# ntp

## 介绍
基于Golang编写，在天泽物联网提供校时服务。

本机地址是：
```text
0x2141000000000401
```

## 服务
服务号|服务
-|-
1|读取时间

### 读取时间服务
- CON请求：空或者带符号的1个字节。

当CON请求为空时，则默认为读取的是北京时间（时区8）。

也可以带1个字节表示时区号。这个字节是有符号的int8。

小技巧，可以使用0x100减去正值即负值。比如8对应的无符号数是0x100-8=248。

- ACK应答：当前时间的字符串

当前时间字符串的格式：2006-01-02 15:04:05 -0700 MST

## 示例
### 读取时间
```go
resp, err := tziot.Call(pipe, 0x2141000000000004, 1, 1000, []uint8{})
fmt.Println("err:", err, "time:", string(resp))
```

输出：
```text
err: 0 time: 2021-03-20 06:34:18 +0800 CST
```

### 读取时区为2的时间
```go
resp, err := tziot.Call(pipe, 0x2141000000000004, 1, 1000, []uint8{2})
fmt.Println("err:", err, "time:", string(resp))
```

输出：
```text
err: 0 time: 2021-03-20 00:36:31 +0200 CST
```

### 读取时区为-6的时间
```go
resp, err := tziot.Call(pipe, 0x2141000000000004, 1, 1000, []uint8{0x100-6})
fmt.Println("err:", err, "time:", string(resp))
```

输出：
```text
err: 0 time: 2021-03-19 16:36:06 -0600 CST
```