install:
	go install github.com/a-h/templ/cmd/templ@latest
	npm install -D tailwindcss
	npx tailwindcss init
	npm install -D daisyui@latest


build:
	npx tailwindcss -i view/css/app.css -o public/styles.css
	templ generate view
	go build -o bin/mengzhao main.go

run: build
	./bin/mengzhao

up:
	go run cmd/migrate/main.go up

down:
	go run cmd/migrate/main.go down

drop:
	go run cmd/drop/main.go up

migrate:
	migrate create -ext sql -dir cmd/migrate/migrates $(filter-out $@,$(MAKECMDGOALS))

gen:
	go run cmd/generate/main.go

seed:
	go run cmd/seed/main.go

tidy:
	go mod tidy
	go mod vendor
	#go mod download
