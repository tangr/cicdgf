// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// CicdPipeline is the golang structure for table cicd_pipeline.
type CicdPipeline struct {
	Id           uint   `json:"id"           orm:"id"            description:""` //
	PipelineName string `json:"pipelineName" orm:"pipeline_name" description:""` //
	GroupId      int    `json:"groupId"      orm:"group_id"      description:""` //
	AgentId      int    `json:"agentId"      orm:"agent_id"      description:""` //
	Concurrency  int    `json:"concurrency"  orm:"concurrency"   description:""` //
	Body         string `json:"body"         orm:"body"          description:""` //
	Author       string `json:"author"       orm:"author"        description:""` //
	UpdatedAt    int64  `json:"updatedAt"    orm:"updated_at"    description:""` //
}
