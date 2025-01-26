set PGPASSWORD=0000
set PGPATH=D:\Program Files\PostgreSQL\16.4-5.1C\bin

rar x rsdata.rar

"%PGPATH%\psql.exe" -U postgres -c "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE datname = 'rsdata';"
"%PGPATH%\psql.exe" -U postgres -c "drop database if exists rsdata"
"%PGPATH%\psql.exe" -U postgres -c "create database rsdata"
"%PGPATH%\psql.exe" --dbname=rsdata -U postgres < "backup.sql"

del backup.sql
