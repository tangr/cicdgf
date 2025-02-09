// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CicdPipelineDao is the data access object for the table cicd_pipeline.
type CicdPipelineDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of the current DAO.
	columns CicdPipelineColumns // columns contains all the column names of Table for convenient usage.
}

// CicdPipelineColumns defines and stores column names for the table cicd_pipeline.
type CicdPipelineColumns struct {
	Id           string //
	PipelineName string //
	GroupId      string //
	AgentId      string //
	Concurrency  string //
	Body         string //
	Author       string //
	UpdatedAt    string //
}

// cicdPipelineColumns holds the columns for the table cicd_pipeline.
var cicdPipelineColumns = CicdPipelineColumns{
	Id:           "id",
	PipelineName: "pipeline_name",
	GroupId:      "group_id",
	AgentId:      "agent_id",
	Concurrency:  "concurrency",
	Body:         "body",
	Author:       "author",
	UpdatedAt:    "updated_at",
}

// NewCicdPipelineDao creates and returns a new DAO object for table data access.
func NewCicdPipelineDao() *CicdPipelineDao {
	return &CicdPipelineDao{
		group:   "default",
		table:   "cicd_pipeline",
		columns: cicdPipelineColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CicdPipelineDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CicdPipelineDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CicdPipelineDao) Columns() CicdPipelineColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CicdPipelineDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CicdPipelineDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *CicdPipelineDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
