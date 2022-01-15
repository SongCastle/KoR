import { useRequest } from './common';

export const useUsers = () => {
  const { data, error } = useRequest('/v1/users');

  return {
    users: data,
    isLoading: !error && !data,
    isError: error,
  };
};
