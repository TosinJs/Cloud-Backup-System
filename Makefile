include app.env

migrate-up:	
	migrate -database ${DSN} -path internal/setup/database/migrations up

migrate-down:	
	migrate -database ${DSN} -path internal/setup/database/migrations down

run: 
	go run main.go