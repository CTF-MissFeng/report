package api

import (
	"fmt"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/os/gtime"

	"assets/app/model"
	"assets/app/service"
	"assets/library/response"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
)

var Assert = new(apiAssets)

type apiAssets struct{}

// ManagerAdd 添加安全管理员
func (a *apiAssets) ManagerAdd(r *ghttp.Request){
	var data *model.RequestAssetsManagerAdd
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.Assets.ManagerAdd(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "安全管理员", fmt.Sprintf("添加安全管理员:[%s]", data.ManagerName))
		response.JsonExit(r, 200, "ok")
	}
}

// ManagerDelete 删除安全管理员
func (a *apiAssets) ManagerDelete(r *ghttp.Request){
	service.User.IsUserAdmin(r,r.Context())
	var data *model.RequestAssetsManagerAdd
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.Assets.ManagerDelete(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "安全管理员", fmt.Sprintf("删除管理员:[%s]", data.ManagerName))
		response.JsonExit(r, 200, "ok")
	}
}

// SearchManager 安全管理员模糊搜索
func (a *apiAssets) SearchManager(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.SearchManager(r.GetInt("page"), r.GetInt("limit"), r.Get("searchParams")))
}


// GroupAssetsType 返回主机资产厂商组
func (a *apiAssets) GroupAssetsType(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsType(r.GetInt("page"), r.GetInt("limit"), r.Get("TypeName")))
}

// GroupAssetsAttribution 返回主机资产应用系统组
func (a *apiAssets) GroupAssetsAttribution(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsAttribution(r.GetInt("page"), r.GetInt("limit"), r.Get("Attribution")))
}

// GroupAssetsDepartment 返回主机资产部门组
func (a *apiAssets) GroupAssetsDepartment(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsDepartment(r.GetInt("page"), r.GetInt("limit"), r.Get("Department")))
}

// SearchType 主机资产模糊分页查询
func (a *apiAssets) SearchType(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.SearchType(r.GetInt("page"), r.GetInt("limit"), r.Get("searchParams")))
}

// TypeAdd 添加/修改主机资产
func (a *apiAssets) TypeAdd(r *ghttp.Request){
	var data *model.RequestAssetsTypeAdd
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.Assets.TypeAdd(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "主机资产", fmt.Sprintf("添加/修改主机资产-系统名称:[%s]", data.AttriBution))
		response.JsonExit(r, 200, "ok")
	}
}

// TypeDelete 删除主机资产
func (a *apiAssets) TypeDelete(r *ghttp.Request){
	service.User.IsUserAdmin(r,r.Context())
	var data *model.RequestAssetsTypeDelete
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.Assets.TypeDelete(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "主机资产", fmt.Sprintf("删除主机资产 ID:[%s]", data.ID))
		response.JsonExit(r, 200, "ok")
	}
}

// TypeShow 查看指定ID的主机资产信息
func (a *apiAssets) TypeShow(r *ghttp.Request){
	var data *model.RequestAssetsTypeDelete
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if result,err := service.Assets.TypeShow(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		r.Response.WriteJson(g.Map{"code":200,"msg":"ok","data":result})
	}
}

// ExportType 导出主机资产表
func (a *apiAssets) ExportType(r *ghttp.Request){
	result,err := service.Assets.ExportType()
	if err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "主机资产", "导出主机资产")
	taskName := "主机资产表"
	r.Response.Header().Set("Content-Type", "application/octet-stream")
	r.Response.Header().Set("Content-Disposition", "attachment; filename="+taskName+".xlsx")
	r.Response.Header().Set("Content-Transfer-Encoding", "binary")
	r.Response.WriteExit(result)
}

// DownloadTypeTemplate 下载主机资产模板
func (a *apiAssets) DownloadTypeTemplate(r *ghttp.Request){
	result,err := service.Assets.DownloadTypeTemplate()
	if err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	taskName := "主机资产模板"
	r.Response.Header().Set("Content-Type", "application/octet-stream")
	r.Response.Header().Set("Content-Disposition", "attachment; filename="+taskName+".xlsx")
	r.Response.Header().Set("Content-Transfer-Encoding", "binary")
	r.Response.WriteExit(result)
}

// ImportType 导入主机资产表
func (a *apiAssets) ImportType(r *ghttp.Request){
	file := r.GetUploadFile("file")
	if file == nil {
		response.JsonExit(r, 201, "导入主机资产失败,文件不存在")
	}
	filename, err := file.Save("./upload/tmp/")
	if err != nil{
		response.JsonExit(r, 201, "导入主机资产表失败,保存文件失败")
	}
	filename = "./upload/tmp/" + filename
	msg,err := service.Assets.ImportType(filename)
	gfile.Remove(filename)
	if err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "主机资产", msg)
	response.JsonExit(r, 200, msg)
}

// EmptyType 清空主机资产表
func (a *apiAssets) EmptyType(r *ghttp.Request){
	service.User.IsUserAdmin(r,r.Context())
	err := service.Assets.EmptyType()
	if err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "主机资产", "清空主机资产")
	response.JsonExit(r, 200, "ok")
}


// GroupAssetsPcDepartment 返回终端资产部门组
func (a *apiAssets) GroupAssetsPcDepartment(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsPcDepartment(r.GetInt("page"), r.GetInt("limit"), r.Get("Department")))
}

// GroupAssetsPcDepartmentSub 返回终端资产二级部门组
func (a *apiAssets) GroupAssetsPcDepartmentSub(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsPcDepartmentSub(r.GetInt("page"), r.GetInt("limit"), r.Get("DepartmentSub")))
}

// SearchPc 终端资产所属模糊分页查询
func (a *apiAssets) SearchPc(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.SearchPc(r.GetInt("page"), r.GetInt("limit"), r.Get("searchParams")))
}

// PcAdd 添加修改终端资产
func (a *apiAssets) PcAdd(r *ghttp.Request){
	var data *model.RequestAssetsPcAdd
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.Assets.PcAdd(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "终端资产", fmt.Sprintf("添加/修改终端资产-员工姓名:[%s]", data.PersonName))
		response.JsonExit(r, 200, "ok")
	}
}

// PcShow 查看指定ID的终端资产信息
func (a *apiAssets) PcShow(r *ghttp.Request){
	var data *model.RequestAssetsTypeDelete
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if result,err := service.Assets.PcShow(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		r.Response.WriteJson(g.Map{"code":200,"msg":"ok","data":result})
	}
}

// PcDelete 删除终端资产
func (a *apiAssets) PcDelete(r *ghttp.Request){
	var data *model.RequestAssetsTypeDelete
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.Assets.PcDelete(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "终端资产", fmt.Sprintf("删除终端资产 ID:[%s]", data.ID))
		response.JsonExit(r, 200, "ok")
	}
}

// PcEmpty 清空终端资产
func (a *apiAssets) PcEmpty(r *ghttp.Request){
	service.User.IsUserAdmin(r,r.Context())
	err := service.Assets.PcEmpty()
	if err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "终端资产", "清空终端资产")
	response.JsonExit(r, 200, "ok")
}

// DownloadPcTemplate 下载终端资产模板
func (a *apiAssets) DownloadPcTemplate(r *ghttp.Request){
	result,err := service.Assets.DownloadPcTemplate()
	if err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	taskName := "终端资产模板"
	r.Response.Header().Set("Content-Type", "application/octet-stream")
	r.Response.Header().Set("Content-Disposition", "attachment; filename="+taskName+".xlsx")
	r.Response.Header().Set("Content-Transfer-Encoding", "binary")
	r.Response.WriteExit(result)
}

// ExportPc 导出终端资产
func (a *apiAssets) ExportPc(r *ghttp.Request){
	result,err := service.Assets.ExportPc()
	if err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "终端资产", "导出终端资产")
	taskName := "终端资产表"
	r.Response.Header().Set("Content-Type", "application/octet-stream")
	r.Response.Header().Set("Content-Disposition", "attachment; filename="+taskName+".xlsx")
	r.Response.Header().Set("Content-Transfer-Encoding", "binary")
	r.Response.WriteExit(result)
}

// ImportPc 导入终端资产表
func (a *apiAssets) ImportPc(r *ghttp.Request){
	file := r.GetUploadFile("file")
	if file == nil {
		response.JsonExit(r, 201, "导入终端资产失败,文件不存在")
	}
	filename, err := file.Save("./upload/tmp/")
	if err != nil{
		response.JsonExit(r, 201, "导入终端资产失败,保存文件失败")
	}
	filename = "./upload/tmp/" + filename
	msg,err := service.Assets.ImportPc(filename)
	gfile.Remove(filename)
	if err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "终端资产", msg)
	response.JsonExit(r, 200, msg)
}

// IchanganUsers 加载员工数据
func (a *apiAssets) IchanganUsers(r *ghttp.Request){
	service.User.IsUserAdmin(r,r.Context())
	err := service.Assets.IchanganUsers()
	if err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "终端资产", "加载员工数据成功")
	response.JsonExit(r, 200, "ok")
}

// PcShowUser 查看指定ID的终端和员工信息
func (a *apiAssets) PcShowUser(r *ghttp.Request){
	var data *model.RequestAssetsTypeDelete
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if result,err := service.Assets.PcShowUser(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		r.Response.WriteJson(g.Map{"code":200,"msg":"ok","data":result})
	}
}


// GroupAssetsManagerName 返回安全管理员Group分组
func (a *apiAssets) GroupAssetsManagerName(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsManagerName(r.GetInt("page"), r.GetInt("limit"), r.Get("ManagerName")))
}

// GroupAssetsWebManagerName 返回Web资产安全管理员Group分组
func (a *apiAssets) GroupAssetsWebManagerName(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsWebManagerName(r.GetInt("page"), r.GetInt("limit"), r.Get("ManagerName")))
}

// GroupAssetsWebAssetsName 返回web应用业务系统名Group分组
func (a *apiAssets) GroupAssetsWebAssetsName(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsWebAssetsName(r.GetInt("page"), r.GetInt("limit"), r.Get("AssetsName")))
}

// GroupAssetsWebAttribution 返回Web资产应用系统组
func (a *apiAssets) GroupAssetsWebAttribution(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsWebAttribution(r.GetInt("page"), r.GetInt("limit"), r.Get("Attribution")))
}

// SearchWeb 应用系统模糊分页查询
func (a *apiAssets) SearchWeb(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.SearchWeb(r.GetInt("page"), r.GetInt("limit"), r.Get("searchParams")))
}

// UploadAssetsImg 上传业务系统截图
func (a *apiAssets) UploadAssetsImg(r *ghttp.Request){
	file := r.GetUploadFile("file")
	if file == nil {
		response.JsonExit(r, 201, "上传失败,文件不存在")
	}
	exts := gfile.Ext(file.Filename)
	if exts != ".jpg" && exts != ".png" && exts != ".jpeg"{
		response.JsonExit(r, 201, "上传失败,只允许jpg|png|jpeg文件")
	}
	md5File,err := gmd5.EncryptString(gtime.Datetime())
	if err != nil {
		response.JsonExit(r, 201, "上传失败,生成文件名错误")
	}
	file.Filename = md5File + exts
	filename, err := file.Save("./public/upload/assetsImg/")
	if err != nil{
		response.JsonExit(r, 201, "上传失败,保存文件失败")
	}
	response.JsonExit(r, 200, "/upload/assetsImg/" + filename)
}

// UploadAssetsFile 上传业务系统资产附件
func (a *apiAssets) UploadAssetsFile(r *ghttp.Request){
	file := r.GetUploadFile("file")
	if file == nil {
		response.JsonExit(r, 201, "上传失败,文件不存在")
	}
	exts := gfile.Ext(file.Filename)
	if len(exts) < 3{
		response.JsonExit(r, 201, "上传失败,文件后缀不正确")
	}
	md5File,err := gmd5.EncryptString(gtime.Datetime())
	if err != nil {
		response.JsonExit(r, 201, "上传失败,生成文件名错误")
	}
	file.Filename = md5File + exts
	filename, err := file.Save("./public/upload/assetsFile/")
	if err != nil{
		response.JsonExit(r, 201, "上传失败,保存文件失败")
	}
	response.JsonExit(r, 200, "/upload/assetsFile/" + filename)
}

// WebAdd 添加修改业务系统资产
func (a *apiAssets) WebAdd(r *ghttp.Request){
	var data *model.RequestAssetsWebAdd
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.Assets.WebAdd(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "业务系统", fmt.Sprintf("添加/修改业务系统:%s", data.AssetsName))
		response.JsonExit(r, 200, "ok")
	}
}

// WebDelete 删除业务系统
func (a *apiAssets) WebDelete(r *ghttp.Request){
	service.User.IsUserAdmin(r,r.Context())
	var data *model.RequestAssetsTypeDelete
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.Assets.WebDelete(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "业务系统", fmt.Sprintf("删除业务系统 ID:[%s]", data.ID))
		response.JsonExit(r, 200, "ok")
	}
}

// ExportWeb 导出业务系统
func (a *apiAssets) ExportWeb(r *ghttp.Request){
	result,err := service.Assets.ExportWeb()
	if err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "业务系统", "导出业务系统资产")
	taskName := "业务系统资产表"
	r.Response.Header().Set("Content-Type", "application/octet-stream")
	r.Response.Header().Set("Content-Disposition", "attachment; filename="+taskName+".xlsx")
	r.Response.Header().Set("Content-Transfer-Encoding", "binary")
	r.Response.WriteExit(result)
}

// UploadReportFile 上传渗透测试报告附件
func (a *apiAssets) UploadReportFile(r *ghttp.Request){
	file := r.GetUploadFile("file")
	if file == nil {
		response.JsonExit(r, 201, "上传失败,文件不存在")
	}
	exts := gfile.Ext(file.Filename)
	if len(exts) < 2{
		response.JsonExit(r, 201, "上传失败,文件后缀不正确")
	}
	md5File,err := gmd5.EncryptString(gtime.Datetime())
	if err != nil {
		response.JsonExit(r, 201, "上传失败,生成文件名错误")
	}
	file.Filename = md5File + exts
	filename, err := file.Save("./public/upload/report/")
	if err != nil{
		response.JsonExit(r, 201, "上传失败,保存文件失败")
	}
	response.JsonExit(r, 200, "/upload/report/" + filename)
}

// WebExportAdd 增加渗透测试漏洞
func (a *apiAssets) WebExportAdd(r *ghttp.Request){
	rjson,err := r.GetJson()
	if err != nil {
		response.JsonExit(r, 201, err.Error())
	}
	if err = service.Assets.WebExportAdd(rjson); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	response.JsonExit(r, 200, "ok")
}

// SearchWebReport 渗透测试报告模糊分页查询
func (a *apiAssets) SearchWebReport(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.SearchWebReport(r.GetInt("page"), r.GetInt("limit"), r.GetInt("assetsWebId")))
}

// WebShow 查看指定ID的业务系统信息
func (a *apiAssets) WebShow(r *ghttp.Request){
	var data *model.RequestAssetsTypeDelete
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if result,err := service.Assets.WebShow(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		r.Response.WriteJson(g.Map{"code":200,"msg":"ok","data":result})
	}
}


// GroupAssetsReportLevelName 返回渗透测试报告漏洞名Group分组
func (a *apiAssets) GroupAssetsReportLevelName(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsReportLevelName(r.GetInt("page"), r.GetInt("limit"), r.Get("level_name")))
}

// GroupAssetsReportAttribution 返回渗透测试报告应用系统组
func (a *apiAssets) GroupAssetsReportAttribution(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsReportAttribution(r.GetInt("page"), r.GetInt("limit"), r.Get("Attribution")))
}

// GroupAssetsReportManagerName 返回安全管理员Group分组
func (a *apiAssets) GroupAssetsReportManagerName(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsReportManagerName(r.GetInt("page"), r.GetInt("limit"), r.Get("ManagerName")))
}

// GroupAssetsReportAssetsName 返回web应用业务系统名Group分组
func (a *apiAssets) GroupAssetsReportAssetsName(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.GroupAssetsReportAssetsName(r.GetInt("page"), r.GetInt("limit"), r.Get("AssetsName")))
}

// SearchReport 渗透测试报告模糊分页查询
func (a *apiAssets) SearchReport(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.SearchReport(r.GetInt("page"), r.GetInt("limit"), r.Get("searchParams")))
}

// ReportDelete 删除漏洞
func (a *apiAssets) ReportDelete(r *ghttp.Request){
	service.User.IsUserAdmin(r,r.Context())
	var data *model.RequestAssetsTypeDelete
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.Assets.ReportDelete(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "渗透测试报告", fmt.Sprintf("删除漏洞 ID:[%s]", data.ID))
		response.JsonExit(r, 200, "ok")
	}
}

// ExportReport 导出漏洞
func (a *apiAssets) ExportReport(r *ghttp.Request){
	result,err := service.Assets.ExportReport()
	if err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "渗透测试报告", "导出所有漏洞")
	taskName := "业务系统漏洞统计表"
	r.Response.Header().Set("Content-Type", "application/octet-stream")
	r.Response.Header().Set("Content-Disposition", "attachment; filename="+taskName+".xlsx")
	r.Response.Header().Set("Content-Transfer-Encoding", "binary")
	r.Response.WriteExit(result)
}

// ReportShow 查看指定ID的漏洞报告
func (a *apiAssets) ReportShow(r *ghttp.Request){
	var data *model.RequestAssetsTypeDelete
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if result,err := service.Assets.ReportShow(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		r.Response.WriteJson(g.Map{"code":200,"msg":"ok","data":result})
	}
}

// ReportEdit 修改漏洞报告
func (a *apiAssets) ReportEdit(r *ghttp.Request){
	var data *model.RequestReport
	if err := r.Parse(&data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}
	if err := service.Assets.ReportEdit(data); err != nil{
		response.JsonExit(r, 201, err.Error())
	}else{
		service.User.UserAddOperation(r.Context(), r.GetRemoteIp(), "渗透测试报告", fmt.Sprintf("修改-漏洞名:[%s]", data.LevelName))
		response.JsonExit(r, 200, "ok")
	}
}


// TongJiCount 资产数据统计
func (a *apiAssets) TongJiCount(r *ghttp.Request){
	r.Response.WriteJson(g.Map{"code":200,"msg":"ok","data":service.Assets.TongJiCount()})
}

// EchartsInfo Echarts图标统计信息
func (a *apiAssets) EchartsInfo(r *ghttp.Request){
	r.Response.WriteJson(service.Assets.EchartsInfo())
}