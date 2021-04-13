package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"seckillsystem/test/partUpDown/client/process"
	"seckillsystem/test/partUpDown/common"
	"strconv"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("listen failed! error: ", err)
		return
	}
	defer listener.Close()

	go WaitExitSignal()

	for {
		log.Println("waiting for connecting from client...")
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept connect failed! error: ", err)
			break
		}

		log.Println("a new client connect to server")
		go ProcessConn(conn)
	}
}

func ProcessConn(conn net.Conn) {
	defer conn.Close()

	tf := &common.Transfer{
		Conn: conn,
	}

	upMgr := process.NewUploadManager()

	closeConn := false
	for {
		log.Println("waiting for message from client...")
		mess, err := tf.ReadPkg()
		if err != nil {
			log.Println("read data from client failed! error: ", err)

			// 判断错误是否是客户端关闭了连接
			if err == io.EOF || strings.Contains(err.Error(), "close") {
				closeConn = true
				log.Println("Client has closed the connection")
			} else {
				continue
			}
		}

		if closeConn {
			log.Println("Server will close the connection")
			break
		}

		if mess.Type == common.MessFileUp {
			upMess, err := upMgr.ParseUpMess(mess.Data)
			if err != nil {
				log.Println("Parse upload message failed! error: ", err)
				break
			}

			if upMess.Type == common.SingleUpload {
				// 文件较小，直接一次性上传所有数据
			} else if upMess.Type == common.MultipleUpload {

				if upMess.Status == common.MultipleUploadConfirm {

					upInfo, err := upMgr.GenerateUpInfo(upMess)
					if err != nil {
						log.Println("generate upload info failed! error: ", err)
						break
					}

					info, err := json.Marshal(upInfo)
					if err != nil {
						log.Println("json marshal upload info failed! error: ", err)
						break
					}
					mess := &common.Message{
						Type:    common.MessFileUpConfirmRes,
						Data:    string(info),
						AddTime: "",
					}
					data, err := json.Marshal(mess)
					if err != nil {
						log.Println("json marshal upload message failed! error: ", err)
						break
					}

					err = tf.WritePkg(data)
					if err != nil {
						log.Println("response message to client failed! error: ", err)
						break
					}

					// 创建临时目录 用来临时存储分片的文件数据
					dirPath := common.DefaultUploadTmpBasePath + upInfo.UploadID

					_, err = os.Stat(dirPath)
					if err != nil {
						err = os.Mkdir(dirPath, os.ModeAppend | os.ModeDir)
						if err != nil {
							log.Println("mkdir upload dir failed! error: ", err)
						} else {
							log.Println("mkdir upload dir success!")
						}
					}

					log.Println("response message to client success!")

				} else if upMess.Status == common.MultipleUploading {

					upRes := &common.MultipleUploadPartRes{
						UploadID: upMess.UploadID,
						PartNo:   upMess.Data.PartNo,
						UpStatus: false,
					}

					err := SaveFilePartDataToTmpFile(upMess)
					if err == nil {
						upRes.UpStatus = true
					}

					upResData, err := json.Marshal(upRes)
					if err != nil {
						log.Printf("json marshal failed! part[%d]\n", upMess.Data.PartNo)
						break
					}

					mess := common.Message{
						Type:    common.MessFileUploadingRes,
						Data:    string(upResData),
						AddTime: "",
					}

					data, err := json.Marshal(mess)
					if err != nil {
						log.Printf("json marshal failed! part[%d]\n", upMess.Data.PartNo)
						break
					}

					err = tf.WritePkg(data)
					if err != nil {
						log.Printf("response result to client failed! part[%d]\n", upMess.Data.PartNo)
						break
					}

					if upRes.UpStatus {
						log.Printf("file-data write to tmp-file success! part[%d]\n", upMess.Data.PartNo)
					}

				} else if upMess.Status == common.MultipleUploadMerge {
					mergeRes := &common.MultipleUploadMergeRes{
						UploadID:    upMess.UploadID,
						MergeStatus: false,
					}

					// 合并文件
					err = MergePartFile(upMess)
					if err == nil {
						mergeRes.MergeStatus = true
						log.Println("merge file success!")

						// 合并成功后删除所有相关的临时文件包括目录
						err = clearTmpFiles(upMess)
					} else {
						log.Println("merge file failed! error: ", err)
					}

					mergeData, err := json.Marshal(mergeRes)
					if err != nil {
						log.Printf("json marshal failed! uploadID[%s]\n", upMess.UploadID)
						break
					}

					mess := &common.Message{
						Type:    common.MessFileUploadMergeRes,
						Data:    string(mergeData),
						AddTime: "",
					}
					data, err := json.Marshal(mess)
					if err != nil {
						log.Printf("json marshal failed! uploadID[%s]\n", upMess.UploadID)
						break
					}

					err = tf.WritePkg(data)
					if err != nil {
						log.Printf("write data to client failed! uploadID[%s]\n", upMess.UploadID)
						break
					}

				} else {
					log.Println("unknown upload status: ", upMess.Status)
					break
				}

			} else {
				log.Println("unknown upload type: ", upMess.Type)
				break
			}
		}

		//fmt.Printf("data from client: %v\n", mess)
	}
}

func MergePartFile(upMess *common.FileUpMessage) (err error) {
	split := strings.Split(upMess.FilePath, ".")
	saveFilePath := common.DefaultUploadBasePath + "up_" + upMess.UploadID + "." + split[len(split) - 1]
	outFile, err := os.OpenFile(saveFilePath, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("create output file failed! UploadID[%s]; error: %v\n", upMess.UploadID, err)
		return
	}
	defer outFile.Close()

	tmpDir := common.DefaultUploadTmpBasePath + upMess.UploadID + "/"
	dirs, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		log.Printf("read tmp dir failed! UploadID[%s]\n", upMess.UploadID)
		return
	}

	for _, f := range dirs {
		log.Println(f.Name())
		tmpFile, err := os.Open(tmpDir + f.Name())
		if err != nil {
			log.Printf("open tmp file failed! UploadID[%s]\n", upMess.UploadID)
			tmpFile.Close()
			return err
		}

		var buf = make([]byte, common.DefaultUpChunkSize)
		n1, err := tmpFile.Read(buf)
		//bytes, err := ioutil.ReadAll(tmpFile)
		if err != nil {
			log.Printf("read tmp file failed! UploadID[%s]\n", upMess.UploadID)
			tmpFile.Close()
			return err
		}
		tmpFile.Close()

		n2, err := outFile.Write(buf[:n1])
		if err != nil || n2 != n1 {
			log.Printf("append data to output-file failed! UploadID[%s]\n", upMess.UploadID)
			return err
		}
	}
	return
}

// 删除临时目录
func clearTmpFiles(upMess *common.FileUpMessage) (err error) {
	tmpDir := common.DefaultUploadTmpBasePath + upMess.UploadID
	return os.RemoveAll(tmpDir)
}

func SaveFilePartDataToTmpFile(upMess *common.FileUpMessage) (err error) {
	if len(upMess.Data.ByteData) != upMess.Data.Len {
		log.Printf("receive file data failed! part[%d] error: file data length is not eq.\n", upMess.Data.PartNo)
		return errors.New("upload file data is wrong")
	}

	filePath := common.DefaultUploadTmpBasePath + upMess.UploadID + "/tmp_" + strconv.FormatInt(int64(upMess.Data.PartNo), 10)
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("create tmp file failed! part[%d], error: %v\n", upMess.Data.PartNo, err)
		return
	}
	defer file.Close()

	n, err := file.Write(upMess.Data.ByteData)
	if err != nil || n != upMess.Data.Len {
		log.Printf("write file-data to tmp-file failed! part[%d]\n", upMess.Data.PartNo)
		return
	}
	return
}

func WaitExitSignal() {
	ch := make(chan os.Signal)

	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch

	log.Println("got exit signal!")

	// 退出前的处理工作

	log.Println("process will exit!")
	os.Exit(1)
}
