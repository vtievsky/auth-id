run:  ## Запустить сервер
	bash -c 'set -a; . ./build/local/.env; set +a; go run cmd/auth-id/main.go'
