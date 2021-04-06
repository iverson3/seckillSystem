package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"path/filepath"
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

	err = ProcessConnWithServer(conn)
	log.Println("client will exit!")
}

func ProcessConnWithServer(conn net.Conn) error {
	defer conn.Close()
	tf := &common.Transfer{
		Conn: conn,
	}

	filePath := "./aaa"
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
		FilePath: filePath,
		FileSize: int(lf.Length),
		FileHash: "",
	}

	data, err := json.Marshal(fum)
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
