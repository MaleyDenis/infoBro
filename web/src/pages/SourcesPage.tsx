import React from 'react';
import { useRunConnector } from '../hooks/useNews';

const sources = [
  {
    id: 'reddit',
    name: 'Reddit',
    type: 'reddit',
    description: 'News from popular technology subreddits',
    color: 'bg-orange-100 text-orange-800',
    icon: (
      <svg className="w-8 h-8 text-orange-500" viewBox="0 0 24 24" fill="currentColor">
        <path d="M12 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0zm5.01 4.744c.688 0 1.25.561 1.25 1.249a1.25 1.25 0 0 1-2.498.056l-2.597-.547-.8 3.747c1.824.07 3.48.632 4.674 1.488.308-.309.73-.491 1.207-.491.968 0 1.754.786 1.754 1.754 0 .716-.435 1.333-1.01 1.614a3.111 3.111 0 0 1 .042.52c0 2.694-3.13 4.87-7.004 4.87-3.874 0-7.004-2.176-7.004-4.87 0-.183.015-.366.043-.534A1.748 1.748 0 0 1 4.028 12c0-.968.786-1.754 1.754-1.754.463 0 .898.196 1.207.49 1.207-.883 2.878-1.43 4.744-1.487l.885-4.182a.342.342 0 0 1 .14-.197.35.35 0 0 1 .238-.042l2.906.617a1.214 1.214 0 0 1 1.108-.701zM9.25 12C8.561 12 8 12.562 8 13.25c0 .687.561 1.248 1.25 1.248.687 0 1.248-.561 1.248-1.249 0-.688-.561-1.249-1.249-1.249zm5.5 0c-.687 0-1.248.561-1.248 1.25 0 .687.561 1.248 1.249 1.248.688 0 1.249-.561 1.249-1.249 0-.687-.562-1.249-1.25-1.249zm-5.466 3.99a.327.327 0 0 0-.231.094.33.33 0 0 0 0 .463c.842.842 2.484.913 2.961.913.477 0 2.105-.056 2.961-.913a.361.361 0 0 0 .029-.463.33.33 0 0 0-.464 0c-.547.533-1.684.73-2.512.73-.828 0-1.979-.196-2.512-.73a.326.326 0 0 0-.232-.095z" />
      </svg>
    ),
    sources: [
      { id: 'golang', name: 'r/golang' },
      { id: 'programming', name: 'r/programming' },
      { id: 'rust', name: 'r/rust' },
      { id: 'machinelearning', name: 'r/MachineLearning' },
    ],
  },
  {
    id: 'telegram',
    name: 'Telegram',
    type: 'telegram',
    description: 'News from tech Telegram channels',
    color: 'bg-blue-100 text-blue-800',
    icon: (
      <svg className="w-8 h-8 text-blue-500" viewBox="0 0 24 24" fill="currentColor">
        <path d="M12 0C5.374 0 0 5.373 0 12c0 6.627 5.374 12 12 12 6.628 0 12-5.373 12-12 0-6.627-5.372-12-12-12zm3.224 17.871c.188.133.43.131.618-.002a.468.468 0 0 0 .28-.417c.242-2.349 1.257-8.292 1.588-10.434.021-.14-.012-.266-.91-.352a.678.678 0 0 0-.628.13c-1.355 1.116-7.291 5.715-7.291 5.715l-3.102 1.033c-.252.083-.394.222-.363.472.025.211.208.344.45.3l2.863-.674 1.73 1.283c.226.167.461.152.62-.055 0 0 .662-2.932.662-2.932.13-.043.22.015.22.015l2.88 1.87c.13.257.414.387.673.296z" />
      </svg>
    ),
    sources: [
      { id: 'golang_news', name: 'Golang News' },
      { id: 'rustlang', name: 'Rust Language' },
      { id: 'python', name: 'Python Insider' },
    ],
  },
  {
    id: 'rss',
    name: 'RSS',
    type: 'rss',
    description: 'News from RSS feeds',
    color: 'bg-orange-100 text-orange-800',
    icon: (
      <svg className="w-8 h-8 text-orange-500" viewBox="0 0 24 24" fill="currentColor">
        <path d="M6.503 20.752c0 1.794-1.456 3.248-3.251 3.248-1.796 0-3.252-1.454-3.252-3.248 0-1.794 1.456-3.248 3.252-3.248 1.795.001 3.251 1.454 3.251 3.248zm-6.503-12.572v4.811c6.05.062 10.96 4.966 11.022 11.009h4.817c-.062-8.71-7.118-15.758-15.839-15.82zm0-3.368c10.58.046 19.152 8.594 19.183 19.188h4.817c-.03-13.231-10.755-23.954-24-24v4.812z"/>
      </svg>
    ),
    sources: [
      { id: 'hackernews', name: 'Hacker News' },
      { id: 'theverge', name: 'The Verge' },
      { id: 'devto', name: 'DEV Community' },
    ],
  },
];

const SourcesPage: React.FC = () => {
  const runConnectorMutation = useRunConnector();

  const handleRunConnector = (id: string) => {
    runConnectorMutation.mutate(id);
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-8">News Sources</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {sources.map((source) => (
          <div key={source.id} className="bg-white rounded-lg shadow-sm overflow-hidden">
            <div className="p-6">
              <div className="flex items-start">
                <div className="flex-shrink-0">{source.icon}</div>
                <div className="ml-4">
                  <h3 className="text-lg font-medium text-gray-900">{source.name}</h3>
                  <p className="mt-1 text-sm text-gray-500">{source.description}</p>
                </div>
              </div>

              <div className="mt-4">
                <div className="flex items-center justify-between">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${source.color}`}>
                    {source.type}
                  </span>
                  <button
                    onClick={() => handleRunConnector(source.id)}
                    disabled={runConnectorMutation.isLoading && runConnectorMutation.variables === source.id}
                    className="btn btn-outline text-sm"
                  >
                    {runConnectorMutation.isLoading && runConnectorMutation.variables === source.id
                      ? 'Fetching...'
                      : 'Fetch Now'}
                  </button>
                </div>
              </div>

              <div className="mt-5">
                <h4 className="text-sm font-medium text-gray-900 mb-2">Available Sources</h4>
                <ul className="space-y-1">
                  {source.sources.map((subSource) => (
                    <li key={subSource.id} className="text-sm text-gray-600">
                      • {subSource.name}
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default SourcesPage;