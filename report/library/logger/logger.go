package logger

import (
	"github.com/gogf/gf/os/glog"
)

var WebLog *glog.Logger

// init 初始化web日志
func init()  {
	logs := glog.New()
	logs.SetPath("logs")
	logs.SetLevelStr("all")
	WebLog = logs
}