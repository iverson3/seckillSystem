package common

type Message struct {
	Type string `json:"type"`  // 消息类型
	Data string `json:"data"`  // 消息内容
	AddTime string `sql:"add_time"`  // 演示 sql tag的用法
}

type MultipleUploadInfo struct {
	FileHash   string `json:"file_hash"`
	FileSize   int    `json:"file_size"`
	UploadID   string `json:"upload_id"`
	ChunkSize  int    `json:"chunk_size"`
	ChunkCount int    `json:"chunk_count"`
}

// SumConfig 计算文件摘要值配置
type SumConfig struct {
	IsMD5Sum      bool
	IsSliceMD5Sum bool
	IsCRC32Sum    bool
}