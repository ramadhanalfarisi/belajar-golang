#!/bin/bash
set -e
export PGPASSWORD=konohagakure;
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "toko_kocak" <<-EOSQL
  CREATE DATABASE toko_kocak;
  GRANT ALL PRIVILEGES ON DATABASE toko_kocak TO "postgres";
EOSQL