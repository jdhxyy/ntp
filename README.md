# ntp

## 介绍
基于Golang编写，在海萤物联网提供校时服务。

本机地址是：
```text
0x2141000000000404
```

## 服务
服务号|服务
---|---
1|读取时间1
2|读取时间2.返回的是结构体

### 读取时间服务1
- CON请求：空或者带符号的1个字节。

当CON请求为空时，则默认为读取的是北京时间（时区8）。

也可以带1个字节表示时区号。这个字节是有符号的int8。

小技巧，可以使用0x100减去正值即负值。比如8对应的无符号数是0x100-8=248。

- ACK应答：当前时间的字符串

当前时间字符串的格式：2006-01-02 15:04:05 -0700 MST

### 读取时间服务2.返回的是结构体
- CON请求：格式与读取时间服务1一致

- ACK应答：
```c
struct {
    // 时区
    uint8 TimeZone
    uint16 Year
    uint8 Month
    uint8 Day
    uint8 Hour
    uint8 Minute
    uint8 Second
    // 星期
    uint8 Weekday
}
```

### 自定义错误码
错误码|含义
---|---
0x40|内部错误
0x41|接收格式错误

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

### 读取结构体格式时间
```go
resp, _ := tziot.Call(pipe, 0x2141000000000004, 2, 3000, []uint8{8})
var ack AckRidGetTime2
_ = dcom.BytesToStruct(resp, &ack)
fmt.Println(ack)
```

输出：
```text
{8 2021 3 30 12 11 28 2}
```