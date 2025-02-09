// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CicdPipeline is the golang structure of table cicd_pipeline for DAO operations like Where/Data.
type CicdPipeline struct {
	g.Meta       `orm:"table:cicd_pipeline, do:true"`
	Id           interface{} //
	PipelineName interface{} //
	GroupId      interface{} //
	AgentId      interface{} //
	Concurrency  interface{} //
	Body         interface{} //
	Author       interface{} //
	UpdatedAt    interface{} //
}
