import React from 'react';
import Slider from 'react-slick';
import '../styles/slick.css';

export const ReactSlick = () => {
  const settings = {
    dots: true,
    infinite: true,
    speed: 500,
    slidesToShow: 3,
    slidesToScroll: 1,
  };
  return (
    <Slider {...settings}>
      <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
      <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
      <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
      <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
      <div className='w-72 h-64 border rounded-lg bg-white mr-2'></div>
    </Slider>
  );
};
