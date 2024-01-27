build:
	go build -o gitlab-ci-component-linter ./cmd/main.go

clean:
	rm -f gitlab-ci-component-linter

run:
	go run ./cmd/main.go --workdir /home/sapher/projects/fleurimont/devops/ci-cd-catalog/nodejs
