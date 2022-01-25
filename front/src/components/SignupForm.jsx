import React from 'react';

export const SingupForm = () => (
  <div className='flex justify-center pt-6'>
    <form className='w-1/3 border-2 rounded-lg shadow-sm p-5'>
      <div className='mb-5'>
        <label className='flex flex-col mb-2' htmlFor='username'>
          Username
          <input
            id='username'
            type='text'
            placeholder='username'
            className='border-2 rounded p-2'
          />
        </label>
      </div>
      <div className='mb-5'>
        <label className='flex flex-col mb-2' htmlFor='email'>
          Email
          <input id='email' type='text' placeholder='email' className='border-2 rounded p-2' />
        </label>
      </div>
      <div className='mb-7'>
        <label className='flex flex-col mb-2' htmlFor='password'>
          Password
          <input
            id='password'
            type='text'
            placeholder='password'
            className='border-2 rounded p-2'
          />
        </label>
      </div>
      <input type='submit' className='bg-purple-400 border-2 rounded-lg px-2 py-1' />
    </form>
  </div>
);
