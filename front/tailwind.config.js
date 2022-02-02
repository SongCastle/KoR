const tailwindConfig = {
  purge: {
    enabled: true,
    content: ['./src/**/*.jsx'],
  },
  darkMode: false,
  theme: {
    extend: {},
  },
  variants: {
    extend: {
      backgroundColor: ['responsive', 'hover', 'focus', 'active'],
      transitionProperty: ['hover', 'focus'],
    },
  },
  plugins: [],
};

module.exports = { tailwindConfig };
