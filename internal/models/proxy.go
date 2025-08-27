package models

// 请求服务列表
type Service struct {
	Type     string `json:"type"`
	Services []struct {
		ID   string `json:"id"`
		Url  string `json:"url"`
		Name string `json:"name"`
	} `json:"services"`
}
