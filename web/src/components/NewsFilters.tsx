import React, { useState } from 'react';

interface FiltersProps {
  onFilterChange: (filters: FilterValues) => void;
}

interface FilterValues {
  source_type?: string;
  source_id?: string;
  query?: string;
  from_date?: string;
  to_date?: string;
}

const NewsFilters: React.FC<FiltersProps> = ({ onFilterChange }) => {
  const [filters, setFilters] = useState<FilterValues>({
    source_type: '',
    source_id: '',
    query: '',
    from_date: '',
    to_date: '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFilters(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Clean up empty filters
    const cleanFilters = Object.fromEntries(
      Object.entries(filters).filter(([_, value]) => value !== '')
    );
    
    onFilterChange(cleanFilters);
  };

  const handleReset = () => {
    setFilters({
      source_type: '',
      source_id: '',
      query: '',
      from_date: '',
      to_date: '',
    });
    
    onFilterChange({});
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-4 mb-6">
      <h3 className="text-lg font-medium text-gray-900 mb-4">Filter News</h3>
      
      <form onSubmit={handleSubmit}>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div>
            <label htmlFor="source_type" className="block text-sm font-medium text-gray-700 mb-1">
              Source Type
            </label>
            <select
              id="source_type"
              name="source_type"
              value={filters.source_type}
              onChange={handleChange}
              className="input w-full rounded-md"
            >
              <option value="">All Sources</option>
              <option value="reddit">Reddit</option>
              <option value="telegram">Telegram</option>
              <option value="rss">RSS</option>
            </select>
          </div>
          
          <div>
            <label htmlFor="query" className="block text-sm font-medium text-gray-700 mb-1">
              Keyword Search
            </label>
            <input
              type="text"
              id="query"
              name="query"
              value={filters.query}
              onChange={handleChange}
              placeholder="Search in titles and content"
              className="input w-full rounded-md"
            />
          </div>
          
          <div>
            <label htmlFor="from_date" className="block text-sm font-medium text-gray-700 mb-1">
              From Date
            </label>
            <input
              type="date"
              id="from_date"
              name="from_date"
              value={filters.from_date}
              onChange={handleChange}
              className="input w-full rounded-md"
            />
          </div>
          
          <div>
            <label htmlFor="to_date" className="block text-sm font-medium text-gray-700 mb-1">
              To Date
            </label>
            <input
              type="date"
              id="to_date"
              name="to_date"
              value={filters.to_date}
              onChange={handleChange}
              className="input w-full rounded-md"
            />
          </div>
        </div>
        
        <div className="mt-4 flex justify-end space-x-3">
          <button
            type="button"
            onClick={handleReset}
            className="btn btn-outline"
          >
            Reset
          </button>
          
          <button
            type="submit"
            className="btn btn-primary"
          >
            Apply Filters
          </button>
        </div>
      </form>
    </div>
  );
};

export default NewsFilters;