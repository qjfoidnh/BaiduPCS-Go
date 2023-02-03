# BaiduPCS-Go 百度网盘客户端(加强版)


仿 Linux shell 文件处理命令的百度网盘命令行客户端.

iikira/BaiduPCS-Go was largely inspired by [GangZhuo/BaiduPCS](https://github.com/GangZhuo/BaiduPCS) and this project was largely based on iikira/BaiduPCS-Go

## 注意

此版本基于iikira原版BaiduPCS-Go v3.6.2继续开发, 并添加了转存功能.

本软件不提供超出官方客户端的下载提速, 普通用户和SVIP的配置建议参见 [显示和修改程序配置项](#显示和修改程序配置项)

<!-- toc -->
## 目录

- [特色](#特色)
- [版本更新](#版本更新)
- [编译/交叉编译 说明](#编译交叉编译-说明)
- [下载/运行 说明](#下载运行-说明)
  * [安装](#安装)
  * [Windows](#windows)
  * [Linux / macOS](#linux--macos)
  * [Android / iOS](#android--ios)
- [命令列表及说明](#命令列表及说明)
  * [注意 ! ! !](#注意---)
  * [检测程序更新](#检测程序更新)
  * [登录百度帐号](#登录百度帐号)
  * [列出帐号列表](#列出帐号列表)
  * [获取当前帐号](#获取当前帐号)
  * [切换百度帐号](#切换百度帐号)
  * [退出百度帐号](#退出百度帐号)
  * [获取网盘配额](#获取网盘配额)
  * [切换工作目录](#切换工作目录)
  * [输出工作目录](#输出工作目录)
  * [列出目录](#列出目录)
  * [列出目录树形图](#列出目录树形图)
  * [获取文件/目录的元信息](#获取文件目录的元信息)
  * [搜索文件](#搜索文件)
  * [下载文件/目录](#下载文件目录)
  * [上传文件/目录](#上传文件目录)
  * [获取下载直链](#获取下载直链)
  * [手动秒传文件](#手动秒传文件)
  * [修复文件MD5](#修复文件MD5)
  * [获取本地文件的秒传信息](#获取本地文件的秒传信息)
  * [导出文件/目录](#导出文件目录)
  * [创建目录](#创建目录)
  * [删除文件/目录](#删除文件目录)
  * [拷贝文件/目录](#拷贝文件目录)
  * [移动/重命名文件/目录](#移动重命名文件目录)
  * [转存文件/目录](#转存文件目录)
  * [分享文件/目录](#分享文件目录)
    + [设置分享文件/目录](#设置分享文件目录)
    + [列出已分享文件/目录](#列出已分享文件目录)
    + [取消分享文件/目录](#取消分享文件目录)
  * [离线下载](#离线下载)
    + [添加离线下载任务](#添加离线下载任务)
    + [精确查询离线下载任务](#精确查询离线下载任务)
    + [查询离线下载任务列表](#查询离线下载任务列表)
    + [取消离线下载任务](#取消离线下载任务)
    + [删除离线下载任务](#删除离线下载任务)
  * [回收站](#回收站)
    + [列出回收站文件列表](#列出回收站文件列表)
    + [还原回收站文件或目录](#还原回收站文件或目录)
    + [删除回收站文件或目录/清空回收站](#删除回收站文件或目录清空回收站)
  * [显示和修改程序配置项](#显示和修改程序配置项)
  * [测试通配符](#测试通配符)
  * [工具箱](#工具箱)
- [初级使用教程](#初级使用教程)
  * [1. 查看程序使用说明](#1-查看程序使用说明)
  * [2. 登录百度帐号 (必做)](#2-登录百度帐号-必做)
  * [3. 切换网盘工作目录](#3-切换网盘工作目录)
  * [4. 网盘内列出文件和目录](#4-网盘内列出文件和目录)
  * [5. 下载文件](#5-下载文件)
  * [6. 设置下载最大并发量](#6-设置下载最大并发量)
  * [7. 恢复默认配置](#7-恢复默认配置)
  * [8. 退出程序](#8-退出程序)
- [已知问题](#已知问题)
- [TODO](#todo)
- [交流反馈](#交流反馈)

<!-- tocstop -->

# 特色

多平台支持, 支持 Windows, macOS, linux, 移动设备等.

百度帐号多用户支持;

通配符匹配网盘路径和 Tab 自动补齐命令和路径, [通配符_百度百科](https://baike.baidu.com/item/通配符);

[下载](#下载文件目录)网盘内文件, 支持多个文件或目录下载, 支持断点续传和单文件并行下载;

[上传](#上传文件目录)本地文件, 支持上传大文件(>2GB), 支持多个文件或目录上传;

[转存](#转存文件目录)其他用户分享的文件, 支持公开、带密码的分享链接, 支持常见的几种秒传链接;

[导出](#导出文件目录)网盘内的文件秒传链接, 可选导出BaiduPCS-Go原生格式或通用格式;

[离线下载](#离线下载), 支持http/https/ftp/电驴/磁力链协议.

# 版本更新
**2022.12.04** v3.9.0:
- 优化转存错误提示
- fix #239
- update go version to 1.18

**2022.11.25** v3.8.9:
- fix #234, 继续修复无法转存文件

**2022.11.12** v3.8.8:
- fix #234, 修复无法转存文件

**2022.2.18** v3.8.7:
- fix #175, 在正式上传前即进行文件大小检测

**2022.2.14** v3.8.6:
- fix #160 #173, 修复上传出现空文件的bug
- fix #165, 支持自带提取码的转存链接
- fix #175, upload增加-policy=rsync策略, 配合--norapid使用, 只跳过大小未发生改变的文件
- 鉴于 #172, 建议下载线程数最大不超过12

**2022.1.1** v3.8.5:
#### 该版本存在已知问题将导致上传文件失败及出现空文件，建议跳过更新
- 2022新年好, 本次更新增加较多特性, 欢迎测试
- fix #146, 提前fail和skip上传策略中重名文件的检测环节（存在问题）
- fix #158, config可配置关闭文件名合法性检测
- fix #141, download增加--mtime选项可保持文件修改时间
- fix #130, config可配置force_login_username, 强制登录指定用户名
- 首条下载链接不可用时自动切换, 增加下载成功率

**2021.10.6** v3.8.4:
- fix 登录时可能出现内存溢出
- 上传文件名允许包含单引号

**2021.8.27** v3.8.3:
- fix 更换默认panUA解决svip限速
- fix 移除失效的秒传修复功能 
- 优化秒传逻辑, 提高成功率
- 优化秒传导出逻辑, 提高新文件的导出成功率  

**2021.7.20** v3.8.2:
- fix 读取大量文件信息容易超时
- fix 秒传链接文件名带"#"时解析错误
- share list增加分享下载数显示  
- config增加配置: 上传的同名文件处理策略

**2021.6.9** v3.8.1:
- fix 部分旧链接无法转存
- 增加上传同名文件自动跳过选项

**2021.5.21** v3.8.0:
- fix 上传到100M左右自动回滚（待测试）
- fix 个别正常的秒传链接无法转存
- fix 文件名含有百分号导出异常
- 优化上传重试策略（待测试）

**2021.4.14** v3.7.9:
- fix 上传时异常退出导致无法加载断点信息
- fix 上传偶发出现0B/s卡住
- 上传时预先检查文件名合法性
- 在线更新使用镜像源加速

**2021.3.20** v3.7.8:

- 优化了上传的输出信息格式
- 优化了上传逻辑，提升上传速度
- transfer增加--fix参数，可转存被屏蔽的秒传链接（inspired by [dupan-rapid-extract](https://github.com/mengzonefire/dupan-rapid-extract)）

**2021.3.11** v3.7.7:

- fix 移动和重命名文件时末尾```/```导致报错
- fix 3.7.2版本后在线升级无效
- fix 转存误报缺少STOKEN

**2021.2.23** v3.7.6:

- fix 下载文件报```x509: certificate is valid```错误
- 完善了下载错误的捕获种类
- download增加--fullpath参数，本地目录保留网盘从根目录开始的完整结构

**2021.2.8** v3.7.5:

- fix 某些时候误报stoken缺失
- fix windows平台上秒传链接转存失败
- fix 某些时候pcs请求缺少Host
- 当分享链接包含多文件/目录时，可选归档到第一个文件命名的目录里（不支持秒传）

**2021.1.31** v3.7.4:

- fix 下载目录会丢失目录结构
- fix 分享列表状态信息显示错误
- 支持自定义文件上传服务器

**2021.1.22** v3.7.3:

- 分享支持自定义分享码和有效天数
- 转存支持转存完毕后自动下载到默认目录
- 增加恢复默认配置功能
- tree命令支持指定输出最大层数和带fsid输出

**2021.1.9** v3.7.2:

- 基本修复了登录验证失效问题([#15](https://github.com/qjfoidnh/BaiduPCS-Go/issues/15))
- 优化下载模块的实现策略, 保证稳定性同时进一步提升下载速度 (需按[显示和修改程序配置项](#显示和修改程序配置项)中建议修改)
- update 功能恢复, 以后可以在线升级了
- 支持导出秒传链接不写文件, 直接输出到控制台; 支持通用秒传格式导出, 具体参见export --help
- 其他bug修正

**2021.1.2** v3.7.1:

- 支持了多文件并发上传，文件并发数和单文件分片数可在配置中指定
- 修复了最大同时下载文件数配置不生效的问题
- 修正了部分显示和帮助的错误

**2020.12.19** v3.7.0:

* 替换了iikira版本的失效仓库
* 转存功能支持旧的短链接
* 默认关闭下载文件校验，配置文件可设置开启
* 修复了关闭校验时会误报下载失败的问题
* 转存功能除了cookies方式登录，现已支持用户名密码登录和bduss登录；bduss登录需同时指定stoken

**2020.11.08** v3.6.3:

* 修复转存失败
* 修复分享文件失败


# 编译/交叉编译 说明
设置好 GOOS 和 GOARCH 环境变量,

运行 go tool dist list 查看所有支持的 GOOS/GOARCH

## Linux/Darwin 例子: 编译 Windows 下的 64 位程序
```
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build
```
## Windows 例子: 编译 Linux 下的 32 位程序
```
set GOOS=linux
set GOARCH=386
set CGO_ENABLED=0
go build
```

# 下载/运行 说明

Go语言程序, 常用几种平台的已编译程序可直接在[蓝奏云](https://wws.lanzoui.com/b01berebe)下载使用. 密码:4pix

如果程序运行时输出乱码, 请检查下终端的编码方式是否为 `UTF-8`.

使用本程序之前, 建议学习一些 linux 基础知识 和 基础命令.

如果未带任何参数运行程序, 程序将会进入仿Linux shell系统用户界面的cli交互模式, 可直接运行相关命令.

cli交互模式下, 光标所在行的前缀应为 `BaiduPCS-Go >`, 如果登录了百度帐号则格式为 `BaiduPCS-Go:<工作目录> <百度ID>$ `

程序会提供相关命令的使用说明.

## 安装

[Homebrew](https://brew.sh/) 用户可以用下面的命令行安装应用：

```sh
brew install baidupcs-go
```

## Windows

程序应在 命令提示符 (Command Prompt) 或 PowerShell 中运行, 在 mintty (例如: GitBash) 可能会有显示问题.

也可直接双击程序运行, 具体使用方法请参见 [命令列表及说明](#命令列表及说明) 和 [初级使用教程](#初级使用教程).

## Linux / macOS

程序应在 终端 (Terminal) 运行.

具体使用方法请参见 [命令列表及说明](#命令列表及说明) 和 [初级使用教程](#初级使用教程).

## Android / iOS

> Android / iOS 移动设备操作比较麻烦, 不建议在移动设备上使用本程序. 移动设备不可直接使用预编译的Linux arm64版本, 使用者需下载源码自行交叉编译.

安卓, 建议使用 [Termux](https://termux.com) 或 [NeoTerm](https://github.com/NeoTerm/NeoTerm) 或 终端模拟器, 以提供终端环境.

示例: [Android 运行本项目程序参考示例](https://web.archive.org/web/20190820154934/https://github.com/iikira/BaiduPCS-Go/wiki/Android-%E8%BF%90%E8%A1%8C%E6%9C%AC%E9%A1%B9%E7%9B%AE%E7%A8%8B%E5%BA%8F%E5%8F%82%E8%80%83%E7%A4%BA%E4%BE%8B), 有兴趣的可以参考一下.

苹果iOS, 需要越狱, 在 Cydia 搜索下载并安装 MobileTerminal, 或者其他提供终端环境的软件.

示例: [iOS 运行本项目程序参考示例](https://web.archive.org/web/20190820155025/https://github.com/iikira/BaiduPCS-Go/wiki/iOS-%E8%BF%90%E8%A1%8C%E6%9C%AC%E9%A1%B9%E7%9B%AE%E7%A8%8B%E5%BA%8F%E5%8F%82%E8%80%83%E7%A4%BA%E4%BE%8B), 有兴趣的可以参考一下.

具体使用方法请参见 [命令列表及说明](#命令列表及说明) 和 [初级使用教程](#初级使用教程).

# 命令列表及说明

## 注意 ! ! !

命令的前缀 `BaiduPCS-Go` 为指向程序运行的全路径名 (ARGv 的第一个参数)

直接运行程序时, 未带任何其他参数, 则程序进入cli交互模式, 运行以下命令时, 要把命令的前缀 `BaiduPCS-Go` 去掉!

cli交互模式已支持按tab键自动补全命令和路径.

## 检测程序更新
```
BaiduPCS-Go update
```

## 登录百度帐号

### 常规登录百度帐号

支持在线验证绑定的手机号或邮箱,
```
BaiduPCS-Go login
```

### 使用百度 BDUSS 来登录百度帐号

[关于 获取百度 BDUSS](https://blog.csdn.net/ykiwmy/article/details/103730962)

```
BaiduPCS-Go login -bduss=<BDUSS>
```

### 使用百度 BDUSS 和 百度网盘 STOKEN 来登录百度账号

STOKEN 获取方式与 BDUSS 基本相同。注意 STOKEN 必须在百度网盘页面获取，否则无效.

```
BaiduPCS-Go login -bduss=<BDUSS> -stoken=<STOKEN>
```

### 使用百度 Cookies 来登录百度账号

[关于 获取百度 Cookies](https://jingyan.baidu.com/article/5553fa829a6a9e65a23934b0.html)
教程中为百度经验的Cookies获取, 这里换成百度网盘首页即可.

```
BaiduPCS-Go login -cookies=<Cookies>
```

#### 例子
```
BaiduPCS-Go login -bduss=1234567
```
```
BaiduPCS-Go login
请输入百度用户名(手机号/邮箱/用户名), 回车键提交 > 1234567
```
```
BaiduPCS-Go login -cookies="BAIDUID=50949C0890YG9735EA6Q3870AFE38:FG=1; BIDUPSID=112335C0ACCAFFJW675EA69A870AFE38; PSTM=1981928511; BDORZ=D6745EBF6F3SW24E515D22A1598; PANWEB=1; BDUSS=ASAYUGFHSTFKGBGSU; STOKEN=gfsdge9gisfgspig34254d7879eee5756b10sgeyrw5vyw342td510ffc9414d32251; SCRC=cwrywec5evyetra26bvvehefvfg6a8; BDCLND=C%4sfgGysrZ%2BML6; PANPSC=wreyewygdfhdggedhsdfg4353"
```

## 列出帐号列表

```
BaiduPCS-Go loglist
```

列出所有已登录的百度帐号

## 获取当前帐号

```
BaiduPCS-Go who
```

## 切换百度帐号

切换已登录的百度帐号
```
BaiduPCS-Go su <uid>
```
```
BaiduPCS-Go su

请输入要切换帐号的 # 值 >
```

## 退出百度帐号

退出当前登录的百度帐号
```
BaiduPCS-Go logout
```

程序会进一步确认退出帐号, 防止误操作.

## 获取网盘配额

```
BaiduPCS-Go quota
```
获取网盘的总储存空间, 和已使用的储存空间

## 切换工作目录
```
BaiduPCS-Go cd <目录>
```

### 切换工作目录后自动列出工作目录下的文件和目录
```
BaiduPCS-Go cd -l <目录>
```

#### 例子
```
# 切换 /我的资源 工作目录
BaiduPCS-Go cd /我的资源

# 切换 上级目录
BaiduPCS-Go cd ..

# 切换 根目录
BaiduPCS-Go cd /

# 切换 /我的资源 工作目录, 并自动列出 /我的资源 下的文件和目录
BaiduPCS-Go cd -l 我的资源

# 使用通配符
BaiduPCS-Go cd /我的*
```

## 输出工作目录
```
BaiduPCS-Go pwd
```

## 列出目录

列出当前工作目录的文件和目录或指定目录
```
BaiduPCS-Go ls
```
```
BaiduPCS-Go ls <目录>
```

### 可选参数
```
-asc: 升序排序
-desc: 降序排序
-time: 根据时间排序
-name: 根据文件名排序
-size: 根据大小排序
```

#### 例子
```
# 列出 我的资源 内的文件和目录
BaiduPCS-Go ls 我的资源

# 绝对路径
BaiduPCS-Go ls /我的资源

# 降序排序
BaiduPCS-Go ls -desc 我的资源

# 按文件大小降序排序
BaiduPCS-Go ls -size -desc 我的资源

# 使用通配符
BaiduPCS-Go ls /我的*
```

## 列出目录树形图

列出当前工作目录的文件和目录或指定目录的树形图
```
BaiduPCS-Go tree <目录>

# 默认获取工作目录元信息
BaiduPCS-Go tree
```

## 获取文件/目录的元信息
```
BaiduPCS-Go meta <文件/目录1> <文件/目录2> <文件/目录3> ...

# 默认获取工作目录元信息
BaiduPCS-Go meta
```

#### 例子
```
BaiduPCS-Go meta 我的资源
BaiduPCS-Go meta /
```

## 搜索文件

按文件名搜索文件（不支持查找目录）。

默认在当前工作目录搜索.

```
BaiduPCS-Go search [-path=<需要检索的目录>] [-r] <关键字>
```

#### 例子
```
# 搜索根目录的文件
BaiduPCS-Go search -path=/ 关键字

# 搜索当前工作目录的文件
BaiduPCS-Go search 关键字

# 递归搜索当前工作目录的文件
BaiduPCS-Go search -r 关键字
```

## 下载文件/目录
```
BaiduPCS-Go download <网盘文件或目录的路径1> <文件或目录2> <文件或目录3> ...
BaiduPCS-Go d <网盘文件或目录的路径1> <文件或目录2> <文件或目录3> ...
```

### 可选参数
```
  --test          测试下载, 此操作不会保存文件到本地
  --ow            overwrite, 覆盖已存在的文件
  --status        输出所有线程的工作状态
  --save          将下载的文件直接保存到当前工作目录
  --saveto value  将下载的文件直接保存到指定的目录
  -x              为文件加上执行权限, (windows系统无效)
  --mode value    下载模式, 可选值: pcs, stream, locate, 默认为 locate, 相关说明见上面的帮助 (default: "locate")
  -p value        指定下载线程数 (default: 0)
  -l value        指定同时进行下载文件的数量 (default: 0)
  --retry value   下载失败最大重试次数 (default: 3)
  --nocheck       下载文件完成后不校验文件

```

下载的文件默认保存到 **程序所在目录** 的 download/ 目录, 支持设置指定目录, 重名的文件会自动跳过!

下载的文件默认保存到, **程序所在目录**的 **download/** 目录.

通过 `BaiduPCS-Go config set -savedir <savedir>`, 自定义保存的目录.

支持多个文件或目录下载.

支持下载完成后自动校验文件, 但并不是所有的文件都支持校验!

自动跳过下载重名的文件!


#### 下载模式说明

* pcs: 通过百度网盘的 PCS API 下载(不建议使用)

* stream: 通过百度网盘的 PCS API, 以流式文件的方式下载, 效果同 pcs(不建议使用)

* locate: 默认的下载模式。从百度网盘 Android 客户端, 获取下载链接的方式来下载

#### 例子
```
# 设置保存目录, 保存到 D:\Downloads
# 注意区别反斜杠 "\" 和 斜杠 "/" !!!
BaiduPCS-Go config set -savedir D:/Downloads

# 下载 /我的资源/1.mp4
BaiduPCS-Go d /我的资源/1.mp4

# 下载 /我的资源 整个目录!!
BaiduPCS-Go d /我的资源

# 下载网盘内的全部文件!!
BaiduPCS-Go d /
BaiduPCS-Go d *
```

## 上传文件/目录
```
BaiduPCS-Go upload <本地文件/目录的路径1> <文件/目录2> <文件/目录3> ... <目标目录>
BaiduPCS-Go u <本地文件/目录的路径1> <文件/目录2> <文件/目录3> ... <目标目录>
```

* 上传默认采用分片上传的方式, 上传的文件将会保存到, <目标目录>.

* 遇到同名文件将会自动覆盖!!

* 当上传的文件名和网盘的目录名称相同时, 不会覆盖目录, 防止丢失数据.


#### 注意:

* 分片上传之后, 服务器可能会记录到错误的文件md5, 可使用 fixmd5 命令尝试修复文件的MD5值, 修复md5不一定能成功, 但文件的完整性是没问题的.

fixmd5 命令使用方法:
```
BaiduPCS-Go fixmd5 -h
```

* 禁用分片上传可以保证服务器记录到正确的md5.

* 禁用分片上传时只能使用单线程上传, 指定的单个文件上传最大线程数将会无效.

#### 例子:
```
# 将本地的 C:\Users\Administrator\Desktop\1.mp4 上传到网盘 /视频 目录
# 注意区别反斜杠 "\" 和 斜杠 "/" !!!
BaiduPCS-Go upload C:/Users/Administrator/Desktop/1.mp4 /视频

# 将本地的 C:\Users\Administrator\Desktop\1.mp4 和 C:\Users\Administrator\Desktop\2.mp4 上传到网盘 /视频 目录
BaiduPCS-Go upload C:/Users/Administrator/Desktop/1.mp4 C:/Users/Administrator/Desktop/2.mp4 /视频

# 将本地的 C:\Users\Administrator\Desktop 整个目录上传到网盘 /视频 目录
BaiduPCS-Go upload C:/Users/Administrator/Desktop /视频
```

## 获取下载直链
```
BaiduPCS-Go locate <文件1> <文件2> ...
```

#### 注意

若该功能无法正常使用, 提示`user is not authorized, hitcode:xxx`, 尝试更换 User-Agent 为 `netdisk;2.2.51.6;netdisk;10.0.63;PC;android-android`:
```
BaiduPCS-Go config set -user_agent "netdisk;2.2.51.6;netdisk;10.0.63;PC;android-android"
```

## 手动秒传文件
```
BaiduPCS-Go rapidupload -length=<文件的大小> -md5=<文件的md5值> -slicemd5=<文件前256KB切片的md5值(可选)> -crc32=<文件的crc32值(可选)> <保存的网盘路径, 需包含文件名>
BaiduPCS-Go ru -length=<文件的大小> -md5=<文件的md5值> -slicemd5=<文件前256KB切片的md5值(可选)> -crc32=<文件的crc32值(可选)> <保存的网盘路径, 需包含文件名>
```

注意: 使用此功能秒传文件, 前提是知道文件的大小, md5, 前256KB切片的 md5 (可选), crc32 (可选), 且百度网盘中存在一模一样的文件.

上传的文件将会保存到网盘的目标目录.

遇到同名文件将会自动覆盖!

可能无法秒传 20GB 以上的文件!!

#### 例子:
```
# 如果秒传成功, 则保存到网盘路径 /test
BaiduPCS-Go rapidupload -length=56276137 -md5=fbe082d80e90f90f0fb1f94adbbcfa7f -slicemd5=38c6a75b0ec4499271d4ea38a667ab61 -crc32=314332359 /test
```


## 修复文件MD5
```
BaiduPCS-Go fixmd5 <文件1> <文件2> <文件3> ...
```

尝试修复文件的MD5值, 以便于校验文件的完整性和导出文件.

使用分片上传文件, 当文件分片数大于1时, 百度网盘服务端最终计算所得的md5值和本地的不一致, 这可能是百度网盘的bug.

不过把上传的文件下载到本地后，对比md5值是匹配的, 也就是文件在传输中没有发生损坏.

对于MD5值可能有误的文件, 程序会在获取文件的元信息时, 给出MD5值 "可能不正确" 的提示, 表示此文件可以尝试进行MD5值修复.

修复文件MD5不一定能成功, 原因可能是服务器未刷新, 可过几天后再尝试.

修复文件MD5的原理为秒传文件, 即修复文件MD5成功后, 文件的**创建日期, 修改日期, fs_id, 版本历史等信息**将会被覆盖, 修复的MD5值将覆盖原先的MD5值, 但不影响文件的完整性.

注意: 无法修复 **20GB** 以上文件的 md5!!

#### 例子:
```
# 修复 /我的资源/1.mp4 的 MD5 值
BaiduPCS-Go fixmd5 /我的资源/1.mp4
```

## 获取本地文件的秒传信息
```
BaiduPCS-Go sumfile <本地文件的路径>
BaiduPCS-Go sf <本地文件的路径>
```

获取本地文件的大小, md5, 前256KB切片的 md5, crc32, 可用于秒传文件.

#### 例子:
```
# 获取 C:\Users\Administrator\Desktop\1.mp4 的秒传信息
BaiduPCS-Go sumfile C:/Users/Administrator/Desktop/1.mp4
```

## 导出文件/目录
```
BaiduPCS-Go export <文件/目录1> <文件/目录2> ...
BaiduPCS-Go ep <文件/目录1> <文件/目录2> ...
```

导出网盘内的文件或目录, 原理为秒传文件, 此操作会生成导出文件或目录的命令.

#### 注意

**无法导出 20GB 以上的文件!!**

**无法导出文件的版本历史等数据!!**

**以通用秒传格式导出会丢失文件路径信息!!**

并不是所有的文件都能导出成功, 程序会列出无法导出的文件列表

#### 例子:
```
# 导出当前工作目录:
BaiduPCS-Go export

# 导出所有文件和目录, 并设置新的根目录为 /root
BaiduPCS-Go export -root=/root /

# 导出 /我的资源
BaiduPCS-Go export /我的资源

# 导出 /我的资源 格式为通用秒传链接格式
BaiduPCS-Go export /我的资源 --link
```

## 创建目录
```
BaiduPCS-Go mkdir <目录>
```

#### 例子
```
BaiduPCS-Go mkdir 123
```

## 删除文件/目录
```
BaiduPCS-Go rm <网盘文件或目录的路径1> <文件或目录2> <文件或目录3> ...
```

注意: 删除多个文件和目录时, 请确保每一个文件和目录都存在, 否则删除操作会失败.

被删除的文件或目录可在网盘文件回收站找回.

#### 例子
```
# 删除 /我的资源/1.mp4
BaiduPCS-Go rm /我的资源/1.mp4

# 删除 /我的资源/1.mp4 和 /我的资源/2.mp4
BaiduPCS-Go rm /我的资源/1.mp4 /我的资源/2.mp4

# 删除 /我的资源 内的所有文件和目录, 但不删除该目录
BaiduPCS-Go rm /我的资源/*

# 删除 /我的资源 整个目录 !!
BaiduPCS-Go rm /我的资源
```

## 拷贝文件/目录
```
BaiduPCS-Go cp <文件/目录> <目标 文件/目录>
BaiduPCS-Go cp <文件/目录1> <文件/目录2> <文件/目录3> ... <目标目录>
```

注意: 拷贝多个文件和目录时, 请确保每一个文件和目录都存在, 否则拷贝操作会失败.

#### 例子
```
# 将 /我的资源/1.mp4 复制到 根目录 /
BaiduPCS-Go cp /我的资源/1.mp4 /

# 将 /我的资源/1.mp4 和 /我的资源/2.mp4 复制到 根目录 /
BaiduPCS-Go cp /我的资源/1.mp4 /我的资源/2.mp4 /
```

## 移动/重命名文件/目录
```
# 移动:
BaiduPCS-Go mv <文件/目录1> <文件/目录2> <文件/目录3> ... <目标目录>
# 重命名:
BaiduPCS-Go mv <文件/目录> <重命名的文件/目录>
```

注意: 移动多个文件和目录时, 请确保每一个文件和目录都存在, 否则移动操作会失败.

#### 例子
```
# 将 /我的资源/1.mp4 移动到 根目录 /
BaiduPCS-Go mv /我的资源/1.mp4 /

# 将 /我的资源/1.mp4 重命名为 /我的资源/3.mp4
BaiduPCS-Go mv /我的资源/1.mp4 /我的资源/3.mp4
```

## 转存文件/目录
```
# 转存分享链接里的文件到当前目录:
BaiduPCS-Go transfer <分享链接> <提取码>
# 转存通用秒传链接里的文件到当前目录:
BaiduPCS-Go transfer <秒传链接>
```

注意: 公开分享链接不需输入提取码, 支持多个文件/目录; 只支持包含单个文件的秒传链接.

转存文件保存到当前工作目录下, 不支持指定.

#### 例子
```
# 将 https://pan.baidu.com/s/12L_ZZVNxz5f_2CccoyyVrW (提取码edv4) 转存到当前目录
BaiduPCS-Go transfer https://pan.baidu.com/s/12L_ZZVNxz5f_2CccoyyVrW edv4

# 将 E7E7B8613854379642F70230B179F37A#FA690D0AB7C8BC6A62WD1B6B3FC5248F#128859362#test.7z 转存到当前目录
BaiduPCS-Go transfer E7E7B8613854379642F70230B179F37A#FA690D0AB7C8BC6A62WD1B6B3FC5248F#128859362#test.7z
```

## 分享文件/目录
```
BaiduPCS-Go share
```

### 设置分享文件/目录
```
BaiduPCS-Go share set <文件/目录1> <文件/目录2> ...
BaiduPCS-Go share s <文件/目录1> <文件/目录2> ...
```

### 列出已分享文件/目录
```
BaiduPCS-Go share list
BaiduPCS-Go share l
```

### 取消分享文件/目录
```
BaiduPCS-Go share cancel <shareid_1> <shareid_2> ...
BaiduPCS-Go share c <shareid_1> <shareid_2> ...
```

目前只支持通过分享id (shareid) 来取消分享.

## 离线下载
```
BaiduPCS-Go offlinedl
BaiduPCS-Go clouddl
BaiduPCS-Go od
```

离线下载支持http/https/ftp/电驴/磁力链协议

离线下载同时进行的任务数量有限, 超出限制的部分将无法添加.

### 添加离线下载任务
```
BaiduPCS-Go offlinedl add -path=<离线下载文件保存的路径> 资源地址1 地址2 ...
```

添加任务成功之后, 返回离线下载的任务ID.

### 精确查询离线下载任务
```
BaiduPCS-Go offlinedl query 任务ID1 任务ID2 ...
```

### 查询离线下载任务列表
```
BaiduPCS-Go offlinedl list
```

### 取消离线下载任务
```
BaiduPCS-Go offlinedl cancel 任务ID1 任务ID2 ...
```

### 删除离线下载任务
```
BaiduPCS-Go offlinedl delete 任务ID1 任务ID2 ...

# 清空离线下载任务记录, 程序不会进行二次确认, 谨慎操作!!!
BaiduPCS-Go offlinedl delete -all
```

#### 例子
```
# 将百度和腾讯主页, 离线下载到根目录 /
BaiduPCS-Go offlinedl add -path=/ http://baidu.com http://qq.com

# 添加磁力链接任务
BaiduPCS-Go offlinedl add magnet:?xt=urn:btih:xxx

# 查询任务ID为 12345 的离线下载任务状态
BaiduPCS-Go offlinedl query 12345

# 取消任务ID为 12345 的离线下载任务
BaiduPCS-Go offlinedl cancel 12345
```

## 回收站
```
BaiduPCS-Go recycle
```

回收站操作.

### 列出回收站文件列表
```
BaiduPCS-Go recycle list
```

#### 可选参数
```
  --page value  回收站文件列表页数 (default: 1)
```

### 还原回收站文件或目录
```
BaiduPCS-Go recycle restore <fs_id 1> <fs_id 2> <fs_id 3> ...
```

根据文件/目录的 fs_id, 还原回收站指定的文件或目录.

### 删除回收站文件或目录/清空回收站
```
BaiduPCS-Go recycle delete [-all] <fs_id 1> <fs_id 2> <fs_id 3> ...
```

根据文件/目录的 fs_id 或 -all 参数, 删除回收站指定的文件或目录或清空回收站.

#### 例子
```
# 从回收站还原两个文件, 其中的两个文件的 fs_id 分别为 1013792297798440 和 643596340463870
BaiduPCS-Go recycle restore 1013792297798440 643596340463870

# 从回收站删除两个文件, 其中的两个文件的 fs_id 分别为 1013792297798440 和 643596340463870
BaiduPCS-Go recycle delete 1013792297798440 643596340463870

# 清空回收站, 程序不会进行二次确认, 谨慎操作!!!
BaiduPCS-Go recycle delete -all
```

## 显示程序环境变量
```
BaiduPCS-Go env
```

BAIDUPCS_GO_CONFIG_DIR: 配置文件路径,

BAIDUPCS_GO_VERBOSE: 是否启用调试.

## 显示和修改程序配置项
```
# 显示配置
BaiduPCS-Go config

# 设置配置
BaiduPCS-Go config set
```

注意: v3.5 以后, 程序对配置文件储存路径的寻找做了调整, 配置文件所在的目录可以是程序本身所在目录, 也可以是家目录.

配置文件所在的目录为家目录的情况:

Windows: `%APPDATA%\BaiduPCS-Go`

其他操作系统: `$HOME/.config/BaiduPCS-Go`

可通过设置环境变量 `BAIDUPCS_GO_CONFIG_DIR`, 指定配置文件存放的目录.

谨慎修改 `appid`, `user_agent`, `pcs_ua`, `pan_ua` 的值, 否则访问网盘服务器时, 可能会出现错误.

上传速度慢的海外用户可尝试修改 `pcs_addr` 值, 选择速度较快的服务器, 目前已知的地址有:

```
pcs.baidu.com
c.pcs.baidu.com
c2.pcs.baidu.com
c3.pcs.baidu.com
c4.pcs.baidu.com
c5.pcs.baidu.com
d.pcs.baidu.com
```

`cache_size` 的值支持可选设置单位了, 单位不区分大小写, `b` 和 `B` 均表示字节的意思, 如 `64KB`, `1MB`, `32kb`, `65536b`, `65536`.

`max_download_rate`, `max_upload_rate` 的值支持可选设置单位了, 单位为每秒的传输速率, 后缀`/s` 可省略, 如 `2MB/s`, `2MB`, `2m`, `2mb` 均为一个意思.

普通用户请将`max_parallel`和`max_download_load`都设置为1, 调大线程数只会在短时间内提升下载速度, 且极易很快触发限速, 导致几小时至几天内账号在各客户端都接近0速. 本软件不支持普通用户提速.

SVIP用户建议`max_parallel`设置为10以上, 根据实际带宽可调大, 但不建议超过20, `max_download_load`设置为1 - 2, 实验表明可以稳定满速下载.

#### 例子
```
# 显示所有可以设置的值
BaiduPCS-Go config -h
BaiduPCS-Go config set -h

# 设置下载文件的储存目录
BaiduPCS-Go config set -savedir D:/Downloads

# 设置下载最大并发量为 150
BaiduPCS-Go config set -max_parallel 150

# 组合设置
BaiduPCS-Go config set -max_parallel 150 -savedir D:/Downloads
```

## 测试通配符
```
BaiduPCS-Go match <通配符表达式>
```

测试通配符匹配路径, 操作成功则输出所有匹配到的路径.

#### 例子
```
# 匹配 /我的资源 目录下所有mp4格式的文件
BaiduPCS-Go match /我的资源/*.mp4
```

## 工具箱
```
BaiduPCS-Go tool
```

目前工具箱支持加解密文件等.

# 初级使用教程

新手建议: **双击运行程序**, 进入仿 Linux shell 的 cli 交互模式;

cli交互模式下, 光标所在行的前缀应为 `BaiduPCS-Go >`, 如果登录了百度帐号则格式为 `BaiduPCS-Go:<工作目录> <百度ID>$ `

以下例子的命令, 均为 cli交互模式下的命令

运行命令的正确操作: **输入命令, 按一下回车键 (键盘上的 Enter 键)**, 程序会接收到命令并输出结果

## 1. 查看程序使用说明

cli交互模式下, 运行命令 `help`

## 2. 登录百度帐号 (必做)

cli交互模式下, 运行命令 `login -h` (注意空格) 查看帮助

cli交互模式下, 运行命令 `login` 程序将会提示你输入百度用户名(手机号/邮箱/用户名)和密码, 必要时还可以在线验证绑定的手机号或邮箱

## 3. 切换网盘工作目录

cli交互模式下, 运行命令 `cd /我的资源` 将工作目录切换为 `/我的资源` (前提: 该目录存在于网盘)

目录支持通配符匹配, 所以你也可以这样: 运行命令 `cd /我的*` 或 `cd /我的??` 将工作目录切换为 `/我的资源`, 简化输入.

将工作目录切换为 `/我的资源` 成功后, 运行命令 `cd ..` 切换上级目录, 即将工作目录切换为 `/`

为什么要这样设计呢, 举个例子,

假设 你要下载 `/我的资源` 内名为 `1.mp4` 和 `2.mp4` 两个文件, 而未切换工作目录, 你需要依次运行以下命令:

```
d /我的资源/1.mp4
d /我的资源/2.mp4
```

而切换网盘工作目录之后, 依次运行以下命令:

```
cd /我的资源
d 1.mp4
d 2.mp4
```

这样就达到了简化输入的目的

## 4. 网盘内列出文件和目录

cli交互模式下, 运行命令 `ls -h` (注意空格) 查看帮助

cli交互模式下, 运行命令 `ls` 来列出当前所在目录的文件和目录

cli交互模式下, 运行命令 `ls /我的资源` 来列出 `/我的资源` 内的文件和目录

cli交互模式下, 运行命令 `ls ..` 来列出当前所在目录的上级目录的文件和目录

## 5. 下载文件

说明: 下载的文件默认保存到 download/ 目录 (文件夹)

cli交互模式下, 运行命令 `d -h` (注意空格) 查看帮助

cli交互模式下, 运行命令 `d /我的资源/1.mp4` 来下载位于 `/我的资源/1.mp4` 的文件 `1.mp4` , 该操作等效于运行以下命令:

```
cd /我的资源
d 1.mp4
```

现在已经支持目录 (文件夹) 下载, 所以, 运行以下命令, 会下载 `/我的资源` 内的所有文件 (违规文件除外):

```
d /我的资源
```

## 6. 设置下载最大并发量

cli交互模式下, 运行命令 `config set -h` (注意空格) 查看设置帮助以及可供设置的值

cli交互模式下, 运行命令 `config set -max_parallel 2` 将下载最大并发量设置为 2

注意：普通用户下载最大并发量的值超过1将导致账号被限速; SVIP同样不宜设置过高, 建议10~20

## 7. 恢复默认配置

cli交互模式下, 运行命令 `config reset`

## 8. 退出程序

运行命令 `quit` 或 `exit` 或 组合键 `Ctrl+C` 或 组合键 `Ctrl+D`

# 已知问题

* 分片上传文件时, 当文件分片数大于1, 网盘端最终计算所得的md5值和本地的不一致, 这可能是百度网盘的bug, 测试把上传的文件下载到本地后，对比md5值是匹配的. 可通过秒传的原理来修复md5值.
* 开启MD5校验下载时可能有 check MD5 不通过, 但文件其实并未出错的情况, 使用--no-check下载或配置中启用no_check即可(3.7版本默认已启用).
* 用户名登录时图片验证码至少要输入两次, 第一次的输入无效
* 登录出现手机/邮箱验证时要输入至少4次图片验证码


# TODO

* 转存文件数量绕过单次限制

# 交流反馈

提交Issue: [Issues](https://github.com/qjfoidnh/BaiduPCS-Go/issues)
