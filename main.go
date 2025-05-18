package main

import (
	"context"
	"fmt"
	pgrepo "forms/pg"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	db, err := pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5432/file-upload?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(ctx); err != nil {
		panic(err)
	}

	defer db.Close(ctx)

	q := pgrepo.New(db)

	app := NewApp(q)

	fmt.Println("Running on http://localhost:3030")
	if err := http.ListenAndServe(":3030", app); err != nil {
		panic(err)
	}
}
