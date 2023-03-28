---
sort: 1
---

# prometheus

## 简介

[chainer](https://github.com/232425wxy/chainer/tree/main/common/metrics/prometheus) 基于 `prometheus` 开发自己的指标 (metrics)，`prometheus` 是一种时间序列数据库，下面将给出一些基本定义，以帮助理解消化 `chainer` 的 `metrics`。

**样本**

`Prometheus` 会将所有采集到的样本数据以时间序列（time-series）的方式保存在内存数据库中，并且定时保存到硬盘上。`time-series` 是按照时间戳和值的序列顺序存放的，我们称之为向量 (vector)。每条`time-series` 通过指标名称 (metrics name) 和一组标签集 (labelset) 命名。如下所示，可以将`time-series` 理解为一个以时间为 `Y` 轴的数字矩阵：

```go
  ^
  │   . . . . . . . . . . . . . . . . .   . .   node_cpu{cpu="cpu0",mode="idle"}
  │     . . . . . . . . . . . . . . . . . . .   node_cpu{cpu="cpu0",mode="system"}
  │     . . . . . . . . . .   . . . . . . . .   node_load1{}
  │     . . . . . . . . . . . . . . . .   . .  
  v
    <------------------ 时间 ---------------->
```

在 `time-series` 中的每一个点称为一个样本（sample），样本由以下三部分组成：

- 指标 (metric)：指标名称和描述当前样本特征的标签集;
- 时间戳(timestamp)：一个精确到毫秒的时间戳;
- 样本值(value)： 一个 `float64` 的浮点型数据表示当前样本的值。

```go
<--------------- metric ---------------------><-timestamp -><-value->
http_request_total{status="200", method="GET"}@1434417560938 => 94355
http_request_total{status="200", method="GET"}@1434417561287 => 94334

http_request_total{status="404", method="GET"}@1434417560938 => 38473
http_request_total{status="404", method="GET"}@1434417561287 => 38544

http_request_total{status="200", method="POST"}@1434417560938 => 4748
http_request_total{status="200", method="POST"}@1434417561287 => 4785
```

上面的例子给出了三组指标名称与时间戳相同、标签集不同的时间序列，可以看到不同时间序列采样的值是不一样的。

**指标 (metrics)**

在形式上，所有的指标 (metrics) 都通过如下格式标示：

```markdown
<metrics name>{<label name>=<label value>, ...}
```

指标的名称 (metrics name) 可以反映被监控样本的含义（比如，http_request_total - 表示当前系统接收到的 HTTP 请求总量）。

标签 (label) 反映了当前样本的特征维度，通过这些维度 `prometheus` 可以对样本数据进行过滤，聚合等。标签的值则可以包含任何Unicode编码的字符。