# tea-urls
A link shortening service

## Features
- Shorten URLs trough a simple web interface
- Redirect them to the original destination

## Stack
- [Go](https://go.dev/)
- [Gin](https://gin-gonic.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- Docker
- Docker Compose

## Running with Docker
1. Clone the repository:
```bash
git clone https://github.com/ZelenyMK/tea-urls.git
cd tea-urls
```
2. Create a .env file with such contents:
```bash
touch .env
```
```
REDIS_PASSWORD=<password>
POSTGRES_PASSWORD=<password>
```

3. Start the application with Docker compose:
```bash
docker compose up --build
```
4. Open the web interface in your browser
```
http://localhost:8080/
```

## Stopping the service
To stop the service run:
```bash
docker compose down
```

You can add ```-v``` flag to the end of the command to remove the database volume

## HTTP routes
| Method  | Route    | Description                   |
|---------|----------|-------------------------------|
| GET     | /        | Display the web interface     |
| POST    | /links   | Shorten the URL               |
| GET     | /:alias  | Redirects to the original URL |