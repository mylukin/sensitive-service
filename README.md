
# Sensitive Service

Sensitive Service 是一个用于管理和过滤敏感词的 HTTP 服务，使用 Golang 和 Echo 框架开发，并打包为 Docker 镜像。

## 功能

1. 添加词库表
2. 删除词库表
3. 查询文本中的敏感词
4. 替换文本中的敏感词
5. 检测文本中的敏感词
6. 找到所有匹配的敏感词
7. 验证文本是否合法
8. 支持添加全局屏蔽词库表

## 接口

### 添加词库表

通过字典 URL 添加一个新的词库表。

**接口:** `GET /add`

**参数:**
- `table` (字符串): 词库表的名称。
- `dict` (字符串): 字典的 URL。

**示例:**

```bash
curl -X GET "http://localhost:8080/add?table=table1&dict=https://raw.githubusercontent.com/mylukin/sensitive/master/dict/dict.txt"
```

### 删除词库表

删除现有的词库表。

**接口:** `GET /del`

**参数:**
- `table` (字符串): 词库表的名称。

**示例:**

```bash
curl -X GET "http://localhost:8080/del?table=table1"
```

### 过滤文本

使用指定的词库表过滤文本中的敏感词。

**接口:** `GET /filter`

**参数:**
- `table` (字符串): 词库表的名称。
- `text` (字符串): 要过滤的文本。

**示例:**

```bash
curl -X GET "http://localhost:8080/filter?table=table1&text=你好吗？我支持习大大，他的名字叫做习近平"
```

### 替换敏感词

使用指定的字符替换文本中的敏感词。

**接口:** `GET /replace`

**参数:**
- `table` (字符串): 词库表的名称。
- `text` (字符串): 要处理的文本。
- `to` (字符串): 用于替换敏感词的字符。

**示例:**

```bash
curl -X GET "http://localhost:8080/replace?table=table1&text=你好吗？我支持习大大，他的名字叫做习近平&to=*"
```

### 检测敏感词

使用指定的词库表检查文本中是否有敏感词。

**接口:** `GET /findin`

**参数:**
- `table` (字符串): 词库表的名称。
- `text` (字符串): 要检查的文本。

**示例:**

```bash
curl -X GET "http://localhost:8080/findin?table=table1&text=你好吗？我支持习大大，他的名字叫做习近平"
```

### 找到所有敏感词

使用指定的词库表找到文本中所有的敏感词。

**接口:** `GET /findall`

**参数:**
- `table` (字符串): 词库表的名称。
- `text` (字符串): 要检查的文本。

**示例:**

```bash
curl -X GET "http://localhost:8080/findall?table=table1&text=你好吗？我支持习大大，他的名字叫做习近平"
```

### 验证文本

使用指定的词库表验证文本是否清洁。

**接口:** `GET /validate`

**参数:**
- `table` (字符串): 词库表的名称。
- `text` (字符串): 要验证的文本。

**示例:**

```bash
curl -X GET "http://localhost:8080/validate?table=table1&text=你好吗？我支持习大大，他的名字叫做习近平"
```

## 全局词典

### 添加全局词典

添加一个适用于所有词库表的全局词典。

**接口:** `GET /add?table=*`

**参数:**
- `dict` (字符串): 字典的 URL。

**示例:**

```bash
curl -X GET "http://localhost:8080/add?table=*&dict=https://raw.githubusercontent.com/mylukin/sensitive/master/dict/dict.txt"
```

### 删除全局词典

删除全局词典。

**接口:** `GET /del?table=*`

**参数:**
- `table` (字符串): 使用 `*` 指定全局词典。

**示例:**

```bash
curl -X GET "http://localhost:8080/del?table=*"
```

## 构建与运行

### 本地运行

确保你已经安装了 Go 环境。

1. 安装依赖

```bash
go mod tidy
```

2. 运行服务

```bash
go run main.go
```

### 使用 Docker

1. 构建 Docker 镜像

```bash
docker build -t sensitive-service .
```

2. 运行 Docker 容器

```bash
docker run -d -p 8080:8080 --name sensitive-service sensitive-service
```

3. 停止 Docker 容器

```bash
docker stop sensitive-service
```

4. 重启 Docker 容器

```bash
docker restart sensitive-service
```

5. 删除 Docker 容器

```bash
docker rm sensitive-service
```

## 贡献

欢迎提交问题和拉取请求。

## 许可证

该项目使用 MIT 许可证。