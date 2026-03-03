/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx,tsx}',
    './App.{js,ts,jsx,tsx}',
    './components/**/*.{js,ts,jsx,tsx}',
    './contexts/**/*.{js,ts,jsx,tsx}',
    './hooks/**/*.{js,ts,jsx,tsx}',
    './pages/**/*.{js,ts,jsx,tsx}',
    './services/**/*.{js,ts,jsx,tsx}',
    './utils/**/*.{js,ts,jsx,tsx}'
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: [
          'Noto Sans SC',
          'Noto Sans CJK SC',
          'Microsoft YaHei',
          'PingFang SC',
          'Hiragino Sans GB',
          'sans-serif'
        ]
      },
      colors: {
        primary: '#2563eb',
        secondary: '#475569'
      }
    }
  },
  plugins: []
};
