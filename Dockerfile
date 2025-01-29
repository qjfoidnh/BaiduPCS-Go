FROM golang:1.20-alpine

WORKDIR /usr/src/app

VOLUME [ "/root/.config" "/root/Downloads"]
ENV username="admin" password="adminadmin"

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY BaiduPCS-Go /usr/local/bin/app

LABEL author="wuzhican"

EXPOSE 8080
ENTRYPOINT app serve -auth -username ${username} -password ${password}