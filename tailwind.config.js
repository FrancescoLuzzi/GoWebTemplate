/** @type {import('tailwindcss').Config} */
module.exports = {
  // https://github.com/saadeghi/daisyui/discussions/640#discussioncomment-6600595
  darkMode: ['class', '[data-theme="dark"]'],
  content: ['./app/views/**/*.{templ,html,go}'],
  theme: {
    container: {
      center: true,
      padding: '2rem',
      screens: {
        '2xl': '1400px',
      },
    },
    extend: {
      colors: {
        //extra default theme
        shadow: {
          DEFAULT: 'hsl(var(--shadow))',
        },
      },
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
      },
      keyframes: {
        'accordion-down': {
          from: { height: 0 },
          to: { height: 'var(--radix-accordion-content-height)' },
        },
        'accordion-up': {
          from: { height: 'var(--radix-accordion-content-height)' },
          to: { height: 0 },
        },
      },
      animation: {
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
        'spin-slow': 'spin 1.5s linear infinite',
      },
    },
  },
  daisyui: {
    themes: ['light', 'dark'],
  },
  plugins: [
    require('tailwindcss-animate'),
    require('@tailwindcss/typography'),
    require('daisyui'),
  ],
};
