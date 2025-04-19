package model

type WsAgentSendMap struct {
	AgentId   int    `json:"agentId"`
	AgentName string `json:"agentName"`
	JobId     int    `json:"jobId"`
	JobStatus string `json:"jobStatus"`
	JobOutput string `json:"jobOutput"`
}

type WsAgentSend struct {
	Items      []WsAgentSendMap `json:"items"`
	TimeoutSec int              `json:"timeoutSec"`
}

type WsServerSend struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    []WsServerSendMap `json:"data"`
}

type WsServerSendMap struct {
	AgentId   int               `json:"agentId"`
	AgentName string            `json:"agentName"`
	JobId     int               `json:"jobId"`
	JobStatus string            `json:"jobStatus"`
	Body      string            `json:"scriptBody"`
	Envs      map[string]string `json:"scriptEnvs"`
	Args      string            `json:"scriptArgs"`

	// ErrMsg    string            `json:"errmsg"`
}
