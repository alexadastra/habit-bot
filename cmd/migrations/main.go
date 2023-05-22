package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/alexadastra/habit_bot/migrations"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const envFile = "values.env" // TODO: set this as argument

func main() {
	if err := godotenv.Load(envFile); err != nil {
		panic(err)
	}
	if len(os.Args) == 1 {
		log.Fatal("Missing options: new, up or down")
	}
	option := os.Args[1]

	if err := mongoConnect(); err != nil {
		log.Fatal(err)
	}

	switch option {
	case "new":
		name, err := new()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("New migration created: %s\n", name)
	case "up":
		if err := up(); err != nil {
			log.Fatal(err)
		}
		log.Println("Successfully run all migrations!")
	case "down":
		if err := down(); err != nil {
			log.Fatal(err)
		}
		log.Println("Successfully rollback one migration!")
	case "down-all":
		if err := downAll(); err != nil {
			log.Fatal(err)
		}
		log.Println("Successfully rollback all migrations!")
	default:
		log.Fatalf("unknown command: %s", option)
	}
}

func mongoConnect() error {
	opt := options.Client().ApplyURI(os.Getenv("MONGODB_ADMIN_HOST"))
	client, err := mongo.NewClient(opt)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	db := client.Database(os.Getenv("MONGODB_ADMIN_DATABASE"))
	migrate.SetDatabase(db)

	return nil
}

func new() (string, error) {
	if len(os.Args) != 3 {
		log.Fatal("Should be: new description-of-migration")
	}
	fName := fmt.Sprintf("./migrations/%s_%s.go", time.Now().Format("20060102150405"), os.Args[2])
	from, err := os.Open("./migrations/template.go")
	if err != nil {
		return "", errors.Wrap(err, "failed to open 'template' file")
	}
	defer from.Close()

	to, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return "", errors.Wrap(err, "failed to copy from 'template' to destination")
	}
	return fName, nil
}

func up() error {
	return migrate.Up(migrate.AllAvailable)
}

func down() error {
	return migrate.Down(1)
}

func downAll() error {
	return migrate.Down(migrate.AllAvailable)
}
