#!/bin/sh

name="BaiduPCS-Go"
version=$1

if [ "$1" = "" ]; then
  version=v3.7.1
fi

output="out"

old_golang() {
  GOROOT=/usr/local/go
  go=$GOROOT/bin/go
}

new_golang() {
  GOROOT=/usr/local/go
  go=$GOROOT/bin/go
}

Build() {
  old_golang
  goarm=$4
  if [ "$4" = "" ]; then
    goarm=7
  fi

  echo "Building $1..."
  export GOOS=$2 GOARCH=$3 GO386=sse2 CGO_ENABLED=0 GOARM=$4
  if [ $2 = "windows" ]; then
    goversioninfo -o=resource_windows_386.syso
    goversioninfo -64 -o=resource_windows_amd64.syso
    $go build -ldflags "-X main.Version=$version -s -w" -o "$output/$1/$name.exe"
    RicePack $1 $name.exe
  else
    $go build -ldflags "-X main.Version=$version -s -w" -o "$output/$1/$name"
    RicePack $1 $name
  fi

  Pack $1
}

AndroidBuild() {
  new_golang
  echo "Building $1..."
  export GOOS=$2 GOARCH=$3 GOARM=$4 CGO_ENABLED=1
  go build -ldflags "-X main.Version=$version -s -w -linkmode=external -extldflags=-pie" -o "$output/$1/$name"

  RicePack $1 $name
  Pack $1
}

IOSBuild() {
  old_golang
  echo "Building $1..."
  mkdir -p "$output/$1"
  cd "$output/$1"
  export CC=/usr/local/go/misc/ios/clangwrap.sh GOOS=darwin GOARCH=arm GOARM=7 CGO_ENABLED=1
  $go build -ldflags "-X main.Version=$version -s -w" -o "armv7" github.com/qjfoidnh/BaiduPCS-Go
  jtool --sign --inplace --ent ../../entitlements.xml "armv7"
  export GOARCH=arm64
  $go build -ldflags "-X main.Version=$version -s -w" -o "arm64" github.com/qjfoidnh/BaiduPCS-Go
  jtool --sign --inplace --ent ../../entitlements.xml "arm64"
  lipo -create "armv7" "arm64" -output $name # merge
  rm "armv7" "arm64"
  cd ../..
  RicePack $1 $name
  Pack $1
}

# zip 打包
Pack() {
  cp README.md "$output/$1"

  cd $output
  zip -q -r "$1.zip" "$1"

  # 删除
  rm -rf "$1"

  cd ..
}

# rice 打包静态资源
RicePack() {
  return # 已取消web功能
  rice -i github.com/qjfoidnh/BaiduPCS-Go/internal/pcsweb append --exec "$output/$1/$2"
}

touch ./vendor/golang.org/x/sys/windows/windows.s

# Android
#export NDK_INSTALL=$ANDROID_NDK_ROOT/bin
# CC=$NDK_INSTALL/arm-linux-androideabi-4.9/bin/arm-linux-androideabi-gcc AndroidBuild $name-$version"-android-16-armv5" android arm 5
# CC=$NDK_INSTALL/arm-linux-androideabi-4.9/bin/arm-linux-androideabi-gcc AndroidBuild $name-$version"-android-16-armv6" android arm 6
#CC=$NDK_INSTALL/arm-linux-androideabi-4.9/bin/arm-linux-androideabi-gcc AndroidBuild $name-$version"-android-16-armv7" android arm 7
#CC=$NDK_INSTALL/aarch64-linux-android-4.9/bin/aarch64-linux-android-gcc AndroidBuild $name-$version"-android-21-arm64" android arm64 7
#CC=$NDK_INSTALL/i686-linux-android-4.9/bin/i686-linux-android-gcc AndroidBuild $name-$version"-android-16-386" android 386 7
#CC=$NDK_INSTALL/x86_64-linux-android-4.9/bin/x86_64-linux-android-gcc AndroidBuild $name-$version"-android-21-amd64" android amd64 7

# iOS
#IOSBuild $name-$version"-darwin-ios-arm"

# OS X / macOS
Build $name-$version"-darwin-osx-amd64" darwin amd64
Build $name-$version"-darwin-osx-arm64" darwin arm64
# Build $name-$version"-darwin-osx-386" darwin 386

# Windows
Build $name-$version"-windows-x86" windows 386
Build $name-$version"-windows-x64" windows amd64
Build $name-$version"-windows-arm" windows arm

# Linux
Build $name-$version"-linux-386" linux 386
Build $name-$version"-linux-amd64" linux amd64
#Build $name-$version"-linux-armv5" linux arm 5
Build $name-$version"-linux-arm" linux arm
Build $name-$version"-linux-arm64" linux arm64
GOMIPS=softfloat Build $name-$version"-linux-mips" linux mips
Build $name-$version"-linux-mips64" linux mips64
GOMIPS=softfloat Build $name-$version"-linux-mipsle" linux mipsle
Build $name-$version"-linux-mips64le" linux mips64le
# Build $name-$version"-linux-ppc64" linux ppc64
# Build $name-$version"-linux-ppc64le" linux ppc64le
# Build $name-$version"-linux-s390x" linux s390x

# Others
# Build $name-$version"-solaris-amd64" solaris amd64
Build $name-$version"-freebsd-386" freebsd 386
Build $name-$version"-freebsd-amd64" freebsd amd64
# Build $name-$version"-freebsd-arm" freebsd arm
# Build $name-$version"-netbsd-386" netbsd	386
# Build $name-$version"-netbsd-amd64" netbsd amd64
# Build $name-$version"-netbsd-arm" netbsd	arm
# Build $name-$version"-openbsd-386" openbsd 386
# Build $name-$version"-openbsd-amd64" openbsd	amd64
# Build $name-$version"-openbsd-arm" openbsd arm
# Build $name-$version"-plan9-386" plan9 386
# Build $name-$version"-plan9-amd64" plan9 amd64
# Build $name-$version"-plan9-arm" plan9 arm
# Build $name-$version"-nacl-386" nacl 386
# Build $name-$version"-nacl-amd64p32" nacl amd64p32
# Build $name-$version"-nacl-arm" nacl arm
# Build $name-$version"-dragonflybsd-amd64" dragonfly amd64
