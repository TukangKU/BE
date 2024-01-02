go test -cover ./features/users/services/ -coverprofile=cover.out && go tool cover -func=cover.out
go test -cover ./features/transaction/services/ -coverprofile=cover.out && go tool cover -func=cover.out
go test -cover ./features/skill/services/ -coverprofile=cover.out && go tool cover -func=cover.out
go test -cover ./features/jobs/services/ -coverprofile=cover.out && go tool cover -func=cover.out