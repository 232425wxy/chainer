---
sort: 1
---

# metrics 包

`chainer` 提供了三种指标类型来区分不同的监控指标：`Counter`、`Gauge`、`Histogram`。

## Counter：只增不减的计数器

Counter 类型的指标其工作方式和计数器一样，只增不减（除非系统发生重置）。常见的监控指标，如http_requests_total，[^1]node_cpu 都是 Counter 类型的监控指标。 一般在定义 Counter 类型指标的名称时推荐使用 `_total` 作为后缀。




[namer](https://232425wxy.github.io/chainer/Chinese/packages/common/metrics/2.namer.html)

[^1]: cpu 的累积使用时间