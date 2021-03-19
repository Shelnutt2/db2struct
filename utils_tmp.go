package db2struct

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ReadFileMethod1 使用ioutil.ReadFile 直接从文件读取到 []byte中
func ReadFile(fileName string) string {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("读取文件失败: %#v \n", err)
		return ""
	}
	return string(f)
}

func CopyFile(sourceFile string, destinationFile string) error {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return err
	}
	return nil
}

type TableParam struct {
	ProjectName    string
	TableName      string
	StructName     string
	TableNote      string
	PkgName        string
	TmpPath        string
	JsonAnnotation bool
	GormAnnotation bool
	GureguTypes    bool
	HandlerUrl     string
}

func CreateHandleFileFromTemp(tp *TableParam) error {
	tmpFile := "./tmp_file/tmp_handle.go.tmp"
	// tmpFile := fmt.Sprintf("%s/tmp_file/tmp_handle.go.tmp", tp.TmpPath)
	srcPath := fmt.Sprintf("/tmp/%s/handle", tp.ProjectName)
	err := os.MkdirAll(srcPath, 0666)
	if err != nil {
		return err
	}
	newFile := fmt.Sprintf("%s/%s_handler.go", srcPath, tp.TableName)
	CopyFile(newFile, newFile+".bak") // 备份文件

	text := ReadFile(tmpFile)
	text = ReplaceAllFun(text, tp)

	err = ioutil.WriteFile(newFile, []byte(text), 0644)
	if err != nil {
		fmt.Println("Error creating", newFile)
		fmt.Println(err)
		return err
	}
	return nil
}

func ReplaceAllFun(str string, tp *TableParam) string {
	text := str
	m := strings.ToLower(string(tp.StructName[0]))
	text = strings.ReplaceAll(text, "pLock", m)
	text = strings.ReplaceAll(text, "handler", tp.PkgName)
	text = strings.ReplaceAll(text, "MisLock", tp.StructName)
	text = strings.ReplaceAll(text, "mis_lock", tp.TableName)
	text = strings.ReplaceAll(text, "mis-lock", strings.ReplaceAll(tp.TableName, "_", "-")) // 下划线转-
	text = strings.ReplaceAll(text, "程序包", tp.TableNote)
	return text
}

// 从tmp_model.go.tmp读取model模板
func GetModelFromTmp(tp *TableParam) string {
	// tmpFile := fmt.Sprintf("%s/tmp_file/tmp_model.go.tmp", tp.TmpPath)
	// tmpFile := "/Users/linhaicheng/go/src/db2struct/tmp_file/tmp_model.go.tmp"
	tmpFile := "./tmp_file/tmp_model.go.tmp"
	srcPath := fmt.Sprintf("/tmp/%s/model", tp.ProjectName)
	err := os.MkdirAll(srcPath, 0666)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	newFile := fmt.Sprintf("%s/%s.go", srcPath, tp.TableName)
	CopyFile(newFile, newFile+".bak") // 备份文件
	text := ReadFile(tmpFile)
	text = ReplaceAllFun(text, tp)
	return text
}
