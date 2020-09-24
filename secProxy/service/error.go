package service

// 错误码
const (
	ErrRequestSuccess      = 1000
	ErrInvalidRequest      = 1001 // 无效的请求
	ErrNotFoundActivityId  = 1002 // 没有找到activity_id对应的活动数据
	ErrParamDeletion       = 1003 // 缺少参数
	ErrParamTypeError      = 1004 // 参数类型错误
	ErrCookieParamDeletion = 1005 // 缺少cookie信息
	ErrUserAuthCheckFailed = 1006 // 用户权限校验未通过
	ErrServiceBusy         = 1007 // 服务器繁忙
)

