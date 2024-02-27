package main

import (
	"os"
	"testing"
)

func readTestData(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile(("../../testdata/" + name))
	if err!= nil {
		t.Errorf("Could not read %v",name)
	}
	return content
}