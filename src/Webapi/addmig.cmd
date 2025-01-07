dotnet tool update --global dotnet-ef

rem dotnet ef migrations remove -f

dotnet ef migrations add MIG-250107-1852 --output-dir Data\Migrations

rem dotnet ef migrations list

