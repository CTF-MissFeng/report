package api

import (
	"fmt"

	"assets/app/model"
	"assets/app/service"
	"assets/library/response"

	"github.com/gogf/gf/net/ghttp"
)

var Users = new(apiUser)

type apiUser struct{}

// Login 用户登录接口
func (a *apiUser) Login(r *ghttp.Request) {
	var data *model.RequestUsersLogin
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.User.Login(r.Context(), data, r.GetRemoteIp(), r.GetHeader("User-Agent")); err != nil {
		response.JsonExit(r, 202, err.Error())
	}else {
		response.JsonExit(r, 200, "ok")
	}
}

// Register 添加用户接口
func (a *apiUser) Register(r *ghttp.Request) {
	service.User.IsUserAdmin(r,r.Context())
	var data *model.RequestUsersRegister
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.User.Register(data); err != nil{
		response.JsonExit(r, 202, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "添加用户", fmt.Sprintf("添加[%s]用户", data.Username))
		response.JsonExit(r, 200, "ok")
	}
}

// UserInfo 获取用户信息接口
func (a *apiUser) UserInfo(r *ghttp.Request) {
	response.JsonExit(r, 200, "ok", service.User.UserInfo(r.Context()))
}

// UserDel 删除用户接口
func (a *apiUser) UserDel(r *ghttp.Request){
	service.User.IsUserAdmin(r,r.Context())
	var data *model.RequestUserDel
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.User.UserDel(data); err != nil{
		response.JsonExit(r, 202, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "删除用户", fmt.Sprintf("删除[%s]用户", data.Username))
		response.JsonExit(r, 200, "ok")
	}
}

// ChangePassword 用户修改密码接口
func (a *apiUser) ChangePassword(r *ghttp.Request){
	var data *model.RequestUserChangePassword
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err:= service.User.ChangePassword(r.Context(), data); err != nil{
		response.JsonExit(r, 202, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "密码修改", "修改成功")
		response.JsonExit(r, 200, "ok")
	}
}

// SetUserInfo 用户修改资料接口
func (a *apiUser) SetUserInfo(r *ghttp.Request){
	var data *model.RequestUserSetInfo
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err:= service.User.SetUserInfo(r.Context(), data); err != nil{
		response.JsonExit(r, 202, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "资料修改", fmt.Sprintf("成功修改 [%s]", data.NickName))
		response.JsonExit(r, 200, "ok")
	}
}

// LoginOut 用户注销接口
func (a *apiUser) LoginOut(r *ghttp.Request) {
	service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "注销用户", "注销成功")
	if err := service.User.LoginOut(r.Context()); err != nil {
		response.JsonExit(r, 201, err.Error())
	}
	response.JsonExit(r, 200, "ok")
}

// SearchUser 用户管理模糊查询分页接口
func (a *apiUser) SearchUser(r *ghttp.Request){
	r.Response.WriteJson(service.User.SearchUser(r.GetInt("page"), r.GetInt("limit"), r.Get("searchParams")))
}

// SearchUserLockIp ip锁定详情分页接口
func (a *apiUser) SearchUserLockIp(r *ghttp.Request){
	r.Response.WriteJson(service.User.SearchUserLockIp(r.GetInt("page"), r.GetInt("limit")))
}

// UserLockIpRest 解锁ip接口
func (a *apiUser) UserLockIpRest(r *ghttp.Request){
	service.User.IsUserAdmin(r,r.Context())
	var data *model.RequestUserLockIp
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.User.UserLockIpRest(data); err != nil{
		response.JsonExit(r, 202, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "解锁IP", fmt.Sprintf("解锁[%s]", data.Ip))
		response.JsonExit(r, 200, "ok")
	}
}

// SearchUserLoginLogs 登录日志管理模糊查询分页接口
func (a *apiUser) SearchUserLoginLogs(r *ghttp.Request){
	r.Response.WriteJson(service.User.SearchUserLoginLogs(r.GetInt("page"), r.GetInt("limit"), r.Get("searchParams")))
}

// SearchUserOperation 用户操作日志模糊查询分页接口
func (a *apiUser) SearchUserOperation(r *ghttp.Request){
	r.Response.WriteJson(service.User.SearchUserOperation(r.GetInt("page"), r.GetInt("limit"), r.Get("searchParams")))
}

// OperationEmpty 清空用户操作日志
func (a *apiUser) OperationEmpty(r *ghttp.Request){
	service.User.IsUserAdmin(r,r.Context())
	err := service.User.OperationEmpty()
	if err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "操作日志", "清空用户操作日志")
	response.JsonExit(r, 200, "ok")
}

// Menu 菜单接口
func (a *apiUser) Menu(r *ghttp.Request) {
	r.Response.WriteJson(model.ModuleInit())
}