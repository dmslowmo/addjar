package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"os"
	"os/exec"
	"bytes"
)

func main() {
	//in the same dir, list all the jar files
	//for each jar file xxx:
	//	if exists xxx.pom:
	//		mvn org.apache.maven.plugins:maven-install-plugin:2.5.2:install-file -Dfile=xxx.jar -DpomFile=xxx.pom
	//	else:
	//		mvn org.apache.maven.plugins:maven-install-plugin:2.5.2:install-file -Dfile=xxx.jar

	files, _ := ioutil.ReadDir("./")
	cmdArgs := []string{"org.apache.maven.plugins:maven-install-plugin:2.5.2:install-file"}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".jar") {
			fmt.Println(f.Name())
			cmdArgs = append(cmdArgs, "-Dfile=" + f.Name())
			pom := strings.Split(f.Name(), ".jar")[0] + ".pom"

			if _, err := os.Stat(pom); err == nil {
				fmt.Println(pom)
				cmdArgs = append(cmdArgs, "-DpomFile=" + pom)
			} else {
				fmt.Println(f.Name() + " has no corresponding POM file")
			}

			cmd := exec.Command("mvn", cmdArgs...)
			fmt.Println(cmd.Args)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(out.String())
		}
	}
}
