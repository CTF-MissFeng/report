// ==========================================================================
// This is auto-generated by gf cli tool. DO NOT EDIT THIS FILE MANUALLY.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/os/gtime"
)

// UserOperation is the golang structure for table user_operation.
type UserOperation struct {
	Id       uint        `orm:"id,primary" json:"id"`        //
	Username string      `orm:"username"   json:"username"`  //
	Ip       string      `orm:"ip"         json:"ip"`        //
	Theme    string      `orm:"theme"      json:"theme"`     //
	Content  string      `orm:"content"    json:"content"`   //
	CreateAt *gtime.Time `orm:"create_at"  json:"create_at"` //
}
