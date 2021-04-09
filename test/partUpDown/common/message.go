package common

type Message struct {
	Type string `json:"type"`  // 消息类型
	Data string `json:"data"`  // 消息内容
	AddTime string `sql:"add_time"`  // 演示 sql tag的用法
}

type FileUpMessage struct {
	Type string     `json:"type"`
	Status string   `json:"status"`
	FilePath string `json:"file_path"`
	FileSize int    `json:"file_size"`
	FileHash string `json:"file_hash"`
	UploadID string `json:"upload_id"`
	Data FileUpData
}

type FileUpData struct {
	PartNo   int    `json:"part_no"`
	Len      int    `json:"len"`
	ByteData []byte `json:"byte_data"`
}

type MultipleUploadInfo struct {
	FilePath   string `json:"file_path"`
	FileHash   string `json:"file_hash"`
	FileSize   int    `json:"file_size"`
	UploadID   string `json:"upload_id"`
	ChunkSize  int    `json:"chunk_size"`
	ChunkCount int    `json:"chunk_count"`
	LocalFile *LocalFile `json:"local_file"`
	UpFilePart []FilePart `json:"done_file_part"`
}

//type MessPartFileMerge struct {
//	UploadID string `json:"upload_id"`
//}
type MultipleUploadPartRes struct {
	UploadID string `json:"upload_id"`
	PartNo int      `json:"file_part_no"`
	UpStatus bool   `json:"up_status"`
}
type MultipleUploadMergeRes struct {
	UploadID string  `json:"upload_id"`
	MergeStatus bool `json:"merge_status"`
}

//filePart 文件分片
type FilePart struct {
	Index int    // 文件分片的序号
	From  int    // 开始byte
	To    int    // 结束byte
	Done  bool   // 是否已经上传成功
}

// SumConfig 计算文件摘要值配置
type SumConfig struct {
	IsMD5Sum      bool
	IsSliceMD5Sum bool
	IsCRC32Sum    bool
}

