# rest-service
## start db
```
docker-compose up
```
## start service
```
go run ./cmd/service
```
## start worker
```
go run ./cmd/worker
```
## start tester
```
go run tester_service.go
```
## benchmark
```
go test ./... -bench=. -benchtime=10s
```
I had **10000** requests with **100** ids per request in **~1sec**.

My local setup: i5-9400F 2.90GHz, 16GB

Usage was: 70% CPU, 7GB RAM

You can configure this benchmark in ./app/object_processor_test.go
