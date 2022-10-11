#!/bin/bash
set -e
export PGPASSWORD=konohagakure;
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "toko_kocak_test" <<-EOSQL
  CREATE DATABASE toko_kocak_test;
  GRANT ALL PRIVILEGES ON DATABASE toko_kocak_test TO "postgres";
EOSQL