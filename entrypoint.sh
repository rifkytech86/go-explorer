#!/bin/sh


# Wait for the MySQL container to be ready
/app/wait-for-it.sh mysql:3306 -t 60

# Wait for the MySQL container to be ready (you might need to customize this)
until mysql -h mysql -u root -pZXCasdqwe123! -e ";" ; do
    echo "MySQL is not yet available. Sleeping..."
    sleep 1
done

echo "Running migrations..."
migrate -database "mysql://root:ZXCasdqwe123!@tcp(bubbme-backend-mysqldb-1:3306)/bubbme" -path db_migrations up
echo "Migrations completed!"

# Start your application
exec /app/bubbme-backend "$@"