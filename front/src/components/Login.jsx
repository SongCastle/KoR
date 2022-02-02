import React from 'react';
import { useState } from 'react';

export const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleChange = (event) => {
    switch (event.target.name) {
      case 'email':
        setEmail(event.target.value);
        break;
      case 'password':
        setPassword(event.target.value);
        break;
      default:
        alert('no params');
    }
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    const params = {
      email: email,
      password: password,
    };
    alert(Object.values(params));
  };

  return (
    <div className='flex justify-center pt-6 bg-gray-800 h-screen'>
      <form
        className='sm:w-1/3 w-5/6 h-2/5 sm:h-3/5 border-2 rounded-lg shadow-sm p-5 bg-white'
        onSubmit={handleSubmit}
      >
        <div className='mb-5'>
          <label className='flex flex-col mb-2' htmlFor='email'>
            Email
            <input
              id='email'
              type='text'
              name='email'
              placeholder='email'
              className='border-2 rounded p-2'
              onChange={handleChange}
            />
          </label>
        </div>
        <div className='mb-7'>
          <label className='flex flex-col mb-2' htmlFor='password'>
            Password
            <input
              id='password'
              type='password'
              name='password'
              placeholder='password'
              className='border-2 rounded p-2'
              onChange={handleChange}
            />
          </label>
        </div>
        <input
          type='submit'
          className='bg-purple-400 border-2 rounded-lg px-2 py-1'
          value='Login'
        />
      </form>
    </div>
  );
};
