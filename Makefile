.PHONY: coverage

coverage:
	go test -v -tags=integration -coverprofile=cover.out
	go tool cover -func=cover.out