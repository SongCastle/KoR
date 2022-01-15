import axios from 'axios';
import useSWR from 'swr';

const backendHost = process.env.BACK_HOST;
const base = axios.create({
  baseURL: backendHost,
});
const fetcher = (url) => base.get(url).then((res) => res.data);

export const useRequest = (url) => useSWR(url, fetcher);
