import React from 'react';
import Slider from 'react-slick';
import '../styles/slick.css';

const settings = {
  dots: true,
  infinite: true,
  speed: 500,
  slidesToShow: 3,
  slidesToScroll: 1,
};

export const ReactSlick = () => {
  const { dots, infinite, speed, slidesToShow, slidesToScroll } = settings;
  return (
    <Slider
      dots={dots}
      infinite={infinite}
      speed={speed}
      slidesToShow={slidesToShow}
      slidesToScroll={slidesToScroll}
    >
      <div className='w-72 h-64 border rounded-lg bg-white mr-2' />
      <div className='w-72 h-64 border rounded-lg bg-white mr-2' />
      <div className='w-72 h-64 border rounded-lg bg-white mr-2' />
      <div className='w-72 h-64 border rounded-lg bg-white mr-2' />
      <div className='w-72 h-64 border rounded-lg bg-white mr-2' />
    </Slider>
  );
};
