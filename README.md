
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

## 接口

### 添加词库表

添加一个词库表。

**请求**

- 方法: GET
- 路径: `/add`
- 参数:
  - `table`: 表名称
  - `dict`: 词库文件的 URL

**示例**

```bash
curl -X GET "http://localhost:8080/add?table=table1&dict=https://raw.githubusercontent.com/mylukin/sensitive/master/dict/dict.txt"
```

### 删除词库表

删除一个词库表。

**请求**

- 方法: GET
- 路径: `/del`
- 参数:
  - `table`: 表名称

**示例**

```bash
curl -X GET "http://localhost:8080/del?table=table1"
```

### 查询敏感词

查询文本中的敏感词。

**请求**

- 方法: GET
- 路径: `/filter`
- 参数:
  - `table`: 表名称
  - `text`: 需要查询的文本

**示例**

```bash
curl -X GET "http://localhost:8080/filter?table=table1&text=你好吗？我支持习大大，他的名字叫做习近平"
```

### 替换敏感词

替换文本中的敏感词。

**请求**

- 方法: GET
- 路径: `/replace`
- 参数:
  - `table`: 表名称
  - `text`: 需要查询的文本
  - `to`: 替换的字符

**示例**

```bash
curl -X GET "http://localhost:8080/replace?table=table1&text=你好吗？我支持习大大，他的名字叫做习近平&to=*"
```

### 检测敏感词

检测文本中的敏感词。

**请求**

- 方法: GET
- 路径: `/findin`
- 参数:
  - `table`: 表名称
  - `text`: 需要检测的文本

**示例**

```bash
curl -X GET "http://localhost:8080/findin?table=table1&text=你好吗？我支持习大大，他的名字叫做习近平"
```

### 找到所有匹配的敏感词

找到文本中所有匹配的敏感词。

**请求**

- 方法: GET
- 路径: `/findall`
- 参数:
  - `table`: 表名称
  - `text`: 需要查询的文本

**示例**

```bash
curl -X GET "http://localhost:8080/findall?table=table1&text=你好吗？我支持习大大，他的名字叫做习近平"
```

### 验证文本是否合法

验证文本是否合法。

**请求**

- 方法: GET
- 路径: `/validate`
- 参数:
  - `table`: 表名称
  - `text`: 需要验证的文本

**示例**

```bash
curl -X GET "http://localhost:8080/validate?table=table1&text=你好吗？我支持习大大，他的名字叫做习近平"
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
make run
```

### 使用 Docker

1. 构建 Docker 镜像

```bash
make docker-build
```

2. 推送 Docker 镜像到仓库

```bash
make docker-push
```

3. 运行 Docker 容器

```bash
make docker-run
```

4. 停止 Docker 容器

```bash
make docker-stop
```

5. 重启 Docker 容器

```bash
make docker-restart
```

## 贡献

欢迎提交问题和拉取请求。

## 许可证

该项目使用 MIT 许可证。
