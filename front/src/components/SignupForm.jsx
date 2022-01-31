import React from 'react';
import { useState } from 'react';
import { postUser } from '../api/userPost';

export const SingupForm = () => {
  const [username, setUserame] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleChange = (event) => {
    switch(event.target.name) {
      case 'username':
        setUserame(event.target.value)
        break;
      case 'email':
        setEmail(event.target.value)
        break;
      case 'password':
        setPassword(event.target.value)
        break;
      default:
        console.log('key not found')
    }
  }

  const handleSubmit = (event) => {
    event.preventDefault();
    const params = {
      username: username,
      email: email,
      password: password
    }
    alert(Object.values(params));
    // postUser(params);
  } 

  return(
    <div className='flex justify-center pt-6 bg-gray-800 h-screen'>
      <form className='sm:w-1/3 w-5/6 h-3/5 sm:h-3/5 border-2 rounded-lg shadow-sm p-5 bg-white' onSubmit={handleSubmit}>
        <div className='mb-5'>
          <label className='flex flex-col mb-2' htmlFor='username'>
            Username
            <input
              id='username'
              type='text'
              name='username'
              placeholder='username'
              className='border-2 rounded p-2'
              onChange={handleChange}
            />
          </label>
        </div>
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
        <div className='mb-7'>
          <label className='flex flex-col mb-2' htmlFor='password'>
            Confirm Password
            <input
              id='password-confirm'
              type='password'
              name='password-confirm'
              placeholder='password'
              className='border-2 rounded p-2'
            />
          </label>
        </div>
        <input type='submit' className='bg-purple-400 border-2 rounded-lg px-2 py-1' value='Singup' />
      </form>
    </div>
  );
};
