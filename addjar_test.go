package main

import (
	"testing"
	"os"
	//"fmt"
)

func TestExtractPomFromJar(t *testing.T) {
	jarName := "some.jar"
	f, err := os.OpenFile(jarName, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}

	pomXmlFile, err := extractPomFromJar(f)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if pomXmlFile.Name() != "some.pom" {
		t.Error("pom.xml file not found")
		t.Fail()
	}

	// uncomment if you want to test installing the jar with mvn
	//cmdArgs := []string{"org.apache.maven.plugins:maven-install-plugin:2.5.2:install-file"}
	//cmdArgs = append(cmdArgs, "-Dfile=" + f.Name())
	//cmdArgs = append(cmdArgs, "-DpomFile="+pomXmlFile.Name())
	//out, err := mvnInstallFile(cmdArgs)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(out.String())
}
