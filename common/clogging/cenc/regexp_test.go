package cenc

import (
	"regexp"
	"testing"
)

var regexpFormatFeiBuHuo = regexp.MustCompile(`Windows(?:95|98|NT|2000)`)

var regexpFormaBuHuo = regexp.MustCompile(`Windows(95|98|NT|2000)`)

func TestBuHuoAndFeiBuHuo(t *testing.T) {
	str := "Windows95 Windows98 WindowsNT Windows2000 WindowsXP"

	matchesBuHuo := regexpFormaBuHuo.FindAllStringSubmatchIndex(str, -1)

	t.Log("捕获结果：", matchesBuHuo)

	matchesFeiBuHuo := regexpFormatFeiBuHuo.FindAllStringSubmatchIndex(str, -1)

	t.Log("非捕获结果：", matchesFeiBuHuo)
}

func BenchmarkBuhuo(b *testing.B) {
	str := "Windows95 Windows98 WindowsNT Windows2000 WindowsXP"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		regexpFormaBuHuo.FindAllStringSubmatchIndex(str, -1)
	}
	b.StopTimer()
}

func BenchmarkFeiBuhuo(b *testing.B) {
	str := "Windows95 Windows98 WindowsNT Windows2000 WindowsXP"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		regexpFormatFeiBuHuo.FindAllStringSubmatchIndex(str, -1)
	}
	b.StopTimer()
}
