package router

import (
	"assets/app/api"
	"assets/app/service"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()

	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware.Ctx)

		group.POST("/login", api.Users.Login)

		group.Group("/user", func(group *ghttp.RouterGroup) {
			group.Middleware(service.Middleware.Auth)
			group.GET("/index", api.Users.UserInfo)
			group.POST("/index", api.Users.SetUserInfo)
			group.PUT("/index", api.Users.ChangePassword)
			group.DELETE("/index", api.Users.UserDel)
			group.POST("/add", api.Users.Register)
			group.GET("/menu", api.Users.Menu)
			group.GET("/manager", api.Users.SearchUser)
			group.GET("/lock", api.Users.SearchUserLockIp)
			group.DELETE("/lock", api.Users.UserLockIpRest)
			group.GET("/log", api.Users.SearchUserLoginLogs)
			group.GET("/logs", api.Users.SearchUserOperation)
			group.GET("/out", api.Users.LoginOut)
			group.POST("/operation",api.Users.OperationEmpty)
		})

		group.Group("/assets", func(group *ghttp.RouterGroup) {
			group.Middleware(service.Middleware.Auth)
			group.PUT("/manager", api.Assert.ManagerAdd)
			group.GET("/manager", api.Assert.SearchManager)
			group.DELETE("/manager", api.Assert.ManagerDelete)

			group.GET("/group/type", api.Assert.GroupAssetsType)
			group.GET("/group/attribution", api.Assert.GroupAssetsAttribution)
			group.GET("/group/department",api.Assert.GroupAssetsDepartment)
			group.GET("/type", api.Assert.SearchType)
			group.PUT("/type", api.Assert.TypeAdd)
			group.DELETE("/type",api.Assert.TypeDelete)
			group.POST("/type",api.Assert.EmptyType)
			group.POST("/typeinfo", api.Assert.TypeShow)
			group.GET("/type/export",api.Assert.ExportType)
			group.GET("/type/download",api.Assert.DownloadTypeTemplate)
			group.POST("/type/upload", api.Assert.ImportType)

			group.GET("/group/pc/department",api.Assert.GroupAssetsPcDepartment)
			group.GET("/group/pc/departmentsub",api.Assert.GroupAssetsPcDepartmentSub)
			group.GET("/pc", api.Assert.SearchPc)
			group.PUT("/pc", api.Assert.PcAdd)
			group.DELETE("/pc",api.Assert.PcDelete)
			group.POST("/pc",api.Assert.PcEmpty)
			group.POST("/pcinfo", api.Assert.PcShow)
			group.GET("/pc/download",api.Assert.DownloadPcTemplate)
			group.GET("/pc/export",api.Assert.ExportPc)
			group.POST("/pc/upload", api.Assert.ImportPc)
			group.GET("/pc/users",api.Assert.IchanganUsers)
			group.POST("/pcuserinfo", api.Assert.PcShowUser)

			group.GET("/group/web/managername",api.Assert.GroupAssetsManagerName)
			group.GET("/group/web/managernames",api.Assert.GroupAssetsWebManagerName)
			group.GET("/group/web/assetsname",api.Assert.GroupAssetsWebAssetsName)
			group.GET("/group/web/attribution",api.Assert.GroupAssetsWebAttribution)
			group.GET("/web",api.Assert.SearchWeb)
			group.PUT("/web", api.Assert.WebAdd)
			group.POST("/web/imgupload",api.Assert.UploadAssetsImg)
			group.POST("/web/fileupload",api.Assert.UploadAssetsFile)
			group.POST("/web/reportupload",api.Assert.UploadReportFile)
			group.DELETE("/web",api.Assert.WebDelete)
			group.GET("/web/export",api.Assert.ExportWeb)
			group.PUT("/web/export", api.Assert.WebExportAdd)
			group.GET("/web/exports",api.Assert.SearchWebReport)
			group.POST("/webinfo", api.Assert.WebShow)

			group.GET("/group/report/levelname",api.Assert.GroupAssetsReportLevelName)
			group.GET("/group/report/attribution",api.Assert.GroupAssetsReportAttribution)
			group.GET("/group/report/managername",api.Assert.GroupAssetsReportManagerName)
			group.GET("/group/report/assetsname",api.Assert.GroupAssetsReportAssetsName)
			group.GET("/report",api.Assert.SearchReport)
			group.DELETE("/report",api.Assert.ReportDelete)
			group.GET("/report/export",api.Assert.ExportReport)
			group.POST("/reportinfo", api.Assert.ReportShow)
			group.POST("/report/edit",api.Assert.ReportEdit)

			group.GET("/home/count",api.Assert.TongJiCount)
			group.GET("/home/echarts",api.Assert.EchartsInfo)
		})
	})
}

