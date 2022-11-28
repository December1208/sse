package api_error

type APIError struct {
	ErrorCode int
	ErrorMsg  string
}

type APIOldError struct {
	ErrorCode int
	ErrorMsg  string
}

var (
	// Basic errors
	NotDefined          = APIError{ErrorCode: 1000, ErrorMsg: "未定义的错误"}
	PermissionDenied    = APIError{ErrorCode: 1001, ErrorMsg: "您没有权限访问"}
	RequiredParamMissed = APIError{ErrorCode: 1002, ErrorMsg: "缺少必要参数"}
	ObjectNotFound      = APIError{ErrorCode: 1004, ErrorMsg: "未找到对象"}
	RequestThrolled     = APIError{ErrorCode: 1005, ErrorMsg: "请求过于频繁"}
	RequestFail         = APIError{ErrorCode: 1006, ErrorMsg: "请求失败"}
	InvalidParam        = APIError{ErrorCode: 1007, ErrorMsg: "无效的请求参数"}
	NotAuthenticated    = APIError{ErrorCode: 1008, ErrorMsg: "未登录"}
	WrongRequestAction  = APIError{ErrorCode: 1009, ErrorMsg: "错误的请求行为"}
	InvalidPage         = APIError{ErrorCode: 1011, ErrorMsg: "无效页"}

	FileLoadFail = APIError{ErrorCode: 2001, ErrorMsg: "文件加载失败"}
)
