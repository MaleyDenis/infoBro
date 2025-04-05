import { useQuery, useMutation, useQueryClient } from 'react-query';
import { 
  fetchNewsList, 
  fetchNewsItem, 
  runConnector, 
  runAllConnectors, 
  NewsFilters 
} from '../services/api';

export const useNewsList = (filters: NewsFilters = {}) => {
  return useQuery(
    ['newsList', filters],
    () => fetchNewsList(filters),
    {
      keepPreviousData: true,
      staleTime: 5 * 60 * 1000, // 5 minutes
    }
  );
};

export const useNewsItem = (id: string) => {
  return useQuery(
    ['newsItem', id],
    () => fetchNewsItem(id),
    {
      enabled: !!id,
      staleTime: 5 * 60 * 1000, // 5 minutes
    }
  );
};

export const useRunConnector = () => {
  const queryClient = useQueryClient();
  
  return useMutation(
    (name: string) => runConnector(name),
    {
      onSuccess: () => {
        // Invalidate news list queries to refresh the data
        queryClient.invalidateQueries('newsList');
      }
    }
  );
};

export const useRunAllConnectors = () => {
  const queryClient = useQueryClient();
  
  return useMutation(
    () => runAllConnectors(),
    {
      onSuccess: () => {
        // Invalidate news list queries to refresh the data
        queryClient.invalidateQueries('newsList');
      }
    }
  );
};