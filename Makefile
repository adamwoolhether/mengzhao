install:
	go install github.com/bokwoon95/wgo@latest
	go install github.com/a-h/templ/cmd/templ@latest
	npm install -D tailwindcss
	npx tailwindcss init
	npm install -D daisyui@latest

build:
	npx tailwindcss -i view/css/app.css -o public/styles.css
	templ generate view
	go build -o bin/mengzhao main.go

run:
	wgo -file=.go -file=.templ -file=.js -file=.css -xfile=_templ.go templ generate :: npx tailwindcss -i view/css/app.css -o public/styles.css :: go run main.go

up:
	go run cmd/migrate/main.go up

reset:
	go run cmd/reset/main.go

down:
	go run cmd/migrate/main.go down

drop:
	go run cmd/migrate/main.go up

migration:
	migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

gen:
	go run cmd/generate/main.go

seed:
	go run cmd/seed/main.go

tidy:
	go mod tidy
	go mod vendor
	#go mod download
