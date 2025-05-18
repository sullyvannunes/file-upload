run:
	templ generate
	go build -o bin/app .
	bin/app

dump:
	pg_dump -f db/schema.sql -U postgres -d file-upload -h localhost
