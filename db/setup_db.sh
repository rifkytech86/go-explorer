#!/bin/bash
sql_slave_user='CREATE USER "root"@"%" IDENTIFIED BY "ZXCasdqwe123!"; GRANT REPLICATION SLAVE ON *.* TO "root"@"%"; FLUSH PRIVILEGES; RESET master'
docker exec mysqldb sh -c "mysql -u root -pZXCasdqwe123! -e '$sql_slave_user'"
MS_STATUS=`docker exec mysqldb sh -c 'mysql -u root -pZXCasdqwe123! -e "SHOW MASTER STATUS"'`
CURRENT_LOG=`echo $MS_STATUS | awk '{print $6}'`
CURRENT_POS=`echo $MS_STATUS | awk '{print $7}'`
sql_set_master="CHANGE MASTER TO MASTER_HOST='mysqldb',MASTER_USER='root',MASTER_PASSWORD='ZXCasdqwe123!',MASTER_LOG_FILE='$CURRENT_LOG',MASTER_LOG_POS=$CURRENT_POS; STOP SLAVE; RESET SLAVE; RESET MASTER; START SLAVE;"
start_slave_cmd='mysql -u root -pZXCasdqwe123! -e "'
start_slave_cmd+="$sql_set_master"
start_slave_cmd+='"'
docker exec mysqldb-replica sh -c "$start_slave_cmd"
docker exec mysqldb-replica sh -c "mysql -u root -pZXCasdqwe123! -e 'SHOW SLAVE STATUS \G'"