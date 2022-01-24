import React from "react";

export const SingupForm = () => {
  return (
    <div className="flex justify-center pt-6">
      <form className="w-1/2 border-2 rounded-lg shadow-sm p-5">
        <div className="flex flex-col mb-5">
          <label className="mb-2">Username</label>
          <input type="text" placeholder="username" className="border-2 rounded pt-2"></input>
        </div>
        <div className="flex flex-col mb-5">
          <label className="mb-2">Email</label>
          <input type="text" placeholder="email" className="border-2 rounded pt-2"></input>
        </div>
        <div className="flex flex-col mb-7">
          <label className="mb-2">Password</label>
          <input type="text" placeholder="password" className="border-2 rounded pt-2"></input>
        </div>
        <input type="submit" className="bg-purple-400 border-2 rounded-lg px-2 py-1" />
      </form>
    </div>
  );
};