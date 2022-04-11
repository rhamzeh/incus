//go:build linux && cgo && !agent

package db

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"database/sql"
	"fmt"

	"github.com/lxc/lxd/lxd/db/cluster"
	"github.com/lxc/lxd/lxd/db/query"
	"github.com/lxc/lxd/shared/api"
)

var _ = api.ServerEnvironment{}

var operationObjects = cluster.RegisterStmt(`
SELECT operations.id, operations.uuid, nodes.address AS node_address, operations.project_id, operations.node_id, operations.type
  FROM operations JOIN nodes ON operations.node_id = nodes.id
  ORDER BY operations.id, operations.uuid
`)

var operationObjectsByNodeID = cluster.RegisterStmt(`
SELECT operations.id, operations.uuid, nodes.address AS node_address, operations.project_id, operations.node_id, operations.type
  FROM operations JOIN nodes ON operations.node_id = nodes.id
  WHERE operations.node_id = ? ORDER BY operations.id, operations.uuid
`)

var operationObjectsByID = cluster.RegisterStmt(`
SELECT operations.id, operations.uuid, nodes.address AS node_address, operations.project_id, operations.node_id, operations.type
  FROM operations JOIN nodes ON operations.node_id = nodes.id
  WHERE operations.id = ? ORDER BY operations.id, operations.uuid
`)

var operationObjectsByUUID = cluster.RegisterStmt(`
SELECT operations.id, operations.uuid, nodes.address AS node_address, operations.project_id, operations.node_id, operations.type
  FROM operations JOIN nodes ON operations.node_id = nodes.id
  WHERE operations.uuid = ? ORDER BY operations.id, operations.uuid
`)

var operationCreateOrReplace = cluster.RegisterStmt(`
INSERT OR REPLACE INTO operations (uuid, project_id, node_id, type)
 VALUES (?, ?, ?, ?)
`)

var operationDeleteByUUID = cluster.RegisterStmt(`
DELETE FROM operations WHERE uuid = ?
`)

var operationDeleteByNodeID = cluster.RegisterStmt(`
DELETE FROM operations WHERE node_id = ?
`)

// GetOperations returns all available operations.
// generator: operation GetMany
func (c *ClusterTx) GetOperations(filter OperationFilter) ([]Operation, error) {
	var err error

	// Result slice.
	objects := make([]Operation, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []any

	if filter.UUID != nil && filter.ID == nil && filter.NodeID == nil {
		stmt = c.stmt(operationObjectsByUUID)
		args = []any{
			filter.UUID,
		}
	} else if filter.NodeID != nil && filter.ID == nil && filter.UUID == nil {
		stmt = c.stmt(operationObjectsByNodeID)
		args = []any{
			filter.NodeID,
		}
	} else if filter.ID != nil && filter.NodeID == nil && filter.UUID == nil {
		stmt = c.stmt(operationObjectsByID)
		args = []any{
			filter.ID,
		}
	} else if filter.ID == nil && filter.NodeID == nil && filter.UUID == nil {
		stmt = c.stmt(operationObjects)
		args = []any{}
	} else {
		return nil, fmt.Errorf("No statement exists for the given Filter")
	}

	// Dest function for scanning a row.
	dest := func(i int) []any {
		objects = append(objects, Operation{})
		return []any{
			&objects[i].ID,
			&objects[i].UUID,
			&objects[i].NodeAddress,
			&objects[i].ProjectID,
			&objects[i].NodeID,
			&objects[i].Type,
		}
	}

	// Select.
	err = query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"operations\" table: %w", err)
	}

	return objects, nil
}

// CreateOrReplaceOperation adds a new operation to the database.
// generator: operation CreateOrReplace
func (c *ClusterTx) CreateOrReplaceOperation(object Operation) (int64, error) {
	args := make([]any, 4)

	// Populate the statement arguments.
	args[0] = object.UUID
	args[1] = object.ProjectID
	args[2] = object.NodeID
	args[3] = object.Type

	// Prepared statement to use.
	stmt := c.stmt(operationCreateOrReplace)

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"operations\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"operations\" entry ID: %w", err)
	}

	return id, nil
}

// DeleteOperation deletes the operation matching the given key parameters.
// generator: operation DeleteOne-by-UUID
func (c *ClusterTx) DeleteOperation(uuid string) error {
	stmt := c.stmt(operationDeleteByUUID)
	result, err := stmt.Exec(uuid)
	if err != nil {
		return fmt.Errorf("Delete \"operations\": %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query deleted %d rows instead of 1", n)
	}

	return nil
}

// DeleteOperations deletes the operation matching the given key parameters.
// generator: operation DeleteMany-by-NodeID
func (c *ClusterTx) DeleteOperations(nodeID int64) error {
	stmt := c.stmt(operationDeleteByNodeID)
	result, err := stmt.Exec(nodeID)
	if err != nil {
		return fmt.Errorf("Delete \"operations\": %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	return nil
}
