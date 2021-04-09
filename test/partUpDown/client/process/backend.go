package process

import (
	"encoding/json"
	"log"
	"seckillsystem/test/partUpDown/common"
)

// 在后台保持跟服务端之间的通讯
// 随时准备接收服务端发送过来的数据 (比如用户上线下线提醒 消息推送提醒等等)
func WaitServerMess(tf *common.Transfer)  {
	for {
		mess, err := tf.ReadPkg()
		if err != nil {
			log.Println("get message from server failed! error: ", err)
			break
		}

		// 消息类型是 server端对于上传确认请求的回复，里面包含了上传相关的信息
		if mess.Type == common.MessFileUpConfirmRes {
			go processFileUp(tf, mess)
		} else if mess.Type == common.MessFileUploadingRes {
			go UpdateUploadProgress(mess, tf)
		} else if mess.Type == common.MessFileUploadMergeRes {
			var res common.MultipleUploadMergeRes
			err = json.Unmarshal([]byte(mess.Data), &res)
			if err != nil {
				log.Println("unmarshal result data failed! error: ", err)
				return
			}


			if !res.MergeStatus {
				log.Println("server merge file failed!")
				return
			}

			log.Println("server merge file success!")
		} else {

		}
	}
}

func processFileUp(tf *common.Transfer, message common.Message) {
	var upInfo common.MultipleUploadInfo
	err := json.Unmarshal([]byte(message.Data), &upInfo)
	if err != nil {
		log.Println("unmarshal upload info failed! error: ", err)
		return
	}

	upInfo.LocalFile = &common.LocalFile{
		LocalFileMeta: common.LocalFileMeta{
			Path: upInfo.FilePath,
		},
	}
	upInfo.LocalFile.BufSize = upInfo.ChunkSize

	err = upInfo.LocalFile.OpenPath()
	if err != nil {
		log.Println("open file failed! error: ", err)
		return
	}

	err = UploadMgr.SaveUpInfo(&upInfo)
	err = UploadMgr.LoopReadFileData(tf, upInfo)
	if err != nil {
		log.Println("send file data to server failed! error: ", err)
		return
	}

	log.Println("send file data to server success!")
}

// 更新上传进度
func UpdateUploadProgress(mess common.Message, tf *common.Transfer) {

	var res common.MultipleUploadPartRes
	err := json.Unmarshal([]byte(mess.Data), &res)
	if err != nil {
		log.Println("json Unmarshal failed! error: ", err)
		return
	}

	if !res.UpStatus {
		log.Printf("file part data upload failed! part[%d], error: %v\n", res.PartNo, err)
		return
	}

	// 标记该part上传完成
	UploadMgr.UploadingFileMap[res.UploadID].UpFilePart[res.PartNo].Done = true

	completeNum := 0
	for _, part := range UploadMgr.UploadingFileMap[res.UploadID].UpFilePart {
		if part.Done {
			completeNum++
		}
	}

	log.Printf("update upload progress success! cur progress: %d/%d\n", completeNum, UploadMgr.UploadingFileMap[res.UploadID].ChunkCount)

	// 所有文件分片上传完成
	if completeNum == UploadMgr.UploadingFileMap[res.UploadID].ChunkCount {
		log.Println("notify server to merge file...")

		err = NotifyServerMergeFile(res, tf)
		if err != nil {
			log.Println("notify server to merge file failed! error: ", err)
		}
	}
}

func NotifyServerMergeFile(res common.MultipleUploadPartRes, tf *common.Transfer) (err error) {
	mergeMess := &common.FileUpMessage{
		Type: common.MultipleUpload,
		Status: common.MultipleUploadMerge,
		UploadID: res.UploadID,
		FilePath: UploadMgr.UploadingFileMap[res.UploadID].FilePath,
	}
	mergeData, err := json.Marshal(mergeMess)
	if err != nil {
		return
	}

	mess := &common.Message{
		Type:    common.MessFileUp,
		Data:    string(mergeData),
		AddTime: "",
	}
	data, err := json.Marshal(mess)
	if err != nil {
		return
	}

	return tf.WritePkg(data)
}