dotnet tool update --global dotnet-ef

rem dotnet ef migrations remove -f

dotnet ef migrations add MIG-250103-1341 --output-dir Data\Migrations

rem dotnet ef migrations list

