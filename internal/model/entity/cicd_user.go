// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// CicdUser is the golang structure for table cicd_user.
type CicdUser struct {
	Id        uint   `json:"id"        orm:"id"         description:""` //
	Email     string `json:"email"     orm:"email"      description:""` //
	Password  string `json:"password"  orm:"password"   description:""` //
	GroupId   string `json:"groupId"   orm:"group_id"   description:""` //
	UpdatedAt int64  `json:"updatedAt" orm:"updated_at" description:""` //
}
