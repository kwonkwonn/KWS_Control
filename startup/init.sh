#!/bin/bash

echo "Creating databases..."

mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "
CREATE DATABASE IF NOT EXISTS \`$DB_NAME\`;
CREATE DATABASE IF NOT EXISTS \`$GUAC_DB_NAME\`;
CREATE DATABASE IF NOT EXISTS core_base;
GRANT ALL PRIVILEGES ON \`$DB_NAME\`.* TO 'root'@'%';
GRANT ALL PRIVILEGES ON \`$GUAC_DB_NAME\`.* TO 'root'@'%';
GRANT ALL PRIVILEGES ON core_base.* TO 'root'@'%';
FLUSH PRIVILEGES;
"

mysql -u root -p"$MYSQL_ROOT_PASSWORD" core_base -e "
CREATE TABLE IF NOT EXISTS subnet (
    id INT PRIMARY KEY,
    last_subnet VARCHAR(64) NOT NULL
);
INSERT IGNORE INTO subnet (id, last_subnet) VALUES (1, '10.0.0.1');

CREATE TABLE IF NOT EXISTS inst_info (
    uuid VARCHAR(128) PRIMARY KEY,
    inst_ip VARCHAR(64),
    guac_pass VARCHAR(256),
    inst_mem INT,
    inst_vcpu INT,
    inst_disk INT
);

CREATE TABLE IF NOT EXISTS inst_loc (
    uuid VARCHAR(128) PRIMARY KEY,
    core INT
);
"

echo "Databases created successfully!"
