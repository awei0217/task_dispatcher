package common

/**
远程调用返回的结果体
*/
type Result struct {
	Code int `json:"code"`

	Message string `json:"message"`
}
