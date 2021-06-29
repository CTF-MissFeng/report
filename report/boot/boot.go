package boot

import (
	"github.com/gogf/gf/os/gfile"
	"time"

	"assets/app/dao"
	"assets/app/model"
	"assets/library/logger"

	_ "assets/library/logger"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"golang.org/x/crypto/bcrypt"
)

// init 初始化Web
func init() {
	server := g.Server()
	if err := server.SetConfigWithMap(g.Map{
		"address": g.Cfg().Get("server.Address"), // web服务器监听地址
		"serverAgent": "assetr", // web服务器server信息
		"serverRoot": "public", // 静态文件服务的目录根路径
		"SessionMaxAge": 300 * time.Minute, // session最大超时时间
		"SessionIdName": "assert", // session会话ID名称
		"SessionCookieOutput": true, // 指定是否将会话ID自动输出到cookie
	}); err != nil{
		logger.WebLog.Fatalf("web服务器配置有误，程序运行失败:%s", err.Error())
	}

	// 静态文件路由设置
	server.SetRewriteMap(g.MapStrStr{
		"/" : "./html/users/login.html",
		"/home": "./html/users/index.html",
		"/user/manager": "./html/users/manager.html",
		"/user/userip": "./html/users/userip.html",
		"/user/loginlog": "./html/users/login_log.html",
		"/user/operation": "./html/users/operation.html",
		"/user/manager/add": "./html/users/add.html",
		"/user/setting": "./html/users/setting.html",
		"/user/password": "./html/users/password.html",

		"/assets/manager": "./html/assets/manager.html",
		"/assets/manager/add": "./html/assets/managerAdd.html",
		"/assets/type": "./html/assets/type.html",
		"/assets/type/add": "./html/assets/typeAdd.html",
		"/assets/type/edit": "./html/assets/typeEditing.html",
		"/assets/pc": "./html/assets/pc.html",
		"/assets/pc/add": "./html/assets/pcAdd.html",
		"/assets/pc/edit": "./html/assets/pcEdit.html",
		"/assets/pc/show": "./html/assets/pcShow.html",

		"/assets/web": "./html/assets/web.html",
		"/assets/web/add": "./html/assets/webAdd.html",
		"/assets/web/report/add": "./html/assets/webReportAdd.html",
		"/assets/web/edit": "./html/assets/webEdit.html",
		"/assets/web/show": "./html/assets/webShow.html",

		"/assets/report": "./html/assets/report.html",
		"/assets/report/edit": "./html/assets/reportEdit.html",

		"/tongji" :"./html/home/home.html",
	})

	//自定义403、404等
	server.BindStatusHandler(404, func(r *ghttp.Request){
		r.Response.RedirectTo("/")
	})
	server.BindStatusHandler(403, func(r *ghttp.Request){
		r.Response.RedirectTo("/")
	})

	// 创建用户
	createAdmin()

	// 创建所需目录
	createDir()
}

// createAdmin 创建默认admin账户
func createAdmin(){
	if i, err := dao.Users.FindCount("username=?", "admin"); err != nil{
		logger.WebLog.Warningf("[创建默认账户] 查询数据库错误:%s", err.Error())
		return
	}else if i != 0{
		return
	}else{
		passwd,err := bcrypt.GenerateFromPassword([]byte("admin888@A"), bcrypt.DefaultCost)
		if err != nil {
			logger.WebLog.Warningf("[创建默认账户] 加密密码错误:%s", err.Error())
			return
		}else{
			users := model.RequestUsersRegister{}
			users.Username = "admin"
			users.Password = string(passwd)
			users.NickName = "管理员"
			users.Email = "admin@qq.com"
			users.Phone = "13888888888"
			users.Remark = "管理员账户"
			if _, err := dao.Users.Insert(users); err != nil {
				logger.WebLog.Warningf("[创建默认账户] 数据库错误:%s", err.Error())
				return
			}else{
				logger.WebLog.Warningf("[创建默认账户成功] 用户名:admin 密码:admin888@A")
			}
		}
	}
}

// createDir 创建程序所需目录
func createDir(){
	if !gfile.Exists("./upload/tmp"){
		gfile.Mkdir("./upload/tmp")
	}
	if !gfile.Exists("./public/upload/assetsImg"){
		gfile.Mkdir("./public/upload/assetsImg")
	}
	if !gfile.Exists("./public/upload/assetsFile"){
		gfile.Mkdir("./public/upload/assetsFile")
	}
	if !gfile.Exists("./public/upload/report"){
		gfile.Mkdir("./public/upload/report")
	}
}
