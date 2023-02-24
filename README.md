# habit-bot

## About
Habit-bot is a Telegram bot for tracking habits. Currently these are 'check-in' for some activity and 'gratitude' for daily being grateful practice. The bot uses mongoDB as a storage for user activities and runs inside Docker container.

## Getting started
### Prerequisites
To run the app locally inside the Docker, you will need only Docker or any other container software with the same API (e.g. Podman). You will also need golang v1.20 to build the app locally or run tests.
### Environment
To handle database secrets and bot token, it is recommended to create `values.env` file in the directory you want to run the app from.
These are the table of variables and their meaning:
|Variable|Meaning|
|:---:|:---:|
`BOT_TOKEN`|Token fetched from Bot Father
`MONGO_INITDB_ROOT_USERNAME`|MongoDB root username set while creating MongoDB instance
`MONGO_INITDB_ROOT_PASSWORD`|MongoDB root password set while creating MongoDB instance
`MONGODB_ADMIN_HOST`|MongoDB uri in the `mongodb://user:pass@localhost:27017/my-db` format
`MONGODB_ADMIN_DATABASE`|the title of the database where migrations will be rolled up
`MONGODB_USER`|App username for MongoDB
`MONGODB_PASSWORD`|App password for MongoDB
`MONGODB_DATABASE`|actually `MONGODB_ADMIN_DATABASE`
`MONGODB_HOST`|App host:port for MongoDB

### Commands
All the command to run the app are inside the Makefile instructions
#### TL;DR
To add and run everything from the scratch, run:
```bash
make dev-run
```
in the terminal.
#### Commands table
|Command|Meaning|
|:---:|:---:|
`build`|Build the app to the binary
`run`|Run the binary locally
`test`|Run the tests
`network-create/rm`|Create/remove the Docker network where containers will run
`volume-create/rm`|Create/remove the Docker volume for MongoDB persistency
`mongo-run/stop/rm`|Create/stop/remove MongoDB container
`migrations-up`|Rollup migrations for MongoDB
`app-build/run/stop/rm`|Build app Docker image/run/stop/rm docker container
`dev-run/stop/clean`|Create and run/stop/remove the entire environment (network, volume, db, app)