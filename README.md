# go-auth-clean-arch

## Description

Golang Backend Application for management user, role, and autentication using Go, Echo, GORM, Goose, MySQL with Paseto and Seserion Versionatign

## Project Structure

```
go-aut-clean-arch/
├── config.yml
├── go.mod, go.sum
├── LICENSE
├── main.go
├── README.md
├── internal/
│   ├── cmd/
│   │   ├── migrate.go
│   │   ├── root.go
│   │   └── server.go
│   ├── controller/
│   │   ├── controller.go
│   │   └── user_controller.go
│   ├── helper/
│   │   ├── common_helper.go
│   │   ├── crypto_helper.go
│   │   └── helper.go
│   ├── middleware/
│   │   ├── db_middleware.go
│   │   └── middleware.go
│   ├── repository/
│   │   ├── repository.go
│   │   └── user_repository.go
│   ├── routes/
│   │   ├── routes.go
│   │   └── user_routes.go
│   ├── service/
│   │   ├── service.go
│   │   └── user_service.go
├── migrations/        # Goose migration files
├── package/
│   ├── external/
│   └── library/
│       ├── common.go
│       ├── config.go
│       ├── db.go
│       ├── echo.go
│       ├── library.go
│       └── logger.go
├── resource/
│   ├── constants/
│   │   ├── channel.go
│   │   └── default.go
│   ├── errors/
│   │   └── error.go
│   ├── model/
│   │   ├── common.go
│   │   ├── helper.go
│   │   ├── role.go
│   │   ├── user.go
│   │   └── validator.go
│   ├── rbac/
│   └── response/
│       └── common.go
```

## Installation

1. Clone repository:
   ```
   git clone https://github.com/petershaan12/go-auth-clean-arch
   cd go-auth-clean-arch
   ```
2. Install dependencies:
   ```
   go mod tidy
   ```
3. Copy file config.yml dari config.yml.example ke root project, lalu edit isinya sesuai kebutuhan (misal: koneksi DB, dsb).

4. Jalankan server:
   ```
   go run main.go serve
   ```
5. Buat migration baru:
   ```
   goose -dir ./migrations create your_sql_name sql
   ```
6. Jalankan migrasi:

   ```
   goose -dir ./migrations mysql "user:password@tcp(localhost:3306)/dbname" up
   ```

   or using migrate.go

   ```
   go run main.go migrate --direction=up
   ```

   or if you want to rollback

   ```
   go run main.go migrate --direction=down
   ```
