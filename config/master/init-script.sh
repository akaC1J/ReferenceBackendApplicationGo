#!/bin/bash

psql -U $POSTGRES_USER --dbname $POSTGRES_DB <<-EOSQL
    CREATE ROLE $POSTGRES_REPLICA_USER WITH REPLICATION LOGIN PASSWORD '$POSTGRES_REPLICA_PASSWORD';
    GRANT SELECT ON ALL TABLES IN SCHEMA public TO $POSTGRES_REPLICA_USER;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO $POSTGRES_REPLICA_USER;
EOSQL
