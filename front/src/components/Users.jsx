import React from 'react';
import { useUsers } from '../api';

export const Users = () => {
  const { users, isLoading, isError } = useUsers();

  if (isLoading) return <p>Loading ...</p>;
  if (isError) return <p>Error</p>;
  console.log(users);

  return users.map(({ id, login, email }) => (
    <p key={id}>{`id: ${id}, login: ${login}, email: ${email}`}</p>
  ));
};
