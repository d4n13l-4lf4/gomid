COV_FILE=$(PWD)/build/coverage.out

coverage:
	sh $(PWD)/scripts/go-coverage-test.sh --check-coverage false
	go tool cover -html=$(COV_FILE)

lint: 
	find . -name \*.go ! -path "./vendor/*" -exec gofmt -w -l {} \;
	find . -name \*.go ! -path "./vendor/*" -exec goimports -w {} \;
