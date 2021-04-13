package process

import (
	"encoding/json"
	"errors"
	"log"
	"seckillsystem/test/partUpDown/common"
	"time"
)

var UploadMgr *UploadManager

type UploadManager struct {
	UploadingFileMap map[string]*common.MultipleUploadInfo   // 记录着所有还未上传结束的文件上传信息
}

func NewUploadManager() *UploadManager {
	UploadMgr = &UploadManager{
		UploadingFileMap: make(map[string]*common.MultipleUploadInfo),
	}
	return UploadMgr
}

func (upm *UploadManager) ParseUpMess(data string) (mess *common.FileUpMessage, err error) {
	err = json.Unmarshal([]byte(data), &mess)
	if err != nil {
		return
	}
	return
}

func (upm *UploadManager) GenerateUpInfo(mess *common.FileUpMessage) (info *common.MultipleUploadInfo, err error) {
	upId := upm.GenerateUploadID(mess)

	chunkSize  := common.DefaultUpChunkSize
	chunkCount := mess.FileSize / chunkSize + 1

	info = &common.MultipleUploadInfo{
		FilePath:   mess.FilePath,
		FileHash:   "",
		FileSize:   mess.FileSize,
		UploadID:   upId,
		ChunkSize:  chunkSize,
		ChunkCount: chunkCount,
	}
	return
}

func (upm *UploadManager) GenerateUploadID(mess *common.FileUpMessage) string {
	return common.Md5String(mess.FilePath + time.Now().String())
}

func (upm *UploadManager) GetUpInfo(uploadId string) *common.MultipleUploadInfo {
	info, ok := upm.UploadingFileMap[uploadId]
	if !ok {
		return nil
	}
	return info
}

func (upm *UploadManager) SaveUpInfo(info *common.MultipleUploadInfo) error {
	info.UpFilePart = make([]common.FilePart, info.ChunkCount)
	for i := 1; i <= info.ChunkCount; i++ {
		fp := common.FilePart{
			Index: i,
			From:  (i - 1) * info.ChunkSize,
			To:    i * info.ChunkSize - 1,
			Done:  false,
		}
		if i == info.ChunkCount {
			fp.To = int(info.LocalFile.Length - 1)
		}
		info.UpFilePart = append(info.UpFilePart, fp)
	}
	upm.UploadingFileMap[info.UploadID] = info
	return nil
}

func (upm *UploadManager) LoopReadFileData(tf *common.Transfer, upInfo *common.MultipleUploadInfo) (err error) {
	if upInfo.LocalFile.File == nil {
		return errors.New("lf.File is nil")
	}
	upInfo.LocalFile.InitBuf()

	//var fileInfo *common.FileUpMessage
	//var infoData []byte
	var i int64 = 1
	var n = 0
	for ; i <= int64(upInfo.ChunkCount); i++ {
		_, _ = upInfo.LocalFile.File.Seek((i - 1)*int64(upInfo.ChunkSize), 0)

		n, err = upInfo.LocalFile.File.Read(upInfo.LocalFile.Buf)
		if err != nil {
			log.Printf("forloop [%d], read file data failed! error: %v\n", i, err)
			break
		}

		fileData := &common.FileUpData{
			PartNo:   int(i),
			Len:      n,
			ByteData: upInfo.LocalFile.Buf[:n],
		}

		//go sendFileDataToServer(tf, int(i), infoData)
		err = sendFileDataToServer(tf, int(i), upInfo, fileData)
		if err != nil {
			log.Printf("send file data[%d] to server failed! error: %v\n", i, err)
			break
		}
	}

	//fileInfo = nil
	UploadMgr.UploadingFileMap[upInfo.UploadID].LocalFile.Close()
	return err

	//var (
	//	begin int64
	//	n int
	//	no int
	//)
	//no = 1
	//
	//// 读文件
	//for {
	//	log.Println("============")
	//	log.Println(len(lf.Buf))
	//	n, err = lf.File.ReadAt(lf.Buf, begin)
	//	if err != nil {
	//		if err == io.EOF {
	//			go sendFileDataToServer(tf, lf, upInfo, no, n)
	//		} else {
	//			fmt.Println("read file failed! error: ", err)
	//			return
	//		}
	//		begin += int64(n)
	//		break
	//	}
	//	go sendFileDataToServer(tf, lf, upInfo, no, n)
	//	begin += int64(n)
	//	no++
	//}
}

func sendFileDataToServer(tf *common.Transfer, partNo int, upInfo *common.MultipleUploadInfo, fileData *common.FileUpData) (err error) {
	fileInfo := &common.FileUpMessage{
		Type:     common.MultipleUpload,
		Status:   common.MultipleUploading,
		FilePath: upInfo.LocalFile.Path,
		FileSize: int(upInfo.LocalFile.Length),
		FileHash: "",
		UploadID: upInfo.UploadID,
		Data: *fileData,
	}
	infoData, err := json.Marshal(fileInfo)
	if err != nil {
		log.Println("json marshal failed! error: ", err)
		return
	}

	mess := &common.Message{
		Type:    common.MessFileUp,
		Data:    string(infoData),
		AddTime: "",
	}
	data, err := json.Marshal(mess)
	if err != nil {
		return
	}

	err = tf.WritePkg(data)
	if err != nil {
		log.Printf("write file data[%d] to server failed! error: %v\n", partNo, err)
		return
	}

	log.Printf("write file data to server success! part[%d]\n", partNo)
	return
}

//func (um *UploadMgr) NewLocalFileInfo(localPath string, bufSize int) *model.LocalFile {
//	return &model.LocalFile{
//		LocalFileMeta: model.LocalFileMeta{
//			Path: localPath,
//		},
//		BufSize: bufSize,
//	}
//}
