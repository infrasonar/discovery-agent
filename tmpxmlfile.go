package main

import (
	"os"
)

func createTmpXmlFile() error {
	f, err := os.CreateTemp("", "infrasonar-scan-*.xml")
	if err != nil {
		return err
	}
	defer f.Close()

	os.Setenv("TMP_XML_FILE", f.Name())
	return nil
}
