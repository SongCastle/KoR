import React from 'react';
import { useUsers } from '../api';
import { postDelete } from '../api/userDelete';

export const Admin = () => {
  const { users, isLoading, isError } = useUsers();
  if (isLoading) return <p>Loading ...</p>;
  if (isError) return <p>Error</p>;
  
  const handleDouble = (event) => {
    // console.log(event.currentTarget.querySelector('#id').dataset.param)
    if (confirm('削除しますか？')) {
      postDelete();
    }
  }

  return (
    <div id='admin-contents' className='flex flex-col pt-6 bg-gray-800 h-screen items-start'>
      {users.map(({ id, login, email }) => (
        <div
          key={id}
          className='w-4/5 border rounded bg-white shadow p-6 mt-0 mb-1 mx-auto cursor-pointer transition duration-500 ease-in-out transform hover:-translate-y-1 hover:scale-105'
          onDoubleClick={handleDouble}
        >
          <p id='id' data-param={id}>{`id: ${id}`}</p>
          <p className='text-2xl'>{`${login}`}</p>
          <p>{`${email}`}</p>
        </div>
      ))}
    </div>
  );
};
