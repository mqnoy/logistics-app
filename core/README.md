# Core
## backend 

### Get Started
1. build image `make build-image`
1. Run with docker compose `make run`
1. Create database `docker exec -it core-logistics_db-1 mysql -uroot -p12345678 -e "CREATE DATABASE IF NOT EXISTS logistics_app_db;"`