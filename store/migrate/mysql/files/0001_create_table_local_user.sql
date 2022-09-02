-- name: create-table-local-user

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

-- name: create-default-local-user
INSERT INTO local_user (id, username, password, name,role, level, created, updated) VALUES ('403519a6-3d89-419e-b424-8b238962025b', 'admin', 'ce1c1cdc2fac8e1167f22cd4bd88d324', '管理员','admin', '0', '2020-01-01 00:00:00', '2020-01-01 00:00:00');