package common


/**
	分页对象
 */
type Page struct {
	/**
		状态码 200 成功
	 */
	Code int `json:"code"`
	/**
		执行结果描述，比如成功，失败
	 */
	Message string `json:"message"`
	/**
		总数
	 */
	Total int64  `json:"total"`
	/**
		页数
	 */
	Page int `json:"page"`
	/**
		每页多少条
	 */
	Limit int  `json:"limit"`

	Data interface{} `json:"data"`
}



