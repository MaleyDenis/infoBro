import React from 'react';
import { Link } from 'react-router-dom';
import { formatDistanceToNow } from 'date-fns';
import { NewsItem } from '../services/api';

interface NewsCardProps {
  item: NewsItem;
}

const NewsCard: React.FC<NewsCardProps> = ({ item }) => {
  const publishedDate = new Date(item.published_at);
  const timeAgo = formatDistanceToNow(publishedDate, { addSuffix: true });
  
  // Get a proper icon for the source type
  const getSourceIcon = (sourceType: string) => {
    switch (sourceType) {
      case 'reddit':
        return (
          <svg className="w-4 h-4 text-orange-500" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0zm5.01 4.744c.688 0 1.25.561 1.25 1.249a1.25 1.25 0 0 1-2.498.056l-2.597-.547-.8 3.747c1.824.07 3.48.632 4.674 1.488.308-.309.73-.491 1.207-.491.968 0 1.754.786 1.754 1.754 0 .716-.435 1.333-1.01 1.614a3.111 3.111 0 0 1 .042.52c0 2.694-3.13 4.87-7.004 4.87-3.874 0-7.004-2.176-7.004-4.87 0-.183.015-.366.043-.534A1.748 1.748 0 0 1 4.028 12c0-.968.786-1.754 1.754-1.754.463 0 .898.196 1.207.49 1.207-.883 2.878-1.43 4.744-1.487l.885-4.182a.342.342 0 0 1 .14-.197.35.35 0 0 1 .238-.042l2.906.617a1.214 1.214 0 0 1 1.108-.701zM9.25 12C8.561 12 8 12.562 8 13.25c0 .687.561 1.248 1.25 1.248.687 0 1.248-.561 1.248-1.249 0-.688-.561-1.249-1.249-1.249zm5.5 0c-.687 0-1.248.561-1.248 1.25 0 .687.561 1.248 1.249 1.248.688 0 1.249-.561 1.249-1.249 0-.687-.562-1.249-1.25-1.249zm-5.466 3.99a.327.327 0 0 0-.231.094.33.33 0 0 0 0 .463c.842.842 2.484.913 2.961.913.477 0 2.105-.056 2.961-.913a.361.361 0 0 0 .029-.463.33.33 0 0 0-.464 0c-.547.533-1.684.73-2.512.73-.828 0-1.979-.196-2.512-.73a.326.326 0 0 0-.232-.095z" />
          </svg>
        );
      case 'telegram':
        return (
          <svg className="w-4 h-4 text-blue-500" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0C5.374 0 0 5.373 0 12c0 6.627 5.374 12 12 12 6.628 0 12-5.373 12-12 0-6.627-5.372-12-12-12zm3.224 17.871c.188.133.43.131.618-.002a.468.468 0 0 0 .28-.417c.242-2.349 1.257-8.292 1.588-10.434.021-.14-.012-.266-.91-.352a.678.678 0 0 0-.628.13c-1.355 1.116-7.291 5.715-7.291 5.715l-3.102 1.033c-.252.083-.394.222-.363.472.025.211.208.344.45.3l2.863-.674 1.73 1.283c.226.167.461.152.62-.055 0 0 .662-2.932.662-2.932.13-.043.22.015.22.015l2.88 1.87c.13.257.414.387.673.296z" />
          </svg>
        );
      case 'rss':
        return (
          <svg className="w-4 h-4 text-orange-500" viewBox="0 0 24 24" fill="currentColor">
            <path d="M6.503 20.752c0 1.794-1.456 3.248-3.251 3.248-1.796 0-3.252-1.454-3.252-3.248 0-1.794 1.456-3.248 3.252-3.248 1.795.001 3.251 1.454 3.251 3.248zm-6.503-12.572v4.811c6.05.062 10.96 4.966 11.022 11.009h4.817c-.062-8.71-7.118-15.758-15.839-15.82zm0-3.368c10.58.046 19.152 8.594 19.183 19.188h4.817c-.03-13.231-10.755-23.954-24-24v4.812z"/>
          </svg>
        );
      default:
        return (
          <svg className="w-4 h-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
          </svg>
        );
    }
  };

  return (
    <div className="card hover:shadow-lg transition-shadow duration-200">
      <div className="p-4">
        <div className="flex items-center space-x-2 mb-1">
          {getSourceIcon(item.source_type)}
          <span className="text-sm font-medium text-gray-600">{item.source_name}</span>
          <span className="text-sm text-gray-400">•</span>
          <span className="text-xs text-gray-500">{timeAgo}</span>
        </div>
        
        <Link to={`/news/${item.id}`} className="block">
          <h3 className="text-lg font-semibold mb-2 text-gray-900 hover:text-primary-600">{item.title}</h3>
          {item.content_preview && (
            <p className="text-gray-600 text-sm line-clamp-2 mb-3">{item.content_preview}</p>
          )}
        </Link>
        
        <div className="flex justify-between items-center pt-2 border-t border-gray-100">
          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
            {item.source_type}
          </span>
          <a 
            href={item.url} 
            target="_blank" 
            rel="noopener noreferrer"
            className="text-sm text-primary-600 hover:text-primary-700 font-medium"
          >
            Read original →
          </a>
        </div>
      </div>
    </div>
  );
};

export default NewsCard;