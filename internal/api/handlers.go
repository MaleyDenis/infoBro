package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dzianismalei/infoBro/internal/connectors"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// API handles HTTP requests for the news dashboard
type API struct {
	connectorService *connectors.ConnectorService
	newsStorage      NewsStorage
}

// NewAPI creates a new API handler
func NewAPI(connectorService *connectors.ConnectorService, newsStorage NewsStorage) *API {
	return &API{
		connectorService: connectorService,
		newsStorage:      newsStorage,
	}
}

// NewsStorage interface for retrieving processed news
type NewsStorage interface {
	GetNewsList(filters map[string]interface{}, page, pageSize int) (*NewsListResult, error)
	GetNewsById(id string) (*NewsItem, error)
}

// NewsItem represents a processed news item for API responses
type NewsItem struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	ContentPreview string    `json:"content_preview,omitempty"`
	Content       string    `json:"content,omitempty"`
	SourceType    string    `json:"source_type"`
	SourceID      string    `json:"source_id"`
	SourceName    string    `json:"source_name"`
	SourceURL     string    `json:"source_url"`
	URL           string    `json:"url"`
	PublishedAt   time.Time `json:"published_at"`
	ProcessedAt   time.Time `json:"processed_at"`
}

// NewsListResult represents paginated news results
type NewsListResult struct {
	Items      []NewsItem  `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// RegisterRoutes registers all API routes
func (a *API) RegisterRoutes(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		// News endpoints
		r.Get("/news", a.GetNewsList)
		r.Get("/news/{id}", a.GetNewsById)
		
		// Connector endpoints
		r.Post("/connectors/run/{name}", a.RunConnector)
		r.Post("/connectors/run-all", a.RunAllConnectors)
	})
}

// GetNewsList handles requests for filtered news lists
func (a *API) GetNewsList(w http.ResponseWriter, r *http.Request) {
	// Parse filters from query parameters
	filters := make(map[string]interface{})
	
	if sourceType := r.URL.Query().Get("source_type"); sourceType != "" {
		filters["source_type"] = sourceType
	}
	
	if sourceID := r.URL.Query().Get("source_id"); sourceID != "" {
		filters["source_id"] = sourceID
	}
	
	if query := r.URL.Query().Get("query"); query != "" {
		filters["query"] = query
	}
	
	if fromDate := r.URL.Query().Get("from_date"); fromDate != "" {
		if date, err := time.Parse(time.RFC3339, fromDate); err == nil {
			filters["from_date"] = date
		}
	}
	
	if toDate := r.URL.Query().Get("to_date"); toDate != "" {
		if date, err := time.Parse(time.RFC3339, toDate); err == nil {
			filters["to_date"] = date
		}
	}
	
	// Parse pagination parameters
	page := 1
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if pageNum, err := strconv.Atoi(pageParam); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	
	pageSize := 20
	if pageSizeParam := r.URL.Query().Get("page_size"); pageSizeParam != "" {
		if size, err := strconv.Atoi(pageSizeParam); err == nil && size > 0 && size <= 100 {
			pageSize = size
		}
	}
	
	// Get news from storage
	result, err := a.newsStorage.GetNewsList(filters, page, pageSize)
	if err != nil {
		a.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve news: "+err.Error())
		return
	}
	
	a.respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    result,
	})
}

// GetNewsById handles requests for a specific news item
func (a *API) GetNewsById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	// Validate ID format
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		a.respondWithError(w, http.StatusBadRequest, "Invalid news ID format")
		return
	}
	
	// Get news item from storage
	news, err := a.newsStorage.GetNewsById(id)
	if err != nil {
		a.respondWithError(w, http.StatusNotFound, "News item not found")
		return
	}
	
	a.respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    news,
	})
}

// RunConnector handles requests to run a specific connector
func (a *API) RunConnector(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	
	count, err := a.connectorService.RunConnector(r.Context(), name)
	if err != nil {
		a.respondWithError(w, http.StatusInternalServerError, "Failed to run connector: "+err.Error())
		return
	}
	
	a.respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"processed":  count,
			"connector":  name,
		},
	})
}

// RunAllConnectors handles requests to run all connectors
func (a *API) RunAllConnectors(w http.ResponseWriter, r *http.Request) {
	results, err := a.connectorService.RunAllConnectors(r.Context())
	if err != nil {
		a.respondWithError(w, http.StatusInternalServerError, "Failed to run connectors: "+err.Error())
		return
	}
	
	a.respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"results": results,
		},
	})
}

// respondWithJSON sends a JSON response
func (a *API) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError sends an error response
func (a *API) respondWithError(w http.ResponseWriter, code int, message string) {
	a.respondWithJSON(w, code, Response{
		Success: false,
		Error:   message,
	})
}