package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"os"
	"os/exec"
	"bytes"
	"archive/zip"
	"errors"
	"io"
)

func extractPomFromJar(archive *os.File) (*os.File, error) {
	reader, err := zip.OpenReader(archive.Name())
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "pom.xml") {
			fmt.Println("Found " + file.Name)
			fileReader, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer fileReader.Close()

			pomName := strings.Split(archive.Name(), ".jar")[0] + ".pom"
			targetFile, err := os.OpenFile(pomName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return nil, err
			}
			defer targetFile.Close()

			if _, err := io.Copy(targetFile, fileReader); err != nil {
				return nil, err
			}
			return targetFile, nil
		}
	}

	return nil, errors.New("pom.xml not found in " + archive.Name())
}

func mvnInstallFile(cmdArgs []string) (bytes.Buffer, error) {
	cmd := exec.Command("mvn", cmdArgs...)
	fmt.Println(cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out, err
}

func main() {
	//in the same dir, list all the jar files
	//for each jar file xxx:
	//	if exists xxx.pom:
	//		mvn org.apache.maven.plugins:maven-install-plugin:2.5.2:install-file -Dfile=xxx.jar -DpomFile=xxx.pom
	//	else:
	//		mvn org.apache.maven.plugins:maven-install-plugin:2.5.2:install-file -Dfile=xxx.jar

	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		cmdArgs := []string{"org.apache.maven.plugins:maven-install-plugin:2.5.2:install-file"}
		if strings.HasSuffix(f.Name(), ".jar") {
			fmt.Println(f.Name())
			cmdArgs = append(cmdArgs, "-Dfile=" + f.Name())

			jarFile, err := os.OpenFile(f.Name(), os.O_RDONLY, 0666)
			if err != nil {
				panic(err)
			}

			if pomXmlFile, err := extractPomFromJar(jarFile); err == nil {
				cmdArgs = append(cmdArgs, "-DpomFile="+pomXmlFile.Name())
			} else {
				fmt.Println(f.Name() + " has no corresponding POM file")
			}

			if err := jarFile.Close(); err != nil {
				panic(err)
			}

			out, err := mvnInstallFile(cmdArgs)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(out.String())
		}
	}
}
