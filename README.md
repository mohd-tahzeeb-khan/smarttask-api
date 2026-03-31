\# 🧠 SmartTask AI API



> \*\*Production-grade, AI-powered Task \& Productivity API\*\* built with Go, Gin, GORM, SQLite and OpenAI.



!\[Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat\&logo=go)

!\[Gin](https://img.shields.io/badge/Gin-Framework-blue?style=flat)

!\[License](https://img.shields.io/badge/License-MIT-yellow.svg)

!\[Deploy](https://img.shields.io/badge/Deployed-Render-46E3B7?style=flat\&logo=render)



🌍 \*\*Live API:\*\* https://smarttask-api.onrender.com  

📖 \*\*Swagger Docs:\*\* https://smarttask-api.onrender.com/swagger/index.html  

💻 \*\*GitHub:\*\* https://github.com/mohd-tahzeeb-khan/smarttask-api



\---



\## ✨ Features



| Feature | Details |

|---|---|

| 🔐 JWT Auth | Signup, Login, bcrypt password hashing |

| 🤖 AI Analysis | Priority suggestion + time estimation via OpenAI (smart mock fallback) |

| ✅ Task CRUD | Full create/read/update/delete with ownership checks |

| 🔍 Advanced Filtering | Filter by priority, status, deadline (overdue/today/week) |

| 📄 Pagination | Page + limit + sort + order on all list endpoints |

| 📊 Analytics | Productivity score, weekly insights, priority breakdown |

| 🚦 Rate Limiting | 100 req/min per IP |

| 🛡️ Middleware | Logger, error handler, auth, panic recovery |

| 📖 Swagger UI | Auto-generated interactive API documentation |

| 🐳 Docker | Single-command deploy with Docker |

| ☁️ Cloud Ready | Deployed on Render free tier |



\---



\## 🗂️ Project Structure

```

smarttask/

├── cmd/

│   └── main.go                       # Entry point + DI wiring

├── internal/

│   ├── config/

│   │   └── config.go                 # Env var loader

│   ├── models/

│   │   ├── models.go                 # GORM models + DTOs

│   │   └── database.go               # DB init + migrations

│   ├── repository/

│   │   ├── user\_repository.go        # User DB operations

│   │   └── task\_repository.go        # Task DB + analytics queries

│   ├── services/

│   │   ├── auth\_service.go           # JWT + bcrypt logic

│   │   └── task\_service.go           # Business logic + AI calls

│   ├── controllers/

│   │   ├── auth\_controller.go        # Auth HTTP handlers

│   │   └── task\_controller.go        # Task + AI + analytics handlers

│   ├── middleware/

│   │   ├── auth.go                   # JWT middleware

│   │   ├── logger.go                 # Colorized request logger

│   │   └── ratelimit.go              # In-memory rate limiter

│   └── routes/

│       └── routes.go                 # All route definitions

├── pkg/

│   └── ai/

│       └── service.go                # OpenAI + smart mock fallback

├── docs/                             # Auto-generated Swagger docs

├── .env.example

├── Dockerfile

├── docker-compose.yml

├── render.yaml

└── README.md

```



\---



\## 🚀 Quick Start



\### 1. Clone \& Configure

```bash

git clone https://github.com/mohd-tahzeeb-khan/smarttask-api.git

cd smarttask-api

cp .env.example .env

\# Optionally add your OPENAI\_API\_KEY in .env

```



\### 2. Run Locally

```bash

go mod tidy

go run ./cmd/main.go

```



Server runs at \*\*http://localhost:8080\*\*



\### 3. Docker

```bash

docker-compose up --build

```



\---



\## 📡 API Endpoints



\### Auth

| Method | Endpoint | Auth | Description |

|---|---|---|---|

| POST | `/api/v1/auth/signup` | ❌ | Register new user |

| POST | `/api/v1/auth/login` | ❌ | Login, returns JWT |

| GET | `/api/v1/auth/me` | ✅ | Current user profile |



\### Tasks

| Method | Endpoint | Auth | Description |

|---|---|---|---|

| POST | `/api/v1/tasks` | ✅ | Create task (auto AI-analyzed) |

| GET | `/api/v1/tasks` | ✅ | List with filters + pagination |

| GET | `/api/v1/tasks/:id` | ✅ | Get single task |

| PUT | `/api/v1/tasks/:id` | ✅ | Update task |

| DELETE | `/api/v1/tasks/:id` | ✅ | Delete task |



\### AI

| Method | Endpoint | Auth | Description |

|---|---|---|---|

| POST | `/api/v1/ai/analyze-task` | ✅ | Analyze task → priority + ETA |



\### Analytics

| Method | Endpoint | Auth | Description |

|---|---|---|---|

| GET | `/api/v1/analytics` | ✅ | Full productivity dashboard |



\### System

| Method | Endpoint | Auth | Description |

|---|---|---|---|

| GET | `/health` | ❌ | Health check |



\---



\## 🔍 Filter \& Query Examples

```bash

\# Filter by priority

GET /api/v1/tasks?priority=high



\# Filter by status

GET /api/v1/tasks?status=pending



\# Overdue tasks

GET /api/v1/tasks?deadline=overdue



\# This week sorted by deadline

GET /api/v1/tasks?deadline=week\&sort=deadline\&order=asc



\# Pagination

GET /api/v1/tasks?page=2\&limit=5



\# Combined

GET /api/v1/tasks?priority=high\&status=pending\&sort=deadline\&order=asc

```



\---



\## 🤖 AI Analysis Example



\*\*Request:\*\*

```json

POST /api/v1/ai/analyze-task

{

&#x20; "title": "Prepare MBA presentation on sales trends",

&#x20; "description": "Deck for the board meeting covering Q3-Q4 results"

}

```



\*\*Response:\*\*

```json

{

&#x20; "success": true,

&#x20; "message": "AI analysis complete",

&#x20; "data": {

&#x20;   "priority": "high",

&#x20;   "estimated\_time\_hours": 3.0,

&#x20;   "reasoning": "Detected high-urgency keyword 'presentation'",

&#x20;   "confidence": 0.88

&#x20; }

}

```



\---



\## 📊 Analytics Response Example

```json

{

&#x20; "success": true,

&#x20; "data": {

&#x20;   "total\_tasks": 24,

&#x20;   "completed\_tasks": 18,

&#x20;   "pending\_tasks": 4,

&#x20;   "in\_progress\_tasks": 2,

&#x20;   "productivity\_score": 72.5,

&#x20;   "overdue\_tasks": 1,

&#x20;   "priority\_breakdown": {

&#x20;     "high": 8,

&#x20;     "medium": 12,

&#x20;     "low": 4

&#x20;   },

&#x20;   "weekly\_insights": \[

&#x20;     { "week": "Mar 04", "completed": 5, "created": 7 },

&#x20;     { "week": "Mar 11", "completed": 6, "created": 5 },

&#x20;     { "week": "Mar 18", "completed": 4, "created": 6 },

&#x20;     { "week": "Mar 25", "completed": 3, "created": 6 }

&#x20;   ]

&#x20; }

}

```



\---



\## 🔧 Environment Variables



| Variable | Default | Description |

|---|---|---|

| `PORT` | `8080` | Server port |

| `GIN\_MODE` | `debug` | `debug` or `release` |

| `JWT\_SECRET` | — | Required — strong random string |

| `JWT\_EXPIRY\_HOURS` | `24` | Token lifetime in hours |

| `OPENAI\_API\_KEY` | — | Optional — mock AI used if empty |

| `OPENAI\_MODEL` | `gpt-3.5-turbo` | OpenAI model |

| `DB\_PATH` | `./smarttask.db` | SQLite file path |

| `RATE\_LIMIT\_REQUESTS` | `100` | Max requests per window |

| `RATE\_LIMIT\_PERIOD` | `1m` | Rate limit window |



\---



\## 🏗️ Architecture

```

HTTP Request

&#x20;    │

&#x20;    ▼

\[Rate Limiter] → \[Logger] → \[Auth Middleware]

&#x20;    │

&#x20;    ▼

\[Controller]      ← validates request DTOs

&#x20;    │

&#x20;    ▼

\[Service]         ← business logic + AI calls

&#x20;    │

&#x20;    ▼

\[Repository]      ← GORM queries

&#x20;    │

&#x20;    ▼

\[SQLite DB]

```



\---



\## 🛠️ Built With



\- \[Gin](https://github.com/gin-gonic/gin) — HTTP framework

\- \[GORM](https://gorm.io) — ORM

\- \[go-openai](https://github.com/sashabaranov/go-openai) — OpenAI client

\- \[golang-jwt](https://github.com/golang-jwt/jwt) — JWT

\- \[godotenv](https://github.com/joho/godotenv) — Env loading

\- \[swaggo](https://github.com/swaggo/swag) — Swagger docs



\---



\## 📝 License



MIT — free to use, fork, and ship.



\---



\## 👤 Author



\*\*Mohd Tahzeeb Khan\*\*  

\[GitHub](https://github.com/mohd-tahzeeb-khan)

