// Package enums 定义了应用程序中使用的枚举类型
package enums

// ProxyType 表示代理服务的类型
type ProxyType string

const (
	// ServiceDocker Docker 镜像加速服务
	ServiceDocker ProxyType = "docker"

	// ServiceGit Git 仓库加速服务
	ServiceGit ProxyType = "git"
)

// String 返回代理类型的字符串表示
func (s ProxyType) String() string {
	return string(s)
}

// IsValid 检查代理类型是否有效
func (s ProxyType) IsValid() bool {
	switch s {
	case ServiceDocker, ServiceGit:
		return true
	default:
		return false
	}
}

// GetAllTypes 返回所有支持的代理类型
func GetAllTypes() []ProxyType {
	return []ProxyType{ServiceDocker, ServiceGit}
}
