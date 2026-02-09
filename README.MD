### Environment
[Backend Environment Example](backend/.env.example)

### Commands
Start Database
```
docker compose up
```
Down Migration
``` 
go run backend/cmd/db/main.go -down
```
Up Migration
```
go run backend/cmd/db/main.go
```
Start CronJob
```
go run backend/cmd/cron/main.go
```
Start Development Server
```
go run backend/main.go
```

Start Production Server
``` 
go run backend/main.go -release
```
