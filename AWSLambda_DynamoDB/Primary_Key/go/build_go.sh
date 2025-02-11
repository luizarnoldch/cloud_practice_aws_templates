GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc
zip go_function.zip bootstrap
rm -rf bootstrap
