package services

import (
	"database/sql"
	"fmt"
	"go-data-migration/models"
	"log"
	"strings"
)

type MigrationService struct {
	sourceDB *sql.DB
	destDB   *sql.DB
	config   models.MigrationConfig
}

func NewMigrationService(sourceDB, destDB *sql.DB, config models.MigrationConfig) *MigrationService {
	return &MigrationService{
		sourceDB: sourceDB,
		destDB:   destDB,
		config:   config,
	}
}

func (m *MigrationService) MigrateData() error {

	// Migrate main table
	mainTableIds, err := m.migrateMainTable()
	if err != nil {
		return err
	}
	// Migrate related tables
	for _, relatedTable := range m.config.RelatedTables {
		err = m.migrateRelatedTable(relatedTable, mainTableIds)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MigrationService) migrateMainTable() ([]int, error) {
	table := m.config.MainTable
	exists, err := TableExists(m.destDB, table.Name)
	if err != nil {
		return nil, err
	}
	log.Printf("Table %s exists: %v", table.Name, exists)
	if !exists {
		err = CreateTable(m.destDB, m.sourceDB, table.Name)
		if err != nil {
			return nil, err
		}
	}

	quotedColumns := make([]string, len(table.Columns))
	for i, col := range table.Columns {
		quotedColumns[i] = fmt.Sprintf("`%s`", col)
	}
	columns := strings.Join(quotedColumns, ", ")

	// Build the actual SQL query with the value directly (for this specific case)
	query := fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` = '%s'",
		columns,
		table.Name,
		table.Filter.Column,
		table.Filter.Value,
	)

	// Print the complete SQL query
	fmt.Printf("\n=== EXECUTING SQL QUERY ===\n%s\n==========================\n", query)

	// Execute the query
	rows, err := m.sourceDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("select query failed: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		values := make([]interface{}, len(table.Columns))
		valuePtrs := make([]interface{}, len(table.Columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		placeholders := make([]string, len(table.Columns))
		for i := range placeholders {
			placeholders[i] = "?"
		}

		insertQuery := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)",
			table.Name,
			columns,
			strings.Join(placeholders, ", "),
		)

		// Log the insert query and values
		log.Printf("Executing INSERT query: %s with values: %v", insertQuery, values)

		_, err := m.destDB.Exec(insertQuery, values...)
		if err != nil {
			return nil, fmt.Errorf("insert failed: %w", err)
		}

		// Store ID for related tables
		if id, ok := values[0].(int64); ok {
			ids = append(ids, int(id))
		} else {
			log.Printf("Warning: ID type conversion failed for value: %v", values[0])
		}
	}

	return ids, nil
}

func (m *MigrationService) migrateRelatedTable(table models.TableConfig, mainTableIds []int) error {
	exists, err := TableExists(m.destDB, table.Name)
	if err != nil {
		return err
	}

	if !exists {
		err = CreateTable(m.destDB, m.sourceDB, table.Name)
		if err != nil {
			return err
		}
	}

	// Build query for related data
	quotedColumns := make([]string, len(table.Columns))
	for i, col := range table.Columns {
		quotedColumns[i] = fmt.Sprintf("`%s`", col)
	}
	columns := strings.Join(quotedColumns, ", ")

	// Convert ids to strings for the IN clause
	idStrings := make([]string, len(mainTableIds))
	for i, id := range mainTableIds {
		idStrings[i] = fmt.Sprintf("%d", id)
	}

	query := fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` IN (%s)",
		columns,
		table.Name,
		table.ForeignKey,
		strings.Join(idStrings, ","),
	)

	// Log the query
	log.Printf("Executing related table SELECT query: %s", query)

	// Get related data
	rows, err := m.sourceDB.Query(query)
	if err != nil {
		return fmt.Errorf("related table select failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		values := make([]interface{}, len(table.Columns))
		valuePtrs := make([]interface{}, len(table.Columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err = rows.Scan(valuePtrs...)
		if err != nil {
			return fmt.Errorf("related table scan failed: %w", err)
		}

		placeholders := make([]string, len(table.Columns))
		for i := range placeholders {
			placeholders[i] = "?"
		}

		insertQuery := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)",
			table.Name,
			columns,
			strings.Join(placeholders, ", "),
		)

		// Log the insert query and values
		log.Printf("Executing related table INSERT query: %s with values: %v", insertQuery, values)

		_, err = m.destDB.Exec(insertQuery, values...)
		if err != nil {
			return fmt.Errorf("related table insert failed: %w", err)
		}
	}

	return nil
}

func createPlaceholders(n int) string {
	placeholders := make([]string, n)
	for i := range placeholders {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}
