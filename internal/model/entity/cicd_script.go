// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// CicdScript is the golang structure for table cicd_script.
type CicdScript struct {
	Id         uint   `json:"id"         orm:"id"          description:""` //
	ScriptName string `json:"scriptName" orm:"script_name" description:""` //
	ScriptBody string `json:"scriptBody" orm:"script_body" description:""` //
	Author     string `json:"author"     orm:"author"      description:""` //
	UpdatedAt  uint64 `json:"updatedAt"  orm:"updated_at"  description:""` //
}
