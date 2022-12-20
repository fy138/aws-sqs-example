set GOOS=linux
set GOARCH=amd64

go build -o main sqstiggerlambda.go
build-lambda-zip.exe -output main.zip main
del main
set GOOS=windows