---
sort: 2
---

# namer.go

时间序列数据库 (prometheus) 通过指标名称 (metrics name) 以及对应的一组标签 (labels) 唯一定义一条时间序列。指标名称 (metrics name) 反映了监控样本的基本表示，而标签 (labels) 则在这个基本特征上为采集到的数据提供了多种特征维度。用户可以基于这些特征维度过滤、聚合、统计从而产生新的计算后的一条时间序列。所以，标签对与监控样本采集来说非常重要。

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

`Namer` 包含了指标名称 (namespace、subsystem、name) 和标签集 (labelNames)。`Namer` 的 `FullyQualifiedName` 方法用来获取完整的指标名称 (namespace.subsystem.name or namespace.name or subsystem.name or name)，该方法的定义如下所示：

**FullyQualifiedName() string**
```go
func (n *Namer) FullyQualifiedName() string {
	switch {
	case n.namespace != "" && n.subsystem != "":
		return strings.Join([]string{n.namespace, n.subsystem, n.name}, ".")
	case n.namespace != "":
		return strings.Join([]string{n.namespace, n.name}, ".")
	case n.subsystem != "":
		return strings.Join([]string{n.subsystem, n.name}, ".")
	default:
		return n.name
	}
}
```

另外，一个样本的完整名称不光包括指标名称，还包括标签名，`Namer` 的 [Format](https://github.com/232425wxy/chainer/blob/main/common/metrics/namer/namer.go#L74) 方法可以获取样本的完整名称，下面给一个例子：

比如说，我们声明一个 `Namer`：
```go
&Namer{
	namespace:  "namespace",
	subsystem:  "subsystem",
	name:       "name",
	nameFormat: "%{#namespace}.%{p2p}",
	labelNames: map[string]struct{}{"p2p": {}},
}
```

然后指定一个标签集：
```
labelValues := []string{"p2p", "gossip"}
```

现在将 `labelValues` 作为 `Format` 方法的输入，我们可以得到字符串 `namespace.gossip`。在上述过程中，我们还需要利用正则表达式去匹配 `Namer` 的 `nameFormat` 字段，这部分在[笔记](https://232425wxy.github.io/chainer/Chinese/packages/common/metrics/namer/notes.html)里进行介绍。