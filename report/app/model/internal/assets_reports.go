// ==========================================================================
// This is auto-generated by gf cli tool. DO NOT EDIT THIS FILE MANUALLY.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/os/gtime"
)

// AssetsReports is the golang structure for table assets_reports.
type AssetsReports struct {
	Id          uint        `orm:"id,primary"   json:"id"`           //
	Attribution string      `orm:"attribution"  json:"attribution"`  //
	ManagerName string      `orm:"manager_name" json:"manager_name"` //
	AssetsName  string      `orm:"assets_name"  json:"assets_name"`  //
	Level       int         `orm:"level"        json:"level"`        //
	LevelName   string      `orm:"level_name"   json:"level_name"`   //
	LevelStatus int         `orm:"level_status" json:"level_status"` //
	FilePath    string      `orm:"file_path"    json:"file_path"`    //
	FileDate    *gtime.Time `orm:"file_date"    json:"file_date"`    //
	CreateAt    *gtime.Time `orm:"create_at"    json:"create_at"`    //
}