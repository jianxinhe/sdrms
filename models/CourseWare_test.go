package models

import (
	"fmt"
	"testing"
)

func TestTraverseFolder(t *testing.T) {
	xxx, err := TraverseFolder()
	fmt.Println(err)
	fmt.Println(xxx)
}

func TestRemovePrefix(t *testing.T) {
	path := "D:/IDEWorkSpace/go-workspace/src/github.com/jianxinhe/coursesuite"
	fmt.Println(removePrefix(path))
}
