---
sort: 2
---

# namer.go

`namer.go` 文件里定义了 `Namer` 结构体，该结构体是 `namer.go` 文件的核心，定义如下：

```go
type Namer struct {
	namespace  string
	subsystem  string
	name       string
	nameFormat string
	labelNames map[string]struct{}
}
```

时间序列数据库 (prometheus) 通过指标名称 (metrics name) 以及对应的一组标签 (labels) 唯一定义一条时间序列。指标名称 (metrics name) 反映了监控样本的基本表示，而标签 (labels) 则在这个基本特征上为采集到的数据提供了多种特征维度。用户可以基于这些特征维度过滤、聚合、统计从而产生新的计算后的一条时间序列。


