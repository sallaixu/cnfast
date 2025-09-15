package services

import (
	"bufio"
	"cnfast/internal/models"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var reHost = regexp.MustCompile(`^https?://([^/]+)`)

func extractHost(rawURL string) string {
	m := reHost.FindStringSubmatch(rawURL)
	if len(m) < 2 {
		return ""
	}
	return m[1] // 含端口，如 ghproxy.com:8080
}

func GitProxy(proxyList []models.ProxyItem) {
	// 检查git命令合法性
	if len(os.Args) < 3 {
		fmt.Printf("args lenght less than 2 \n")
		os.Exit(1)
	}

	// 获取代理地址
	fmt.Printf("use max score proxy: %s\n", proxyList[0].ProxyUrl)
	proxyPrefix = proxyList[0].ProxyUrl
	proxyHost := extractHost(proxyPrefix)
	fmt.Println(proxyHost)
	// 保留 git 子命令
	newArgs := []string{}
	supportCmd := []string{"clone", "pull", "fetch"}

	found := false
	for _, s := range supportCmd {
		if s == os.Args[2] {
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("not support command %s \n", newArgs[1])
		fmt.Printf("supported command %s \n", strings.Join(supportCmd, ", "))
		os.Exit(1)
	}
	for _, arg := range os.Args[2:] {
		if strings.HasPrefix(arg, "https://github.com/") ||
			strings.HasPrefix(arg, "http://github.com/") {
			arg = proxyPrefix + "/" + arg
		}
		newArgs = append(newArgs, arg)
	}

	cmd := exec.Command(os.Args[1], newArgs...)
	pr, pw := io.Pipe()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = pw

	// ① 先启动读 goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		sc := bufio.NewScanner(pr)
		for sc.Scan() {
			line := sc.Text()
			// 过滤掉含代理 host 的行
			// if strings.Contains(line, proxyHost) {
			// 	continue
			// }
			fmt.Fprintln(os.Stderr, line)
		}
	}()

	// ② 再运行命令
	if err := cmd.Run(); err != nil {
		pw.Close()
		<-done
		os.Exit(1)
	}
	pw.Close()
	<-done
}
