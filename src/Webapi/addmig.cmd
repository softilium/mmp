dotnet tool update --global dotnet-ef

rem dotnet ef migrations remove -f

dotnet ef migrations add MIG-250108-2359 --output-dir Data\Migrations

rem dotnet ef migrations list

