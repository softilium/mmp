dotnet tool update --global dotnet-ef

rem dotnet ef migrations remove -f

dotnet ef migrations add MIG-250114-2304 --output-dir Data\Migrations

rem dotnet ef migrations list
