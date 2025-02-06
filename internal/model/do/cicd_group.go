// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CicdGroup is the golang structure of table cicd_group for DAO operations like Where/Data.
type CicdGroup struct {
	g.Meta    `orm:"table:cicd_group, do:true"`
	Id        interface{} //
	GroupName interface{} //
	ParentId  interface{} //
}
