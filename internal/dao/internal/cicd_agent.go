// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CicdAgentDao is the data access object for the table cicd_agent.
type CicdAgentDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of the current DAO.
	columns CicdAgentColumns // columns contains all the column names of Table for convenient usage.
}

// CicdAgentColumns defines and stores column names for the table cicd_agent.
type CicdAgentColumns struct {
	Id        string //
	AgentName string //
	Ipaddr    string //
	UpdatedAt string //
}

// cicdAgentColumns holds the columns for the table cicd_agent.
var cicdAgentColumns = CicdAgentColumns{
	Id:        "id",
	AgentName: "agent_name",
	Ipaddr:    "ipaddr",
	UpdatedAt: "updated_at",
}

// NewCicdAgentDao creates and returns a new DAO object for table data access.
func NewCicdAgentDao() *CicdAgentDao {
	return &CicdAgentDao{
		group:   "default",
		table:   "cicd_agent",
		columns: cicdAgentColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CicdAgentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CicdAgentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CicdAgentDao) Columns() CicdAgentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CicdAgentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CicdAgentDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *CicdAgentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
