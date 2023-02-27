module github.com/alexadastra/habit_bot

go 1.20

replace github.com/alexadastra/habit_bot/migrations => ./migrations

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/joho/godotenv v1.5.1
	github.com/pkg/errors v0.9.1
	github.com/xakep666/mongo-migrate v0.2.1
	go.mongodb.org/mongo-driver v1.11.1
	golang.org/x/sync v0.1.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/text v0.3.7 // indirect
)
