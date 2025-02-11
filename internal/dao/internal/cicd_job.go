// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CicdJobDao is the data access object for the table cicd_job.
type CicdJobDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of the current DAO.
	columns CicdJobColumns // columns contains all the column names of Table for convenient usage.
}

// CicdJobColumns defines and stores column names for the table cicd_job.
type CicdJobColumns struct {
	Id          string //
	PipelineId  string //
	AgentId     string //
	Concurrency string //
	JobType     string //
	JobStatus   string //
	Script      string //
	Comment     string //
	Author      string //
	CreatedAt   string //
}

// cicdJobColumns holds the columns for the table cicd_job.
var cicdJobColumns = CicdJobColumns{
	Id:          "id",
	PipelineId:  "pipeline_id",
	AgentId:     "agent_id",
	Concurrency: "concurrency",
	JobType:     "job_type",
	JobStatus:   "job_status",
	Script:      "script",
	Comment:     "comment",
	Author:      "author",
	CreatedAt:   "created_at",
}

// NewCicdJobDao creates and returns a new DAO object for table data access.
func NewCicdJobDao() *CicdJobDao {
	return &CicdJobDao{
		group:   "default",
		table:   "cicd_job",
		columns: cicdJobColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CicdJobDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CicdJobDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CicdJobDao) Columns() CicdJobColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CicdJobDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CicdJobDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *CicdJobDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
