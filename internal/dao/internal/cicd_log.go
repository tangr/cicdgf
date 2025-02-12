// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CicdLogDao is the data access object for the table cicd_log.
type CicdLogDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of the current DAO.
	columns CicdLogColumns // columns contains all the column names of Table for convenient usage.
}

// CicdLogColumns defines and stores column names for the table cicd_log.
type CicdLogColumns struct {
	Id         string //
	PipelineId string //
	AgentId    string //
	JobType    string //
	JobId      string //
	TaskStatus string //
	Ipaddr     string //
	UpdatedAt  string //
	Output     string //
}

// cicdLogColumns holds the columns for the table cicd_log.
var cicdLogColumns = CicdLogColumns{
	Id:         "id",
	PipelineId: "pipeline_id",
	AgentId:    "agent_id",
	JobType:    "job_type",
	JobId:      "job_id",
	TaskStatus: "task_status",
	Ipaddr:     "ipaddr",
	UpdatedAt:  "updated_at",
	Output:     "output",
}

// NewCicdLogDao creates and returns a new DAO object for table data access.
func NewCicdLogDao() *CicdLogDao {
	return &CicdLogDao{
		group:   "default",
		table:   "cicd_log",
		columns: cicdLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CicdLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CicdLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CicdLogDao) Columns() CicdLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CicdLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CicdLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *CicdLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
