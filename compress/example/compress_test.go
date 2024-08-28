package main

import (
	"github.com/ameise84/pi_common/str_conv"
	"github.com/ameise84/pi_net/compress"
	"log"
	"slices"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func TestLz4(t *testing.T) {
	r := "1dkfjnadsfjnioawjhrfijq3nadnfjahdfnuianfansfaidfaionetrajdnfaenfqwierfnakfnczvnaiweoaiwfansdkfaweitnadvxcvsdgfq"
	out := make([]byte, 1024)
	src := make([]byte, 1024)
	x1, _ := compress.Zip(str_conv.ToBytes(r), out)
	log.Println(x1)
	x2 := slices.Clone(x1)
	x3, _ := compress.Unzip(x2, src)
	log.Println(str_conv.ToString(x3))
	log.Println(r)
}
