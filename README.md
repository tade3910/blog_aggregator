# blog_aggregator

## Setting up local database

### Installing postgress server
- brew install postgresql@15
- brew services start postgresql : Start services in the background
- psql postgres
- CREATE DATABASE blogator;

**NOTE** : You may need to export the postgress path
**NOTE**: In future iterations it may be easier to run a postgress docker image instead
### Install command line tool kits
- go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
- go install github.com/pressly/goose/v3/cmd/goose@latest

### Setting up tables
1. cd sql/schema
2. goose postgres postgres://username:@localhost:5432/blogator up