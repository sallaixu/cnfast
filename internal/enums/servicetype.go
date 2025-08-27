package enums

type ProxyType string

const (
	ServiceDocker ProxyType = "docker"
	ServiceGit    ProxyType = "git"
)

func (s ProxyType) String() string {
	return string(s)
}
