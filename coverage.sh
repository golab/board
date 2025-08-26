dir="$1"
cd $dir
go test ./... -coverprofile=cover.out
go tool cover -html=cover.out
rm cover.out
cd -
