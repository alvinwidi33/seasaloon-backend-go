# 💈 Seasaloon Backend
 
REST API backend for the Seasaloon application, built with **Go** + **Gin** + **PostgreSQL**, containerized with Docker.
 
---
 
## 🚀 Tech Stack
 
- **Language**: Go
- **Framework**: Gin
- **Database**: PostgreSQL 17
- **Container**: Docker + Docker Compose
 
---
 
## 🐳 Running with Docker
 
```bash
# Start all services
docker compose up -d
 
# Stop all services
docker compose down
 
# View logs
docker logs backend-app
```
 
> The app runs on port **8081** by default.
 
---
 
## 📡 API Endpoints
 
Base URL: `http://apin-devops.my.id:8081`
 
### 🔐 Auth
 
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/register` | Register new user |
| POST | `/api/login` | Login |
| GET | `/api/activate` | Activate user account |
 
---
 
### 👤 User
 
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/me` | Get current logged-in user |
| PATCH | `/api/user-profile/:id` | Update user profile |
 
---
 
### 💈 Saloon
 
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/saloon/:id` | Get saloon by ID |
| POST | `/api/saloon` | Create new saloon |
| PUT | `/api/saloon/:id` | Update saloon |
| PATCH | `/api/saloon/:id/delete` | Soft delete saloon |
 
---
 
### 📅 Reservation
 
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/reservation` | Get all reservations |
| GET | `/api/reservation/:id` | Get reservations by customer ID |
| POST | `/api/reservation` | Create new reservation |
| PATCH | `/api/reservation/:id/cancel` | Cancel a reservation |
| PATCH | `/api/reservation/:id/done` | Mark reservation as done |
 
---
 
### 👥 Customers & Users
 
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users` | Get all customers |
| PATCH | `/api/users/:id/member` | Set customer membership |
| GET | `/api/saloon/customer` | Get all saloon customers |
 
---
 
### 🛡️ Admin
 
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin` | Register admin |
| GET | `/api/admin/saloon` | Get all saloons (admin) |
 
---
 
### 🩺 Doctor
 
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/doctor` | Register doctor |
| GET | `/api/doctor` | Get all doctors |
 
---
 
## 🔄 CI/CD
 
This project uses **GitHub Actions** with `appleboy/ssh-action` to auto-deploy on push:
 
1. SSH into VPS
2. Pull latest Docker image
3. Restart containers via `docker compose up -d --remove-orphans`
---
 <img width="1906" height="1009" alt="image" src="https://github.com/user-attachments/assets/6bd5e7a7-1563-4adb-9325-5b312f0e6ece" />

## 📁 Project Structure
 
```
seasaloon-backend-go/
├── controllers/       # Route handlers
├── database/
│   └── connections/   # DB connection & migration
├── main.go            # Entry point
├── docker-compose.yml
└── .env
```
