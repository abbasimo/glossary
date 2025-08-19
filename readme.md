# Go Glossary — Step 1


## Requirements
- Go 1.22+
- PostgreSQL 13+


## Setup
1. Create DB and user (example):
```sql
create database glossary_db;
create user glossary_user with password 'glossary_pass';
grant all privileges on database glossary_db to glossary_user;
```


2. Copy env:
```bash
cp .env.example .env
```


3. Run the server:
```bash
go mod tidy
go run ./cmd/server
```


4. Verify:
- Open http://localhost:8080/healthz → should return `ok`.
- The server applies `db/schema.sql` automatically on start.