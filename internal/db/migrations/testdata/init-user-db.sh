#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-'EOSQL'
    DROP USER IF EXISTS test_user;
    CREATE USER test_user WITH PASSWORD 'test_password' SUPERUSER;
EOSQL
