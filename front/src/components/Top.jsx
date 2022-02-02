import React from 'react';
import { Link } from 'react-router-dom';
import { MyComponent } from './sampleTransition';

export const Top = () => (
  <>
    <div className='w-full h-screen bg-gray-800'>
      <div className='text-white text-4xl text-center pt-3 animate-pulse'>
        Welcome to our coding community. Lets join us!
      </div>
      <div className='w-2/3 h-1/2 pt-8 mx-auto mt-3 flex flex-wrap justify-center'>
        <div className='h-5 border border-blue-800 rounded-2xl text-white text-5xl bg-blue-600 px-4 py-2'>
          <Link to='/signup' className='hover:text-gray-200'>
            Get started
          </Link>
        </div>
        <div className='h-5 border border-gray-500 rounded-2xl text-white text-5xl bg-gray-400 px-4 py-2'>
          <Link to='/' className='hover:text-gray-200'>
            Document
          </Link>
        </div>
      </div>
    </div>
    <div className='w-full h-screen bg-white'>
      <div className='text-center text-7xl'>Usage</div>
      <p>ここに説明を入れる。少し動きをつけたい。headless ui</p>
      <MyComponent />
    </div>
  </>
);
