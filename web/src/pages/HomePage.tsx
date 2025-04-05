import React, { useState } from 'react';
import NewsFilters from '../components/NewsFilters';
import NewsList from '../components/NewsList';
import SourceStats from '../components/SourceStats';

interface FilterValues {
  source_type?: string;
  source_id?: string;
  query?: string;
  from_date?: string;
  to_date?: string;
}

const HomePage: React.FC = () => {
  const [filters, setFilters] = useState<FilterValues>({});
  const [currentPage, setCurrentPage] = useState(1);

  // Reset page when filters change
  const handleFilterChange = (newFilters: FilterValues) => {
    setFilters(newFilters);
    setCurrentPage(1);
  };

  // Sample source stats data - in a real app, this would come from API
  const sourceStatsData = [
    { name: 'Reddit', value: 45, color: '#FF4500' },
    { name: 'Telegram', value: 30, color: '#0088CC' },
    { name: 'RSS', value: 25, color: '#FF8C00' },
  ];

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">Latest Tech News</h1>
          
          <NewsFilters onFilterChange={handleFilterChange} />
          
          <NewsList 
            filters={filters} 
            page={currentPage} 
            setPage={setCurrentPage} 
          />
        </div>
        
        <div className="lg:col-span-1 space-y-6">
          <SourceStats data={sourceStatsData} />
          
          <div className="bg-white rounded-lg shadow-sm p-4">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Popular Tags</h3>
            <div className="flex flex-wrap gap-2">
              {['go', 'programming', 'javascript', 'react', 'mongodb', 'ai', 'python', 'cybersecurity', 'cloud', 'kubernetes'].map((tag) => (
                <span 
                  key={tag}
                  className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-gray-100 text-gray-800 hover:bg-gray-200 cursor-pointer"
                  onClick={() => handleFilterChange({ query: tag })}
                >
                  {tag}
                </span>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default HomePage;