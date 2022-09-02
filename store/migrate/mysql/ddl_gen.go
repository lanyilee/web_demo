package mysql

import (
	"database/sql"
)

var migrations = []struct {
	name string
	stmt string
}{
	{
		name: "create-table-local-user",
		stmt: createTableLocalUser,
	},
	{
		name: "create-default-local-user",
		stmt: createDefaultLocalUser,
	},

}

// Migrate performs the database migration. If the migration fails
// and error is returned.
func Migrate(db *sql.DB) error {
	if err := createTable(db); err != nil {
		return err
	}
	completed, err := selectCompleted(db)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, migration := range migrations {
		if _, ok := completed[migration.name]; ok {

			continue
		}

		if _, err := db.Exec(migration.stmt); err != nil {
			return err
		}
		if err := insertMigration(db, migration.name); err != nil {
			return err
		}

	}
	return nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(migrationTableCreate)
	return err
}

func insertMigration(db *sql.DB, name string) error {
	_, err := db.Exec(migrationInsert, name)
	return err
}

func selectCompleted(db *sql.DB) (map[string]struct{}, error) {
	migrations := map[string]struct{}{}
	rows, err := db.Query(migrationSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations[name] = struct{}{}
	}
	return migrations, nil
}

//
// migration table ddl and sql
//

var migrationTableCreate = `
CREATE TABLE IF NOT EXISTS migrations (
 name VARCHAR(255)
,UNIQUE(name)
)
`

var migrationInsert = `
INSERT INTO migrations (name) VALUES (?)
`

var migrationSelect = `
SELECT name FROM migrations
`

//
// 0001_create_table_local_user.sql
//

var createTableLocalUser = `
CREATE TABLE IF NOT EXISTS local_user (
	id VARCHAR(40) NOT NULL,
	parent_id VARCHAR(40),
	username VARCHAR(255),
	password VARCHAR(255),
	name VARCHAR(255),
	role VARCHAR(255),
	level VARCHAR(255),
	phone VARCHAR(20),
	email VARCHAR(255),
	created DateTime,
	updated DateTime,
	PRIMARY KEY (id)
);
`

var createDefaultLocalUser = `
INSERT INTO local_user (id, username, password, name,role, level, created, updated) VALUES ('403519a6-3d89-419e-b424-8b238962025b', 'admin', 'ce1c1cdc2fac8e1167f22cd4bd88d324', '管理员','admin', '0', '2020-01-01 00:00:00', '2020-01-01 00:00:00');
`
