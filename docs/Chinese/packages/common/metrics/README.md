---
sort: 1
---

# metrics 包

`chainer` 提供了三种指标类型来区分不同的监控指标：`Counter`、`Gauge`、`Histogram`。

## Counter：只增不减的计数器

Counter 类型的指标其工作方式和计数器一样，只增不减（除非系统发生重置）。常见的监控指标，如http_requests_total，[^1]node_cpu 都是 Counter 类型的监控指标。 一般在定义 Counter 类型指标的名称时推荐使用 `_total` 作为后缀。

我们可以在应用程序中利用 Counter 记录某些事件发生的次数，通过以时序的形式存储这些数据，我们可以轻松的了解该事件产生速率的变化。

## Gauge：可增可减的仪表盘

与 Counter 不同，Gauge 类型的指标侧重于反应系统的当前状态。因此这类指标的样本数据可增可减。常见指标如：node_memory_MemFree（主机当前空闲的内容大小）、node_memory_MemAvailable（可用内存大小）都是 Gauge 类型的监控指标。

## Histogram：反映数据分布情况的工具

在大多数情况下人们都倾向于使用某些量化指标的平均值，例如CPU的平均使用率、页面的平均响应时间。这种方式的问题很明显，以系统API调用的平均响应时间为例：如果大多数API请求都维持在100ms的响应时间范围内，而个别请求的响应时间需要5s，那么就会导致某些WEB页面的响应时间落到中位数的情况，而这种现象被称为长尾问题。

为了区分是平均的慢还是长尾的慢，最简单的方式就是按照请求延迟的范围进行分组。例如，统计延迟在0~10ms之间的请求数有多少而10~20ms之间的请求数又有多少。通过这种方式可以快速分析系统慢的原因。Histogram 就是能够解决这样问题的存在，通过 Histogram 监控指标，我们可以快速了解监控样本的分布情况。

比如，我们调用了三次 `Histogram` 的 `Observe` 功能：

```go
histogram.With("http", "www.baidu.com", "https", "github.com").Observe(0.8) // 大于 0.8 的 bucket 只有 {1, 2.5, 5, 10}
	histogram.With("http", "www.baidu.com", "https", "github.com").Observe(2.8) // 大于 2.8 的 bucket 只有 {5, 10}
	histogram.With("http", "www.baidu.com", "https", "github.com").Observe(3.1) // 大于 3.1 的 bucket 只有 {5, 10}
```

得到如下结果：

```shell
histogram_namespace_histogram_subsystem_histogram_name_bucket{http="www.baidu.com",https="github.com",le="0.005"} 0
histogram_namespace_histogram_subsystem_histogram_name_bucket{http="www.baidu.com",https="github.com",le="0.01"} 0
histogram_namespace_histogram_subsystem_histogram_name_bucket{http="www.baidu.com",https="github.com",le="0.025"} 0
histogram_namespace_histogram_subsystem_histogram_name_bucket{http="www.baidu.com",https="github.com",le="0.05"} 0
histogram_namespace_histogram_subsystem_histogram_name_bucket{http="www.baidu.com",https="github.com",le="0.1"} 0
histogram_namespace_histogram_subsystem_histogram_name_bucket{http="www.baidu.com",https="github.com",le="0.25"} 0
histogram_namespace_histogram_subsystem_histogram_name_bucket{http="www.baidu.com",https="github.com",le="0.5"} 0
histogram_namespace_histogram_subsystem_histogram_name_bucket{http="www.baidu.com",https="github.com",le="1"} 1
histogram_namespace_histogram_subsystem_histogram_name_bucket{http="www.baidu.com",https="github.com",le="2.5"} 1
histogram_namespace_histogram_subsystem_histogram_name_bucket{http="www.baidu.com",https="github.com",le="5"} 3
histogram_namespace_histogram_subsystem_histogram_name_bucket{http="www.baidu.com",https="github.com",le="10"} 3
histogram_namespace_histogram_subsystem_histogram_name_bucket{http="www.baidu.com",https="github.com",le="+Inf"} 3
histogram_namespace_histogram_subsystem_histogram_name_sum{http="www.baidu.com",https="github.com"} 6.699999999999999
histogram_namespace_histogram_subsystem_histogram_name_count{http="www.baidu.com",https="github.com"} 3
```

其中，Histogram 类型的样本以 `_count` 作为后缀的条目反映了记录的总数，以 `_sum` 作为后缀的条目反映了值的总量。其他以 `_bucket` 作为后缀的条目则直接反映了在不同区间内样本的个数。

[^1]: cpu 的累积使用时间