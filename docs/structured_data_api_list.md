# 结构化数据API列表



## 更新通知：

2013.7.2   修改“创建table”接口，请求参数增加“sk”

## 创建table

### 功能

    创建一个表，定义索引，其中包括对唯一索引的支持。 

**注意：** 



*   一个应用最多创建5个表，一个表上最多创建5个索引；
*   关于表和索引的创建规则，您可以参考“结构化数据表基本概念”；

*   为保证一致性，创建表后，可能需要等待一段时间才能用describe table接口查看到；

*   创建一张表必须带有该表所属app的密匙sk，用于之后的psstoken鉴权使用。


**关于“唯一索引”的说明：** 



*   可以建联合的唯一索引；

*   一个table上唯一索引数量同一般索引限制，5个；

*   唯一索引不支持在表创建后在增加，所以需要在表设计的时候尽量考虑；（如确实需要，可联系我们。）

*   insert的时候必须带上唯一索引的所有字段，否则会失败。


### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/table

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：create。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>开发者的应用所对应的access_token。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>表名。
</td></tr>
<tr>
<td>sk
</td><td>string
</td><td>是
</td><td>该表所属app的密匙（secret key），用于psstoken鉴权使用。
</td></tr>
<tr>
<td>column
</td><td>json
</td><td>否
</td><td>列描述。
</td></tr>
<tr>
<td>index
</td><td>json
</td><td>否
</td><td>索引描述：

*   1：表示升序索引

*   -1：表示降序索引
</td></tr></tbody></table>
 
返回参数 
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码, 如果不出错, 则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示, 如果不出错, 则返回值没有该字段。
</td></tr>
<tr>
<td>app_id
</td><td>int
</td><td>应用对应的ID。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>表名。
</td></tr>
<tr>
<td>request_id
</td><td>int
</td><td>请求ID号。
</td></tr></tbody></table>
 
### 示例

请求示例: 

#### 1. 创建一般索引

<pre style="font-family:monospace;">$ cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_create_table <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"column"</span> <span style="color: #339933;">:</span>
    <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"id"</span> <span style="color: #339933;">:</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"description"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">""</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"type"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"int"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"required"</span> <span style="color: #339933;">:</span> <span style="color: #003366; font-weight: bold;">true</span>
        <span style="color: #009900;">}</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"index"</span><span style="color: #339933;">:</span>
    <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"id_index"</span> <span style="color: #339933;">:</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"column"</span> <span style="color: #339933;">:</span> <span style="color: #009900;">{</span><span style="color: #3366CC;">"id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #009900;">}</span>
        <span style="color: #009900;">}</span>
    <span style="color: #009900;">}</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
$ curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST
<span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/table?method=create&amp;access_token=2.b06c3e00010fdb879d12345dcd5f8545.2587600.134819999.1175746697-238347&amp;sk=cRgk8uMGX098yMfmttoVYswcv3XKBLGX"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_create_table"</span></pre> 

#### 2. 创建唯一索引

<pre style="font-family:monospace;">$ cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_create_table <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"column"</span> <span style="color: #339933;">:</span>
    <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"id"</span> <span style="color: #339933;">:</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"description"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">""</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"type"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"int"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"required"</span> <span style="color: #339933;">:</span> <span style="color: #003366; font-weight: bold;">true</span>
        <span style="color: #009900;">}</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"index"</span><span style="color: #339933;">:</span>
    <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"id_index"</span> <span style="color: #339933;">:</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"column"</span> <span style="color: #339933;">:</span> <span style="color: #009900;">{</span><span style="color: #3366CC;">"id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #009900;">}</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"unique"</span> <span style="color: #339933;">:</span> <span style="color: #003366; font-weight: bold;">true</span>
        <span style="color: #009900;">}</span>
    <span style="color: #009900;">}</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
$ curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST
<span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/table?method=create&amp;access_token=2.b06c3e00010fdb879d12345dcd5f8545.2587600.134819999.1175746697-238347&amp;sk=cRgk8uMGX098yMfmttoVYswcv3XKBLGX "</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_create_table"</span></pre>

### 注意

    unique字段为true，表示唯一索引；为false，则表示一般索引；不指定则默认为一般索引。

正确响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">200</span> OK
&nbsp;
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"app_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"table"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">3728395580</span>
<span style="color: #009900;">}</span></pre>
出错响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">400</span> Bad Request
&nbsp;
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"error_code"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">31472</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"error_msg"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"table already exist"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">9085631045</span>
<span style="color: #009900;">}</span></pre>



## 修改table

### 功能

    修改一个表，添加或者删除索引

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/table

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：alter。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>开发者的access_token。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>表名。
</td></tr>
<tr>
<td>add_index
</td><td>json
</td><td>否
</td><td>增加的索引。
</td></tr>
<tr>
<td>drop_index
</td><td>json
</td><td>否
</td><td>删除的索引。
</td></tr></tbody></table>
 
### 返回参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码，如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示，如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>appid
</td><td>int
</td><td>开发者App ID。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>表名。
</td></tr>
<tr>
<td>column
</td><td>json
</td><td>列描述。
</td></tr>
<tr>
<td>index
</td><td>json
</td><td>索引描述：

*   1：表示升序索引；

*   -1：表示降序索引。
</td></tr></tbody></table>
 
### 示例

请求示例: <pre style="font-family:monospace;">$ cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_alter_table <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"add_index"</span> <span style="color: #339933;">:</span>
    <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"direction"</span> <span style="color: #339933;">:</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"column"</span> <span style="color: #339933;">:</span> <span style="color: #009900;">{</span><span style="color: #3366CC;">"direction"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #009900;">}</span>
        <span style="color: #009900;">}</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"drop_index"</span><span style="color: #339933;">:</span>
    <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"id_index"</span> <span style="color: #339933;">:</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"column"</span> <span style="color: #339933;">:</span> <span style="color: #009900;">{</span><span style="color: #3366CC;">"id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #009900;">}</span>
        <span style="color: #009900;">}</span>
    <span style="color: #009900;">}</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
$ curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST
<span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/table?method=alter&amp;access_token=2.b06c3e00010fdb879d12345dcd5f8545.2587600.134819999.1175746697-238347"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_alter_table"</span>
&nbsp;</pre> 
正确响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">200</span> OK
&nbsp;
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"app_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"table"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"lastindex"</span> <span style="color: #339933;">:</span>
    <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"direction"</span> <span style="color: #339933;">:</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"column"</span> <span style="color: #339933;">:</span>
            <span style="color: #009900;">{</span>
                <span style="color: #3366CC;">"direction"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span>
            <span style="color: #009900;">}</span>
        <span style="color: #009900;">}</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">6402295586</span>
<span style="color: #009900;">}</span></pre>
出错响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">400</span> Bad Request
&nbsp;
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"error_code"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">31409</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"error_msg"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"table not exist"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">9360693908</span>
<span style="color: #009900;">}</span></pre>



## 删除table

### 功能

    删除一个table

### 注意 

    如果drop到回收站（默认情况），则drop后该表处于不可访问状态，不能再创建同名的table。

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/table

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：drop。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>开发者的App对应的access_token。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>表名。
</td></tr>
<tr>
<td>op
</td><td>string
</td><td>否
</td><td>值为recycled: drop到回收站，可用restore接口恢复。
</td></tr></tbody></table>
 
返回参数 
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>app_id
</td><td>int
</td><td>App对应的ID。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>表名。
</td></tr>
<tr>
<td>request_id
</td><td>int
</td><td>请求ID号。
</td></tr></tbody></table>
 
### 示例

请求示例: <pre style="font-family:monospace;"> $ cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_drop_table <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
$ curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST
<span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/table?method=drop&amp;access_token=2.b06c3e00010fdb879d12345dcd5f8545.2587600.134819999.1175746697-238347"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_drop_table"</span></pre> 
正确响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">200</span> OK
&nbsp;
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"app_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"table"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">3728395580</span>
<span style="color: #009900;">}</span>
&nbsp;
<span style="color: #339933;">&lt;/</span>pre<span style="color: #339933;">&gt;</span>
出错响应示例：<span style="color: #339933;">&lt;</span>javascript<span style="color: #339933;">&gt;</span>HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">400</span> Bad Request
&nbsp;
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"error_code"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">31409</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"error_msg"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"table not exist"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">9085631045</span>
<span style="color: #009900;">}</span></pre>



## 从回收站恢复table

### 功能

    恢复一个在回收站中的表。

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/table

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：restore。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>开发者的App对应的access_token。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>表名。
</td></tr></tbody></table>
 
返回参数 
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>app_id
</td><td>int
</td><td>App对应的ID。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>表名。
</td></tr>
<tr>
<td>request_id
</td><td>int
</td><td>请求ID号。
</td></tr></tbody></table>
 
### 示例

请求示例: <pre style="font-family:monospace;">$ cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_restore_table <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
$ curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST
<span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/table?method=restore&amp;access_token=2.b06c3e00010fdb879d12345dcd5f8545.2587600.134819999.1175746697-238347"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_restore_table"</span>
&nbsp;</pre> 
正确响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">200</span> OK
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"app_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"table"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">3728395580</span>
<span style="color: #009900;">}</span></pre>
出错响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">400</span> Bad Request
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"error_code"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">31474</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"error_msg"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"table not drop, cannot restore"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">9085631045</span>
<span style="color: #009900;">}</span></pre>



## 查看table创建信息

### 功能

    查看一个表的创建信息。

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/table

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号， 默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：describe。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>开发者的App对应的access_token。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>表名。
</td></tr></tbody></table>
 
返回参数 
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>app_id
</td><td>int
</td><td>App对应的ID。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>表名。
</td></tr>
<tr>
<td>request_id
</td><td>int
</td><td>请求ID号。
</td></tr>
<tr>
<td>column
</td><td>json
</td><td>表的列描述。
</td></tr>
<tr>
<td>index
</td><td>json
</td><td>表的索引描述。
</td></tr>
<tr>
<td>quota
</td><td>int
</td><td>该表单个用户最大的条目数限制。
</td></tr>
<tr>
<td>auth_code
</td><td>string
</td><td>第三方应用请忽略此参数。
</td></tr></tbody></table>
 
### 示例

请求示例: <pre style="font-family:monospace;"> $ cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_describe_table <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
$ curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST
<span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/table?method=describe&amp;access_token=2.b06c3e00010fdb879d12345dcd5f8545.2587600.134819999.1175746697-238347"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_describe_table"</span></pre> 
正确响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">200</span> OK
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"appid"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"table"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists008"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"status"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">0</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"ctime"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1347417209</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"mtime"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1347417209</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"cluster"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"cluster0"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"subtablenum"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"column"</span> <span style="color: #339933;">:</span>
    <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"id"</span> <span style="color: #339933;">:</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"description"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">""</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"type"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"number"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"required"</span> <span style="color: #339933;">:</span> <span style="color: #003366; font-weight: bold;">true</span>
        <span style="color: #009900;">}</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"index"</span><span style="color: #339933;">:</span>
    <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"id_index"</span> <span style="color: #339933;">:</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"column"</span> <span style="color: #339933;">:</span> <span style="color: #009900;">{</span><span style="color: #3366CC;">"id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #009900;">}</span>
        <span style="color: #009900;">}</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"quota"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">10000</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"auth_code"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"e3725dd9a7cbd0a5e3eb7928ab922d33"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">2707788940</span>
<span style="color: #009900;">}</span></pre>
出错响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">400</span> Bad Request
&nbsp;
<span style="color: #009900;">{</span>
  <span style="color: #3366CC;">"error_code"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">31409</span><span style="color: #339933;">,</span>
  <span style="color: #3366CC;">"error_msg"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"table not exist"</span><span style="color: #339933;">,</span>
  <span style="color: #3366CC;">"request_id"</span> <span style="color: #339933;">:</span> <span style="color: #CC0000;">5574355722</span>
<span style="color: #009900;">}</span></pre>



## 添加record

### 功能

    新增record，每次调用都会新增传入的record。

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/data

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：insert。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>用户的access_token，HTTPS调用时必须使用。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>要插入的目标表名。
</td></tr>
<tr>
<td>records
</td><td>json array
</td><td>是
</td><td>需要插入的record JSON对象构成的数组。
</td></tr></tbody></table>
 
返回参数 
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>records
</td><td>json array
</td><td>返回服务器端已经处理的records(_key, _mtime, _ctime)列表，顺序与输入顺序一致；如果一个请求包含多个record，遇到第一个出错record即中止，返回的records只包含已处理成功的key。
</td></tr></tbody></table>
 
### 示例

请求示例: 
<pre style="font-family:monospace;">cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_insert_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"id"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">85617</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"刘德华"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"type"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"男歌手"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"intro"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"香港著名歌手、演员"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"add_time"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">1340949289</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"language"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"国语"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"粤语"</span><span style="color: #009900;">]</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"tags"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"香港电影金像奖"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"四大天王"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"东亚唱片"</span><span style="color: #009900;">]</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"top_song"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span>
                <span style="color: #3366CC;">"id"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">3</span><span style="color: #339933;">,</span> 
                <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"爱你一万年"</span>
            <span style="color: #009900;">}</span>
        <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"id"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">85618</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"凤凰传奇"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"type"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"组合"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"intro"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"中国大陆具有广泛知名度的男女二人音乐组合"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"add_time"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">1340949289</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"language"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"国语"</span><span style="color: #009900;">]</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"tags"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"月亮之上"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"最炫民族风"</span><span style="color: #009900;">]</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"top_song"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span>
                <span style="color: #3366CC;">"id"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">5</span><span style="color: #339933;">,</span> 
                <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"月亮之上"</span>
            <span style="color: #009900;">}</span>
         <span style="color: #009900;">}</span>
    <span style="color: #009900;">]</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
$ curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST <span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=insert&amp;access_token=2.b06c3e86610fdb879d12345dcd5f8545.2587600.134819999.1175746697-238347"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_select1_request"</span></pre> 
正确响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">200</span> OK
&nbsp;
<span style="color: #009900;">{</span>
  <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span><span style="color: #009900;">[</span>
    <span style="color: #009900;">{</span>
      <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"f44603de003c57d5-1346066442"</span><span style="color: #339933;">,</span>
      <span style="color: #3366CC;">"_mtime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1346066442</span><span style="color: #339933;">,</span>
      <span style="color: #3366CC;">"_ctime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1346066442</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #009900;">{</span>
      <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"1aaef0010c012db7-1346066442"</span><span style="color: #339933;">,</span>
      <span style="color: #3366CC;">"_mtime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1346066442</span><span style="color: #339933;">,</span>
      <span style="color: #3366CC;">"_ctime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1346066442</span>
    <span style="color: #009900;">}</span>
  <span style="color: #009900;">]</span><span style="color: #339933;">,</span>
  <span style="color: #3366CC;">"request_id"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">3728395580</span>
<span style="color: #009900;">}</span></pre>
出错响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">400</span> Bad Request
&nbsp;
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"error_code"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">31430</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"error_msg"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"bad record"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">0</span>
    <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
           <span style="color: #009900;">{</span><span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"7d4febca4a68e763-1344915172"</span><span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #009900;">]</span>
<span style="color: #009900;">}</span></pre>



## 更新record

### 功能

    根据_key更新record；支持批量更新，但只能更新非回收站的record。

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/data

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否

必需 

</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：update。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>用户的access_token，HTTPS调用时必须使用。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>要更新的目标表名。
</td></tr>
<tr>
<td>records
</td><td>json array
</td><td>是
</td><td>需要更新的record。
</td></tr>
<tr>
<td>op
</td><td>string
</td><td>否
</td><td>

*   当值为“merge”时，请求中record不带的column，保持旧值（默认值）；

*   当值为“replace”时，参数中传的record将全量替换整个旧的record。
</td></tr></tbody></table>

说明：

#### 其中records是一个数组，其数组成员结构如下：

<table>

<tbody><tr>
<th scope="col">名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>record
</td><td>json
</td><td>是
</td><td>需要更新的record，只能是一个record，并且必须指定_key。
</td></tr>
<tr>
<td>if-match
</td><td>string
</td><td>否
</td><td>条件更新，防止写操作覆盖了其它client的数据值；为上次获取该item时返回的_mtime属性，只有server端保存的_mtime和用户携带的_mtime一致时，才会进行更新操作。
</td></tr></tbody></table>
 
返回参数 
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>records
</td><td>json array
</td><td>返回服务器端已经处理的records(_key, _mtime)列表，顺序与输入顺序一致；如果一个请求包含多个record，遇到第一个出错record即中止，返回的records只包含已经处理成功的key。
</td></tr></tbody></table>
 
### 示例

请求示例: <pre style="font-family:monospace;">$ cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_update_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
        <span style="color: #009900;">{</span>
                    <span style="color: #3366CC;">"record"</span><span style="color: #339933;">:</span>
                    <span style="color: #009900;">{</span>
                    <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"f44603de003c57d5-1346066442"</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"id"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">85617</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"刘德华"</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"type"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"男歌手"</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"intro"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"香港著名歌手、演员"</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"add_time"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">1340949289</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"language"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"国语"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"粤语"</span><span style="color: #009900;">]</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"tags"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"香港电影金像奖"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"四大天王"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"东亚唱片"</span><span style="color: #009900;">]</span>
                <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
                <span style="color: #3366CC;">"if-match"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">1346066442</span>
         <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
        <span style="color: #009900;">{</span>
                <span style="color: #3366CC;">"record"</span><span style="color: #339933;">:</span>
                <span style="color: #009900;">{</span>
                    <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"1aaef0010c012db7-1346066442"</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"id"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">85618</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"凤凰传奇"</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"type"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"组合"</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"intro"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"中国大陆具有广泛知名度的男女二人音乐组合"</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"add_time"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">1340949289</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"language"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"国语"</span><span style="color: #009900;">]</span><span style="color: #339933;">,</span>
                    <span style="color: #3366CC;">"tags"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"月亮之上"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"最炫民族风"</span><span style="color: #009900;">]</span>
                <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
                <span style="color: #3366CC;">"if-match"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">1346066442</span>
         <span style="color: #009900;">}</span>
    <span style="color: #009900;">]</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
$ curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST
<span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=update&amp;access_token=2.b06c3e00010fdb879d12345dcd5f8545.2587600.134819999.1175746697-238347"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_update_request"</span>
&nbsp;</pre> 
正确响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">200</span> OK
&nbsp;
<span style="color: #009900;">{</span>
  <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span><span style="color: #009900;">[</span>
    <span style="color: #009900;">{</span>
      <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"f44603de003c57d5-1346066442"</span><span style="color: #339933;">,</span>
      <span style="color: #3366CC;">"_mtime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1346066823</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #009900;">{</span>
      <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"1aaef0010c012db7-1346066442"</span><span style="color: #339933;">,</span>
      <span style="color: #3366CC;">"_mtime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1346066824</span>
    <span style="color: #009900;">}</span>
  <span style="color: #009900;">]</span><span style="color: #339933;">,</span>
  <span style="color: #3366CC;">"request_id"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">9201162933</span>
<span style="color: #009900;">}</span></pre>
出错响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">400</span> Bad Request
&nbsp;
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"error_code"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">31430</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"error_msg"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"bad record"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">0</span>
    <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
    <span style="color: #009900;">]</span>
<span style="color: #009900;">}</span></pre>

## 删除record

### 功能

    根据_key，删除record。

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/data

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：delete。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>用户access_token，HTTPS调用时必须使用。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>要删除的目标表名。
</td></tr>
<tr>
<td>records
</td><td>json array
</td><td>是
</td><td>需要删除的record _key 数组。
</td></tr>
<tr>
<td>op
</td><td>string
</td><td>否
</td><td>

*   当值为“permanent”时，永久删除record；无论是普通record还是回收record。

*   当值为“recycled”时，将普通record放进回收站；缺省情况为放进回收站。
</td></tr></tbody></table>

### 说明：

#### 其中records是一个数组，其数组成员结构如下：

<table>

<tbody><tr>
<th scope="col">名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>_key
</td><td>string
</td><td>是
</td><td>需要更新的record _key字段的值。
</td></tr>
<tr>
<td>if-match
</td><td>string
</td><td>否
</td><td>类似update中的条件更新值为上次获取该item时返回的_mtime属性；只有server端保存的_mtime和用户携带的_mtime一致时，才会发生delete。
</td></tr></tbody></table>
 
返回参数 
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误消息。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>records
</td><td>json array
</td><td>返回服务器端已经处理的records(_key, _mtime)列表，顺序与输入顺序一致，如果一个请求包含多个record，遇到第一个出错record即中止，返回的records 只包含已经处理成功的key。
</td></tr></tbody></table>
 
### 示例

请求示例: <pre style="font-family:monospace;">$ cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_delete_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"f44603de003c57d5-1346066442"</span><span style="color: #339933;">,</span>
             <span style="color: #3366CC;">"if-match"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">1346066823</span>
        <span style="color: #009900;">}</span>
    <span style="color: #009900;">]</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
$ curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST
<span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=delete&amp;access_token=2.b06c3e00010fdb879d12345dcd5f8545.2587600.134819999.1175746697-238347"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_delete_request"</span>
&nbsp;</pre> 
正确响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">200</span> OK
&nbsp;
<span style="color: #009900;">{</span>
  <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span><span style="color: #009900;">[</span>
    <span style="color: #009900;">{</span>
      <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"f44603de003c57d5-1346066442"</span><span style="color: #339933;">,</span>
      <span style="color: #3366CC;">"_mtime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1346066823</span>
    <span style="color: #009900;">}</span>
  <span style="color: #009900;">]</span><span style="color: #339933;">,</span>
  <span style="color: #3366CC;">"request_id"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">9494352006</span>
<span style="color: #009900;">}</span></pre>
出错响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">400</span> Bad Request
&nbsp;
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"error_code"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">31430</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"error_msg"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"bad record"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">0</span>
    <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
           <span style="color: #009900;">{</span><span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"7d4febca4a68e763-1344915172"</span><span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #009900;">]</span>
<span style="color: #009900;">}</span></pre>



## 查询record

### 功能

    通过一定条件查询record，只能select非回收站的record。

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/data

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：select。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>用户的access_token，HTTPS调用时必须使用。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>查询的目标表名。
</td></tr>
<tr>
<td>condition
</td><td>json
</td><td>是
</td><td>查询条件，参见查询条件描述。
</td></tr>
<tr>
<td>projection
</td><td>array
</td><td>否
</td><td>指定需要哪些字段，_key为默认返回值。
</td></tr>
<tr>
<td>order_by
</td><td>array
</td><td>否
</td><td>排序字段。
</td></tr>
<tr>
<td>start
</td><td>number
</td><td>否
</td><td>分页用，默认为“0”，范围要求&gt;=0。
</td></tr>
<tr>
<td>limit
</td><td>number
</td><td>否
</td><td>分页用，默认为“100”，范围要求[1, 10000]。
</td></tr></tbody></table>

所支持的查询条件如下表所示：

<table>

<tbody><tr>
<th scope="col">查询条件
</th><th scope="col">类型
</th><th scope="col">表达查询条件
</th><th scope="col">示例
</th></tr>
<tr>
<td>'='
</td><td>number/string
</td><td>表示范围查询=
</td><td>"name": {"=": "刘德华"}
</td></tr>
<tr>
<td>'&lt;'
</td><td>number/string
</td><td>表示范围查询&lt;
</td></tr>
<tr>
<td>'&gt;'
</td><td>number/string
</td><td>表示范围查询&gt;
</td><td>"add_time": {"&gt;": 1340949589}
</td></tr>
<tr>
<td>'&lt;='
</td><td>number/string
</td><td>表示范围查询&lt;=
</td></tr>
<tr>
<td>'&gt;='
</td><td>number/string
</td><td>表示范围查询&gt;=
</td><td>"add_time": {"&gt;=": 1340949589}
</td></tr>
<tr>
<td>'!= '
</td><td>number/string
</td><td>不等于
</td><td>"add_time": {"!=": 1340949589} (coming soon)
</td></tr>
<tr>
<td>'like'
</td><td>string
</td><td>SQL中like语法(不区分大小写)
* 表示0到多个字符，_ 表示一个字符
</td><td>"message": {"like": "%windows%"}
</td></tr>
<tr>
<td>'like_binary'
</td><td>string
</td><td>SQL中binary like语法（区分大小写）
* 表示0到多个字符， _ 表示一个字符
</td><td>"message": {"like_binary": "%Windows%"}
</td></tr>
<tr>
<td>'contain'
</td><td>string
</td><td>包含在数组中
</td><td>"language": {"contain": "国语"}
</td></tr>
<tr>
<td>'in'
</td><td>array
</td><td>in
</td><td>"_key": {"in": ["_key1", "_key2"]}
</td></tr>
<tr>
<td>‘notin’
</td><td>array
</td><td>不在集合中
</td><td>"_key": {"notin": ["_key1", "_key2"]}
</td></tr></tbody></table>

说明：

（1）当“condition”条件为空时，表示获取所有record。

（2）根据key获取一条record，可使用如下condition条件表达：

    "_key": {"=": "385d24b3baef3290-1344915172"}

（3）order_by表示支持的排序方式，它是一个数组，数组元素信息如下：

<table>

<tbody><tr>
<th scope="col">Key
</th><th scope="col">Value
</th><th scope="col">描述
</th></tr>
<tr>
<td>列名
</td><td>asc/desc
</td><td>将某列按照“asc/desc”排序。
</td></tr></tbody></table>
 
返回参数 (JSON格式)
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码，如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示，如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>count
</td><td>number
</td><td>总条目数。
</td></tr>
<tr>
<td>records
</td><td>json array
</td><td>record 数组。
</td></tr></tbody></table>
 
### 示例

请求示例: 

1.	简单查询

<pre style="font-family:monospace;">cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_select_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"condition"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"and"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>                                                                                                                                              
            <span style="color: #009900;">{</span> <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"="</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"刘德华"</span> <span style="color: #009900;">}</span> <span style="color: #009900;">}</span>
        <span style="color: #009900;">]</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"order_by"</span> <span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
        <span style="color: #009900;">{</span><span style="color: #3366CC;">"add_time"</span> <span style="color: #339933;">:</span> <span style="color: #3366CC;">"desc"</span> <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
        <span style="color: #009900;">{</span><span style="color: #3366CC;">"name"</span>     <span style="color: #339933;">:</span> <span style="color: #3366CC;">"asc"</span> <span style="color: #009900;">}</span>
    <span style="color: #009900;">]</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"start"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">0</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"limit"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">10</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST <span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=select&amp;access_token=2.85e37d20acd37c3a5ebc9726bd5606eb.31536000.1384932826.1175746697-309847"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_select_request"</span> <span style="color: #CC0000;">2</span><span style="color: #339933;">&gt;/</span>dev<span style="color: #339933;">/</span><span style="color: #003366; font-weight: bold;">null</span> </pre>

2.	组合查询

<pre style="font-family:monospace;">cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_select_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"condition"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"and"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>                                                                                                                                              
            <span style="color: #009900;">{</span> <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"="</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"刘德华"</span> <span style="color: #009900;">}</span> <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
            <span style="color: #009900;">{</span> <span style="color: #3366CC;">"tags"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"contain"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"四大天王"</span> <span style="color: #009900;">}</span> <span style="color: #009900;">}</span>
        <span style="color: #009900;">]</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"projection"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"name"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"intro"</span><span style="color: #009900;">]</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST <span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=select&amp;access_token=2.85e37d20acd37c3a5ebc9726bd5606eb.31536000.1384932826.1175746697-309847"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_select_request"</span> <span style="color: #CC0000;">2</span><span style="color: #339933;">&gt;/</span>dev<span style="color: #339933;">/</span><span style="color: #003366; font-weight: bold;">null</span> </pre>

3.	select 支持对嵌套属性的查询

<pre style="font-family:monospace;">cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_select_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"condition"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"and"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>                          
            <span style="color: #009900;">{</span> <span style="color: #3366CC;">"top_song.name"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"="</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"爱你一万年"</span> <span style="color: #009900;">}</span> <span style="color: #009900;">}</span>
        <span style="color: #009900;">]</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"projection"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"name"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"intro"</span><span style="color: #009900;">]</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST <span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=select&amp;access_token=2.85e37d20acd37c3a5ebc9726bd5606eb.31536000.1384932826.1175746697-309847"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_select_request"</span> <span style="color: #CC0000;">2</span><span style="color: #339933;">&gt;/</span>dev<span style="color: #339933;">/</span><span style="color: #003366; font-weight: bold;">null</span> <span style="color: #339933;">|</span>.<span style="color: #339933;">/</span>json_decode</pre>

4.	or 条件支持

<pre style="font-family:monospace;">cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_select_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"condition"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"or"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
            <span style="color: #009900;">{</span> <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"="</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"刘德华"</span> <span style="color: #009900;">}</span> <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
            <span style="color: #009900;">{</span> <span style="color: #3366CC;">"top_song.name"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"="</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"月亮之上"</span> <span style="color: #009900;">}</span> <span style="color: #009900;">}</span>
        <span style="color: #009900;">]</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"projection"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"name"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"intro"</span><span style="color: #009900;">]</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST <span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=select&amp;access_token=2.85e37d20acd37c3a5ebc9726bd5606eb.31536000.1384932826.1175746697-309847"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_select_request"</span> <span style="color: #CC0000;">2</span><span style="color: #339933;">&gt;/</span>dev<span style="color: #339933;">/</span><span style="color: #003366; font-weight: bold;">null</span> </pre>

5.	and/or 混合条件支持

<pre style="font-family:monospace;">cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_select_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"condition"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"or"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
            <span style="color: #009900;">{</span> <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"="</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"刘德华"</span> <span style="color: #009900;">}</span> <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
            <span style="color: #009900;">{</span> <span style="color: #3366CC;">"and"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
                <span style="color: #009900;">{</span> <span style="color: #3366CC;">"top_song.name"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"="</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"月亮之上"</span> <span style="color: #009900;">}</span> <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
                <span style="color: #009900;">{</span> <span style="color: #3366CC;">"tags"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"contain"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"月亮之上"</span> <span style="color: #009900;">}</span> <span style="color: #009900;">}</span>
            <span style="color: #009900;">]</span><span style="color: #009900;">}</span>
        <span style="color: #009900;">]</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"projection"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"name"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"intro"</span><span style="color: #009900;">]</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST <span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=select&amp;access_token=2.85e37d20acd37c3a5ebc9726bd5606eb.31536000.1384932826.1175746697-309847"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_select_request"</span> <span style="color: #CC0000;">2</span><span style="color: #339933;">&gt;/</span>dev<span style="color: #339933;">/</span><span style="color: #003366; font-weight: bold;">null</span> </pre>
响应示例：<pre style="font-family:monospace;"><span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"count"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #339933;">,</span>
     <span style="color: #3366CC;">"start"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">0</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"limit"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">0</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
        <span style="color: #009900;">{</span>
           <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"f44603de003c57d5-1346066442"</span><span style="color: #339933;">,</span>
           <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"刘德华"</span><span style="color: #339933;">,</span>
           <span style="color: #3366CC;">"intro"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"香港著名歌手、演员"</span><span style="color: #339933;">,</span>
         <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #009900;">]</span>
<span style="color: #009900;">}</span></pre>



## record增量更新查询

### 功能

    数据更新增量查询接口。 

### HTTP请求方式

    POST

### URL 

    https://pcs.baidu.com/rest/2.0/structure/data 

### 请求参数

<table>
<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：diff。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>用户的access_token，HTTPS调用时必须使用。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>要更新的目标表名。
</td></tr>
<tr>
<td>cursor
</td><td>string
</td><td>是
</td><td>用于标记更新的游标。第一次调用时设置cursor=null，第二次调用时，使用上一次调用该接口的返回结果中的cursor。
</td></tr>
<tr>
<td>projection
</td><td>array
</td><td>否
</td><td>指定需要哪些字段，“_key”、“_mtime”、“_ctime”是默认返回的。
</td></tr></tbody></table>
 返回参数 
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示。如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>表名。
</td></tr>
<tr>
<td>entries
</td><td>array
</td><td>record 数组。
</td></tr>
<tr>
<td>reset
</td><td>boolean
</td><td>客户端是否需要清空本地所有数据。True：表示服务器通知客户端清理所有本地数据，从头获取一份完整的数据列表。
</td></tr>
<tr>
<td>has_more
</td><td>boolean
</td><td>是否还有更新。

*   True：本次调用diff接口结果无法一次返回，立刻再调用一次diff接口获取剩余结果；
*   False：已返回全部更新，等待一段时间（5分钟）之后再调用该接口查看是否有更新。
</td></tr>
<tr>
<td>cursor
</td><td>string
</td><td>游标，下次调用diff 接口，需要使用该参数
</td></tr></tbody></table>

说明：

（1）其中records是一个record数组，标志从上次调用该接口以来的更新操作： 


*   对于update后的record，会得到一个最新版的record；
*   对于删除的record，得到的record中_isdelete字段为“1”。


（2）常见“reset=true”的场景如下： 


*   服务器端程序升级等，提示客户端重新拉去文件列表等。


（3）**注意：** 


#### diff接口有一定延迟（约10s），客户端不可假设新增record之后马上就会在diff 接口中获得更新。 

 示例

请求示例: 

当创建了“刘德华”和“凤凰传奇”两个record后，第一次调用diff接口，使用cursor为null作为参数： 

<pre style="font-family:monospace;">&nbsp;
$ cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_diff1_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"cursor"</span>    <span style="color: #339933;">:</span> <span style="color: #3366CC;">"null"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"projection"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span><span style="color: #3366CC;">"name"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"intro"</span><span style="color: #009900;">]</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
$ curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST
<span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=diff&amp;access_token=2.b06c3e00010fdb879d12345dcd5f8545.2587600.134819999.1175746697-238347"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_diff1_request"</span></pre> 

返回：

<pre style="font-family:monospace;"><span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"entries"</span><span style="color: #339933;">:</span>
    <span style="color: #009900;">[</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"f44603de003c57d5-1346066442"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"刘德华"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"intro"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"香港著名歌手、演员"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"_ctime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1345786801</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"_mtime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1345786801</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"_isdelete"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">0</span>
        <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"1aaef0010c012db7-1346066442"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"凤凰传奇"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"intro"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"中国大陆具有广泛知名度的男女二人音乐组合"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"_ctime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1345787059</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"_isdelete"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">0</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"_mtime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1345787666</span>
        <span style="color: #009900;">}</span>
    <span style="color: #009900;">]</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"cursor"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"3861315431477246513534776b33573367765677314d595a6571724c69753574426133356f4f71316239342b5937577037766874316330493447465a346172445776504e45793235552b456f39796f4c6b307a30447a6f4e68774233616130362f63356f67586d66647879736f72686a70757a575a5342582b4c4b506479325431486f3937526333514a4a6d72626d7830574a35456d46705153454c4873614f6a6368743948743575386b45765477376a634e453848457737522f756d714235464d7a374372574b5777675134423231366a6f3431673d3d"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">4060333081</span>
<span style="color: #009900;">}</span></pre> 

此时如果删除了“刘德华”，再调用diff接口，应该使用刚才的cursor作为参数调用diff接口：

<pre style="font-family:monospace;">$ cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_diff2_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span><span style="color: #3366CC;">"cursor"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"3861315431477246513534776b33573367765677314d595a6571724c69753574426133356f4f71316239342b5937577037766874316330493447465a346172445776504e45793235552b456f39796f4c6b307a30442b6f6b4f585a7672687647426342783858576a524666316470356473754f6a364356674365366f53647979663646356649744747336e50694a44767a7258627a77473078572b327a366c4f674a4c757a76596e3843454b36496b59474153566c4447514c506632704c64494a5a764f337269546d6d454c6869676765367a4866513d3d"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"projection"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span> <span style="color: #3366CC;">"name"</span><span style="color: #339933;">,</span> <span style="color: #3366CC;">"id"</span> <span style="color: #009900;">]</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
curl  <span style="color: #339933;">-</span>i <span style="color: #339933;">-</span>X POST
<span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=diff&amp;access_token=2.b06c3e00010fdb879d12345dcd5f8545.2587600.134819999.1175746697-238347"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=artists_diff2_request"</span></pre> 

得到的响应如下：

<pre style="font-family:monospace;"><span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"entries"</span><span style="color: #339933;">:</span>
    <span style="color: #009900;">[</span>
        <span style="color: #009900;">{</span>
            <span style="color: #3366CC;">"id"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">85617</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"<span style="color: #000099; font-weight: bold;">\u</span>5218<span style="color: #000099; font-weight: bold;">\u</span>5fb7<span style="color: #000099; font-weight: bold;">\u</span>534e"</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"_ctime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1346066442</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"_mtime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1346067702</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"_isdelete"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1</span><span style="color: #339933;">,</span>
            <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"f44603de003c57d51346066442"</span>
         <span style="color: #009900;">}</span>
    <span style="color: #009900;">]</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"has_more"</span><span style="color: #339933;">:</span><span style="color: #003366; font-weight: bold;">false</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"reset"</span><span style="color: #339933;">:</span><span style="color: #003366; font-weight: bold;">false</span><span style="color: #339933;">,</span><span style="color: #3366CC;">"cursor"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"3861315431477246513534776b33573367765677314d595a6571724c69753574426133356f4f71316239342b5937577037766874316330493447465a346172445776504e45793235552b456f39796f4c6b307a304434436368685971496f4e544d446e326a6341505358767251356b79616244597a52745750436d544a327250703956444f7338593752737a4e32735a4d634233372f34416e454e4f3744474c4b4e6d726d64726f774e6b594a4e6553524a716d65743442553375696b354f585738474d376e4f635973507239636c6e5071585141673d3d"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1242060313</span>
<span style="color: #009900;">}</span></pre> 


## 查询record（回收站）

### 功能

    与select相同，只不过操作对象是回收站中。

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/data

### 说明

    该接口参数与select完全一样, 只不过操作对象是回收站中的records；返回的records中_isdelete为“1”。详细信息，请参考“查询record—select”。



## 从回收站中恢复record

### 功能

    从回收站中恢复文件。

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/data

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：restore。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>用户的access_token，HTTPS调用时必须使用。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>要恢复的目标表名。
</td></tr>
<tr>
<td>records
</td><td>json array
</td><td>是
</td><td>需要恢复的record _key数组。
</td></tr></tbody></table>

说明：

#### 其中records是一个数组，每个数组成员结构如下：

<table>

<tbody><tr>
<th scope="col">名称
</th><th scope="col">类型
</th><th scope="col">是否必需
</th><th scope="col">描述
</th></tr>
<tr>
<td>_key
</td><td>string
</td><td>是
</td><td>需要恢复的record _key字段的值。
</td></tr></tbody></table>
 
返回参数 
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码。如果不出错，则返回值没有该字段
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示。如果不出错，则返回值没有该字段
</td></tr>
<tr>
<td>records
</td><td>json array
</td><td>返回服务器端已经处理的records(_key,_mtime)列表，顺序与输入顺序一致，如果一个请求包含多个record，遇到第一个出错record即中止，返回的records只包含已经处理成功的key。
</td></tr></tbody></table>
 
### 示例

请求示例: <pre style="font-family:monospace;">$ cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_restore_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
        <span style="color: #009900;">{</span><span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"7d4febca4a68e763-1344915172"</span><span style="color: #009900;">}</span><span style="color: #339933;">,</span>
        <span style="color: #009900;">{</span><span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"385d24b3baef3290-1344915172"</span><span style="color: #009900;">}</span>
    <span style="color: #009900;">]</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
$ curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST
<span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=restore&amp;access_token=2.b06c3e00010fdb879d12345dcd5f8545.2587600.134819999.1175746697-238347"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_restore_request"</span>
&nbsp;</pre> 
正确响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">200</span> OK
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
        <span style="color: #009900;">{</span>
          <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"7d4febca4a68e763-1344915172"</span><span style="color: #339933;">,</span>
          <span style="color: #3366CC;">"_mtime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1344927006</span>
        <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
        <span style="color: #009900;">{</span>
          <span style="color: #3366CC;">"_key"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"385d24b3baef3290-1344915172"</span><span style="color: #339933;">,</span>
          <span style="color: #3366CC;">"_mtime"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">1344927006</span>
        <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
&nbsp;
    <span style="color: #009900;">]</span>
<span style="color: #009900;">}</span></pre>
出错响应示例：<pre style="font-family:monospace;">HTTP<span style="color: #339933;">/</span><span style="color: #CC0000;">1.1</span> <span style="color: #CC0000;">400</span> Bad Request
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"error_code"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">31430</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"error_msg"</span><span style="color: #339933;">:</span><span style="color: #3366CC;">"key not exist"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"request_id"</span><span style="color: #339933;">:</span><span style="color: #CC0000;">0</span>
    <span style="color: #3366CC;">"records"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
        <span style="color: #009900;">{</span><span style="color: #3366CC;">'_key'</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"7d4febca4a68e763-1344915172"</span><span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #009900;">]</span>
<span style="color: #009900;">}</span></pre>



## 按条件更新record

### 功能

    对符合一定条件的record 执行更新操作。

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/data

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否

必需 

</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：update
</td></tr>
<tr>
<td>type
</td><td>string
</td><td>是
</td><td>固定值：by-condition。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>用户的access_token，HTTPS调用时必须使用。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>查询的目标表名。
</td></tr>
<tr>
<td>condition
</td><td>json
</td><td>是
</td><td>条件描述, 与select 中的condition一样。
</td></tr>
<tr>
<td>action
</td><td>json
</td><td>是
</td><td>需要对命中的record进行的操作。
</td></tr></tbody></table>

#### 说明：

##### action为一个json字典，其格式为：

    "action": {                                                                     
        column: {action: value}
    }

##### 如: 

    "action": {                                                                     
        "name": {"=": "LiuDeHua"}
    }

#### 其中column 支持嵌套列。

*   所支持的action如下表所示：


<table>

<tbody><tr>
<th scope="col">action
</th><th scope="col">类型
</th><th scope="col">描述
</th><th scope="col">示例
</th></tr>
<tr>
<td>'='
</td><td>number/string
</td><td>表示将目标列设置为value。
</td><td>"name": {"=": "LiuDeHua"},
</td></tr>
<tr>
<td>'+='
</td><td>number
</td><td>表示将目标列的值增加value，如果该列不存在，默认值为0。
</td><td>"age": {"+=":1},
</td></tr>
<tr>
<td>'-='
</td><td>number
</td><td>表示将目标列的值减少value，如果该列不存在，默认值为0。
</td><td>"age": {"-=":1},
</td></tr></tbody></table>
 
返回参数 (JSON格式)
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码，如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示，如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>request_id
</td><td>number
</td><td>请求唯一标识ID。
</td></tr>
<tr>
<td>affected
</td><td>number
</td><td>返回受影响的行数。
</td></tr></tbody></table>
 
### 示例

请求示例: <pre style="font-family:monospace;">cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_update_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"condition"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"and"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
            <span style="color: #009900;">{</span> <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"="</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"刘德华"</span> <span style="color: #009900;">}</span> <span style="color: #009900;">}</span>
        <span style="color: #009900;">]</span>
    <span style="color: #009900;">}</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"action"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span><span style="color: #3366CC;">"="</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"LiuDeHua"</span><span style="color: #009900;">}</span><span style="color: #339933;">,</span> 
        <span style="color: #3366CC;">"age"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"+="</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span> <span style="color: #009900;">}</span>
    <span style="color: #009900;">}</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST <span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=update&amp;type=by-condition&amp;access_token=2.85e37d20acd37c3a5ebc9726bd5606eb.31536000.1384932826.1175746697-309847"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_update_request"</span></pre> 
响应示例：<pre style="font-family:monospace;"><span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"affected"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #339933;">,</span> 
    <span style="color: #3366CC;">"request_id"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">4060311005</span>
<span style="color: #009900;">}</span></pre>



## 按条件删除record

### 功能

    对符合一定条件的record 执行删除操作。

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/data

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否

必需 

</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：delete。
</td></tr>
<tr>
<td>type
</td><td>string
</td><td>是
</td><td>固定值：by-condition。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>用户的access_token，HTTPS调用时必须使用。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>查询的目标表名。
</td></tr>
<tr>
<td>condition
</td><td>json
</td><td>是
</td><td>条件描述，与select 中的condition一样。
</td></tr></tbody></table>
 
返回参数 (JSON格式)
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码，如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示，如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>request_id
</td><td>number
</td><td>请求唯一标识ID。
</td></tr>
<tr>
<td>affected
</td><td>number
</td><td>返回受影响的行数。
</td></tr></tbody></table>
 
### 示例

请求示例: <pre style="font-family:monospace;"> cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_update_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"condition"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"and"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
            <span style="color: #009900;">{</span> <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"="</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"LiuDeHua"</span> <span style="color: #009900;">}</span> <span style="color: #009900;">}</span>
        <span style="color: #009900;">]</span>
    <span style="color: #009900;">}</span>
<span style="color: #009900;">}</span>
DELIM
&nbsp;
curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST <span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=delete&amp;type=by-condition&amp;access_token=2.85e37d20acd37c3a5ebc9726bd5606eb.31536000.1384932826.1175746697-309847"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_update_request"</span> </pre> 
响应示例：<pre style="font-family:monospace;"><span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"affected"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #339933;">,</span> 
    <span style="color: #3366CC;">"request_id"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">4060311005</span>
<span style="color: #009900;">}</span></pre>



## 按条件恢复record

### 功能

    对回收站中符合一定条件的record 执行restore操作。

### HTTP请求方式

    POST

### URL

    https://pcs.baidu.com/rest/2.0/structure/data

### 请求参数 

<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">是否

必需 

</th><th scope="col">描述
</th></tr>
<tr>
<td>v
</td><td>string
</td><td>否
</td><td>版本号，默认为“1.0”。
</td></tr>
<tr>
<td>method
</td><td>string
</td><td>是
</td><td>固定值：restore。
</td></tr>
<tr>
<td>type
</td><td>string
</td><td>是
</td><td>固定值：by-condition。
</td></tr>
<tr>
<td>access_token
</td><td>string
</td><td>是
</td><td>用户的access_token，HTTPS调用时必须使用。
</td></tr>
<tr>
<td>table
</td><td>string
</td><td>是
</td><td>查询的目标表名。
</td></tr>
<tr>
<td>condition
</td><td>json
</td><td>是
</td><td>条件描述, 与select 中的condition一样。
</td></tr></tbody></table>
 
返回参数 (JSON格式)
<table>

<tbody><tr>
<th scope="col">参数名称
</th><th scope="col">类型
</th><th scope="col">描述
</th></tr>
<tr>
<td>error_code
</td><td>number
</td><td>错误码，如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>error_msg
</td><td>string
</td><td>错误提示，如果不出错，则返回值没有该字段。
</td></tr>
<tr>
<td>request_id
</td><td>number
</td><td>请求唯一标识ID。
</td></tr>
<tr>
<td>affected
</td><td>number
</td><td>返回受影响的行数。
</td></tr></tbody></table>
 
### 示例

请求示例: <pre style="font-family:monospace;">cat <span style="color: #339933;">&gt;</span> .<span style="color: #339933;">/</span>artists_update_request <span style="color: #339933;">&lt;&lt;</span>DELIM
<span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"table"</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"artists"</span><span style="color: #339933;">,</span>
    <span style="color: #3366CC;">"condition"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span>
        <span style="color: #3366CC;">"and"</span><span style="color: #339933;">:</span> <span style="color: #009900;">[</span>
            <span style="color: #009900;">{</span> <span style="color: #3366CC;">"name"</span><span style="color: #339933;">:</span> <span style="color: #009900;">{</span> <span style="color: #3366CC;">"="</span><span style="color: #339933;">:</span> <span style="color: #3366CC;">"LiuDeHua"</span> <span style="color: #009900;">}</span> <span style="color: #009900;">}</span>
        <span style="color: #009900;">]</span>
    <span style="color: #009900;">}</span>
<span style="color: #009900;">}</span>
DELIM
curl  <span style="color: #339933;">-</span>v <span style="color: #339933;">-</span>X POST <span style="color: #3366CC;">"http://pcs.baidu.com/rest/2.0/structure/data?method=restore&amp;type=by-condition&amp;access_token=2.85e37d20acd37c3a5ebc9726bd5606eb.31536000.1384932826.1175746697-309847"</span> <span style="color: #339933;">-</span>F <span style="color: #3366CC;">"param=&lt;artists_update_request"</span></pre> 
响应示例：<pre style="font-family:monospace;"><span style="color: #009900;">{</span>
    <span style="color: #3366CC;">"affected"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">1</span><span style="color: #339933;">,</span> 
    <span style="color: #3366CC;">"request_id"</span><span style="color: #339933;">:</span> <span style="color: #CC0000;">4060311005</span>
<span style="color: #009900;">}</span></pre>
