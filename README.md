# SMS TCP Server

一个异步TCP服务器，用于接收SMS消息并通过HTTP转发。服务器接收JSON格式的TCP消息，解析后通过HTTP GET请求转发到指定URL。

## 功能特性

- 异步TCP服务器
- JSON消息解析
- 消息队列处理
- HTTP请求转发
- HTTPS支持
- 优雅关闭
- 详细的日志记录

## 消息格式

服务器接收以下格式的JSON消息：

```json
{
    "txt": "消息内容",
    "num": "手机号码",
    "cmd": "sms",
    "metas": {
        "tz": 32,
        "min": 21,
        "seqNum": 0,
        "refNum": 0,
        "year": 25,
        "sec": 57,
        "maxNum": 0,
        "mon": 4,
        "hour": 15,
        "day": 22
    }
}
```

## 安装

1. 克隆仓库：
```bash
git clone https://github.com/yourusername/sms-tcpserver.git
cd sms-tcpserver
```

2. 编译：
```bash
go build
```

## 使用方法

### 启动服务器

```bash
# 使用默认端口(8080)和指定的HTTP URL
go run main.go -http-url https://test.com/path

# 指定端口和HTTP URL
go run main.go -port 9090 -http-url https://test.com/path

# 指定自定义图标URL
go run main.go -http-url https://test.com/path -icon https://example.com/icon.jpg
```

### 命令行参数

- `-port`: TCP服务器监听端口（默认：8080）
- `-http-url`: HTTP请求的基础URL（必需）
- `-icon`: 图标URL（默认：https://i.ibb.co/WNcWfLJP/unnamed.jpg）

### 测试服务器

使用netcat发送测试消息：

```bash
echo '{"txt":"测试消息","num":"+8613061709786","cmd":"sms","metas":{"tz":32,"min":21,"seqNum":0,"refNum":0,"year":25,"sec":57,"maxNum":0,"mon":4,"hour":15,"day":22}}' | nc localhost 8080
```

## HTTP转发

服务器会将接收到的消息转发到以下格式的URL：
```
{http-url}/{phone-number}/{message-text}?icon={icon-url}
```

例如，如果设置 `-http-url https://test.com/path`，消息将被转发到：
```
https://test.com/path/+8613061709786/测试消息?icon=https://i.ibb.co/WNcWfLJP/unnamed.jpg
```

## 配置说明

### HTTP客户端配置

- 请求超时：120秒
- 连接超时：30秒
- TLS握手超时：10秒
- 空闲连接超时：90秒
- 最大空闲连接数：100
- 支持HTTPS和自签名证书

### 消息队列

- 队列大小：1000条消息
- 异步处理
- 自动重试失败的请求

## 日志输出

服务器会输出以下格式的日志：

1. 消息接收日志：
```
Message [txt=测试消息 num=+8613061709786 cmd=sms tz=32 time=15:21:57 date=2025-04-22 seq=0 ref=0 max=0]
```

2. 消息入队日志：
```
Queued message for HTTP request [txt=测试消息 num=+8613061709786]
```

3. HTTP响应日志：
```
HTTP Response [url=https://test.com/path/+8613061709786/测试消息?icon=https://i.ibb.co/WNcWfLJP/unnamed.jpg status=200]
```

## 优雅关闭

服务器支持优雅关闭：
- 处理SIGINT和SIGTERM信号
- 关闭TCP服务器
- 等待所有HTTP请求完成
- 关闭消息队列

## 错误处理

服务器会处理以下错误情况：
- TCP连接错误
- JSON解析错误
- HTTP请求错误
- URL解析错误

## 性能考虑

- 使用goroutine处理并发连接
- 消息队列缓冲防止内存溢出
- 连接池管理优化
- 超时控制防止资源耗尽

## 贡献

欢迎提交问题和改进建议！

## 许可证

MIT License 