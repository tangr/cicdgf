// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// CicdGroup is the golang structure for table cicd_group.
type CicdGroup struct {
	Id        int    `json:"id"        orm:"id"         description:""` //
	GroupName string `json:"groupName" orm:"group_name" description:""` //
	ParentId  int    `json:"parentId"  orm:"parent_id"  description:""` //
}
