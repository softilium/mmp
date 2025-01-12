dotnet tool update --global dotnet-ef

rem dotnet ef migrations remove -f

dotnet ef migrations add MIG-250112-2112 --output-dir Data\Migrations

rem dotnet ef migrations list
