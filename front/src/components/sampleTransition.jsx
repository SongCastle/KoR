import React from 'react';
import { Transition } from '@headlessui/react';
import { useState } from 'react';

// 上手くrotateしてないので要修正
export const MyComponent = () => {
  const [isShowing, setIsShowing] = useState(false);

  return (
    <>
      <button onClick={() => setIsShowing((isShowing) => !isShowing)}>Toggle</button>
      <Transition
        show={isShowing}
        enter='transform transition duration-75 duration-[400ms]'
        enterFrom='opacity-0 rotate-[-120deg] scale-50'
        enterTo='opacity-100 rotate-0 scale-100'
        leave='transform rotate-[-240deg] transition duration-150 ease-in-out'
        leaveFrom='opacity-100 rotate-0 scale-100'
        leaveTo='opacity-0'
      >
        <div className='m-auto border w-20 h-20'></div>I will fade in and out
      </Transition>
    </>
  );
};
