// Package models 定义了应用程序中使用的数据模型
package models

import (
	"fmt"
	"strings"
)

// ProxyItem 表示一个代理服务项
// 包含代理服务的详细信息，用于加速网络请求
type ProxyItem struct {
	// ID 代理服务的唯一标识符
	ID string `json:"id"`

	// UseType 使用类型，描述代理的用途
	UseType int `json:"useType"`

	// ProxyUrl 代理服务的URL地址
	ProxyUrl string `json:"proxyUrl"`

	// Name 代理服务的名称
	Name string `json:"name"`

	// Score 代理服务的评分，用于选择最优代理
	Score int `json:"score"`

	// ProxyType 代理类型，如 "docker" 或 "git"
	ProxyType string `json:"proxyType"`
}

// IsValid 检查代理项是否有效
func (p *ProxyItem) IsValid() bool {
	return p.ID != "" &&
		p.ProxyUrl != "" &&
		p.Name != "" &&
		p.Score >= 0
}

// GetDisplayName 获取代理项的显示名称
func (p *ProxyItem) GetDisplayName() string {
	if p.Name != "" {
		return p.Name
	}
	return p.ID
}

// GetScoreDescription 获取评分的描述信息
func (p *ProxyItem) GetScoreDescription() string {
	switch {
	case p.Score >= 90:
		return "优秀"
	case p.Score >= 70:
		return "良好"
	case p.Score >= 50:
		return "一般"
	default:
		return "较差"
	}
}

// String 返回代理项的字符串表示
func (p *ProxyItem) String() string {
	return fmt.Sprintf("ProxyItem{ID: %s, Name: %s, URL: %s, Score: %d, Type: %s}",
		p.ID, p.Name, p.ProxyUrl, p.Score, p.ProxyType)
}

// ProxyList 代理服务列表
type ProxyList struct {
	Items []ProxyItem `json:"items"`
	Total int         `json:"total"`
}

// AddItem 添加代理项到列表
func (pl *ProxyList) AddItem(item ProxyItem) {
	pl.Items = append(pl.Items, item)
	pl.Total++
}

// GetBestProxy 获取评分最高的代理
func (pl *ProxyList) GetBestProxy() *ProxyItem {
	if len(pl.Items) == 0 {
		return nil
	}

	best := &pl.Items[0]
	for i := 1; i < len(pl.Items); i++ {
		if pl.Items[i].Score > best.Score {
			best = &pl.Items[i]
		}
	}
	return best
}

// FilterByType 根据类型过滤代理列表
func (pl *ProxyList) FilterByType(proxyType string) *ProxyList {
	filtered := &ProxyList{}
	for _, item := range pl.Items {
		if strings.EqualFold(item.ProxyType, proxyType) {
			filtered.AddItem(item)
		}
	}
	return filtered
}
