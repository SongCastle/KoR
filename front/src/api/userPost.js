import axios from 'axios';

const backendHost = process.env.BACK_HOST;
const base = axios.create({
  baseURL: backendHost,
  headers: {
    "Access-Control-Allow-Headers": "*"
  }
});

// バックエンド側のcors対策がfrontendHostのみになっているので通らない。
// バックエンドのcorsに/signupも含める必要がある。
export const postUser = (params) => base.post('/v1/users', {
  username: params.username,
  email: params.email,
  password: params.password
})
.then((res) => {
  alert('success'); 
});