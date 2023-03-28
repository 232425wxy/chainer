---
sort: 3
---

# 笔记

## 正则表达式：捕获与非捕获

在许多正则表达式里，能够见到 `?:`，这个表达式表示的含义就是非捕获匹配，网上给了很多案例去解释它的用途，但是感觉都解释的不是很明白，下面我们将给出一个利用 `Go` 语言编写的案例，来揭开它神秘的面纱。

我们先给出一个非捕获匹配表达式：`Windows(?:95|98|NT|2000)`。

再给出一个捕获匹配表达式：`Windows(95|98|NT|2000)`。

以上两个表达式都能匹配字符串 `"Windows95 Windows98 WindowsNT Windows2000 WindowsXP"` 里的 `"Windows95"`、`"Windows98"`、`"WindowsNT"` 和 `"Windows2000"`。但是区别在哪里呢？看下面的测试代码：

```go
func TestBuHuoAndFeiBuHuo(t *testing.T) {
	str := "Windows95 Windows98 WindowsNT Windows2000 WindowsXP"

	matchesBuHuo := regexpFormaBuHuo.FindAllStringSubmatchIndex(str, -1)

	t.Log("捕获结果：", matchesBuHuo)

	matchesFeiBuHuo := regexpFormatFeiBuHuo.FindAllStringSubmatchIndex(str, -1)

	t.Log("非捕获结果：", matchesFeiBuHuo)
}
```

上面代码的输出结果如下：
```
regexp_test.go:17: 捕获结果： [[0 9 7 9] [10 19 17 19] [20 29 27 29] [30 41 37 41]]
regexp_test.go:21: 非捕获结果： [[0 9] [10 19] [20 29] [30 41]]
```

分析所得的结果：

- str[0:9] $\rightarrow$ `"Windows95"`，str[7:9] $\rightarrow$ `"95"`
- str[10:19] $\rightarrow$ `"Windows98"`，str[17:19] $\rightarrow$ `"98"`
- str[20:29] $\rightarrow$ `WindowsNT`，str[27:29] $\rightarrow$ `"NT"`
- str[30:41] $\rightarrow$ `"Windows2000"`，str[37:41] $\rightarrow$ `"2000"`

可以看到，非捕获匹配输出的结果比捕获匹配输出的结果从内容上来说少了一些，即非捕获匹配结果没有保留匹配到的 `"95"`、`"98"` 等字符串。

另外，由于非捕获匹配没有额外的存储匹配到的非捕获分组，所以，非捕获的匹配性能要更好：

```go
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
```

```
BenchmarkBuhuo-96                 907278              1293 ns/op             371 B/op          5 allocs/op
BenchmarkFeiBuhuo-96              973719              1263 ns/op             306 B/op          5 allocs/op
```