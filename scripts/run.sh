mkdir -p build
rm build/* 2> /dev/null
go build -o build cmd/*
#. ./env
./build/main
