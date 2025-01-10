pushd ..

pushd src\vueclient
call npm run build
popd

swa deploy --env production

popd

