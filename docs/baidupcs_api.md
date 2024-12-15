baidupcs-go的api说明

# 启动api服务

通过serve子命令开启服务，支持basic auth对api进行验证。命令参数如下

| 参数名   | 说明              | 类型   | 默认值     |
| -------- | ----------------- | ------ | ---------- |
| port     | 服务端口号        | int    | 8080       |
| auth     | 是否开启basic验证 | bool   | false      |
| username | basic auth用户名  | string | admin      |
| password | basicauth密码     | string | adminadmin |

开启api服务会阻塞主进程，如需关闭可以ctrl+c退出程序。

# api 描述

api基本和终端命令一一对应，接受get或post方法。在包含path的子命令中时，只能使用post方法，因为get方法的url code会造成路径中空格和+的混淆。

## ls

- 请求路径：{{baseurl}}/api/ls;

- 请求方式：POST

| 参数名 | 说明               | 类型   | 默认值 |
| ------ | ------------------ | ------ | ------ |
| asc    | 是否顺序排序       | bool   | true   |
| desc   | 是否降序排序       | bool   | false  |
| time   | 是否按照时间排序   | bool   | false  |
| name   | 是否按照文件名排序 | bool   | true   |
| size   | 查看的路径         | string | .      |

## search

- 请求路径：{{baseurl}}/api/search
- 请求方式POST,GET

| 参数名            | 说明         | 类型   | 默认值 |
| ----------------- | ------------ | ------ | ------ |
| rescure           | 是否递归搜索 | bool   | false  |
| keyword（必须值） | 检索的关键词 | string |        |
| path              | 检索的路劲   | string | .      |

## pwd

- 请求路径：{{baseurl}}/api/pwd
- 请求方式：GET，POST

该方法不带任何参数

## meta

- 请求路径：{{baseurl}}/api/meta
- 请求方式：POST

| 参数名       | 说明                         | 类型     | 默认值 |
| ------------ | ---------------------------- | -------- | ------ |
| target_paths | 需要查看元数据的文件路径列表 | []string | null   |

## download

- 请求路劲：{{baseurl}}/api/download
- 请求方式：POST

| 参数名                 | 说明                                      | 类型     | 默认值 |
| ---------------------- | ----------------------------------------- | -------- | ------ |
| save                   | 是否将下载文件保存到当前路径中            | bool     | false  |
| save_to                | 将下载的文件存放到指定目录                | string   | .      |
| mode                   | 下载模式，可选值：pcs,stream,locate       | string   | locate |
| is_test                | 是否为测试下载，测试方式下，不会保存文件  | bool     | false  |
| is_printstatus         | 是否输出所有线程工作状态                  | bool     | false  |
| is_executed_permission | 是否为文件加上执行权限（windows系统无效） | bool     | false  |
| is_overwrite           | 是否覆盖已存在的文件                      | bool     | false  |
| parallel               | 下载线程数                                | int      | 0      |
| load                   | 指定同时进行下载文件的数量                | int      | 1      |
| max_retry              | 下载失败后最大重试次数                    | int      | 3      |
| no_check               | 下载完成后，不校验文件                    | bool     | false  |
| link_prefer            | 使用备选下载链接中的第几个                | int      | 1      |
| modifyMTime            | 是否将本地文件的时间修改为服务器上的时间  | bool     | false  |
| full_path              | 是否以网盘完整路径保存到本地              | bool     | false  |
| paths                  | 待下载路径（必须给定值）                  | []string | null   |

- TODO: 这了考虑给download加上生命周期的webhook接口

## rm

- 请求路劲：{{baseurl}}/api/rm
- 请求方式：POST，GET

| 参数名       | 说明                               | 类型     | 默认值 |
| ------------ | ---------------------------------- | -------- | ------ |
| target_paths | 待移除的路径列表（该参数必须指定） | []string | null   |

## mkdir

- 请求路劲：{{baseurl}}/api/mkdir
- 请求方式：POST，GET

| 参数名      | 说明                               | 类型   | 默认值 |
| ----------- | ---------------------------------- | ------ | ------ |
| target_path | 待创建的目录路径（该参数必须指定） | string | null   |

## cp

- 请求路劲：{{baseurl}}/api/cp
- 请求方式：POST

| 参数名     | 说明                                   | 类型     | 默认值 |
| ---------- | -------------------------------------- | -------- | ------ |
| from_paths | 待拷贝的文件或目录路径(该参数必须指定) | []string | null   |
| to_path    | 目标路径（该参数必须指定）             | string   | null   |

## mv

- 请求路劲：{{baseurl}}/api/cp
- 请求方式：POST

| 参数名     | 说明                                   | 类型     | 默认值 |
| ---------- | -------------------------------------- | -------- | ------ |
| from_paths | 待拷贝的文件或目录路径(该参数必须指定) | []string | null   |
| to_path    | 目标路径（该参数必须指定）             | string   | null   |
