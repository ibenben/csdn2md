package util

import (
	"os"
	"time"
	"fmt"
	"runtime/debug"
	"path"
)

//
//捕捉协程的错就错误并记录
//
func CatchErr() {
	errs := recover()
	if errs == nil {
		return
	}

	now := time.Now()  //获取当前时间
	pid := os.Getpid() //获取进程ID

	time_str := now.Format("20060102150405")                          //设定时间格式
	fname := fmt.Sprintf("./log/%d-%s-dump.log", pid, time_str) //保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）
	CreateFile(path.Dir(fname))
	f, err := os.Create(fname)
	fmt.Println(err, fname)
	if err != nil {
		return
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("%v\r\n", errs)) //输出panic信息
	f.WriteString("========\r\n")

	f.WriteString(string(debug.Stack())) //输出堆栈信息
}
