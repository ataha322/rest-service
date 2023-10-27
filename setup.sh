#!/bin/bash

sudo -i -u postgres psql <<EOF
CREATE DATABASE db_people;
CREATE USER my_app WITH PASSWORD 'pwd';
GRANT ALL PRIVILEGES ON DATABASE db_people TO my_app;
\q
EOF
