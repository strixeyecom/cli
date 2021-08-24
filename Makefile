
source_path=./cmd/strixeye/main.go
# install and run source code security analysis
security:
	-echo "installing requirements"
	-curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s latest
	@./bin/gosec ./...
	@echo "[OK] Go security check was completed!"


build:
	go build -o strixeye ${source_path}
	go install ./...