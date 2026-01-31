module github.com/tukhvatullinsm/golang-project

go 1.24.5

replace github.com/tukhvatullinsm/golang-project/internal/storage => ./internal/storage

replace github.com/tukhvatullinsm/golang-project/internal/handlers => ./internal/handlers

require (
	github.com/tukhvatullinsm/golang-project/internal/handlers v0.0.0-00010101000000-000000000000
	github.com/tukhvatullinsm/golang-project/internal/storage v0.0.0-00010101000000-000000000000
)
