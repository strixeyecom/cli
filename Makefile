

# install and run source code security analysis
security:
	-echo "installing requirements"
	-curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s latest
	@./bin/gosec ./...
	@echo "[OK] Go security check was completed!"