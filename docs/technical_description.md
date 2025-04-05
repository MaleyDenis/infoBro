# Personal Technical News Dashboard

## 1. Task Description
Development of a personal technical news dashboard that aggregates data from various sources (Telegram, RSS, websites, Reddit). In the first phase, the system will collect, store, and display news with filtering capabilities. In subsequent phases, AI analysis for classification and personalization is planned.

## 2. Technologies Used
- **Backend**: Go with Chi web framework (minimalist approach)
- **Database**: MongoDB (document-oriented NoSQL)
- **Message Queues**: Redis (Lists for queues)
- **Data Parsers**: colly (scraping), gofeed (RSS)
- **Frontend**: React + Tailwind CSS
- **Visualization**: Recharts for graphs
- **Deployment**: Docker + Docker Compose
- **Configuration**: YAML files
- **Modular Architecture**: Interfaces for component standardization

## 3. Architectural Solution
### Overall Scheme
```
DATA COLLECTION
  ├─ Telegram Connector
  ├─ RSS Connector
  ├─ Web Scraping Connector
  └─ Reddit Connector
      │
      ▼
  ┌─────────────────┐    ┌─────────────┐
  │ MongoDB (raw)   │◄───┤ Redis Queue │
  └─────────────────┘    └─────┬───────┘
                               │
PROCESSING (PHASE 1)           │
                               ▼
                       ┌───────────────┐
                       │ Basic Processor│
                       └───────┬───────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │ MongoDB (processed) │
                    └──────────┬──────────┘
                               │
PRESENTATION                   │
                               ▼
                       ┌───────────────┐
                       │  API (Chi)    │
                       └───────┬───────┘
                               │
                               ▼
                       ┌───────────────┐
                       │React Frontend │
                       └───────────────┘
```

### Data Flow Description
- **Data Collection**: Independent connectors collect data from different sources, save them to MongoDB (raw), and then place only the document identifiers (ObjectId) into the Redis queue.
- **Processing**: In the first phase, the processor extracts data from the queue, performs basic cleaning and normalization, then saves it to the processed news collection.
- **Presentation**: The Chi API provides access to the data, and the React frontend displays it in a convenient format with filtering capabilities.

## 4. Technical Details

### MongoDB Data Structure
**Raw News Collection**:
```json
{
  "_id": ObjectId,
  "source_type": String,
  "source_id": String,
  "title": String,
  "content": String,
  "url": String,
  "published_at": DateTime,
  "fetched_at": DateTime,
  "metadata": Object
}
```

**Processed News Collection**:
```json
{
  "_id": ObjectId,
  "raw_id": ObjectId,
  "title": String,
  "content": String,
  "source_type": String,
  "source_id": String,
  "source_name": String,
  "source_url": String,
  "url": String,
  "published_at": DateTime,
  "processed_at": DateTime
}
```

**Channel States Collection**:
```json
{
  "_id": ObjectId,
  "channel_id": String,
  "last_message_id": String,
  "last_update_time": DateTime,
  "processed_messages": Integer
}
```

### Redis Queues
- `news:queue` - main queue of raw news
- `news:processing` - items being processed
- `news:failed` - problematic items

### Objects in Redis Queues
In Redis queues, only MongoDB document identifiers are stored as strings:
```
"615a8b2c7d3a2f1a3c9b4d7e"  // Hex representation of MongoDB ObjectId
```

### Interfaces for Connector System

```go
// NewsConnector - interface for all news sources
type NewsConnector interface {
    GetNews(ctx context.Context) ([]RawNews, error)
}

// RawNews - structure for storing news in a standard format
type RawNews struct {
    SourceType  string
    SourceID    string
    SourceName  string
    SourceURL   string
    Title       string
    Content     string
    URL         string
    PublishedAt time.Time
    FetchedAt   time.Time
    Metadata    map[string]interface{}
}

// ChannelState - source state structure
type ChannelState struct {
    ChannelID         string
    LastMessageID     string
    LastUpdateTime    time.Time
    ProcessedMessages int
}

// ChannelStateRepository - interface for storing source states
type ChannelStateRepository interface {
    GetChannelState(ctx context.Context, channelID string) (*ChannelState, error)
    UpdateChannelState(ctx context.Context, state *ChannelState) error
}

// NewsStorage - interface for news storage
type NewsStorage interface {
    SaveRawNews(ctx context.Context, news []RawNews) ([]primitive.ObjectID, error)
}

// NewsQueue - interface for news queue
type NewsQueue interface {
    AddToQueue(ctx context.Context, newsIDs []primitive.ObjectID) error
}
```

### Connector Configuration (connectors.yaml)

```yaml
# connectors.yaml - Separate configuration file for connectors

# Telegram connector
telegram:
  enabled: true
  channels:
    - name: "Golang News"
      url: "https://t.me/golang_news"
    - name: "Rust Language"
      url: "https://t.me/rustlang"
    - name: "Python Insider"
      url: "https://t.me/python"
  credentials:
    api_id: "your_api_id"
    api_hash: "your_api_hash"

# RSS connector
rss:
  enabled: true
  feeds:
    - name: "Hacker News"
      url: "https://news.ycombinator.com/rss"
    - name: "The Verge"
      url: "https://www.theverge.com/rss/index.xml"
    - name: "DEV Community"
      url: "https://dev.to/feed"
  settings:
    timeout: 30s
    user_agent: "NewsAggregator/1.0"
```

### News Deduplication Mechanism
To avoid news duplication, the system:
1. Stores state for each source in MongoDB (last processed ID)
2. When launching a connector, checks the last processed ID
3. Retrieves only new messages after that ID
4. For new sources without history, loads only the latest N messages

### API Endpoints
**GET /api/news**
Parameters: source_type, source_id, query, from_date, to_date, page, page_size
Response:
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "615a8b2c7d3a2f1a3c9b4d7e",
        "title": "Go 1.21 Version Released",
        "content_preview": "The Go development team announced the release of a new version...",
        "source_type": "telegram",
        "source_id": "golang_news",
        "source_name": "Golang News",
        "source_url": "https://t.me/golang_news",
        "url": "https://t.me/golang_news/1234",
        "published_at": "2025-04-02T15:30:42Z",
        "processed_at": "2025-04-02T15:32:10Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_pages": 8,
      "total_items": 150
    }
  }
}
```

**GET /api/news/{id}**
Parameters: id (in URL)
Response:
```json
{
  "success": true,
  "data": {
    "id": "615a8b2c7d3a2f1a3c9b4d7e",
    "title": "Go 1.21 Version Released",
    "content": "Full news text...",
    "source_type": "telegram",
    "source_id": "golang_news",
    "source_name": "Golang News",
    "source_url": "https://t.me/golang_news",
    "url": "https://t.me/golang_news/1234",
    "published_at": "2025-04-02T15:30:42Z",
    "processed_at": "2025-04-02T15:32:10Z"
  }
}
```

**POST /api/connectors/run/{name}**
Manually run a specific connector
Response:
```json
{
  "success": true,
  "data": {
    "processed": 5,
    "connector": "telegram"
  }
}
```

**POST /api/connectors/run-all**
Run all active connectors
Response:
```json
{
  "success": true,
  "data": {
    "results": {
      "telegram": {"status": "success", "processed": 5},
      "rss": {"status": "success", "processed": 12},
      "reddit": {"status": "error", "message": "Auth failed"}
    }
  }
}
```

### Field Specifics
- `source_url` - URL of the news source overall (channel, feed, site)
- `url` - URL of the specific news item or post

## 5. Development Stages
- **Phase 1**: Basic data collection and display
- **Phase 2**: Integration with Claude API for analysis and classification
- **Phase 3**: Personalization and user preferences
- **Phase 4**: Analytical dashboards and visualization

This project provides an excellent opportunity to learn Go, MongoDB, Redis, and work with various APIs, while also offering room for gradual development and functionality enhancement.