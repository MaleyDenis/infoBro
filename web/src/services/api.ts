import axios from 'axios';

// Create an axios instance with base configuration
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  headers: {
    'Content-Type': 'application/json'
  }
});

// Types
export interface NewsItem {
  id: string;
  title: string;
  content?: string;
  content_preview?: string;
  source_type: string;
  source_id: string;
  source_name: string;
  source_url: string;
  url: string;
  published_at: string;
  processed_at: string;
}

export interface Pagination {
  page: number;
  page_size: number;
  total_pages: number;
  total_items: number;
}

export interface NewsListResponse {
  success: boolean;
  data: {
    items: NewsItem[];
    pagination: Pagination;
  };
  error?: string;
}

export interface NewsItemResponse {
  success: boolean;
  data: NewsItem;
  error?: string;
}

export interface ConnectorRunResponse {
  success: boolean;
  data: {
    processed: number;
    connector: string;
  };
  error?: string;
}

export interface ConnectorRunAllResponse {
  success: boolean;
  data: {
    results: {
      [key: string]: {
        status: string;
        processed?: number;
        message?: string;
      };
    };
  };
  error?: string;
}

export interface NewsFilters {
  source_type?: string;
  source_id?: string;
  query?: string;
  from_date?: string;
  to_date?: string;
  page?: number;
  page_size?: number;
}

// API Functions
export const fetchNewsList = async (filters: NewsFilters = {}): Promise<NewsListResponse> => {
  const response = await api.get('/news', { params: filters });
  return response.data;
};

export const fetchNewsItem = async (id: string): Promise<NewsItemResponse> => {
  const response = await api.get(`/news/${id}`);
  return response.data;
};

export const runConnector = async (name: string): Promise<ConnectorRunResponse> => {
  const response = await api.post(`/connectors/run/${name}`);
  return response.data;
};

export const runAllConnectors = async (): Promise<ConnectorRunAllResponse> => {
  const response = await api.post('/connectors/run-all');
  return response.data;
};

export default api;