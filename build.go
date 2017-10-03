package main

import (
	"os/exec"
	"fmt"
	"strings"
	"sync"
	"os"
)

type OSNAME int
type ARCHNAME int

const (
	linux OSNAME = 1 << iota
	windows
	darwin
)

const (
	amd64 ARCHNAME = 1 << iota
	i386
)

func (osname OSNAME) String() string {
	if osname & linux == linux {
		return "linux"
	} else if osname & windows == windows {
		return "windows"
	} else if osname & darwin == darwin {
		return "darwin"
	} else {
		panic("OS not supported")
	}
}

func (archname ARCHNAME) String() string {
	if archname & amd64 == amd64 {
		return "amd64"
	} else if archname & i386 == i386 {
		return "386"
	} else {
		panic ("ARCH not supported")
	}
}

func build(osName OSNAME, archName ARCHNAME, target string, goFile string, wg *sync.WaitGroup) {
	goCmd := "go"
	cmdArgs := []string{"build", "-o", target, goFile}
	cmd := exec.Command(goCmd, cmdArgs...)

	osEnv := "GOOS=" + osName.String()
	archEnv :=  "GOARCH=" + archName.String()
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, osEnv, archEnv)
	defer wg.Done()
	err := cmd.Run()
	if err != nil {
		panic("Failed to build executable: " + err.Error())
	} else {
		fmt.Println("Created " + target)
	}
}

func main() {
	osNames := []OSNAME{linux, darwin, windows}
	goFile := os.Args[2]
	var target []string

	var wg sync.WaitGroup
	for _, o := range osNames {
		arch := amd64
		suffix := o.String()
		if o == windows {
			arch = i386
			suffix += ".exe"
		} else if o == darwin {
			suffix = "mac"
		}
		targetName := strings.Split(goFile, ".go")[0] + "_" + suffix
		wg.Add(1)
		go build(o, arch, targetName, goFile, &wg)
		target = append(target, targetName)
	}
	wg.Wait()

	//check file exists
	//for _, f := range target {
	//	if _, err := os.Stat(f); os.IsExist(err) {
	//		fmt.Println(err)
	//	} else {
	//		fmt.Println("Executable file " + f + " created")
	//	}
	//}
}
