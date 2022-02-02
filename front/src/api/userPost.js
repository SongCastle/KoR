import axios from 'axios';

const backendHost = process.env.BACK_HOST;
const base = axios.create({
  baseURL: backendHost,
  headers: {
    'Access-Control-Allow-Headers': '*',
  },
});

export const postUser = (params) =>
  base
    .post('/v1/users', {
      login: params.username,
      email: params.email,
      password: params.password,
    })
    .then(() => {
      alert('success');
    });
