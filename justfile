build:
  go build

ci:
  just build
  just fmt

docker-build:
  docker-compose build

docker-run:
  docker-compose up && docker-compose down

fmt:
  #!/usr/bin/env sh
  if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
    gofmt -d -s -l .
    exit 1
  fi
  echo "\033[92mgofmt Success\033[0m"

fix-fmt:
  gofmt -w -s .


run:
  go run main.go