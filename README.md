# 简易用户文件系统
## **使用：**

1.在根文件夹中使用以下命令来启动服务程序


**`go run ./cmd/filesys/main.go `**

**！！！第一次运行时需要在根文件夹下依次运行以下两个命令：**

**`go run ./cmd/init_db/main.go `**

**`go run ./cmd/gen/main.go `**

访问运行服务程序主机的8080端口就可以查看并且测试文件系统，不过前端界面只是测试界面，是一个很简陋的测试界面。

## 功能概述

1. 实现一个简易的用户文件系统。
2. 每个用户拥有独立的文件树，节点类型包括文件和文件夹。
3. 每个文件或文件夹都有唯一 ID，根目录文件夹 ID 为 0。
4. 支持文件和文件夹的基本操作接口。
5. 系统默认自带管理员账号 `admin`，仅管理员可创建其他用户。
6. 同一文件夹下文件和文件夹均不可重名，若重名则自动重命名（规则参考 Windows）。
7. 文件支持版本管理，版本号从 1 开始递增，当前版本号最大。
8. ~~高级功能：客户端配合服务端使用 rsync 差分算法实现历史版本增量上传（后续版本实现）。~~
9. 文件存储于工程根目录下的 `data` 目录，需考虑磁盘文件清理。
10. 推荐使用 Postman 进行接口测试，支持简单前端页面展示更佳。

---

## 技术栈

- Web 框架：Gin
- ORM 框架：GORM
- 数据库：SQLite

---

## 接口说明

### 登录接口

- `POST /login`  
  用户登录，返回 `sid`，后续接口需在 Cookie 中携带 `sid`。

### 管理接口（仅管理员）

- `POST /api/user`  
  创建新用户。

### 文件接口

- `POST /api/file/{file_id}/new`  
  新建文件夹，`file_id` 为父目录 ID，返回文件夹信息。

- `POST /api/file/{file_id}/upload`  
  上传文件，`file_id` 为父目录 ID，文件二进制内容放在请求体，返回文件信息。

- `POST /api/file/{file_id}/update`  
  更新文件，文件二进制内容放在请求体，返回文件信息。

- `DELETE /api/file/{file_id}`  
  删除文件或文件夹。

- `POST /api/file/{file_id}/copy`  
  复制文件或文件夹。

- `POST /api/file/{file_id}/move`  
  移动文件或文件夹。

- `POST /api/file/{file_id}/rename`  
  重命名文件或文件夹。

- `GET /api/file/{file_id}`  
  获取文件或文件夹信息。

- `GET /api/file/{file_id}/list`  
  获取文件夹下的文件和文件夹列表。

- `GET /api/file/{file_id}/content`  
  下载文件内容。

- `GET /api/file/{file_id}/version/{ver_num}/content`  
  下载指定历史版本的文件内容。
