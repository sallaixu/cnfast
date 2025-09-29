package util

import "regexp"

// 网址提取出主机名
func ExtractHostFromURL(url string) string {
	re := regexp.MustCompile(`https?://([^/]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
