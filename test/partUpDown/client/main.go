package main

import (
	"github.com/prometheus/common/log"
	"net"
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
		log.Error("connect to server failed! error: ", err)
		return
	}
	defer conn.Close()

	err = ProcessConnWithServer(conn)
	log.Info("client will exit!")
}

func ProcessConnWithServer(conn net.Conn) error {
	defer conn.Close()
	tf := &common.Transfer{
		Conn: conn,
	}

	err := tf.WritePkg([]byte("hello update"))
	if err != nil {
		log.Error("send data to server failed! error: ", err)
		return err
	}

	log.Info("send data to server success!")
	return nil
}
