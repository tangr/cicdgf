// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// CicdJob is the golang structure for table cicd_job.
type CicdJob struct {
	Id          int64  `json:"id"          orm:"id"          description:""` //
	PipelineId  int    `json:"pipelineId"  orm:"pipeline_id" description:""` //
	AgentId     int    `json:"agentId"     orm:"agent_id"    description:""` //
	Concurrency int    `json:"concurrency" orm:"concurrency" description:""` //
	JobType     string `json:"jobType"     orm:"job_type"    description:""` //
	JobStatus   string `json:"jobStatus"   orm:"job_status"  description:""` //
	Script      string `json:"script"      orm:"script"      description:""` //
	Comment     string `json:"comment"     orm:"comment"     description:""` //
	Author      string `json:"author"      orm:"author"      description:""` //
	CreatedAt   int64  `json:"createdAt"   orm:"created_at"  description:""` //
}
