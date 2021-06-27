package service

import (
	"context"
	"errors"
	"fmt"

	"assets/app/dao"
	"assets/app/model"
	"assets/library/logger"
	"assets/library/response"

	"github.com/gogf/gf/encoding/ghtml"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"golang.org/x/crypto/bcrypt"
)

var User = new(serviceUser)

type serviceUser struct{}

// IsSignedIn 判断用户是否已经登录
func (s *serviceUser) IsSignedIn(ctx context.Context) bool {
	if v := Context.Get(ctx); v != nil && v.User != nil {
		return true
	}
	return false
}

// checkPassword 校检密码
func (s *serviceUser) checkPassword(password, newPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(newPassword))
	return err == nil
}

// isLoginIP 判断登录IP是否锁定
func (s *serviceUser) isLoginIP(ip string) bool {
	userIp, err := dao.UserIp.FindOne("ip=?", ip)
	if err != nil{
		logger.WebLog.Warningf("查询Lock锁数据库错误:%s", err.Error())
		return false
	}
	if userIp == nil{
		return false
	}
	if userIp.LockCount >= 5 {
		return true
	}
	return false
}

// addLoginIP 记录登录失败的ip
func (s *serviceUser) addLoginIP(ip string){
	userIp, err := dao.UserIp.FindOne("ip=?", ip)
	if err != nil{
		logger.WebLog.Warningf("查询登录IP数据库错误:%s", err.Error())
		return
	}
	if userIp == nil{ // 如果没有记录则insert插入一条
		if _, err := dao.UserIp.Insert(g.Map{"ip":ip, "lock_count":1}); err != nil {
			logger.WebLog.Warningf("添加登录锁定IP数据库错误:%s", err.Error())
			return
		}
		return
	}
	if _, err := dao.UserIp.Update(g.Map{"lock_count":userIp.LockCount+1}, "ip", ip); err != nil {
		logger.WebLog.Warningf("修改登录锁定IP数据库错误:%s", err.Error())
		return
	}
}

// restLoginIP 登录成功重置lock锁
func (s *serviceUser) restLoginIP(ip string){
	userIp, err := dao.UserIp.FindOne("ip=?", ip)
	if err != nil{
		logger.WebLog.Warningf("重置登录Lock锁数据库错误:%s", err.Error())
		return
	}
	if userIp == nil{
		return
	}
	if _, err = dao.UserIp.Update(g.Map{"lock_count":0}, "ip", ip); err != nil{
		logger.WebLog.Warningf("重置登录Lock锁数据库错误:%s", err.Error())
		return
	}
}

// addLoginInfo 记录用户登录成功日志
func (s *serviceUser) addLoginInfo(username, ip, userAgent string){
	if _, err := dao.UserLog.Insert(g.Map{"username": username, "ip": ip, "user_agent": userAgent}); err != nil {
		logger.WebLog.Warningf("增加登录日志数据库错误:%s", err.Error())
	}
}

// Login 用户登录
func (s *serviceUser) Login(ctx context.Context, r *model.RequestUsersLogin, ip, userAgent string)error{
	if s.isLoginIP(ip){
		logger.WebLog.Debugf("登录接口 username:%s  ip:%s IP已锁定", r.Username, ip)
		return errors.New("该IP已被锁定,请联系管理员解锁")
	}
	user, err := dao.Users.FindOne("username=?", r.Username)
	if err != nil {
		logger.WebLog.Warningf("登录查询用户 数据库错误:%s", err.Error())
		return errors.New("登录失败,数据库错误")
	}
	if user == nil {
		s.addLoginIP(ip)
		logger.WebLog.Debugf("用户登录 [%s] [%s] 用户不存在", r.Username, ip)
		return errors.New("账号或密码错误")
	}
	if !s.checkPassword(user.Password, r.Password){
		s.addLoginIP(ip)
		logger.WebLog.Debugf("用户登录 [%s] [%s]  密码错误", r.Username, ip)
		return errors.New("账号或密码错误")
	}
	if err := Session.SetUser(ctx, user); err != nil {
		logger.WebLog.Debugf("用户登录 [%s] [%s]  Session错误:%s", r.Username, ip, err.Error())
		return errors.New("登录失败,Session设置错误")
	}
	logger.WebLog.Debugf("用户登录 [%s] [%s]  登录成功", r.Username, ip)
	s.restLoginIP(ip)
	s.addLoginInfo(r.Username, ip, userAgent)
	Context.SetUser(ctx, &model.ContextUser{
		Id:       int(user.Id),
		UserName: user.Username,
		Email:    user.Email,
	})
	return nil
}

// setPassword 加密明文密码
func (s *serviceUser) setPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Register 添加用户
func (s *serviceUser) Register(r *model.RequestUsersRegister) error {
	if i, err := dao.Users.FindCount("username=?", r.Username); err != nil{
		logger.WebLog.Warningf("添加用户 数据库错误:%s", err.Error())
		return errors.New("数据库错误,添加用户失败")
	}else if i > 0{
		return errors.New(fmt.Sprintf("账户 %s 已存在", r.Username))
	}
	encPassword, err := s.setPassword(r.Password)
	if err != nil {
		logger.WebLog.Warningf("添加用户 密码加密失败:%s", err.Error())
		return errors.New("添加用户失败,加密密码错误")
	}
	r.Password = encPassword
	r.NickName = ghtml.SpecialChars(r.NickName)
	r.Remark = ghtml.SpecialChars(r.Remark)
	if _, err := dao.Users.Insert(r); err != nil {
		logger.WebLog.Warningf("添加用户 数据库错误:%s", err.Error())
		return errors.New("添加用户失败,数据库错误")
	}
	logger.WebLog.Warningf("添加用户成功:%s", r.Username)
	return nil
}

// UserInfo 返回当前用户信息
func (s *serviceUser) UserInfo(ctx context.Context)*model.Users{
	user := Session.GetUser(ctx)
	user1,err := dao.Users.FindOne("username=?", user.Username)
	if err != nil {
		logger.WebLog.Warningf("获取用户资料 数据库错误:%s", err.Error())
		return &model.Users{}
	}
	user1.Password = "********"
	return user1
}

// UserDel 删除指定用户
func (s *serviceUser) UserDel(r *model.RequestUserDel) error {
	if r.Username == "admin"{
		return errors.New("删除用户失败:不能删除admin内置账户")
	}
	result,err := dao.Users.Delete("username=?", r.Username)
	if err != nil{
		logger.WebLog.Warningf("删除用户 数据库错误:%s", err.Error())
		return errors.New("数据库错误,删除用户失败")
	}
	if result != nil{
		return nil
	}
	return errors.New("删除用户失败,无此用户")
}

// ChangePassword 用户修改密码
func (s *serviceUser) ChangePassword(ctx context.Context, r *model.RequestUserChangePassword) error {
	userinfo := Session.GetUser(ctx)
	currentUser, err := dao.Users.FindOne("username=?", userinfo.Username)
	if err != nil{
		logger.WebLog.Warningf("修改密码 数据库错误:%s", err.Error())
		return errors.New("数据库错误,修改密码失败")
	}
	if !s.checkPassword(currentUser.Password, r.Password){
		return errors.New("修改密码失败,原密码不正确")
	}
	encPassword,err := s.setPassword(r.Password1)
	if err != nil {
		logger.WebLog.Warningf("修改密码 加密密码错误:%s", err.Error())
		return errors.New("密码加密失败,修改密码失败")
	}
	if result,err := dao.Users.Update(g.Map{"password":encPassword}, "username", userinfo.Username); err != nil{
		logger.WebLog.Warningf("修改密码 数据库错误:%s", err.Error())
		return errors.New("数据库错误,修改密码失败")
	}else if result != nil{
		return nil
	}else{
		return errors.New("修改密码失败,无此用户")
	}
}

// SetUserInfo 用户修改资料
func (s *serviceUser) SetUserInfo(ctx context.Context, r *model.RequestUserSetInfo) error {
	r.Remark = ghtml.SpecialChars(r.Remark)
	r.NickName = ghtml.SpecialChars(r.NickName)
	result,err := dao.Users.Update(r, "username", Session.GetUser(ctx).Username)
	if err != nil{
		logger.WebLog.Warningf("修改用户资料 数据库错误:%s", err.Error())
		return errors.New("数据库错误,修改用户资料失败")
	}
	if result == nil{
		return errors.New("修改用户资料失败,无此用户")
	}else{
		logger.WebLog.Warningf("%s:修改用户资料成功", Session.GetUser(ctx).Username)
		return nil
	}
}

// LoginOut 用户注销
func (s *serviceUser) LoginOut(ctx context.Context) error {
	return Session.RemoveUser(ctx)
}

// SearchUser 用户管理模糊分页查询
func (s *serviceUser) SearchUser(page, limit int, search interface{}) *model.ResponeUserManager{
	var resultUser []*model.Users
	UserSearch := dao.Users.Clone()
	searchStr := gconv.String(search)
	if search != ""{
		j := gjson.New(searchStr)
		if gconv.String(j.Get("username")) != ""{
			UserSearch = UserSearch.Where("username like ?", "%"+gconv.String(j.Get("username"))+"%")
		}
		if gconv.String(j.Get("phone")) != ""{
			UserSearch = UserSearch.Where("phone like ?", "%"+gconv.String(j.Get("phone"))+"%")
		}
		if gconv.String(j.Get("email")) != ""{
			UserSearch = UserSearch.Where("email like ?", "%"+gconv.String(j.Get("email"))+"%")
		}
		if gconv.String(j.Get("nickname"))!= ""{
			UserSearch = UserSearch.Where("nick_name like ?", "%"+gconv.String(j.Get("nickname"))+"%")
		}
	}
	count,_ := UserSearch.Count()
	if page > 0 && limit > 0 {
		err := UserSearch.Order("id desc").Limit((page-1)*limit,limit).Scan(&resultUser)
		if err != nil {
			logger.WebLog.Warningf("用户管理分页查询 数据库错误:%s", err.Error())
			return &model.ResponeUserManager{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponeUserManager{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range resultUser{
		resultUser[i].Password = ""
		index++
		resultUser[i].Id = uint(index)
	}
	return &model.ResponeUserManager{Code:0, Msg:"ok", Count:int64(count), Data:resultUser}
}

// SearchUserLockIp ip锁定管理模糊分页查询
func (s *serviceUser) SearchUserLockIp(page, limit int) *model.ResponeUserIpLockIp{
	var resultUserLockIp []*model.UserIp
	if page > 0 && limit > 0 {
		if err := dao.UserIp.Where("lock_count >= ?", 5).Order("id desc").Limit((page-1)*limit,limit).Scan(&resultUserLockIp); err != nil{
			logger.WebLog.Warningf("IP锁定分页查询 数据库错误:%s", err.Error())
			return &model.ResponeUserIpLockIp{Code:201,Msg:"查询失败,数据库错误",Count:0,Data:nil}
		}
	}else{
		return &model.ResponeUserIpLockIp{Code:201,Msg:"查询失败,分页参数有误",Count:0,Data:nil}
	}
	count,_ := dao.UserIp.Count()
	return &model.ResponeUserIpLockIp{Code:0, Msg:"ok", Count:int64(count), Data:resultUserLockIp}
}

// UserLockIpRest 解锁IP
func (s *serviceUser) UserLockIpRest(r *model.RequestUserLockIp) error {
	if result,err := dao.UserIp.Update(g.Map{"lock_count":0}, "ip", r.Ip); err != nil{
		logger.WebLog.Warningf("解锁IP失败 数据库错误:%s", err.Error())
		return errors.New("数据库错误,解锁IP失败")
	}else{
		if result == nil{
			return errors.New("解锁IP失败,该IP不存在")
		}else{
			return nil
		}
	}
}

// SearchUserLoginLogs 登录日志管理模糊分页查询
func (s *serviceUser) SearchUserLoginLogs(page, limit int, search interface{}) *model.ResponeUserOperationLogins{
	var result []*model.UserLog
	SearchModel := dao.UserLog.Clone()
	searchStr := gconv.String(search)
	if search != ""{
		j := gjson.New(searchStr)
		if gconv.String(j.Get("username")) != ""{
			SearchModel = SearchModel.Where("username like ?", "%"+gconv.String(j.Get("username"))+"%")
		}
		if gconv.String(j.Get("lastip")) != ""{
			SearchModel = SearchModel.Where("ip like ?", "%"+gconv.String(j.Get("lastip"))+"%")
		}
	}
	count,_ := SearchModel.Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Order("id desc").Limit((page-1)*limit,limit).Scan(&result)
		if err != nil {
			logger.WebLog.Warningf("用户登录日志分页查询 数据库错误:%s", err.Error())
			return &model.ResponeUserOperationLogins{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponeUserOperationLogins{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].Id = uint(index)
	}
	return &model.ResponeUserOperationLogins{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// SearchUserOperation 操作日志管理模糊分页查询
func (s *serviceUser) SearchUserOperation(page, limit int, search interface{}) *model.ResponeUserOperationLogs{
	var result []*model.UserOperation
	SearchModel := dao.UserOperation.Clone()
	searchStr := gconv.String(search)
	if search != ""{
		j := gjson.New(searchStr)
		if gconv.String(j.Get("username")) != ""{
			SearchModel = SearchModel.Where("username like ?", "%"+gconv.String(j.Get("username"))+"%")
		}
		if gconv.String(j.Get("theme")) != ""{
			SearchModel = SearchModel.Where("theme like ?", "%"+gconv.String(j.Get("theme"))+"%")
		}
		if gconv.String(j.Get("content")) != ""{
			SearchModel = SearchModel.Where("content like ?", "%"+gconv.String(j.Get("content"))+"%")
		}
	}
	count,_ := SearchModel.Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Order("id desc").Limit((page-1)*limit,limit).Scan(&result)
		if err != nil {
			logger.WebLog.Warningf("用户操作日志分页查询 数据库错误:%s", err.Error())
			return &model.ResponeUserOperationLogs{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponeUserOperationLogs{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index ++
		result[i].Id = uint(index)
	}
	return &model.ResponeUserOperationLogs{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// OperationEmpty 清空用户操作日志
func (s *serviceUser) OperationEmpty()error{
	_,err := dao.UserOperation.Where("1=",1).Delete()
	if err != nil {
		logger.WebLog.Warningf("清空用户操作日志,数据库错误:%s", err.Error())
		return errors.New("清空用户操作日志失败,数据库错误")
	}
	return nil
}

// UserAddOperation 添加用户操作记录
func (s *serviceUser) UserAddOperation(ctx context.Context, ip, Theme, Content string){
	dao.UserOperation.Insert(g.Map{
		"Username": Session.GetUser(ctx).Username,
		"Ip": ip,
		"Theme": Theme,
		"Content": Content,
	})
}

// IsUserAdmin 判断是否为管理员用户
func (s *serviceUser)IsUserAdmin(r *ghttp.Request,ctx context.Context){
	user := Session.GetUser(ctx)
	if user.Username != "admin"{
		response.JsonExit(r,201, "非管理员用户无权限使用该功能")
	}
}