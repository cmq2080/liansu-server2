package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	APP_NAME    = "liansu-server"
	APP_NAME_ZH = "liansu-server服务器控制端"
)

var (
	thisExe     string
	thisDir     string
	thisConfDir string

	cfgFilepath string
	cfg         *Config
)

func main() {
	initialize()

	// 开启GUI
	startGUI()
}

func initialize() {
	thisExe, _ = os.Executable()
	thisDir = getRealDir(filepath.Dir(thisExe))
	thisConfDir = getRealDir(thisDir + "/conf")

	// 载入配置信息
	reloadCfg()
}

func isWindows() bool {
	return strings.Contains(runtime.GOOS, "windows")
}

func isLinux() bool {
	return strings.Contains(runtime.GOOS, "linux")
}
