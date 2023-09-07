#!/bin/bash
echo "Executing custom entry point..."
# Additional commands or logic here
migrate \
  -path /migrations \
  -database "mysql://root:ZXCasdqwe123!@tcp(127.0.0.1:3306)/bubbme_db?loc=Asia%2FJakarta&multiStatements=true&parseTime=1" \
  up

# Wait for the first migration to complete before starting the second one
migrate \
  -path /migrations \
  -database "mysql://root:ZXCasdqwe123!@tcp(127.0.0.1:3308)/bubbme_db?loc=Asia%2FJakarta&multiStatements=true&parseTime=1" \
  up