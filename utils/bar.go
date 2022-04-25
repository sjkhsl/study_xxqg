package utils

import (
	"fmt"
	"io"
)

type Bar struct {
	percent int64  // 百分比
	cur     int64  // 当前进度位置
	total   int64  // 总进度
	rate    string // 进度条
	graph   string // 显示符号
	io.Reader
}

func (bar *Bar) Read(p []byte) (n int, err error) {
	n, err = bar.Reader.Read(p)
	if err != nil {
		return 0, err
	}
	bar.Play(bar.cur + int64(n))
	return n, err
}

func (bar *Bar) NewOption(start, total int64, reader io.Reader) {
	bar.Reader = reader
	bar.cur = start
	bar.total = total
	if bar.graph == "" {
		bar.graph = "█"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph // 初始化进度条位置
	}
}

func (bar *Bar) getPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}

func (bar *Bar) NewOptionWithGraph(start, total int64, reader io.Reader, graph string) {
	bar.graph = graph
	bar.NewOption(start, total, reader)
}

func (bar *Bar) Play(cur int64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\r[%-50s]%3d%%  %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}
