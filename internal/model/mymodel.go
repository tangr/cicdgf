package model

type WsAgentSendMap struct {
	AgentId   int    `json:"agentId"`
	AgentName string `json:"agentName"`
	JobId     int    `json:"jobId"`
	JobStatus string `json:"jobStatus"`
	JobOutput string `json:"jobOutput"`

	TaskId     int    `json:"taskId"`
	TaskStatus string `json:"taskStatus"`
	TaskOutput string `json:"taskOutput"`
}

type WsAgentSend struct {
	Items      []WsAgentSendMap `json:"items"`
	TimeoutSec int              `json:"timeoutSec"`
}

type WsAgentSendLogMap struct {
	PipelineId int    `json:"pipelineId"`
	AgentId    int    `json:"agentId"`
	JobType    string `json:"jobType"`
	JobId      int    `json:"jobId"`
	TaskStatus string `json:"taskStatus"`
	Ipaddr     string `json:"ipaddr"`
	UpdatedAt  int64  `json:"updatedAt"`
	Output     string `json:"output"`
}

type WsServerSend struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    []WsServerSendMap `json:"data"`
}

type WsServerSendMap struct {
	AgentId   int    `json:"agentId"`
	AgentName string `json:"agentName"`
	JobId     int    `json:"jobId"`
	JobStatus string `json:"jobStatus"`

	TaskId     int               `json:"taskId"`
	TaskStatus string            `json:"taskStatus"`
	Body       string            `json:"scriptBody"`
	Envs       map[string]string `json:"scriptEnvs"`
	Args       string            `json:"scriptArgs"`

	// ErrMsg    string            `json:"errmsg"`
}

type TaskGetRes struct {
	Message string      `json:"message"`
	Data    TaskInfoMap `json:"data"`
}

type TaskInfoMap struct {
	Id         int    `json:"id"`
	PipelineId int    `json:"pipelineId"`
	AgentId    int    `json:"agentId"`
	JobType    string `json:"jobType"`
	JobId      int    `json:"jobId"`
	TaskStatus string `json:"taskStatus"`
	Ipaddr     string `json:"ipaddr"`
	UpdateAt   int    `json:"updateAt"`
	Output     string `json:"output"`
}

type JobGetRes struct {
	Message string     `json:"message"`
	Data    JobInfoMap `json:"data"`
}

type JobInfoMap struct {
	Id          int    `json:"id"`
	PipelineId  int    `json:"pipelineId"`
	AgentId     int    `json:"agentId"`
	Concurrency int    `json:"concurrency"`
	JobType     string `json:"jobType"`
	JobStatus   string `json:"jobStatus"`
	Script      string `json:"script"`
	Comment     string `json:"comment"`
	Author      string `json:"author"`
	CreatedAt   int    `json:"createdAt"`
}

type Script struct {
	Body string            `json:"scriptBody"`
	Envs map[string]string `json:"scriptEnvs"`
	Args string            `json:"scriptArgs"`
}
