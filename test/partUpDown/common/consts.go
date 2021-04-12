package common

const DefaultUpChunkSize = 8 * 1024 * 1024    // 50MB
//const DefaultUpChunkSize = 200 * 1024 * 1024    // 50MB
//const DefaultUpChunkSize =  8    // 8字节
const DefaultUploadTmpBasePath = "./public/upload/tmp/"
const DefaultUploadBasePath = "./public/upload/"

// 消息类型
const (
	UserRegister = "UserRegister"
	UserLogin    = "UserLogin"
	UserLogout   = "UserLogout"

	MessFileUp   = "MessFileUp"
	MessFileUpConfirmRes = "MessFileUpConfirmRes"
	MessFileUploadingRes = "MessFileUploadingRes"
	MessFileUploadMergeRes = "MessFileUploadMergeRes"

	MessFileDown = "MessFileDown"
)

// 文件上传/下载的类型
const (
	SingleUpload = "SingleUpload"
	MultipleUpload = "MultipleUpload"
	SingleDownload = "SingleDownload"
	MultipleDownload = "MultipleDownload"
)

const (
	MultipleUploadConfirm = "MultipleUploadConfirm"
	MultipleUploading     = "MultipleUploading"
	MultipleUploadMerge   = "MultipleUploadMerge"
)
