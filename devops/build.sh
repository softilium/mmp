pushd ../src/mmp-go
GOOG=linux go build -ldflags "-s -w"
popd

pushd ../src/vueclient-go
npm run build
popd
