import React from "react";
import { Link } from "react-router-dom";


export const Header = () => (
  <header className='bg-gray-800 text-gray-100 flex'>
    <div className='title text-4xl p-5 flex-grow'><Link to='/'>KoR</Link></div>
    <ul className='flex items-center px-8'>
      <li className='mr-3'><Link to='/signup'>Signup</Link></li>
      <li><Link to='/login'>Login</Link></li>
    </ul>
  </header>
)