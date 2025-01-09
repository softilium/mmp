rem c:\Utils\kiota\kiota generate --language CSharp -d http://localhost:5078/openapi/v1.json --namespace-name mmp.Client --output mmpClient
c:\Utils\kiota\kiota generate --clean-output --language "csharp" --openapi "http://localhost:5078/openapi/v1.json" --output "csharp" --namespace-name "mmp.Client"
c:\Utils\kiota\kiota generate --clean-output --language "Typescript" --openapi "http://localhost:5078/openapi/v1.json" --output "typescript" --namespace-name "mmp.Client"

rem dotnet tool install --global Microsoft.OpenApi.Kiota

rem dotnet kiota generate --clean-output --language "csharp" --openapi "http://localhost:5078/openapi/v1.json" --output "ApiClient" --namespace-name "mmp.ApiClient"

rem --class-name ""

rem npm install -g @azure/static-web-apps-cli