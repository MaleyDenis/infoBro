import React from 'react';
import { useParams, Link } from 'react-router-dom';
import { useNewsItem } from '../hooks/useNews';
import { format } from 'date-fns';

const NewsDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { data, isLoading, isError } = useNewsItem(id || '');

  if (isLoading) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="animate-pulse">
          <div className="h-6 bg-gray-200 rounded w-3/4 mb-4"></div>
          <div className="h-4 bg-gray-200 rounded w-1/4 mb-8"></div>
          <div className="h-4 bg-gray-200 rounded mb-2.5"></div>
          <div className="h-4 bg-gray-200 rounded mb-2.5"></div>
          <div className="h-4 bg-gray-200 rounded mb-2.5"></div>
          <div className="h-4 bg-gray-200 rounded w-4/5 mb-2.5"></div>
        </div>
      </div>
    );
  }

  if (isError || !data || !data.success) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="bg-red-50 p-4 rounded-md">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-red-800">Error loading news item</h3>
              <div className="mt-2 text-sm text-red-700">
                <p>The news item could not be found or there was an error loading it.</p>
              </div>
              <div className="mt-4">
                <Link to="/" className="btn btn-outline">
                  Back to Home
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  const newsItem = data.data;
  const publishedDate = new Date(newsItem.published_at);
  const formattedDate = format(publishedDate, 'MMMM d, yyyy h:mm a');

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="mb-6">
        <Link to="/" className="inline-flex items-center text-primary-600 hover:text-primary-700">
          <svg className="mr-2 h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Back to News
        </Link>
      </div>

      <article className="bg-white rounded-lg shadow-md overflow-hidden">
        <div className="p-6 md:p-8">
          <div className="flex items-center mb-4">
            <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
              {newsItem.source_type}
            </span>
            <span className="mx-2 text-gray-500">•</span>
            <span className="text-gray-600">{newsItem.source_name}</span>
            <span className="mx-2 text-gray-500">•</span>
            <span className="text-gray-500">{formattedDate}</span>
          </div>

          <h1 className="text-3xl font-bold text-gray-900 mb-6">{newsItem.title}</h1>

          <div className="prose max-w-none">
            {newsItem.content ? (
              <div dangerouslySetInnerHTML={{ __html: newsItem.content.replace(/\n/g, '<br />') }} />
            ) : (
              <p className="text-gray-600">No content available for this news item.</p>
            )}
          </div>

          <div className="mt-8 pt-4 border-t border-gray-200">
            <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center">
              <a
                href={newsItem.url}
                target="_blank"
                rel="noopener noreferrer"
                className="btn btn-primary mb-4 sm:mb-0"
              >
                Read Original Article
              </a>

              <a
                href={newsItem.source_url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary-600 hover:text-primary-700 font-medium"
              >
                View Source: {newsItem.source_name} →
              </a>
            </div>
          </div>
        </div>
      </article>
    </div>
  );
};

export default NewsDetailPage;