package utils

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func TrimHTML(str string) string {
	if str == "" {
		return str
	}
	re, _ := regexp.Compile(`<[\s\S]+?(>|$)`)
	newstr := re.ReplaceAllString(str, "")
	return newstr
}

func SubStr(str string, start, end int) string {
	if start < 0 {
		log.Panic("start position is wrong!")
	}
	if end > len(str) {
		log.Panic("end positon is wrong!")
	}
	if start > end {
		log.Panic("wrong position!")
	}

	rs := []rune(str)
	return string(rs[start:end])
}

/**
检测文件是否存在 Stat返回fileInfo
*/
func IsExists(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

//复制文件
func CopyFile(src, dst string) (w int64, err error) {
	f, err := os.Open(src)
	if err != nil {
		return
	}
	defer f.Close()
	dstf, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	defer dstf.Close()
	return io.Copy(dstf, f)
}

//递归复制目录以及其文件
func CopyDir(source, dest string) (err error) {
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return &CustomError{"Source is not a directory"}
	}

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}
	entries, err := ioutil.ReadDir(source)
	for _, entry := range entries {
		sfp := source + "/" + entry.Name()
		dfp := dest + "/" + entry.Name()
		if entry.IsDir() {
			err = CopyDir(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		} else {
			_, err = CopyFile(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		}

	}
	return
}

func XmlEscapString(str string) string {
	str = strings.Replace(str, `<pre class="prettyprint linenums">`, "@@PRE_BEGIN", -1)
	str = strings.Replace(str, `</pre>`, "@@PRE_END", -1)
	str = template.HTMLEscapeString(str)
	str = strings.Replace(str, "@@PRE_BEGIN", `<pre class="prettyprint linenums">`, -1)
	str = strings.Replace(str, "@@PRE_END", `</pre>`, -1)
	return str
}

//转义
func Unescaped(str string) interface{} {
	re := regexp.MustCompile(`<pre class="prettyprint linenums">([\s\S]*?)</pre>`)
	str = re.ReplaceAllStringFunc(str, XmlEscapString)
	return template.HTML(str)

}

type CustomError struct {
	msg string
}

func (e *CustomError) Error() string {
	return e.msg
}

func GetDirName(uri string) string {

	es := strings.Split(uri, "/")
	path := es[:len(es)-1]
	return strings.Join(path, "/")
}
