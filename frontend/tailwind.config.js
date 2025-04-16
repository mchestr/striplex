module.exports = {
  content: [
    './src/**/*.{js,jsx,ts,tsx}',
    './public/index.html'
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#4b6bfb',
          hover: '#3557fa',
          light: '#e6ebfe'
        },
        secondary: {
          DEFAULT: '#e9ecef',
          hover: '#dee2e6'
        },
        cancel: {
          DEFAULT: '#ff7675',
          hover: '#ff6b6b'
        },
        donate: {
          DEFAULT: '#ffb26b',
          hover: '#ffc988',
          text: '#7d5a50'
        }
      },
      fontFamily: {
        sans: ['Inter', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'sans-serif']
      },
    },
  },
  plugins: [],
}
