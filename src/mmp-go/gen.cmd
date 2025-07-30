pushd ..\..\..\elorm-gen
go build
popd
..\..\..\elorm-gen\elorm-gen.exe mmp-go.schema.json dbcontext.go

