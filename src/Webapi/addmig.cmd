dotnet tool update --global dotnet-ef

rem dotnet ef migrations remove -f

dotnet ef migrations add MIG-250124-2053 --output-dir Data\Migrations

rem dotnet ef migrations list
