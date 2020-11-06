# 文件API列表

## 更新通知：

* 2013.6.20 上传、下载新域名正式上线使用，相关接口“上传单个文件”、“分片上传-文件分片上传”、“下载单个文件”及“下载流式文件”相关接口信息更新

* 2013.3.20 去除API权限申请开通相关说明，开发者可通过“管理中心”自行开启。

* 2013.3.1 新增“回收站功能”，新增“还原单个文件或目录”、“还原多个文件或目录”、“获取回收站文件或目录列表”及“清空回收站”等接口

## 基本功能

### 空间配额信息

#### 功能

    获取当前用户空间配额信息。

#### HTTP请求方式

    GET

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/quota

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：info。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |

#### 返回参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| quota | uint64 | 是 | 空间配额，单位为字节。 |
| used | uint64 | 是 | 已使用空间大小，单位为字节。 |

#### 示例

##### 请求示例

    GET https://pcs.baidu.com/rest/2.0/pcs/quota?method=info&access_token=1.54be391000a16ee6a21791d4a8ea04fe.86400.1331206383.67272939-188383

##### 响应示例

    {
        "quota":15000000000,
        "used":5221166,
        "request_id":4043312634
    }

### 上传单个文件

#### 功能

    上传单个文件。
    百度PCS服务目前支持最大2G的单个文件上传。
    如需支持超大文件（>2G）的断点续传，请参考下面的“分片文件上传”方法。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：upload。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| path | string | 是 | 上传文件路径（含上传的文件名称)。<br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|
| file | char[] | 是 | 上传文件的内容。 |
| ondup | string | 是 | * overwrite：表示覆盖同名文件；<br/> * newcopy：表示生成文件副本并进行重命名，命名规则为“文件名_日期.后缀”。 |

#### 返回参数

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| path | string | 是 | 该文件的绝对路径。 |
| size | uint64 | 否 | 文件字节大小。 |
| ctime | uint64 | 否 | 文件创建时间。 |
| mtime | uint64 | 否 | 文件修改时间。 |
| md5 | string | 否 | 文件的md5签名。 |
| fs_id | uint64 | 否 | 文件在PCS的临时唯一标识ID。 |

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=upload&path=%2fapps%2falbum%2f1.JPG&access_token=b778fb598c717c0ad7ea8c97c8f3a46f

##### 响应示例

    {
    　    "path" : "/apps/album/1.jpg",
    　    "size" : 372121,
    　    "ctime" : 1234567890,
    　    "mtime" : 1234567890,
    　    "md5" : "cb123afcc12453543ef",
    　    "fs_id" : 12345,
        　"request_id":4043312669
    }

### 分片上传—文件分片及上传

#### 功能

    百度PCS服务支持每次直接上传最大2G的单个文件。
    如需支持上传超大文件（>2G），则可以通过组合调用分片文件上传的upload方法和createsuperfile方法实现：
    首先，将超大文件分割为2G以内的单文件，并调用upload将分片文件依次上传；
    其次，调用createsuperfile，完成分片文件的重组。
    除此之外，如果应用中需要支持断点续传的功能，也可以通过分片上传文件并调用createsuperfile接口的方式实现。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：upload。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| type | string | 是 | 固定值，tmpfile。 |
| file | char[] | 是 | 上传文件的内容。 |

#### 返回参数

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| md5 | string | 否 | 文件的md5签名。 |

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=upload&access_token=1.54bef000f2416ee4a41791d4a8ea04fe.86400.1331206383.67272939-188383&type=tmpfile

##### 响应示例

    {
        "md5":"a7619410bca74850f985e488c9a0d51e",
        "request_id":3238563823
    }

### 分片上传—合并分片文件

#### 功能

    与分片文件上传的upload方法配合使用，可实现超大文件（>2G）上传，同时也可用于断点续传的场景。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：createsuperfile。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| path | string | 是 | 上传文件路径（含上传的文件名称)。<br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|
| param | string | 是 | block_list数组，数组的取值为子文件内容的MD5；子文件至少两个，最多1024个。 <br/> *本参数必须放在Http Body中进行传输，value示例： <br/> {"block_list":["d41d8cd98f00b204e9800998ecf8427e","89dfb274b42951b973fc92ee7c252166","1c83fe229cb9b1f6116aa745b4ef3c0d"]} |
| ondup | string | 是 | * overwrite：表示覆盖同名文件；<br/> * newcopy：表示生成文件副本并进行重命名，命名规则为“文件名_日期.后缀”。 |

#### 返回参数

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| path | string | 是 | 该文件的绝对路径。 |
| size | uint64 | 否 | 文件大小（以字节为单位）。 |
| ctime | uint64 | 否 | 文件创建时间。 |
| mtime | uint64 | 否 | 文件修改时间。 |
| md5 | string | 否 | 文件的md5签名。 |
| fs_id | uint64 | 否 | 文件在PCS的临时唯一标识ID。 |

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/file?method=createsuperfile&path=%2fapps%2fyunform%2f6ddddd.JPG&access_token=1.9fb09e8cce44c0d000e6787138924a26.86400.1331273905.2600617452-188383

##### 响应示例

    {
        "path":"/apps/yunform/6ddddd.JPG",
        "size":6844,
        "ctime":1331197101,
        "mtime":1331197101,
        "md5":"baa7c379639b74e9bf98c807498e1b64",
        "fs_id":1548308694,
        "request_id":4043313276
    }

### 下载单个文件

#### 功能

    下载单个文件。
    Download接口支持HTTP协议标准range定义，通过指定range的取值可以实现断点下载功能。 例如：
    如果在request消息中指定“Range: bytes=0-99”，那么响应消息中会返回该文件的前100个字节的内容；继续指定“Range: bytes=100-199”，那么响应消息中会返回该文件的第二个100字节内容。

#### HTTP请求方式

    GET

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：download。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| path | string | 是 | 下载文件路径，以/开头的绝对路径。<br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|

#### 返回参数

    无

#### 示例

##### 请求示例

    GET https://pcs.baidu.com/rest/2.0/pcs/file?method=download&access_token=3.d9000194f4b5d2da3fe8b6f850ace082.2592000.1348645419.2233553628-248414&path=%2Fapps%2F%E6%B5%8B%E8%AF%95%E5%BA%94%E7%94%A8%2F%2F01.jpg

##### 响应示例

    文件内容

### 创建目录

#### 功能

    为当前用户创建一个目录。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：mkdir。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| path | string | 是 | 需要创建的目录，以/开头的绝对路径。<br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|

#### 返回参数

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| fs_id | uint64 | 否 | 目录在PCS的临时唯一标识id。 |
| path | string | 否 | 该目录的绝对路径。 |
| ctime | uint64 | 否 | 目录创建时间。 |
| mtime | uint64 | 否 | 目录修改时间。 |

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=mkdir&access_token=1.54bef000f2416ee4a41791d4a8ea04fe.86400.1331206383.67272939-188383&path=%2Fapps%2Fyunform%2Fmusic

##### 响应示例

    {
        "fs_id":1636599174,
        "path":"/apps/yunfom/music",
        "ctime":1331183814,
        "mtime":1331183814,
        "request_id":4043312656
    }

### 获取单个文件/目录的元信息

#### 功能

    获取单个文件或目录的元信息。

#### HTTP请求方式

    GET

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：meta。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| path | string | 是 | 需要获取文件属性的目录，以/开头的绝对路径。如：/apps/album/a/b/c <br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|

#### 返回参数

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| fs_id | uint64 | 否 | 文件或目录在PCS的临时唯一标识ID。 |
| path | string | 否 | 文件或目录的绝对路径。 |
| ctime | uint | 否 | 文件或目录的创建时间。 |
| mtime | uint | 否 | 文件或目录的最后修改时间。 |
| block_list | string | 否 | 文件所有分片的md5数组JSON字符串。 |
| size | uint64 | 否 | 文件大小（byte）。 |
| isdir | uint | 否 | 是否是目录的标识符：<br/> * “0”为文件 <br/> * “1”为目录 |
| ifhassubdir | uint | 否 | 是否含有子目录的标识符：<br/> * “0”表示没有子目录 <br/> * “1”表示有子目录 |

#### 示例

##### 请求示例

    GET https://pcs.baidu.com/rest/2.0/pcs/file?method=meta&access_token=1.5400f91df2416ee4a41791d4a8ea04fe.86400.1331206383.67272939-188383&path=%2Fapps%2Fyunform%2Fmusic%2Fhello 

##### 响应示例

    {
        "list": [{
            "fs_id": 3528850315,
            "path": "/apps/yunform/music/hello",
            "ctime": 1331184269,
            "mtime": 1331184269,
            "block_list": ["59ca0efa9f5633cb0371bbc0355478d8"],
            "size": 13,
            "isdir": 1
        }],
        "request_id": 4043312678
    }

### 批量获取文件/目录的元信息

#### 功能

    批量获取文件或目录的元信息。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：meta。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| param | string | 是 | JSON字符串。<br/> {"list":[{"path":"\/apps\/album\/a\/b\/c"},{"path":"\/apps\/album\/a\/b\/d"}]} <br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|

#### 返回参数

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| fs_id | uint64 | 否 | 文件或目录在PCS的临时唯一标识ID。 |
| path | string | 否 | 文件或目录的绝对路径。 |
| server_filename | string | 否 | 文件或目录的名称。 |
| ctime | uint | 否 | 文件或目录的创建时间。 |
| mtime | uint | 否 | 文件或目录的最后修改时间。 |
| md5 | string | 否 | 文件的md5值。 |
| block_list | string | 否 | 文件所有分片的md5数组JSON字符串。 |
| size | uint64 | 否 | 文件大小（byte）。 |
| isdir | uint | 否 | 是否是目录的标识符：<br/> * “0”为文件 <br/> * “1”为目录 |
| ifhassubdir | uint | 否 | 是否含有子目录的标识符：<br/> * “0”表示没有子目录 <br/> * “1”表示有子目录 |


#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=meta&access_token=1.54b0091ee2416ee4a41791d4a8ea04fe.86400.1331206383.67272939-188383

##### 响应示例

    {
        "list": [{
                "fs_id": 3528850315,
                "path": "/apps/album/a/b/c",
                "ctime": 1331184269,
                "mtime": 1331184269,
                "block_list": ["59ca0efa9f5633cb0371bbc0355478d8"],
                "size": 13,
                "isdir": 0
            },
            {
                "fs_id": 3528850320,
                "path": "/apps/album/a/b/d",
                "ctime": 1331184269,
                "mtime": 1331184269,
                "block_list": ["59ca0efa9f5633cb0371bbc0355478d8"],
                "size": 13,
                "isdir": 0
            }
        ],
        "request_id": 4043312678
    }

### 获取目录下的文件列表

#### 功能

    获取目录下的文件列表。

#### HTTP请求方式

    GET

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：list。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| path | string | 是 | 需要list的目录，以/开头的绝对路径。<br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|
| by | string | 否 | 排序字段，缺省根据文件类型排序：<br/> * time（修改时间）<br/> * name（文件名）<br/> * size（大小，注意目录无大小） |
| order | string | 否 | “asc”或“desc”，缺省采用降序排序。<br/> * asc（升序）<br/> * desc（降序） |
| limit | string | 否 | 返回条目控制，参数格式为：n1-n2。<br/> 返回结果集的[n1, n2)之间的条目，缺省返回所有条目；n1从0开始。 |

#### 返回参数

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| fs_id | uint64 | 否 | 文件或目录在PCS的临时唯一标识ID。 |
| path | string | 否 | 文件或目录的绝对路径。 |
| server_filename | string | 否 | 文件或目录的名称。 |
| ctime | uint | 否 | 文件或目录的创建时间。 |
| mtime | uint | 否 | 文件或目录的最后修改时间。 |
| md5 | string | 否 | 文件的md5值。 |
| block_list | string | 否 | 文件所有分片的md5数组JSON字符串。 |
| size | uint64 | 否 | 文件大小（byte）。 |
| isdir | uint | 否 | 是否是目录的标识符：<br/> * “0”为文件 <br/> * “1”为目录 |

#### 示例

##### 请求示例

    GET https://pcs.baidu.com/rest/2.0/pcs/file?method=list&access_token=1.54bef91002416ee4a41791d4a8ea04fe.86400.1331206383.67272939-188383&path=%2Fapps%2Fyunform%2Fhello 

##### 响应示例

    {
        "list": [{
            "fs_id": 3528850315,
            "path": "/apps/yunform/music/hello",
            "ctime": 1331184269,
            "mtime": 1331184269,
            "block_list": ["59ca0efa9f5633cb0371bbc0355478d8"],
            "size": 13,
            "isdir": 0
        }],
        "request_id": 4043312670
    }

### 移动单个文件/目录

#### 功能

    移动单个文件/目录。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：move。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| from | string | 是 | 源文件地址（包括文件名）。<br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|
| to | string | 是 | 目标文件地址（包括文件名）。<br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|

#### 返回参数

    如果move操作执行成功，那么response会返回执行成功的from/to列表。

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| from | string | 是 | 执行move操作成功的源文件地址。 |
| to | string | 是 | 执行move操作成功的目标文件地址。 |

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=move&from=%2fapps%2f pcstest_oauth%2f test1%2fyyyytestwer.jpg&to=%2fapps%2fpcstest_oauth%2ftest2%2f2.jpg&access_token=b778fb598c717c0ad7ea8c97c8f3a46f 

##### 响应示例

    {
        "extra": {
            "list": [{
                "to": "/apps/pcstest_oauth/test2/2.jpg",
                "from": "/apps/pcstest_oauth/test1/yyyytestwer.jpg"
            }]
        },
        "request_id": 2298812844
    }

### 批量移动文件/目录

#### 功能

    批量移动文件/目录。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：move。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| param | string | 是 | 源文件地址和目标文件地址对应的列表。<br/> {"list":[{"from":"/apps/album/a/b/c","to":"/apps/album/b/b/c"},{"from":"/apps/album/a/b/d","to":"/apps/album/b/b/d"}]} <br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|

#### 返回参数

    返回参数extra由list数组组成，list数组的两个元素分别是“from”和“to”，代表move操作的源地址和目的地址。

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| from | string | 是 | 执行move操作成功的源文件地址。 |
| to | string | 是 | 执行move操作成功的目标文件地址。 |

#### 注意

    调用move接口时，目标文件的名称如果和源文件不相同，将会在move操作时对文件进行重命名。

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=move&from=%2fapps%2f pcstest_oauth%2f test1%2fyyyytestwer.jpg&to=%2fapps%2fpcstest_oauth%2ftest2%2f2.jpg&access_token=b778fb598c717c0ad7ea8c97c8f3a46f 

##### 响应示例

    {
        "extra": {
            "list": [{
                "to": "/apps/pcstest_oauth/test2/2.jpg",
                "from": "/apps/pcstest_oauth/test1/yyyytestwer.jpg"
            }]
        },
        "request_id": 2298812844
    }

### 拷贝单个文件/目录

#### 功能

    拷贝文件(目录)。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：copy。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| from | string | 是 | 源文件地址。<br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|
| to | string | 是 | 目标文件地址。<br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|

#### 返回参数

    如果copy操作执行成功，那么response会返回执行成功的from/to列表。

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| from | string | 是 | 执行copy操作成功的源文件地址。 |
| to | string | 是 | 执行copy操作成功的目标文件地址。 |

#### 注意

    move操作后，源文件被移动至目标地址；copy操作则会保留原文件。

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=copy&from=%2fapps%2fpcstest_oauth%2f test1%2f6.jpg&to=%2fapps%2fpcstest_oauth%2ftest2%2f6.jpg&access_token=b700fb598c717c0ad7ea8c97c8f3a46f 

##### 响应示例

    {
        "extra": {
            "list": [{
                "to": "/apps/pcstest_oauth/test2/6.jpg",
                "from": "/apps/pcstest_oauth/test1/6.jpg"
            }]
        },
        "request_id": 2298812844
    }

### 批量拷贝文件/目录

#### 功能

    批量拷贝文件/目录。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：copy。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| param | string | 是 | 源文件地址和目标文件地址对应的列表。<br/> {"list":[{"from":"/apps/album/a/b/c","to":"/apps/album/b/b/c"},{"from":"/apps/album/a/b/d","to":"/apps/album/b/b/d"}]} <br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|

#### 返回参数

    返回参数extra由list数组组成，list数组的两个元素分别是“from”和“to”，代表copy操作的源地址和目的地址。

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| from | string | 是 | 执行copy操作成功的源文件地址。 |
| to | string | 是 | 执行copy操作成功的目标文件地址。 |

#### 注意

    执行批量copy操作时，param参数通过HTTP Body传递；
    批量执行copy操作时，copy接口一次对请求参数中的每个from/to进行操作；执行失败就会退出，成功就继续，返回执行成功的from/to列表。

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=copy&access_token=b778fb008c717c0ad7ea8c97c8f3a46f

##### 响应示例

    {
        "extra": {
            "list": [{
                    "to": "/apps/pcstest_oauth/test1/6.jpg",
                    "from": "/apps/pcstest_oauth/test2/6.jpg"
                },
                {
                    "to": "/apps/pcstest_oauth/test2/89.jpg",
                    "from": "/apps/pcstest_oauth/89.jpg"
                }
            ]
        },
        "request_id": 2166619191
    }

### 删除单个文件/目录

#### 功能

    删除单个文件/目录。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：delete。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| path | string | 是 | 需要删除的文件或者目录路径。如：/apps/album/a/b/c <br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|

#### 返回参数

    无

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=delete&access_token=1.54bef91002416ee4a41791d4a8ea04fe.86400.1331206383.67272939-188383&path=%2Fapps%2Fyunform%2Fmusic 

##### 响应示例

    {
        "request_id": 4043312866
    }

### 批量删除文件/目录

#### 功能

    批量删除文件/目录。
    注意：
    * 文件/目录删除后默认临时存放在回收站内，删除文件或目录的临时存放不占用用户的空间配额；
    * 存放有效期为10天，10天内可还原回原路径下，10天后则永久删除。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：delete。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| param | string | 是 | 需要删除的文件或者目录路径。如：<br/> {"list":[{"path":"\/apps\/album\/a\/b\/c"},{"path":"\/apps\/album\/a\/b\/d"}]} <br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|

#### 返回参数

    无

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=delete&access_token=1.54bef91002416ee4a41791d400ea04fe.86400.1331206383.67272939-188383 

##### 响应示例

    {
        "request_id": 4043312865
    }

### 搜索

#### 功能

    按文件名搜索文件（不支持查找目录）。

#### HTTP请求方式

    GET

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：search。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| path | string | 是 | 需要检索的目录。 <br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|
| wd | string | 是 | 关键词。 |
| re | string | 否 | 是否递归。 <br/> * “0”表示不递归 <br/> * “1”表示递归 <br/> * 缺省为“0” |

#### 返回参数

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| fs_id | uint64 | 否 | 文件或目录在PCS的临时唯一标识ID。 |
| path | string | 否 | 文件或目录的绝对路径。 |
| server_filename | string | 否 | 文件或目录的名称。 |
| ctime | uint | 否 | 文件或目录的创建时间。 |
| mtime | uint | 否 | 文件或目录的最后修改时间。 |
| md5 | string | 否 | 文件的md5值。 |
| block_list | string | 否 | 文件所有分片的md5数组JSON字符串。 |
| size | uint64 | 否 | 文件大小（byte）。 |
| isdir | uint | 否 | 是否是目录的标识符：<br/> * “0”为文件 <br/> * “1”为目录 |

#### 示例

##### 请求示例

    GET https://pcs.baidu.com/rest/2.0/pcs/file?method=search&access_token=1.54bee00df241eee4a41791d4a8ea04fe.86400.1331206383.67272939-188383&path=%2Fapps%2Fyunform%2Fmusic&wd=hello&re=1 

##### 响应示例

    {
        "list": [{
            "fs_id": 3528850315,
            "path": "/apps/yunform/music/hello",
            "ctime": 1331184269,
            "mtime": 1331184269,
            "block_list": ["59ca0efa9f5633cb0371bbc0355478d8"],
            "size": 13,
            "isdir": 0
        }],
        "request_id": 4043312670
    }

## 高级功能

### 缩略图

#### 功能

    获取指定图片文件的缩略图。

#### HTTP请求方式

    GET

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/thumbnail

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：generate。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| path | string | 是 | 源图片的路径。 <br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|
| quality | int32 | 否 | 缩略图的质量，默认为“100”，取值范围(0,100]。 |
| height | int | 是 | 指定缩略图的高度，取值范围为(0,1600]。 |
| width | int | 是 | 指定缩略图的宽度，取值范围为(0,1600]。 |

#### 返回参数

    无

#### 注意

    有以下限制条件：
    * 原图大小(0, 10M]；
    * 原图类型: jpg、jpeg、bmp、gif、png；
    * 目标图类型:和原图的类型有关；例如：原图是gif图片，则缩略后也为gif图片。

#### 示例

##### 请求示例

    GET https://pcs.baidu.com/rest/2.0/pcs/thumbnail?method=generate&path=%2Fapps%2Fpcstest_oauth%2FSunset.jpg&quality=100&width=1600&height=1600

##### 响应示例

    缩略图文件内容

### 增量更新查询

#### 功能

    文件增量更新操作查询接口。本接口有数秒延迟，但保证返回结果为最终一致。

#### HTTP请求方式

    GET

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：diff。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| cursor | string | 是 | 用于标记更新断点。<br/> * 首次调用cursor=null；<br/> * 非首次调用，使用最后一次调用diff接口的返回结果中的cursor。 |

#### 返回参数

| 参数名称 | 类型 |  描述 |
| :- | :-: | :- |
| entries | array | k-v形式的列表，分为以下两种形式： <br/> 1. key为path，value为path对应的meta值，meta中isdelete=0为更新操作 <br/> * 如果path为文件，则更新path对应的文件；<br/> * 如果path为目录，则更新path对应的目录信息，但不更新path下的文件。<br/> 2. key为path，value为path删除的meta信息，meta中“isdelete!=0”为删除操作。<br/> * isdelete=1 该文件被永久删除；<br/> * isdelete=-1 该文件被放置进回收站；<br/> * 如果path为文件，则删除该path对应的文件；<br/> * 如果path为目录，则删除该path对应的目录和目录下的所有子目录和文件；<br/> * 如果path在本地没有任何记录，则跳过本删除操作。 |
| has_more | boolean | * True： 本次调用diff接口，增量更新结果服务器端无法一次性返回，客户端可以立刻再调用一次diff接口获取剩余结果；<br/> * False： 截止当前的增量更新结果已经全部返回，客户端可以等待一段时间（1-2分钟）之后再diff一次查看是否有更新。 |
| reset | boolean | * True： 服务器通知客户端，服务器端将按时间排序从第一条开始向客户端返回一份完整的数据列表；<br/> * False：返回上次请求返回cursor之后的增量更新结果。 |
| cursor | string | 用于下一次调用diff接口时传入的断点参数。 |

#### 示例

##### 请求示例

    *First time: cursor=null
    GET https://pcs.baidu.com/rest/2.0/pcs/file?method=diff&access_token=1.54bef91df2416ee4a41791d4a8ea04fe.86400.1331206383.67272939-188383&cursor=null

    *Next every time: cursor={cursor from last response}
    GET https://pcs.baidu.com/rest/2.0/pcs/file?method=diff&access_token=1.54bef91df2416ee4a41791d4a8ea04fe.86400.1331206383.67272939-188383&cursor=MxKx6UPi3w2Jt%2B%2BktMKKQpBbnC%2B11aH7Ec9pt%2BfteS%2F%2BknWrp3JIz%2F6fXHccEkZo2kkkSH748hScdRgcA4VCZJuCMQMvNkXAlSmzT5TwqBVc3xwhSxaFkClqbcogAOc8I0k7xtTb9nG6rBJsxNgRFgBV4F695TkrLDHYHRy%2BQ%3D%3D

##### 响应示例

    {
        "entries": {
            "\/baiduapp\/browser": {
                "fs_id": 2427025269,
                "path": "\/baiduapp\/browser",
                "size": 0,
                "isdir": 1,
                "md5": "",
                "mtime": 1336631762,
                "ctime": 1336631762
            }
        },
        "has_more": true,
        "reset": true,
        "cursor": "MxKx6UPie/9WzBkwALPrVWQlyxlmK0LgHG8zutwXp8oyC/ngIdGgS3w2Jt++ktMKKQpBbnC+11aH7Ec9pt+fteS/+knWrp3JIz/6fXHccEkZo2kkkSH748hScdRgcA4VCZJuCMQMvNkXAlSmzT5TwqBVc3xwhSxaFkClqbcogAOc8I0k7xtTb9nG6rBJsxNgRFgBV4F695TkrLDHYHRy+Q==",
        "request_id": 3355443548
    }

### 视频转码

#### 功能

    对视频文件进行转码，实现实时观看视频功能。 可下载支持HLS/M3U8的[媒体云播放器SDK](http://developer.baidu.com/wiki/index.php?title=docs/cplat/media/sdk)配合使用。

#### HTTP请求方式

    GET

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：streaming。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| path | string | 是 | 需要下载的视频文件路径，以/开头的绝对路径，需含源文件的文件名。 <br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|
| type | string | 是 | 目前支持以下格式：<br/> * M3U8_320_240、M3U8_480_224、M3U8_480_360、M3U8_640_480和M3U8_854_480 |

#### 注意

    目前这个接口支持的源文件格式如下：

| 格式名称 | 扩展名 | 备注 |
| :- | :- | :- |
| Apple HTTP Live Streaming | m3u8/m3u | iOS支持的视频格式 |
| ASF | asf | 视频格式 |
| AVI | avi | 视频格式 |
| Flash Video (FLV) | flv | Macromedia Flash视频格式 |
| GIF Animation | gif | 视频格式 |
| Matroska | mkv | Matroska/WebM视频格式 |
| MOV/QuickTime/MP4 | mov/mp4/m4a/3gp/3g2/mj2 | 支持3GP、3GP2、PSP、iPod 之类视频格式 |
| MPEG-PS (program stream) | mpeg | 也就是VOB文件、SVCD DVD格式 |
| MPEG-TS (transport stream) | ts | 即DVB传输流 |
| RealMedia | rm/rmvb | Real视频格式 |
| WebM | webm | Html视频格式 |

#### 返回参数

    无

#### 示例

##### 请求示例

    GET https://pcs.baidu.com/rest/2.0/pcs/file?method=streaming&path=%2fapps%2fvideo%2fv1.mov&access_token=b778fb000c717c0ad7ea8c97c8f3a46f&type=MP4_480P

##### 响应示例

    直接返回文件内容

### 获取流式文件列表

#### 功能

    以视频、音频、图片及文档四种类型的视图获取所创建应用程序下的文件列表。

#### HTTP请求方式

    GET

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/stream

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：list。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| type | string | 是 | 类型分为video、audio、image及doc四种。 |
| start | string | 否 | 返回条目控制起始值，缺省值为0。 |
| limit | string | 否 | 返回条目控制长度，缺省为1000，可配置。 |
| filter_path | string | 否 | 需要过滤的前缀路径，如：/apps/album <br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|

#### 返回参数

| 参数名称 | 类型 |  描述 |
| :- | :-: | :- |
| total | uint | 文件总数。 |
| start | uint | 起始数。 |
| limit | uint | 获取数。 |
| path | string | 获取流式文件的绝对路径。 |
| block_list | string | 分片MD5列表。 |
| size | uint | 流式文件的文件大小（byte）。 |
| mtime | uint | 流式文件在服务器上的修改时间 。 |
| ctime | uint | 流式文件在服务器上的创建时间 。 |
| fs_id | uint64 | 流式文件在PCS中的唯一标识ID 。 |
| isdir | uint | * 0：文件 <br/> * 1：目录 |

#### 示例

##### 请求示例

    GET https://pcs.baidu.com/rest/2.0/pcs/stream?method=list&type=image&start=50&limit=100&access_token=b778fb000c717c0ad7ea8c97c8f3a46f

##### 响应示例

    {
        "total": 13,
        "start": 0,
        "limit": 1,
        "list": [{
            "path": "/apps/album/1.jpg",
            "size": 372121,
            "ctime": 1234567890,
            "mtime": 1234567890,
            "md5": "cb123afcc12453543ef",
            "fs_id": 12345,
            "isdir": 0
        }]
    }

### 下载流式文件

#### 功能

    为当前用户下载一个流式文件。其参数和返回结果与下载单个文件的相同。

#### HTTP请求方式

    GET

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/stream

#### 注意

    1. 兼容原有域名pcs.baidu.com；使用新域名d.pcs.baidu.com，则提供更快、更稳定的下载服务。
    2. 需注意处理好 302 跳转问题。

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：download。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| path | string | 是 | 需要下载的文件路径，以/开头的绝对路径，含文件名。 <br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|

#### 返回参数

    无

#### 示例

##### 请求示例

    GET https://pcs.baidu.com/rest/2.0/pcs/stream?method=download&access_token=b778fb000c717c0ad7ea8c97c8f3a46f&path=%2fapps%2falbum%2f1.jpg

##### 响应示例

    流式文件内容

### 秒传文件

#### 功能

    秒传一个文件。
    注意：
    * 被秒传文件必须大于256KB（即 256*1024 B）。
    * 校验段为文件的前256KB，秒传接口需要提供校验段的MD5。(非强一致接口，上传后请等待1秒后再读取)

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：rapidupload。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| path | string | 是 | 上传文件的全路径名。 <br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B|
| content-length | int | 是 | 待秒传的文件长度。 |
| content-md5 | string | 是 | 待秒传的文件的MD5。 |
| slice-md5 | string | 是 | 待秒传文件校验段的MD5。 |
| content-crc32 | string | 是 | 待秒传文件CRC32 |
| ondup | string | 否 | * overwrite：表示覆盖同名文件；<br/> * newcopy：表示生成文件副本并进行重命名，命名规则为“文件名_日期.后缀”。 |

#### 返回参数

| 参数名称 | 类型 | 描述 |
| :- | :-: | :- |
| path | string | 秒传文件的绝对路径。 |
| size | uint64 | 秒传文件的字节大小 。 |
| ctime | uint64 | 秒传文件的创建时间。 |
| mtime | uint64 | 秒传文件的修改时间 。 |
| md5 | string | 秒传文件的md5签名。 |
| fs_id | uint64 | 秒传文件在PCS的唯一标识ID。 |
| isdir | uint | * 0：文件 <br/> * 1：目录 |

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=rapidupload&access_token=b778fb000c717c0ad7ea8c97c8f3a46f&content-length=1542719&content-md5=3edf3d47292280e0182db6750bd176e5&slice-md5=6fce289cfee3e4414788dcd000a3ddc4&path=%2fa%2fb%2fc

##### 响应示例

    {
        "path": "/apps/album/1.jpg",
        "size": 372121,
        "ctime": 1234567890,
        "mtime": 1234567890,
        "md5": "cb123afcc12453543ef",
        "fs_id": 12345,
        "isdir": 0,
        "request_id": 12314124
    }

### 添加离线下载任务

#### 功能

    添加离线下载任务，实现单个文件离线下载。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/services/cloud_dl

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：add_task。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| expires | int | 否 | 请求失效时间，如果有，则会校验。 |
| save_path | string | 是 | 下载后的文件保存路径。<br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B |
| source_url | string | 是 | 源文件的URL。 |
| rate_limit | int | 否 | 下载限速，默认不限速。 |
| timeout | int | 否 | 下载超时时间，默认3600秒。 |
| callback | string | 否 | 下载完毕后的回调，默认为空。 |

#### 返回参数

    任务ID号

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/service/cloud_dl?method=add_task&access_token=40001fsdjfdjskaf&source_url=http:\/\/dl_dir.qq.com:80\/qqfile\/qq\/QQ2012\/QQ2012.exe

##### 响应示例

    任务ID号成功：
    {"task_id":432432432432432,"request_id":3372220525}
    任务并发太大：
    {"error_code":36013,"error_msg":"too many tasks","request_id":3372220539}

### 精确查询离线下载任务

#### 功能

    根据任务ID号，查询离线下载任务信息及进度信息。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/service/cloud_dl

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：query_task。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| expires | int | 否 | 请求失效时间，如果有，则会校验。 |
| task_ids | string | 是 | 要查询的任务ID信息，如：1,2,3,4 |
| op_type | int | 是 | * 0：查任务信息 <br/> * 1：查进度信息，默认为1 |

#### 返回参数

    无

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/services/cloud_dl?method=query_task&access_token=43000fsdjfdjskaf

##### Body示例

    查询任务信息：
    {
    “request_id”:12394838223,
    “task_info”:
    {
        “123456” : {
            "result" : 0, //0查询成功，结果有效，1要查询的task_id不存在
            "source_url":"http://www.example.com/xxx.zip",//下载数据源地址
            "save_path":"http://xxxx/",//下载完成后的存放地址
            "rate_limit":10,
            "timeout": 3600,
            "callback":"http://XXX",
            "status":1 (0下载成功，1下载进行中 2系统错误，3资源不存在，4下载超时，5资源存在但下载失败 6存储空间不足 7目标地址数据已存在 8任务取消)
            "create_time":"UNIX_TIMESTAMP",//任务创建时间
    }
    “43829483”:
    {
    "result" : 1, //要查询的task_id不存在
    }
    }
    }
    查询进度信息：
    {
    “request_id”:12394838223,
    “task_info”:
    {
        “123456” : {
            "result" : 0, //0查询成功，结果有效，1要查询的task_id不存在
            "status": (0下载成功，1下载进行中 2系统错误，3资源不存在，4下载超时，5资源存在但下载失败 6存储空间不足 7任务取消)
            //其余字段，在status为0、1时有效
            "file_size":1024
            "finished_size":512
            "create_time": 123232132,
            "start_time": 43728943,
            "finish_time": 43728948,
        }
    “43829483”:
    {
    "result" : 1, //要查询的task_id不存在
    }
    }
    }

##### 响应示例

    任务信息或进度信息

### 查询离线下载任务列表

#### 功能

    查询离线下载任务ID列表及任务信息。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/services/cloud_dl

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：list_task。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| expires | int | 否 | 请求失效时间，如果有，则会校验。 |
| start | int | 否 | 查询任务起始位置，默认为0。 |
| limit | int | 否 | 设定返回任务数量，默认为10。 |
| asc | int | 否 | * 0：降序，默认值 <br/> * 1：升序 |
| source_url | string | 否 | 源地址URL，默认为空。 |
| save_path | string | 否 | 文件保存路径，默认为空。<br/>注意：<br/> * 路径长度限制为1000 <br/> * 路径中不能包含以下字符：\\ ? \| " > < : * <br/> * 文件名或路径名开头结尾不能是“.”或空白字符，空白字符包括: \r, \n, \t, 空格, \0, \x0B |
| create_time | int | 否 | 任务创建时间，默认为空。 |
| status | int | 否 | 任务状态，默认为空。 |
| need_task_info | int | 否 | 是否需要返回任务信息: <br/> * 0：不需要 <br/> * 1：需要，默认为1 |

#### 返回参数

    任务信息或任务列表

#### 示例

##### 请求示例

    查询离线下载任务ID列表
    POST https://pcs.baidu.com/rest/2.0/pcs/services/cloud_dl?method=list_task&access_token=43000fsdjfdjskaf&need_task_info=1

##### 响应示例

    任务列表
    {"task_info":[{"task_id":"26"}],"total":"1","request_id":1283164486}
    任务信息
    {"task_info":[{"task_id":"26","source_url":"http:\/\/dl_dir.qq.com:80\/qqfile\/qq\/QQ2012\/QQ2012.exe","save_path":"\/apps\/Slideshow\/wu_jing_test0","rate_limit":"100","timeout":"10000","callback":"http:\/\/www.baidu.com","status":"1","create_time":"1347449048"}],"total":"1","request_id":1285732167}

### 取消离线下载任务

#### 功能

    取消离线下载任务。

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/services/cloud_dl

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：cancel_task。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| expires | int | 否 | 请求失效时间，如果有，则会校验。 |
| task_id | string | 是 | 要取消的任务ID号。 |

#### 返回参数

    任务信息或任务列表

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/services/cloud_dl?method=cancel_task&access_token=43000fsdjfdjskaf&task_id=26

##### 响应示例

    {
        “request_id”:12394838223,
    }

### 回收站

#### 功能

    回收站用于临时存放删除文件，且不占空间配额；但回收站的文件存放具有10天有效期，删除文件默认扔到回收站，10天内可通过回收站找回，逾期永久删除。

    回收站功能，目前支持以下几种操作和接口：

    * 删除文件到回收站： delete
    （目前默认是删除到回收站），API详细说明请参考“删除单个文件或目录”及“批量删除文件或目录”部分
    * 查看回收站文件： listrecycle
    * 还原回收站文件（单个文件或多个文件）： restore
    * 清空回收站：delete

### 查询回收站文件

#### 功能

    获取回收站中的文件及目录列表。

#### HTTP请求方式

    GET

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：listrecycle。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| start | int | 否 | 返回条目的起始值，缺省值为0 |
| limit | int | 否 | 返回条目的长度，缺省值为1000 |

#### 返回参数

    获取成功，返回回收站文件或目录列表信息

| 参数名称 | 描述 |
| :- | :- |
| list | list是一个数组，以JSON串的形式显示所获取到的回收站文件或目录的具体信息，详细信息参考下面的“list数组元素说明”。 |
| request_id | 请求ID，也就是服务器用于追踪错误的的日志ID |

    list数组中的元素说明如下：

| 参数名称 | 类型 | UrlEncode | 描述 |
| :- | :-: | :-: | :- |
| fs_id | uint64 | 否 | 目录在PCS上的临时唯一标识 |
| path | string | 是 | 该目录的绝对路径 |
| ctime | uint | 否 | 文件在服务器上的创建时间 |
| mtime | uint | 否 | 文件在服务器上的修改时间 |
| md5 | string | 否 | 分片MD5 |
| size | uint | 否 | 文件大小（byte） |
| isdir | uint | 否 | 是否是目录的标识符：<br/> * “0”为文件 <br/> * “1”为目录 |

    获取失败，则返回错误信息

| 参数名称 | 描述 |
| :- | :- |
| error_code | 错误码，详见“[文件API错误码列表](https://github.com/iikira/BaiduPCS-Go/blob/master/docs/file_data_apis_error.md)” |
| error_msg | 错误信息，详见“[文件API错误码列表](https://github.com/iikira/BaiduPCS-Go/blob/master/docs/file_data_apis_error.md)” |
| request_id | 请求ID，也就是服务器用于追踪错误的的日志ID |

#### 示例

##### 请求示例

    GET https://pcs.baidu.com/rest/2.0/pcs/file?method=listrecycle&start=50&limit=100&access_token=111f1118c717111a8111c8f3a46f

##### 响应示例

    获取成功: 
    {
    "list":[
    {
    "fs_id":1579174,
    "path":"\/apps\/CloudDriveDemo\/testfile-10.rar",
    "ctime":1361934614,
    "mtime":1361934625,
    "md5":"1131170ac11cfbec411a5e8d4e111769",
    "size":10730431,
    "isdir":0
    },
    {
    "fs_id":304521061,
    "path":"\/apps\/CloudDriveDemo\/testfile-4.rar",
    "ctime":1361934605,
    "mtime":1361934625,
    "md5":"9552bf5e5abdf962e2de94be243bec7c",
    "size":4287611,
    "isdir":0
    }
    ],
    "request_id":3779302504
    }
    获取失败: 
    {"error_code":110,"error_msg":"Access token invalid or no longer valid","request_id":1444638699}

### 还原单个文件或目录

#### 功能

    还原单个文件或目录（非强一致接口，调用后请sleep 1秒读取）

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：restore。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| fs_id | string | 是 | 所还原的文件或目录在PCS的临时唯一标识ID。 |

#### 返回参数

    还原成功时，返回以下参数

| 参数名称 | 描述 |
| :- | :- |
| extra | extra由list数组组成，list数组中包含一个元素fs_id，即文件或目录在PCS的临时唯一标识ID。 |
| request_id | 请求ID，也就是服务器用于追踪错误的的日志ID |

    还原失败，则返回错误信息

| 参数名称 | 描述 |
| :- | :- |
| error_code | 错误码，详见“[文件API错误码列表](https://github.com/iikira/BaiduPCS-Go/blob/master/docs/file_data_apis_error.md)” |
| error_msg | 错误信息，详见“[文件API错误码列表](https://github.com/iikira/BaiduPCS-Go/blob/master/docs/file_data_apis_error.md)” |
| request_id | 请求ID，也就是服务器用于追踪错误的的日志ID |

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=restore&fs_id=123456&access_token=111f1118c717111a8111c8f3a46f

##### 响应示例

    还原成功
    {"extra":{"list":[{"fs_id":"1356099017"}]},"request_id":3775323016}
    还原失败
    {"error_code":31061,"error_msg":"file already exists","request_id":811204199}

### 批量还原文件或目录

#### 功能

    批量还原文件或目录（非强一致接口，调用后请sleep1秒 ）

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值：restore。 |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| param | string | 是 | Body中的JSON串，用于批量处理 |

#### 返回参数

| 参数名称 | 描述 |
| :- | :- |
| error_code | 错误码，详见“[文件API错误码列表](https://github.com/iikira/BaiduPCS-Go/blob/master/docs/file_data_apis_error.md)” |
| error_msg | 错误信息，详见“[文件API错误码列表](https://github.com/iikira/BaiduPCS-Go/blob/master/docs/file_data_apis_error.md)” |
| extra | extra由list数组组成，list数组中包含一个元素fs_id，即文件或目录ID |
| request_id | 请求ID，也就是服务器用于追踪错误的的日志ID |

    说明：
    * 全部还原成功的情况下，返回extra及request_id信息；
    * 还原多个文件或目录时，如果还原某个文件或目录失败，则报错并终止还原操作，返回error_code、error_msg、extra及request_id信息；
        * 其中，如果未成功还原任何文件时，返回的extra中的list数组为空；
        * 部分文件或目录还原成功时，则返回的extra中list数组中显示该部分文件的fs_id。

#### 示例

##### 请求示例

    POST  https://pcs.baidu.com/rest/2.0/pcs/file?method=restore&access_token=111f1118c717111a8111c8f3a46f.2592000.1364548100.3123371436-248414&param={"list":[{"fs_id":"4059450057"},{"fs_id":"2959141864"}]}

##### 响应示例

    全部还原成功
    {"extra":{"list":[{"fs_id":"2959141864"}],"request_id":1359873129}
    全部还原失败
    {"error_code":31061,"error_msg":"file already exists","extra":{"list":[]},"request_id":1342759216}
    部分还原成功
    {"error_code":31061,"error_msg":"file already exists","extra":{"list":[{"fs_id":"2959141864"}]},"request_id":1359873129}

### 清空回收站

#### 功能

    清空回收站

#### HTTP请求方式

    POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/file

#### 请求参数

| 参数名称 | 类型 | 是否必需 | 描述 |
| :- | :-: | :-: | :- |
| method | string | 是 | 固定值为delete |
| access_token | string | 是 | 开发者准入标识，HTTPS调用时必须使用。 |
| type | string | 是 | 固定值为recycle |

#### 返回参数

    清空成功，返回请求ID 

| 参数名称 | 描述 |
| :- | :- |
| request_id | 请求ID，也就是服务器用于追踪错误的的日志ID |

    清空失败，则返回错误信息

| 参数名称 | 描述 |
| :- | :- |
| error_code | 错误码，详见“[文件API错误码列表](https://github.com/iikira/BaiduPCS-Go/blob/master/docs/file_data_apis_error.md)” |
| error_msg | 错误信息，详见“[文件API错误码列表](https://github.com/iikira/BaiduPCS-Go/blob/master/docs/file_data_apis_error.md)” |
| request_id | 请求ID，也就是服务器用于追踪错误的的日志ID |

#### 示例

##### 请求示例

    POST https://pcs.baidu.com/rest/2.0/pcs/file?method=delete&type=recycle&access_token=111f1118c717111a8111c8f3a46f

##### 响应示例

    清空成功
    {"request_id":2307473052}
    清空失败
    {"error_code":31070,"error_msg":"file delete failed","request_id":12345678} 
