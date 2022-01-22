const tailwindcss = require('tailwindcss');
const postcssPresetEnv = require('postcss-preset-env');

const { tailwindConfig } = require('./tailwind.config');

module.exports = {
  plugins: [tailwindcss(tailwindConfig), postcssPresetEnv()],
};
