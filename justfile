# Generate mocks
mock:
    GOFLAGS="-buildvcs=false" go run github.com/vektra/mockery/v2@v2.53.5 --config /dev/null --all --recursive --inpackage --case underscore --with-expecter=true

# Run tests with coverage
test: mock
    go test -cover -bench=. -benchmem -race ./... -coverprofile=coverage.out

# Build zeshion binary to GOPATH/bin
build version="dev":
    go build -buildvcs=false -ldflags "-X 'main.version={{version}}'" -o `go env GOPATH`/bin/zeshion

# Generate man page
man: build
    mkdir -p share/man/man1
    zeshion man > share/man/man1/zeshion.1
