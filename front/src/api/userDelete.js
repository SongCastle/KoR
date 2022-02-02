import axios from 'axios';

const backendHost = process.env.BACK_HOST;
const base = axios.create({
  baseURL: backendHost,
  headers: {
    'Access-Control-Allow-Headers': '*',
  },
});

export const postDelete = (id) =>
  base.delete(`/v1/users/${id}`).then((res) => {
    console.log(res);
  });
