package sniffer_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestCheckIfIp(t *testing.T) {
	s := "8.8.8.8"
	str := strings.Split(s, ".")
	fmt.Println(str)
	for _, v := range str {
		if _, err := strconv.Atoi(v); err != nil {
			fmt.Println(str)
			fmt.Println("false")
		}
	}
	fmt.Println("true")
}

func TestSubdomainList(t *testing.T) {
	path := "/Users/teadove/projects/ekiren/Sublist3r/misis.txt"
	//domain := "misis.ru"
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	fInfo, _ := f.Stat()
	size := fInfo.Size()
	wantByte := make([]byte, size)
	f.Read(wantByte)
	want := string(wantByte)
	//gotList := subdomainList.ListSubdomains(domain)
	//got := strings.Join(gotList, "\n")
	got := "bbb01.misis.ru\nchem.misis.ru\ncommlab.misis.ru\nwww.commlab.misis.ru"
	if !assert.Equal(t, want, got) {
		t.Errorf("results are not equal")
	}
}
