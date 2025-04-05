# 📰 InfoBro

Personal Technical News Dashboard that aggregates data from various sources (Telegram, RSS, websites, Reddit).

## 🔍 Overview

InfoBro is a Go-based application that collects, stores, and displays tech news from multiple sources. 
The system uses MongoDB for storage, Redis for message queues, and provides a React frontend for viewing the news.

For detailed technical specifications, see [Technical Description](docs/technical_description.md).

## ✨ Features

- 🔄 Multi-source news aggregation (Telegram, RSS, Reddit, Web scraping)
- ⚙️ Configurable connectors for each source type
- 🧹 Efficient news deduplication mechanism
- 🌐 REST API with filtering and pagination
- ⚛️ Modern React frontend with Tailwind CSS
- 📱 Responsive UI that works on mobile and desktop
- 🐳 Docker-based deployment

## 📂 Project Structure

```
├── cmd/
│   └── server/               # Main application entry point
├── config/                   # Configuration files
├── docs/                     # Documentation
├── internal/
│   ├── api/                  # API handlers
│   ├── config/               # Configuration loading
│   ├── connectors/           # News source connectors
│   │   ├── reddit/           # Reddit-specific connector
│   │   ├── telegram/         # Telegram-specific connector (coming soon)
│   │   └── rss/              # RSS-specific connector (coming soon)
│   ├── models/               # Common data models
│   ├── queue/                # Message queue implementation
│   └── storage/              # Database storage implementation
├── scripts/                  # Helper scripts
└── web/                      # React frontend
    ├── public/               # Static files
    └── src/
        ├── components/       # Reusable UI components
        ├── hooks/            # Custom React hooks
        ├── pages/            # Page components
        ├── services/         # API service layer
        ├── styles/           # CSS and Tailwind styles
        └── utils/            # Utility functions
```

## 🛠️ Prerequisites

- 🔹 Go 1.21+
- 🔹 MongoDB
- 🔹 Redis
- 🔹 Node.js 18+ and npm (for frontend development)
- 🔹 Docker and Docker Compose (optional, for containerized setup)

## 🚀 Getting Started

### Option 1: Quick Start (All-in-one)

Run everything with a single command:
```
make run-all
```

This will start MongoDB, Redis, the backend API, and the React frontend.

### Option 2: Manual Setup

1. Clone the repository:
   ```
   git clone https://github.com/MaleyDenis/infoBro.git
   cd infoBro
   ```

2. Install Go dependencies:
   ```
   go mod download
   ```

3. Configure the application:
   - Edit `config/connectors.yaml` with your source configurations
   - For Reddit, obtain API credentials from https://www.reddit.com/prefs/apps

4. Start MongoDB and Redis:
   ```
   # Using Docker
   docker run -d -p 27017:27017 --name mongodb mongo
   docker run -d -p 6379:6379 --name redis redis
   ```

5. Build and run the backend server:
   ```
   go build -o infobro ./cmd/server
   ./infobro
   ```

6. Run the frontend (in a separate terminal):
   ```
   cd web
   npm install
   npm start
   ```

7. Access the application:
   - Frontend: http://localhost:3000
   - API: http://localhost:8080/api/news

### Option 3: Docker Compose 🐳

```
docker-compose up -d
```

This will start all components in containers:
- MongoDB: localhost:27017
- Redis: localhost:6379
- Backend API: localhost:8080
- Frontend: localhost:3000

## 💻 Frontend Features

The React frontend provides a modern, responsive user interface:

- 📊 **Homepage Dashboard**: View latest news with filtering and statistics
- 📄 **News Detail View**: Read full articles with source information
- 🔌 **Sources Management**: View and refresh individual news sources
- ⚙️ **Settings Page**: Customize application behavior

The UI is built with:
- 📘 TypeScript for type safety
- 🎨 Tailwind CSS for styling
- 🧭 React Router for navigation
- 🔄 React Query for data fetching and state management
- 📈 Recharts for data visualization

See [UI Preview](docs/ui_preview.md) for screenshots of the interface.

## 🏗️ Development

The project is organized following clean architecture principles with clear separation of concerns:

- **Backend**
  - 🔌 **Connectors**: Each news source has its own connector implementation
  - 💾 **Storage**: MongoDB is used for storing news items and source states
  - 📬 **Queue**: Redis is used for message queueing
  - 🌐 **API**: Chi router provides REST endpoints

- **Frontend**
  - 🧩 **Components**: Reusable UI building blocks
  - 📱 **Pages**: Main application views
  - 🔄 **Services**: API integration layer
  - 🪝 **Hooks**: Custom React hooks for data fetching

## 🌐 API Endpoints

- `GET /api/news` - Get news list with filtering and pagination
- `GET /api/news/{id}` - Get a specific news item
- `POST /api/connectors/run/{name}` - Run a specific connector
- `POST /api/connectors/run-all` - Run all enabled connectors

## 🛠️ Useful Commands

```
# Backend
make build          # Build the backend
make test           # Run tests
make run            # Run the backend

# Frontend
make frontend-install  # Install frontend dependencies
make frontend-start    # Start the frontend dev server
make frontend-build    # Build the frontend for production

# Docker
make docker-build      # Build Docker images
make docker-compose-up # Start all services with Docker Compose
```

## 🔮 Future Plans

- 🔌 Add more connectors (Twitter, HackerNews, etc.)
- 🤖 Implement AI analysis for classification
- 🔍 Add personalization features
- 📊 Create analytical dashboards

## 📝 License

MIT