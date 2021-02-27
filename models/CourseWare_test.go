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
