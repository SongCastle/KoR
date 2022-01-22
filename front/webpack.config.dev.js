const path = require('path');
const { merge } = require('webpack-merge');
const base = require('./webpack.config.base');

const staticPath = path.resolve(__dirname, 'static');

const config = merge(base, {
  mode: 'development',
  devtool: 'inline-source-map',
  devServer: {
    contentBase: staticPath,
    historyApiFallback: true,
    host: '0.0.0.0',
    hot: true,
    port: 8080,
    publicPath: '/assets/',
    watchContentBase: true,
  },
  watchOptions: {
    ignored: /node_modules/,
  },
  target: 'web',
});

module.exports = config;
