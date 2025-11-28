#!/bin/sh
set -e


echo "run db migration"
/usr/src/app/migrate -path /usr/src/app/migration -database "$DB_SOURCE" -verbose up


echo "start the app"
exec "$@"