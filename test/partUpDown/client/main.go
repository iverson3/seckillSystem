package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"seckillsystem/test/partUpDown/client/process"
	"seckillsystem/test/partUpDown/common"
)

type PartUploader struct {

}

type PartDownloader struct {
	FilePath       string
	FileName       string
	FileSize       int
	PartCount      int
	OutputFilePath string
	OutputFileName string
	DoneFileParts  []*FilePart  // 已下载完成的分片
}

type FilePart struct {
	Index int     // 分片序号
	From  int     // 起始字节
	To    int     // 结束字节
	Data  []byte  // 当前分片的具体数据
}

func main() {
	//http.StatusPartialContent

	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatal("connect to server failed! error: ", err)
		return
	}
	defer conn.Close()

	ProcessConnWithServer(conn)
	log.Println("client will exit!")
}

func ProcessConnWithServer(conn net.Conn) {
	defer conn.Close()
	tf := &common.Transfer{
		Conn: conn,
	}

	go process.WaitServerMess(tf)
	_ = process.NewUploadManager()

	fmt.Println("---------------聊天系统主界面----------------")
	fmt.Println("\t\t 1.显示在线用户列表")
	fmt.Println("\t\t 2.群发消息")
	fmt.Println("\t\t 3.个人私聊")

	fmt.Println("\t\t 4.文件上传")
	fmt.Println("\t\t 5.信息列表")
	fmt.Println("\t\t 6.退出登录")

	fmt.Println("\t\t 7.退出系统")

	var key int
	for {
		fmt.Println("\t\t 请选择 (1-7)：")
		// 需要处理用户输入字符串的情况
		_, err := fmt.Scanf("%d\n", &key)
		if err != nil {
			log.Println("input is wrong! error: ", err)
			break
		}

		switch key {
		case 4: {
			// 让用户选择文件 / 输入文件路径
			filePath := "./ekw029.mp4"
			//filePath := "./ccc.txt"
			//filePath := "./bbb.mp4"
			//filePath := "./aaa"
			err = UploadFile(tf, filePath)
			if err != nil {
				log.Println("upload file failed! error: ", err)
				break
			}
		}
		default:
			log.Println("其他功能暂不支持！")
		}
	}
}

func UploadFile(tf *common.Transfer, filePath string) error {
	filePath, _ = filepath.Abs(filePath)
	fmt.Println(filePath)

	lf := &common.LocalFile{
		LocalFileMeta: common.LocalFileMeta{
			Path: filePath,
		},
	}

	err := lf.OpenPath()
	if err != nil {
		log.Fatal("open file failed! error: ", err)
		return err
	}

	fum := &common.FileUpMessage{
		Type:     common.MultipleUpload,
		Status:   common.MultipleUploadConfirm,
		FilePath: filePath,
		FileSize: int(lf.Length),
		FileHash: "",
	}
	data1, err := json.Marshal(fum)
	if err != nil {
		log.Fatal("json marshal message failed! error: ", err)
		return err
	}

	mess := &common.Message{
		Type:    common.MessFileUp,
		Data:    string(data1),
		AddTime: "",
	}

	data, err := json.Marshal(mess)
	if err != nil {
		log.Fatal("json marshal message failed! error: ", err)
		return err
	}

	err = tf.WritePkg(data)
	if err != nil {
		log.Fatal("send data to server failed! error: ", err)
		return err
	}

	log.Println("send data to server success!")
	return nil
}
