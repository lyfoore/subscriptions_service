# Setup
Create a .env file with the following variables:
```
DB_HOST=db
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=mydb
POSTGRES_PORT=5432
PGADMIN_DEFAULT_EMAIL=postgres@gmail.com
PGADMIN_DEFAULT_PASSWORD=postgres
```

# Run
Run the docker compose file:
```bash
docker compose up --build -d
```

# Swagger UI
Open the following URL in your browser:
```
http://localhost:8080/swagger/index.html
```