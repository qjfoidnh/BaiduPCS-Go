# 结构化数据API错误码

结构化数据API以HTTP提供，因此请求者应首先检查HTTP协议级别的响应状态码。 当请求有错误或执行失败时，HTTP协议的返回响应状态码不为“200”，且在Content-body的JSON格式数据中以error_code给出错误码，并可能以error_msg字段提示错误信息。错误码的详细信息，请参考以下错误码列表。 一般来说, 错误分为客户端错误和服务器端错误：

* 客户端错误：
返回状态码4xx表示结构化数据平台认为客户端请求有错误，例如：auth失败，参数不对，超出quota等。这些情况，客户端需要先解决参数问题，再向服务器端重新发起请求。

* 服务器端错误：
返回状态码为5xx，表示结构化数据平台内部发生错误，客户端需要重试。

## 注意

请求成功时，HTTP协议的返回响应状态码为“200”，不会设置error_code和error_msg。

## 结构化数据API错误码

| HTTP状态码 | 错误码 | 错误信息 | 备注 | 是否重试 |
| - | :- | :- | :- | - |
| 500 | 1 | Unknown error | 未知错误 | 是 |
| 500 | 2 | Service temporarily unavailable | 服务暂不可用 | 是 |
| 403 | 6 | No permission to access user data | 无权访问用户数据 | 否 |
| 403 | 7 | No permission to access data for this referer | 无权访问数据 | 否 |
| 400 | 100 | Invalid parameter | 无效参数 | 否 |
| 401 | 101 | Invalid API key | 无效API Key | 否 |
| 401 | 102 | Session key invalid or no longer valid | 会话密钥无效 | 否 |
| 401 | 103 | Invalid/Used call_id parameter | call_id参数无效/已被使用 | 否 |
| 400 | 104 | Incorrect signature | 签名错误 | 否 |
| 400 | 105 | Too many parameters | 参数过多 | 否 |
| 400 | 106 | Unsupported signature method | 不支持此签名方式 | 否 |
| 400 | 107 | Invalid/Used timestamp parameter | 时间戳无效 | 否 |
| 401 | 108 | Invalid user id | 用户ID无效 | 否 |
| 400 | 109 | Invalid user info field | 用户信息字段无效 | 否 |
| 401 | 110 | Access token invalid or no longer valid | Access token无效或已失效 | 否 |
| 401 | 111 | Access token expired | Access token已过期 | 否 |
| 401 | 112 | Session key expired | 会话密钥已过期 | 否 |
| 400 | 114 | Invalid Ip | 无效IP | 否 |
| 400 | 31400 | param error | 参数错误 | 否 |
| 400 | 31401 | malformed json | JSON格式错误 | 否 |
| 400 | 31402 | no "table" in request | 请求中没有“table”字段 | 否 |
| 400 | 31403 | no "records" in request | 请求中没有“records”字段 | 否 |
| 400 | 31405 | too many records in request | 请求中的records 过多，目前限制为500 | 否 |
| 400 | 31406 | bad columnname | 列名非法，请参考API文档 | 否 |
| 400 | 31407 | record too large | record过大，> 1M | 否 |
| 400 | 31408 | bad table name | table名称不合法 | 否 |
| 400 | 31409 | table not exist | table不存在，请先创建 | 否 |
| 400 | 31410 | bad record | record格式错误，请检查JSON | 否 |
| 400 | 31411 | no appid | 请求中没有“app_id”字段 | 否 |
| 400 | 31412 | no userid | 请求中没有“user_id”字段 | 否 |
| 400 | 31420 | bad condition | condition描述错误。 | 否 |
| 400 | 31421 | bad projection | projection描述错误 | 否 |
| 400 | 31422 | bad order_by | order_by描述错误 | 否 |
| 400 | 31423 | bad operator | condition中的operation 非法 | 否 |
| 400 | 31424 | bad start/limit | start/limit 错误 | 否 |
| 400 | 31425 | unsupported operator | 操作符暂未支持，如：or、like、regex等 | 否 |
| 400 | 31430 | no key in record | update/delete 请求，但是record 中没有_key 字段 | 否 |
| 400 | 31431 | record not exist | 符合条件的record不存在，比如if-match不匹配、在回收站等 | 否 |
| 400 | 31432 | unknown op | 参数op非法 | 否 |
| 400 | 31433 | bad key | key非法 | 否 |
| 400 | 31440 | param cursor not set | 参数cursor未设值 | 否 |
| 400 | 31441 | param cursor format error | 参数cursor格式错误 | 否 |
| 400 | 31442 | param cursor appid wrong | 参数cursor appid错误 | 否 |
| 400 | 31443 | param cursor user_id wrong | 参数cursor user_id错误 | 否 |
| 400 | 31450 | exceed quota | 超出配额 | 否 |
| 400 | 31451 | quota size param not exist | 找不到参数quota size | 否 |
| 503 | 31452 | quota info fail | quota info失败 | 是 |
| 400 | 31453 | quota too big | quota过大 | 否 |
| 400 | 31454 | quota size param not numberic | quota size 参数未数值化 | 否 |
| 400 | 31460 | no permission | 未授权 | 否 |
| 400 | 31461 | account not login | 账户为登录，使用bduss认证失败 | 否 |
| 400 | 31462 | access token errro | access token校验失败 | 否 |
| 400 | 31470 | index num too much | index num太多 | 否 |
| 400 | 31472 | table already exist | table已存在 | 否 |
| 400 | 31473 | abnormal table already exist | 异常table已存在 | 否 |
| 400 | 31474 | table not drop, cannot restore | table不在回收站，无法恢复 | 否 |
| 400 | 31475 | engine not support | 不支持此项操作 | 否 |
| 400 | 31480 | param op wrong, should be recycled or permanent | 参数op错误，应为可回收或永久的 | 否 |
| 400 | 31490 | api not support | 调用了错误的API | 否 |
| 500 | 31500 | Internal error (Try Again Later) | 内部错误 | 是 |
| 503 | 31501 | storeengine construct fail | construct失败 | 是 |
| 503 | 31502 | storeengine select fail | 选择操作失败 | 是 |
| 503 | 31503 | storeengine insert fail | 插入操作失败 | 是 |
| 503 | 31504 | storeengine update fail | 更新操作失败 | 是 |
| 503 | 31505 | storeengine delete fail | 删除操作失败 | 是 |
| 503 | 31506 | storeengine count fail | count操作失败 | 是 |
| 503 | 31507 | storeengine ensure index fail | 查询或创建索引失败 | 是 |
| 503 | 31508 | storeengine delete index fail | 删除索引失败 | 是 |
| 503 | 31509 | storeengine drop table fail | 删除table操作失败 | 是 |
| 503 | 31530 | config set num match fail | 配置中num匹配失败 | 是 |
| 503 | 31590 | db query error | db交互出错 | 是 |
| 503 | 31591 | network error | 内部网络交互错误 | 是 |
