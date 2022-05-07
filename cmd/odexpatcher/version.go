package main

import (
	"fmt"
	"runtime"
)

var (
	AppName    string
	Version    string
	CommitHash string
	BuiltTime  string
	OsArch     string
)

func init() {
	AppName = "ODEX Patcher"
	OsArch = runtime.GOOS + "/" + runtime.GOARCH
}

func printVersion() {
	fmt.Println(AppName)
	fmt.Printf(" Version:\t%s\n", Version)
	fmt.Printf(" Commit:\t%s\n", CommitHash)
	fmt.Printf(" Built Time:\t%s\n", BuiltTime)
	fmt.Printf(" OS/Arch:\t%s\n", OsArch)
}
