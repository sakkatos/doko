test:
	@DB_NAME=$(DB_NAME)_test \
	DB_HOST=localhost \
	DB_USER=$(DB_USER) \
	DB_PASSWORD=$(DB_PASSWORD) \
	go test ./...

migration:
	@echo "=== Run migrations and insert seed for main database ==="
	@cd migrations && DB_HOST=localhost go run *.go init;
	@cd migrations && DB_HOST=localhost go run *.go up;
	@cd migrations && DB_HOST=localhost go run *.go seed;
	@echo "=== Run migrations for test database ==="
	@cd migrations && DB_HOST=localhost DB_NAME=$(DB_NAME)_test go run *.go init;
	@cd migrations && DB_HOST=localhost DB_NAME=$(DB_NAME)_test go run *.go up;

truncate:
	@echo "=== Run truncate for main database ==="
	@cd migrations && DB_HOST=localhost go run *.go truncate;
	@echo "=== Run truncate for test database ==="
	@cd migrations && DB_HOST=localhost DB_NAME=$(DB_NAME)_test go run *.go truncate;

tcr:
	@make test || (git reset --hard; echo "----- TCR reverted -----"; exit 1)
	@git add .
	@git commit -am "tcring" | tee /dev/tty | grep -qE "nothing to commit$$" || echo

pull_rebase:
	@set -o pipefail; git pull --rebase | tee /dev/tty > ./limbo.local
	
limbo: tcr pull_rebase
	@grep -qE "up to date\.$$" ./limbo.local && git push || make limbo
