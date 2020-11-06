// AndroidNDKBuild
// go build -ldflags "-X main.APILevel=15 -X main.Arch=x86_64"
// env ANDROID_API_LEVEL NDK ANDROID_NDK_ROOT GOARCH

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

var (
	// NDKPath path to Android NDK
	NDKPath string
	// APILevel Android api level
	APILevel string
	// Arch arch
	Arch string
)

func getNDKPath() string {
	ndkPath, ok := os.LookupEnv("NDK")
	if ok {
		return ndkPath
	}
	ndkPath, ok = os.LookupEnv("ANDROID_NDK_ROOT")
	if ok {
		return ndkPath
	}
	ndkPath, ok = os.LookupEnv("ANDROID_NDK_DIR")
	if ok {
		return ndkPath
	}
	return ""
}

func getAPILevel() string {
	apiLevelStr, ok := os.LookupEnv("ANDROID_API_LEVEL")
	if ok {
		return apiLevelStr
	}
	return "21"
}

func getGoarch() string {
	arch, ok := os.LookupEnv("GOARCH")
	if ok {
		return arch
	}

	return runtime.GOARCH
}

func getArch() string {
	if Arch != "" {
		return Arch
	}
	goarch := getGoarch()
	switch goarch {
	case "386":
		return "x86"
	case "amd64":
		return "x86_64"
	case "arm64":
		return "aarch64"
	}
	return goarch
}

func getPlatformsArch() string {
	arch := getArch()
	switch arch {
	case "aarch64":
		return "arm64"
	}
	return arch
}

func main() {
	if NDKPath == "" {
		NDKPath = getNDKPath()
	}
	if APILevel == "" {
		APILevel = getAPILevel()
	}
	if Arch == "" {
		Arch = getArch()
	}

	lastPattern := "*-gcc"
	if runtime.GOOS == "windows" {
		lastPattern += ".exe"
	}

	gccPaths, err := filepath.Glob(filepath.Join(NDKPath, "toolchains", getArch()+"-*", "prebuilt", runtime.GOOS+"-*", "bin", lastPattern))
	checkErr(err)
	if len(gccPaths) == 0 {
		panic("no match gcc")
	}

	args := make([]string, len(os.Args))
	copy(args[1:], os.Args[1:])
	args[0] = "--sysroot=" + filepath.Join(NDKPath, "platforms", "android-"+APILevel, "arch-"+getPlatformsArch())

	gccExec := exec.Command(gccPaths[0], args...)
	gccExec.Stdout = os.Stdout
	gccExec.Stderr = os.Stderr

	err = gccExec.Run()
	exitError, ok := err.(*exec.ExitError)
	if ok {
		status := exitError.ProcessState.Sys().(syscall.WaitStatus)
		os.Exit(status.ExitStatus())
	}

	if err != nil {
		println(err.Error())
	}

	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
