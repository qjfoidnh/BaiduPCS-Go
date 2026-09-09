#!/bin/bash
# 这个脚本以，菜单形式运行 BaiduPCS-Go 命令
# 菜单结构同 3.9.7 read.me
# 使用前自行修改 BDUSS_FILE 位置
# BDUSS 会定期更换
# bd.BDUSS 文件存放 baidu BDUSS KEY
# 测试在 Ubuntu 24 可以运行


# BDUSS 文件:
BDUSS_FILE="/share/Multimedia/2024-MyProgramFiles/29.QTS_conf_files/bd.BDUSS"

# 清屏
clear_screen() {
    clear
}

# 暂停 
pause() {
    echo
    read -p "按回车键继续..." key
}

# loading BDUSS
read_bduss() {
    if [ ! -f "$BDUSS_FILE" ]; then
        echo "错误: BDUSS文件不存在: $BDUSS_FILE"
        return 1
    fi
    
    BDUSS=$(cat "$BDUSS_FILE")
    if [ -z "$BDUSS" ]; then
        echo "错误: BDUSS文件内容为空"
        return 1
    fi
    
    echo "$BDUSS"
    return 0
}

show_main_menu() {
    clear_screen
    echo "=== BaiduPCS-Go 主菜单 ==="
    echo "1. 账号管理"
    echo "2. 文件操作" 
    echo "3. 分享/转存"
    echo "4. 回收站"
    echo "5. 系统设置"
    echo "0. 退出"
    echo "===================="
}

show_account_menu() {
    clear_screen
    echo "=== 账号管理 ==="
    echo "1. 使用BDUSS登录"
    echo "2. 切换账号"
    echo "3. 退出账号" 
    echo "4. 显示当前账号"
    echo "5. 显示账号列表"
    echo "9. 返回主菜单"
    echo "0. 退出"
    echo "===================="
}

show_file_menu() {
    clear_screen
    echo "=== 文件操作 ==="
    echo "1. 列出文件"
    echo "2. 切换目录"
    echo "3. 下载文件/目录"
    echo "4. 上传文件/目录"
    echo "5. 创建目录"
    echo "6. 删除文件/目录"
    echo "7. 复制文件/目录"
    echo "8. 移动/重命名"
    echo "9. 返回主菜单"
    echo "0. 退出"
    echo "===================="
}

show_share_menu() {
    clear_screen
    echo "=== 分享/转存 ==="
    echo "1. 分享文件/目录"
    echo "2. 列出已分享"
    echo "3. 取消分享"
    echo "4. 转存分享文件"
    echo "9. 返回主菜单"
    echo "0. 退出"
    echo "===================="
}

show_recycle_menu() {
    clear_screen
    echo "=== 回收站 ==="
    echo "1. 列出回收站"
    echo "2. 还原文件/目录"
    echo "3. 清空回收站"
    echo "9. 返回主菜单"
    echo "0. 退出"
    echo "===================="
}

show_config_menu() {
    clear_screen
    echo "=== 系统设置 ==="
    echo "1. 显示配置"
    echo "2. 修改配置"
    echo "3. 恢复默认配置"
    echo "9. 返回主菜单"
    echo "0. 退出"
    echo "===================="
}

account_operations() {
    while true; do
        show_account_menu
        read -p "请选择操作 [0-9]: " choice
        case $choice in
            1) # 使用BDUSS登录
                clear_screen
                echo "正在从文件读取BDUSS: $BDUSS_FILE"
                BDUSS=$(read_bduss)
                if [ $? -eq 0 ]; then
                    echo "成功读取BDUSS，正在登录..."
                    BaiduPCS-Go login -bduss="$BDUSS"
                fi
                pause
                ;;
            2) # 切换账号 
                clear_screen
                BaiduPCS-Go su
                pause
                ;;
            3) # 退出账号
                clear_screen
                BaiduPCS-Go logout
                pause
                ;;
            4) # 显示当前账号
                clear_screen
                BaiduPCS-Go who
                pause
                ;;
            5) # 显示账号列表
                clear_screen
                BaiduPCS-Go loglist
                pause
                ;;
            9) # 返回主菜单
                return
                ;;
            0) # 退出
                exit 0
                ;;
            *)
                echo "无效的选择"
                pause
                ;;
        esac
    done
}

file_operations() {
    while true; do
        show_file_menu
        read -p "请选择操作 [0-9]: " choice
        case $choice in
            1) # 列出文件
                clear_screen
                read -p "请输入要列出的目录路径(直接回车列出当前目录): " path
                if [ -z "$path" ]; then
                    BaiduPCS-Go ls
                else
                    BaiduPCS-Go ls "$path"
                fi
                pause
                ;;
            2) # 切换目录
                clear_screen
                read -p "请输入要切换到的目录路径: " path
                BaiduPCS-Go cd "$path"
                pause
                ;;
            3) # 下载文件/目录
                clear_screen
                read -p "请输入要下载的文件/目录路径: " path
                BaiduPCS-Go download "$path"
                pause
                ;;
            4) # 上传文件/目录
                clear_screen
                read -p "请输入要上传的本地文件/目录路径: " local_path
                read -p "请输入要上传到的网盘目录路径: " remote_path
                BaiduPCS-Go upload "$local_path" "$remote_path"
                pause
                ;;
            5) # 创建目录
                clear_screen
                read -p "请输入要创建的目录路径: " path
                BaiduPCS-Go mkdir "$path"
                pause
                ;;
            6) # 删除文件/目录
                clear_screen
                read -p "请输入要删除的文件/目录路径: " path
                BaiduPCS-Go rm "$path"
                pause
                ;;
            7) # 复制文件/目录
                clear_screen
                read -p "请输入要复制的源文件/目录路径: " src
                read -p "请输入目标路径: " dst
                BaiduPCS-Go cp "$src" "$dst"
                pause
                ;;
            8) # 移动/重命名
                clear_screen
                read -p "请输入要移动/重命名的源文件/目录路径: " src
                read -p "请输入新路径: " dst
                BaiduPCS-Go mv "$src" "$dst"
                pause
                ;;
            9) # 返回主菜单
                return
                ;;
            0) # 退出
                exit 0
                ;;
            *)
                echo "无效的选择"
                pause
                ;;
        esac
    done
}

share_operations() {
    while true; do
        show_share_menu
        read -p "请选择操作 [0-9]: " choice
        case $choice in
            1) # 分享文件/目录
                clear_screen
                read -p "请输入要分享的文件/目录路径: " path
                BaiduPCS-Go share set "$path"
                pause
                ;;
            2) # 列出已分享
                clear_screen
                BaiduPCS-Go share list
                pause
                ;;
            3) # 取消分享
                clear_screen
                read -p "请输入要取消的分享ID: " share_id
                BaiduPCS-Go share cancel "$share_id"
                pause
                ;;
            4) # 转存分享文件
                clear_screen
                read -p "请输入分享链接: " link
                read -p "请输入提取码: " code
                BaiduPCS-Go transfer "$link" "$code"
                pause
                ;;
            9) # 返回主菜单
                return
                ;;
            0) # 退出
                exit 0
                ;;
            *)
                echo "无效的选择"
                pause
                ;;
        esac
    done
}

recycle_operations() {
    while true; do
        show_recycle_menu
        read -p "请选择操作 [0-9]: " choice
        case $choice in
            1) # 列出回收站
                clear_screen
                BaiduPCS-Go recycle list
                pause
                ;;
            2) # 还原文件/目录
                clear_screen
                read -p "请输入要还原的文件/目录fs_id: " fs_id
                BaiduPCS-Go recycle restore "$fs_id"
                pause
                ;;
            3) # 清空回收站
                clear_screen
                echo "警告:该操作将清空回收站!"
                read -p "确认要清空吗?(y/n) " confirm
                if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
                    BaiduPCS-Go recycle delete -all
                fi
                pause
                ;;
            9) # 返回主菜单
                return
                ;;
            0) # 退出
                exit 0
                ;;
            *)
                echo "无效的选择"
                pause
                ;;
        esac
    done
}

config_operations() {
    while true; do
        show_config_menu
        read -p "请选择操作 [0-9]: " choice
        case $choice in
            1) # 显示配置
                clear_screen
                BaiduPCS-Go config
                pause
                ;;
            2) # 修改配置
                clear_screen
                echo "常用配置选项:"
                echo "1) 设置下载目录: config set -savedir <目录路径>"
                echo "2) 设置下载并发数: config set -max_parallel <数值>"
                echo "3) 设置同时下载文件数: config set -max_download_load <数值>"
                echo
                read -p "请输入完整的配置命令: " cmd
                BaiduPCS-Go $cmd
                pause
                ;;
            3) # 恢复默认配置
                clear_screen
                echo "警告:该操作将恢复所有默认配置!"
                read -p "确认要恢复吗?(y/n) " confirm
                if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
                    BaiduPCS-Go config reset
                fi
                pause
                ;;
            9) # 返回主菜单
                return
                ;;
            0) # 退出
                exit 0
                ;;
            *)
                echo "无效的选择"
                pause
                ;;
        esac
    done
}

while true; do
    show_main_menu
    read -p "请选择操作 [0-5]: " choice
    case $choice in
        1) # 账号管理
            account_operations
            ;;
        2) # 文件操作
            file_operations
            ;;
        3) # 分享/转存
            share_operations
            ;;
        4) # 回收站
            recycle_operations
            ;;
        5) # 系统设置
            config_operations
            ;;
        0) # 退出
            echo "谢谢使用,再见!"
            exit 0
            ;;
        *)
            echo "无效的选择"
            pause
            ;;
    esac
done
