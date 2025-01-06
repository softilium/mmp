dotnet tool update --global dotnet-ef

rem dotnet ef migrations remove -f

dotnet ef migrations add MIG-250105-2357 --output-dir Data\Migrations

rem dotnet ef migrations list

