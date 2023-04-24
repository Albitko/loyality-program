build-and-run:
	cd cmd/gophermart
	go build -o gophermart -buildvcs=false
	./gophermart -r localhost:8181
	cd ../accrual
	./accrual_darwin_amd64 -a "localhost:8181"

lint:
	golangci-lint run