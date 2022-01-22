const path = require('path');
const webpack = require('webpack');

const postcssconfig = require('./postcss.config');

const srcPath = path.resolve(__dirname, 'src');
const distPath = path.resolve(__dirname, 'dist');

const config = {
  context: srcPath,
  entry: './index.jsx',
  output: {
    path: distPath,
    filename: 'bundle.js',
    publicPath: '/assets/',
  },
  module: {
    rules: [
      {
        test: /\.jsx?$/,
        include: srcPath,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/env', '@babel/react'],
          },
        },
      },
      {
        test: /\.css$/,
        include: srcPath,
        use: [
          {
            loader: 'style-loader',
          },
          {
            loader: 'css-loader',
            options: {
              importLoaders: 1,
            },
          },
          {
            loader: 'postcss-loader',
            options: {
              postcssOptions: postcssconfig,
            },
          },
        ],
      },
    ],
  },
  resolve: {
    extensions: ['.jsx', '.js', '.css'],
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env.BACK_HOST': JSON.stringify(process.env.BACK_HOST),
    }),
  ],
};

module.exports = config;
