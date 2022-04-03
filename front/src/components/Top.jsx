import React from 'react';
import { Link } from 'react-router-dom';
import { MyComponent } from './sampleTransition';
import { ReactSlick } from './ReactSlick';
import "../styles/top.css";

export const Top = () => (
  <>
    <div className='h-auto bg-gray-800'>
      <div className='container mx-auto py-7 h-full flex'>
        <div className='artcle-index pr-12'>

          <div className='java-contents mx-auto'>
          <p className='text-white'>Java</p>
            <ReactSlick />
          </div>

          <div className='ruby-contents'>
            <p className='text-white'>Ruby</p>
            <div className='flex'>
              <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
              <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
              <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
            </div>
          </div>

          <div className='php-contents'>
            <p className='text-white'>PHP</p>
            <div className='flex'>
              <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
              <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
              <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
            </div>
          </div>
          <div className='php-contents'>
            <p className='text-white'>Python</p>
            <div className='flex'>
              <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
              <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
              <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
            </div>
          </div>
          <div className='php-contents'>
            <p className='text-white'>Go</p>
            <div className='flex'>
              <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
              <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
              <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
            </div>
          </div>
        </div>
        <div className='border-l-4 border-black w-full px-6 pt-6'>
          <div className='border border-gray-700 rounded bg-gray-100 h-96 p-4'>
            <p>タグ(ペンの画像）</p>
            <ul className='text-2xl pt-5'>
              <li className='border border-gray-300 rounded-lg bg-white px-4 inline-block'># Ruby</li>
              <li className='border border-gray-300 rounded-lg bg-white px-4 inline-block'># PHP</li>
              <li className='border border-gray-300 rounded-lg bg-white px-4 inline-block'># Go</li>
              <li className='border border-gray-300 rounded-lg bg-white px-4 inline-block'># Python</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </>
);
