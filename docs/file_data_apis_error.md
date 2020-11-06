# 文件数据API错误码

| HTTP状态码 | 错误码 | 错误信息 | 备注 |
| :- | -: | :- | :- |
| 200 | 0 | no error | 没有错误 |
| 400 | 3 | Unsupported open api | 不支持此接口 |
| 403 | 4 | No permission to do this operation | 没有权限执行此操作 |
| 403 | 5 | Unauthorized client IP address | IP未授权 |
| 503 | 31001 | db query error | 数据库查询错误 |
| 503 | 31002 | db connect error | 数据库连接错误 |
| 503 | 31003 | db result set is empty | 数据库返回空结果 |
| 503 | 31021 | network error | 网络错误 |
| 503 | 31022 | can not access server | 暂时无法连接服务器 |
| 400 | 31023 | param error | 输入参数错误 |
| 400 | 31024 | app id is empty | app id为空 |
| 503 | 31025 | bcs error | 后端存储错误 |
| 403 | 31041 | bduss is invalid | 用户的cookie不是合法的百度cookie |
| 403 | 31042 | user is not login | 用户未登陆 |
| 403 | 31043 | user is not active | 用户未激活 |
| 403 | 31044 | user is not authorized | 用户未授权 |
| 403 | 31045 | user not exists | 用户不存在 |
| 403 | 31046 | user already exists | 用户已经存在 |
| 400 | 31061 | file already exists | 文件已经存在 |
| 400 | 31062 | file name is invalid | 文件名非法 |
| 400 | 31063 | file parent path does not exist | 文件父目录不存在 |
| 403 | 31064 | file is not authorized | 无权访问此文件 |
| 400 | 31065 | directory is full | 目录已满 |
| 403 | 31066 | file does not exist | 文件不存在 |
| 503 | 31067 | file deal failed | 文件处理出错 |
| 503 | 31068 | file create failed | 文件创建失败 |
| 503 | 31069 | file copy failed | 文件拷贝失败 |
| 503 | 31070 | file delete failed | 文件删除失败 |
| 503 | 31071 | get file meta failed | 不能读取文件元信息 |
| 503 | 31072 | file move failed | 文件移动失败 |
| 503 | 31073 | file rename failed | 文件重命名失败 |
| 503 | 31081 | superfile create failed | superfile创建失败 |
| 503 | 31082 | superfile block list is empty | superfile 块列表为空 |
| 503 | 31083 | superfile update failed | superfile 更新失败 |
| 503 | 31101 | tag internal error | tag系统内部错误 |
| 503 | 31102 | tag param error | tag参数错误 |
| 503 | 31103 | tag database error | tag系统错误 |
| 403 | 31110 | access denied to set quota | 未授权设置此目录配额 |
| 400 | 31111 | quota only sopport 2 level directories | 配额管理只支持两级目录 |
| 400 | 31112 | exceed quota | 超出配额 |
| 403 | 31113 | the quota is bigger than one of its parent directorys | 配额不能超出目录祖先的配额 |
| 403 | 31114 | the quota is smaller than one of its sub directorys | 配额不能比子目录配额小 |
| 503 | 31141 | thumbnail failed, internal error | 请求缩略图服务失败 |
| 401 | 110 | Access token invalid or no longer valid | Access Token不正确或者已经过期 |
| 400 | 31201 | signature error | 签名错误 |
| 400 | 31203 | acl put error | 设置acl失败 |
| 400 | 31204 | acl query error | 请求acl验证失败 |
| 400 | 31205 | acl get error | 获取acl失败 |
| 404 | 31079 | File md5 not found, you should use upload API to upload the whole file. | 未找到文件MD5 |，| 请使用上传API上传整个文件。 |
| 404 | 31202 | object not exists | 文件不存在 |
| 404 | 31206 | acl get error | acl不存在 |
| 400 | 31207 | bucket already exists | bucket已存在 |
| 400 | 31208 | bad request | 用户请求错误 |
| 500 | 31209 | baidubs internal error | 服务器错误 |
| 501 | 31210 | not implement | 服务器不支持 |
| 403 | 31211 | access denied | 禁止访问 |
| 503 | 31212 | service unavailable | 服务不可用 |
| 503 | 31213 | service unavailable | 重试出错 |
| 503 | 31214 | put object data error | 上传文件data失败 |
| 503 | 31215 | put object meta error | 上传文件meta失败 |
| 503 | 31216 | get object data error | 下载文件data失败 |
| 503 | 31217 | get object meta error | 下载文件meta失败 |
| 403 | 31218 | storage exceed limit | 容量超出限额 |
| 403 | 31219 | request exceed limit | 请求数超出限额 |
| 403 | 31220 | transfer exceed limit | 流量超出限额 |
| 500 | 31298 | the value of KEY[VALUE] in pcs response headers is invalid | 服务器返回值KEY非法 |
| 500 | 31299 | no KEY in pcs response headers | 服务器返回值KEY不存在 |
