import React from 'react';
import { Link } from 'react-router-dom';

export const Header = () => (
  <header className='bg-gray-900 text-gray-100 flex'>
    <div className='title text-4xl p-5 flex-grow'>
      <Link to='/'>KoR</Link>
    </div>
    <ul className='flex items-center px-8'>
      <li className='mr-3'>
        <Link to='/signup'>Signup</Link>
      </li>
      <li className='mr-3'>
        <Link to='/login'>Login</Link>
      </li>
      <li>
        <Link to='/admin'>Admin</Link>
      </li>
    </ul>
  </header>
);
