package main

import (
	"testing"
	"os"
)

func ExtractPom_test(t *testing.T) {
	f, err := os.OpenFile("some.jar", os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	pomXmlFile, err := ExtractPomFromJar(f)
	if err != nil {
		t.Fail()
	}

	if pomXmlFile.Name != "pom.xml" {
		t.Fail()
	}
}
