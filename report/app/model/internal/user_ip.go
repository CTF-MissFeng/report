// ==========================================================================
// This is auto-generated by gf cli tool. DO NOT EDIT THIS FILE MANUALLY.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/os/gtime"
)

// UserIp is the golang structure for table user_ip.
type UserIp struct {
	Id        uint        `orm:"id,primary" json:"id"`         //
	Ip        string      `orm:"ip,unique"  json:"ip"`         //
	LockCount int         `orm:"lock_count" json:"lock_count"` //
	CreateAt  *gtime.Time `orm:"create_at"  json:"create_at"`  //
}
