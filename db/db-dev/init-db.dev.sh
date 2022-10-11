#!/bin/bash
set -e
export PGPASSWORD=konohagakure;
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "toko_kocak_dev" <<-EOSQL
  CREATE DATABASE toko_kocak_dev;
  GRANT ALL PRIVILEGES ON DATABASE toko_kocak_dev TO "postgres";
EOSQL