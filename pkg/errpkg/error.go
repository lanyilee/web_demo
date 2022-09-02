package errpkg

// ErrorMsg the struct of error message
type ErrorMsg struct {
	Code    int    `json:"err_code"`
	Message string `json:"message"`
}

func (e ErrorMsg) Error() string {
	return e.Message
}

func NewAPIError(code int, err error) ErrorMsg {
	return ErrorMsg{
		Code:    code,
		Message: err.Error(),
	}
}

var (
	APIForbidden = ErrorMsg{
		Code:    403,
		Message: "权限不足，请联系管理员",
	}
	APINotFound = ErrorMsg{
		Code:    400,
		Message: "不存在",
	}
	APIParamError = ErrorMsg{
		Code:    400,
		Message: "参数错误",
	}
)
