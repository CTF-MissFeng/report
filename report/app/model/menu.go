package model

// FrameRoute 框架模块路由所需信息
type FrameRoute struct{
	HomeInfo ModuleRoute `json:"homeInfo"`
	LogoInfo ModuleRoute `json:"logoInfo"`
	MenuInfo []SonModuleRoute `json:"menuInfo"`
}

// ModuleRoute 主框架描述信息
type ModuleRoute struct{
	Title string `json:"title"`
	Image string `json:"image"`
	Href string `json:"href"`
}

// SonModuleRoute 主框架模块路由
type SonModuleRoute struct{
	Title string `json:"title"`
	Icon string `json:"icon"`
	Href string `json:"href"`
	Target string `json:"target"`
	Child []SonModuleRoute `json:"child"`
}

// ModuleInit 生成模块路由
func ModuleInit()FrameRoute{
	frameRoute := FrameRoute{}
	frameRoute.HomeInfo = ModuleRoute{
		Title: "首页",
		Href : "tongji",
	}
	frameRoute.LogoInfo = ModuleRoute{
		Title:"资产",
		Image:"/images/logo.png",
	}
	frameRoute.MenuInfo = []SonModuleRoute{
		AssetsMenu(),
		userMenu(),

	}
	return frameRoute
}

// user菜单
func userMenu()SonModuleRoute{
	return SonModuleRoute{
		Title:"后台管理",
		Icon:"fa fa-cog",
		Target:"_self",
		Child:[]SonModuleRoute{
			SonModuleRoute{
				Title:"用户管理",
				Icon:"fa fa-users",
				Target:"_self",
				Child:[]SonModuleRoute{
					SonModuleRoute{
						Title:"用户管理",
						Icon:"fa fa-user",
						Target:"_self",
						Href: "user/manager",
					},
					SonModuleRoute{
						Title:"IP锁定管理",
						Icon:"fa fa-unlock-alt",
						Target:"_self",
						Href: "user/userip",
					},
				},
			},
			SonModuleRoute{
				Title:"日志管理",
				Icon:"fa fa-book",
				Target:"_self",
				Child:[]SonModuleRoute{
					SonModuleRoute{
						Title:"登录日志",
						Icon:"fa fa-calendar",
						Target:"_self",
						Href: "user/loginlog",
					},
					SonModuleRoute{
						Title:"操作日志",
						Icon:"fa fa-calendar-o",
						Target:"_self",
						Href: "user/operation",
					},
				},
			},
		},
	}
}

// AssetsMenu 资产管理菜单
func AssetsMenu()SonModuleRoute{
	return SonModuleRoute{
		Title:"资产管理",
		Icon:"fa fa-archive",
		Target:"_self",
		Child:[]SonModuleRoute{
			SonModuleRoute{
				Title:"资产管理",
				Icon:"fa fa-dashboard",
				Target:"_self",
				Child:[]SonModuleRoute{
					SonModuleRoute{
						Title:"安全管理员",
						Icon:"fa fa-users",
						Target:"_self",
						Href: "assets/manager",
					},
					SonModuleRoute{
						Title:"主机资产",
						Icon:"fa fa-jsfiddle",
						Target:"_self",
						Href: "assets/type",
					},
					SonModuleRoute{
						Title:"终端资产",
						Icon:"fa fa-desktop",
						Target:"_self",
						Href: "assets/pc",
					},
				},
			},
			SonModuleRoute{
				Title:"业务系统",
				Icon:"fa fa-dashboard",
				Target:"_self",
				Child:[]SonModuleRoute{
					SonModuleRoute{
						Title:"业务系统",
						Icon:"fa fa-globe",
						Target:"_self",
						Href: "assets/web",
					},
					SonModuleRoute{
						Title:"渗透测试报告",
						Icon:"fa fa-user-secret",
						Target:"_self",
						Href: "assets/report",
					},
				},
			},
		},
	}
}

