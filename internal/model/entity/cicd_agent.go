// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// CicdAgent is the golang structure for table cicd_agent.
type CicdAgent struct {
	Id        int    `json:"id"        orm:"id"         description:""` //
	AgentName string `json:"agentName" orm:"agent_name" description:""` //
	Ipaddr    string `json:"ipaddr"    orm:"ipaddr"     description:""` //
	UpdatedAt int64  `json:"updatedAt" orm:"updated_at" description:""` //
}
