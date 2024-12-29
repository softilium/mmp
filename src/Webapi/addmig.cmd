dotnet tool update --global dotnet-ef

rem dotnet ef migrations remove -f

dotnet ef migrations add MIG-241228-1939 --output-dir Data\Migrations

rem dotnet ef migrations list

