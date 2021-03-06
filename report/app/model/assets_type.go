// ==========================================================================
// This is auto-generated by gf cli tool. Fill this file as you wish.
// ==========================================================================

package model

import (
	"assets/app/model/internal"
)

// AssetsType is the golang structure for table assets_type.
type AssetsType internal.AssetsType

// Fill with you ideas below.

// ResponseAssetsTypeGroup 主机资产 返回厂商组信息
type ResponseAssetsTypeGroup struct{
	Code int `json:"code"`
	Msg string `json:"msg"`
	Count int64 `json:"count"`
	Data []ResponseAssetsTypeGroupInfo `json:"data"`
}

// ResponseAssetsTypeGroupInfo 主机资产 返回厂商详细信息
type ResponseAssetsTypeGroupInfo struct{
	ID int `json:"id"`
	TypeName string `json:"type_name"`
}

// ResponseAssetsTypeAttributionGroup 主机资产 返回应用系统组信息
type ResponseAssetsTypeAttributionGroup struct{
	Code int `json:"code"`
	Msg string `json:"msg"`
	Count int64 `json:"count"`
	Data []ResponseAssetsTypeAttributionGroupInfo `json:"data"`
}

// ResponseAssetsTypeAttributionGroupInfo 主机资产 返回应用系统组详细信息
type ResponseAssetsTypeAttributionGroupInfo struct{
	ID int `json:"id"`
	Attribution string `json:"attribution"`
}

// ResponseAssetsTypeDepartmentGroup 主机资产 返回部门组信息
type ResponseAssetsTypeDepartmentGroup struct{
	Code int `json:"code"`
	Msg string `json:"msg"`
	Count int64 `json:"count"`
	Data []ResponseAssetsTypeDepartmentGroupInfo `json:"data"`
}

// ResponseAssetsTypeDepartmentGroupInfo 主机资产 返回部门组详细信息
type ResponseAssetsTypeDepartmentGroupInfo struct{
	ID int `json:"id"`
	Department string `json:"department"`
}

// ResponseType 主机资产管理 模糊分页查询返回数据所需信息
type ResponseType struct{
	Code int `json:"code"`
	Msg string `json:"msg"`
	Count int64 `json:"count"`
	Data []*ResponseTypeInfo `json:"data"`
}

// ResponseTypeInfo 主机资产管理 模糊分页查询返回数据所需详细信息
type ResponseTypeInfo struct{
	ID1 uint `json:"id1"`
	*AssetsType
}

// RequestAssetsTypeAdd 添加主机资产所需信息
type RequestAssetsTypeAdd struct{
	TypeName string `v:"required#厂商不能为空"`
	AttriBution string `v:"required#应用系统名不能为空"`
	Department string
	SubDomain string
	IntranetIp string
	PublicIp string
	AssetsUserName string
	TypeID string
}

// RequestAssetsTypeDelete 删除主机资产信息
type RequestAssetsTypeDelete struct {
	ID string `v:"required|integer#ID标识不能为空|ID值必须为整数"`
}