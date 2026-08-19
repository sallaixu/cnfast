package services

import (
	"cnfast/internal/models"
	"cnfast/internal/pkg/util"
	"fmt"
	"strings"
	"testing"
)

// proxyRow 构建表格一行（与 renderProxyTable 相同的列宽与 pad 规则）
func proxyRow(idx int, proxy models.ProxyItem) string {
	return strings.Join([]string{
		util.Pad(fmt.Sprintf("%d", idx), carTableIdxColWidth),
		util.Pad(util.Truncate(proxy.ProxyUrl, carTableURLColWidth-1), carTableURLColWidth),
		util.Pad(fmt.Sprintf("%d", proxy.Score), carTableScoreColWidth),
		proxy.GetScoreDescription(),
	}, " ")
}

func TestProxyRowAlignment(t *testing.T) {
	cases := []models.ProxyItem{
		{ID: "3e58a08", ProxyUrl: "https://cdn.gh-proxy.org", Score: 10},
		{ID: "2da2189", ProxyUrl: "https://edgeone.gh-proxy.org", Name: "edge", Score: 7},
		{ID: "5c767aa", ProxyUrl: "项目地址： https://github.com/sallaixu/cnfast这是一个非常非常长的描述用于测试", Name: "中文名称超级长", Score: 0},
	}

	// 校验各列起始位置一致
	var scorePositions, ratingPositions []int

	for i, p := range cases {
		urlPos := util.VisualLen(util.Pad(fmt.Sprintf("%d", i+1), carTableIdxColWidth) + " ")
		if urlPos != carTableIdxColWidth+1 {
			t.Errorf("地址列起始不符: %d", urlPos)
		}

		scorePos := util.VisualLen(
			util.Pad(fmt.Sprintf("%d", i+1), carTableIdxColWidth) + " " +
				util.Pad(util.Truncate(p.ProxyUrl, carTableURLColWidth-1), carTableURLColWidth) + " ")
		ratingPos := util.VisualLen(
			util.Pad(fmt.Sprintf("%d", i+1), carTableIdxColWidth) + " " +
				util.Pad(util.Truncate(p.ProxyUrl, carTableURLColWidth-1), carTableURLColWidth) + " " +
				util.Pad(fmt.Sprintf("%d", p.Score), carTableScoreColWidth) + " ")

		scorePositions = append(scorePositions, scorePos)
		ratingPositions = append(ratingPositions, ratingPos)
	}

	if !allEqual(scorePositions) {
		t.Errorf("评分列未对齐: %v", scorePositions)
	}
	if !allEqual(ratingPositions) {
		t.Errorf("评级列未对齐: %v", ratingPositions)
	}
}

func allEqual(nums []int) bool {
	for _, n := range nums {
		if n != nums[0] {
			return false
		}
	}
	return true
}

func TestGetScoreDescription10Scale(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{10, "优秀"},
		{9, "优秀"},
		{8, "良好"},
		{7, "良好"},
		{6, "一般"},
		{5, "一般"},
		{4, "较差"},
		{0, "较差"},
	}
	for _, c := range cases {
		p := models.ProxyItem{Score: c.score}
		if got := p.GetScoreDescription(); got != c.want {
			t.Errorf("score=%d evaluate=%s want=%s", c.score, got, c.want)
		}
	}
}
