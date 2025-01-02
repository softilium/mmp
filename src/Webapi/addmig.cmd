dotnet tool update --global dotnet-ef

rem dotnet ef migrations remove -f

dotnet ef migrations add MIG-250102-1919 --output-dir Data\Migrations

rem dotnet ef migrations list

