import React from 'react';
import { NewsItem } from '../services/api';
import NewsCard from './NewsCard';
import { useNewsList } from '../hooks/useNews';

interface NewsListProps {
  filters?: {
    source_type?: string;
    source_id?: string;
    query?: string;
    from_date?: string;
    to_date?: string;
  };
  page: number;
  setPage: (page: number) => void;
}

const NewsList: React.FC<NewsListProps> = ({ filters = {}, page, setPage }) => {
  const { data, isLoading, isError, error } = useNewsList({
    ...filters,
    page,
    page_size: 12,
  });

  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-12">
        <svg className="animate-spin h-10 w-10 text-primary-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="bg-red-50 p-4 rounded-md">
        <div className="flex">
          <div className="flex-shrink-0">
            <svg className="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
          </div>
          <div className="ml-3">
            <h3 className="text-sm font-medium text-red-800">Error loading news</h3>
            <div className="mt-2 text-sm text-red-700">
              <p>{error instanceof Error ? error.message : 'Unknown error occurred'}</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!data || !data.success || !data.data.items.length) {
    return (
      <div className="bg-white rounded-lg shadow p-6 text-center">
        <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
        </svg>
        <h3 className="mt-2 text-sm font-medium text-gray-900">No news found</h3>
        <p className="mt-1 text-sm text-gray-500">
          Try refreshing or changing your filters.
        </p>
      </div>
    );
  }

  return (
    <div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {data.data.items.map((item: NewsItem) => (
          <NewsCard key={item.id} item={item} />
        ))}
      </div>
      
      {/* Pagination */}
      {data.data.pagination.total_pages > 1 && (
        <div className="flex justify-between items-center mt-8">
          <button
            onClick={() => setPage(Math.max(1, page - 1))}
            disabled={page === 1}
            className="btn btn-outline"
          >
            Previous
          </button>
          
          <span className="text-sm text-gray-700">
            Page {page} of {data.data.pagination.total_pages}
          </span>
          
          <button
            onClick={() => setPage(Math.min(data.data.pagination.total_pages, page + 1))}
            disabled={page === data.data.pagination.total_pages}
            className="btn btn-outline"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
};

export default NewsList;