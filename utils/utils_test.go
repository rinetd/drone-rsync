package utils

import "testing"

var str1 string = "/dir/wr"

// var str2 string = "/dir/wr/"
// var str3 string = "/dir/wr/wear.war"

// var testdata []string =["/dir/wr","/dir/wr/","/dir/wr/ww.war"]
func TestGetDirName(t *testing.T) {
	if GetDirName(str1) != "/dir/" {
		t.Fatalf("err")
	}
	// for i, data := range testdata {
	// 	c := GetDirName(data)
	// 	switch i {
	// 	case 1:
	// 		if c != "/dir/"
	// 		t.Fatalf("err")
	// 	case 2:
	// 		if c != "/dir/wr/"
	// 	}

	// }
}
