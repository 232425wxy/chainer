---
sort: 3
---

# formatter.go

一条日志记录由 `color`、`id`、`level`、`message`、`module`、`shortfunc`、`time` 七个部分组成，其中 `id` 表示日志条目的序号。

`formatter.go` 文件内部则定义了对上述日志条目所包含的七个部分的格式化方法。例如，`time` 采用 `2006-01-02T15:04:05.999Z07:00` 格式进行输出。