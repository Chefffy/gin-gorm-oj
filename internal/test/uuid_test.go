package test

import (
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestGenerateUUid(t *testing.T){
	//生成uuid
	s := uuid.NewV4().String()
	println(s)
}