// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CicdJob is the golang structure of table cicd_job for DAO operations like Where/Data.
type CicdJob struct {
	g.Meta      `orm:"table:cicd_job, do:true"`
	Id          interface{} //
	PipelineId  interface{} //
	AgentId     interface{} //
	Concurrency interface{} //
	JobType     interface{} //
	JobStatus   interface{} //
	Script      interface{} //
	Comment     interface{} //
	Author      interface{} //
	CreatedAt   interface{} //
}
