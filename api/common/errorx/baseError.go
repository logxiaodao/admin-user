package errorx

type CodeError struct {
	Code int         `json:"code"`
	Msg  Msg         `json:"msg"`
	Data interface{} `json:"data"`
}

type CodeErrorResponse struct {
	Code int         `json:"code"`
	Msg  Msg         `json:"msg"`
	Data interface{} `json:"data"`
}

type SuccessResponse struct {
	Code int         `json:"code"`
	Msg  Msg         `json:"msg"`
	Data interface{} `json:"data"`
}

type Msg struct {
	En string `json:"en"`
	Zh string `json:"zh"`
}

func SendCodeError(code int, msg Msg) error {
	return &CodeError{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

type ErrorCustom interface {
	Error() string
	ErrorZh() string
}

// 实现接口，默认使用中文
func (e *CodeError) Error() string {
	return e.Msg.En
}

func (e *CodeError) ErrorZh() string {
	return e.Msg.Zh
}

func (m *Msg) Error() string {
	return m.En
}

func (m *Msg) ErrorZh() string {
	return m.Zh
}

func (e *CodeError) GetData() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
		Data: nil,
	}
}

// SendSuccess 快捷方法：成功
func SendSuccess(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Code: Success,
		Msg:  DefaultCodeMessage[Success],
		Data: data,
	}
}

// SendLogicalError 快捷方法：业务逻辑错误
func SendLogicalError(msg Msg) error {
	code := FindCode(msg)
	if code != -1 {
		return SendCodeError(code, msg)
	} else {
		return SendCodeError(LogicalError, msg)
	}
}

// SendParameterError 快捷方法：参数错误
func SendParameterError(msg Msg) error {
	code := FindCode(msg)
	if code != -1 {
		return SendCodeError(code, msg)
	} else {
		return SendCodeError(ParameterError, msg)
	}
}

// SendPermissionError 快捷方法：权限错误
func SendPermissionError(msg Msg) error {
	code := FindCode(msg)
	if code != -1 {
		return SendCodeError(code, msg)
	} else {
		return SendCodeError(PermissionError, msg)
	}
}

// SendDatabaseError 快捷方法：数据库错误
func SendDatabaseError(msg Msg) error {
	code := FindCode(msg)
	if code != -1 {
		return SendCodeError(code, msg)
	} else {
		return SendCodeError(DatabaseError, msg)
	}
}

// SendLimitError 快捷方法：超出限制
func SendLimitError(msg Msg) error {
	code := FindCode(msg)
	if code != -1 {
		return SendCodeError(code, msg)
	} else {
		return SendCodeError(LimitError, msg)
	}
}

// SendNetworkError 快捷方法：网络错误
func SendNetworkError(msg Msg) error {
	code := FindCode(msg)
	if code != -1 {
		return SendCodeError(code, msg)
	} else {
		return SendCodeError(NetworkError, msg)
	}
}

// SendThirdPartyError 快捷方法：第三方调用错误
func SendThirdPartyError(msg Msg) error {
	code := FindCode(msg)
	if code != -1 {
		return SendCodeError(code, msg)
	} else {
		return SendCodeError(ThirdPartyError, msg)
	}
}

// SendIoError 快捷方法：io错误
func SendIoError(msg Msg) error {
	code := FindCode(msg)
	if code != -1 {
		return SendCodeError(code, msg)
	} else {
		return SendCodeError(IoError, msg)
	}
}

// SendServiceError 快捷方法：服务错误
func SendServiceError(msg Msg) error {
	code := FindCode(msg)
	if code != -1 {
		return SendCodeError(code, msg)
	} else {
		return SendCodeError(ServiceError, msg)
	}
}

// FindCode 通过msg找code， 找到则使用默认或自定义的code
func FindCode(msg Msg) (Code int) {

	// 匹配默认
	for k, v := range DefaultCodeMessage {
		if msg == v {
			return k
		}
	}

	// 匹配自定义
	for k, v := range CodeMessage {
		if msg == v {
			return k
		}
	}

	return -1
}

// FindCodeByMsg 通过message找code， 找到则使用默认或自定义的code
func FindCodeByMsg(msg string) (Code int) {

	// 匹配默认
	for k, v := range DefaultCodeMessage {
		if msg == v.Error() || msg == v.ErrorZh() {
			return k
		}
	}

	// 匹配自定义
	for k, v := range CodeMessage {
		if msg == v.Error() || msg == v.ErrorZh() {
			return k
		}
	}

	return -1
}
