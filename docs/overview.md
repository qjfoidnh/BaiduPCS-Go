# 概述

百度开放云平台为广大开发者提供了访问PCS资源的系列接口，目前开放的接口主要分两个部分：

* 文件API:

    主要提供文件上传、下载、拷贝、删除、搜索、断点续传及缩略图等功能。

* 结构化数据API:

    主要提供结构数据存储、查询、删除及同步等功能。

通过对这些API的组合调用，开发者可以实现基本的用户文件操作以及结构数据存储和管理功能，也能够支持用户数据在多种不同终端上的同步，以提供更优质的用户体验。

除了原生的REST（Representational State Transfer，即“表述性状态转移”） API之外，百度开放云平台还提供了多种平台的SDK来帮助开发者缩短开发周期，具体请参考“SDK”部分相关内容。

## PCS REST API使用说明

### 开通PCS API权限

PCS所有REST API都必须经过开通权限才能正常使用。申请的方法请参考“开通PCS API权限”部分相关内容。

注意：PCS未提供分享接口，download等接口仅供个人获取数据使用。 access_token不能泄露，否则会直接封禁应用。

### API请求方式说明

目前所有的提交类接口仅支持POST方式，查询类接口同时支持POST方式和GET方式。

PCS REST API的所有参数在传入时应当使用：UTF-8编码。

#### HTTP 请求方式

    GET | POST

#### URL

    https://pcs.baidu.com/rest/2.0/pcs/{object_name}?{query_string}

#### 参数说明

| 参数名称 | 描述 |
| :- | :- |
| object_name | PCS REST API操作实体名称，如：quota、file、thumbnail。 |
| query_string | 放在HTTP头部传入的参数，必须经过UrlEncode编码。 |

#### HTTP GET和POST方式使用说明

<table width="600" border="1" cellpadding="1" cellspacing="1">
    <tbody>
        <tr>
            <th scope="row" width="80">请求方式
            </th>
            <th scope="col">GET
            </th>
            <th scope="col">POST
            </th>
        </tr>
        <tr>
            <th scope="row">URL
            </th>
            <td colspan="2">https://pcs.baidu.com/rest/2.0/pcs/{object_name}?{query_string}
            </td>
        </tr>
        <tr>
            <th scope="row">请求参数
            </th>
            <td> 全部携带在 HTTPS 请求头部的 query_string 中。
            </td>
            <td> 既可携带在 query_string 中，也可携带在 HTTP Body 中。
                <dl>
                    <dd>
                        <ul>
                            <li> method 及 access token 等参数必须携带在 query_string 中进行传输，请参考各个API的具体说明；
                            </li>
                            <li> 携带在 query_string 中的参数的值，必须进行 UrlEncode 编码；
                            </li>
                            <li> 携带在 HTTP Body 中的参数，则不需要进行 UrlEncode 编码。
                            </li>
                        </ul>
                    </dd>
                </dl>
                <div style="border:solid 1px #d7d7d7;padding:10px 16px 2px 16px; background-color:#fbfafb;">
                    <div>注</div>
                    HTTP URL 长度有限，若参数值长度过长，建议将参数放在 HTTP Body 中进行传输。</div>
            </td>
        </tr>
        <tr>
            <th scope="row"> HTTP BODY
            </th>
            <td> 不携带HTTP Body
            </td>
            <td> multipart/form-data
            </td>
        </tr>
        <tr>
            <th scope="row">注意
            </th>
            <td colspan="2">如果 HTTP Body 和 query_string 存在相同的参数，则以 query_string 中的参数为准。
            </td>
        </tr>
    </tbody>
</table>

#### 使用示例

1. GET请求：

    用HTTP GET请求方式发送两个参数：key1=value1和key2=value2。

    https://pcs.baidu.com/rest/2.0/pcs/quota?key1=UrlEncode(value1)&key2=UrlEncode(value2)

2. POST请求：

    分别用两种方式使用POST方式发送三个参数：key1=value1、key2=value2和key3=value3；方式一与方式二效果等同。

##### 方式一：

    POST /rest/2.0/pcs/quota?key2=value2&key3=value3 HTTP/1.1

    User-Agent: curl/7.12.1 (x86_64-redhat-linux-gnu) libcurl/7.12.1 OpenSSL/0.9.7a zlib/1.2.1.2 libidn/0.5.6
    Pragma: no-cache
    Accept: */*
    Host:pcs.baidu.com
    Content-Length:123
    Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryS0JIa4uHF7yHd8xJ
    ------WebKitFormBoundaryS0JIa4uHF7yHd8xJ
    Content-Disposition: form-data; name="key1"
    value1
    ------WebKitFormBoundaryS0JIa4uHF7yHd8xJ—

##### 方式二：

    POST /rest/2.0/pcs/quota HTTP/1.1
    User-Agent: curl/7.12.1 (x86_64-redhat-linux-gnu) libcurl/7.12.1 OpenSSL/0.9.7a zlib/1.2.1.2 libidn/0.5.6
    Pragma: no-cache
    Accept: */*
    Host:pcs.baidu.com
    Content-Length:123
    Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryS0JIa4uHF7yHd8xJ
    ------WebKitFormBoundaryS0JIa4uHF7yHd8xJ
    Content-Disposition: form-data; name="key1"

    value1
    ------WebKitFormBoundaryS0JIa4uHF7yHd8xJ
    Content-Disposition: form-data; name="key2"

    value2
    ------WebKitFormBoundaryS0JIa4uHF7yHd8xJ
    Content-Disposition: form-data; name="key3"

    value3
    ------WebKitFormBoundaryS0JIa4uHF7yHd8xJ--

### API响应格式说明

<table>
    <tbody>
        <tr>
            <th scope="col">
            </th>
            <th scope="col">正常请求
            </th>
            <th scope="col">异常请求
            </th>
        </tr>
        <tr>
            <th scope="row">HTTP状态码
            </th>
            <td> 200 OK
            </td>
            <td> 4**&nbsp;: 用户请求错误。<br>5** ：server服务失败。
            </td>
        </tr>
        <tr>
            <th scope="row">HTTP BODY
            </th>
            <td> API响应内容
            </td>
            <td> 异常请求的返回值为JSON字符串。<br>例如：<br>{<br>"error_code":110,<br>"error_msg":"Access token invalid or no longer valid",<br>"request_id":729562373<br>}<br>说明：<br>
                <dl>
                    <dd>- error_code：错误码；<br>
                    </dd>
                    <dd>- error_msg: 错误描述信息；<br>
                    </dd>
                    <dd>- request_id: 请求ID。由server生成，用于追查和定位请求日志。<br>
                    </dd>
                </dl>
            </td>
        </tr>
    </tbody>
</table>
