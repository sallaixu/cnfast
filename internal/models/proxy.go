package models

// 请求服务列表
type ProxyItem struct {
	ID        string `json:"id"`
	UseType   string `json:"useType"`
	ProxyUrl  string `json:"proxyUrl"`
	Name      string `json:"name"`
	Score     int    `json:"score"`
	ProxyType string `json:"proxyType"`
}
