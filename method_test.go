package natsrpc

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	ret, err := parseStruct(&A{})
	if nil != err {
		t.Error(err)
	}
	for _, v := range ret {
		fmt.Println(v)
	}

}
