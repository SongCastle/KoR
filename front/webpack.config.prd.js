const { merge } = require('webpack-merge');
const base = require('./webpack.config.base');

const config = merge(base, {
  mode: 'production',
  target: ['web', 'es5'],
});

module.exports = config;
