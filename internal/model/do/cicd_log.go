// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CicdLog is the golang structure of table cicd_log for DAO operations like Where/Data.
type CicdLog struct {
	g.Meta     `orm:"table:cicd_log, do:true"`
	Id         interface{} //
	PipelineId interface{} //
	AgentId    interface{} //
	JobType    interface{} //
	JobId      interface{} //
	TaskStatus interface{} //
	Ipaddr     interface{} //
	UpdatedAt  interface{} //
	Output     interface{} //
}
