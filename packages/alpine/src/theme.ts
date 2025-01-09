type Theme = 'light' | 'dark';
const themeKey = 'customTheme';

function loadTheme(): Theme {
  let localTheme = localStorage.getItem(themeKey);
  if (localTheme === 'light' || localTheme === 'dark') {
    return localTheme as Theme;
  }
  return window.matchMedia('(prefers-color-scheme: dark)').matches
    ? 'dark'
    : 'light';
}

export const theme = () => ({
  theme: loadTheme(),
  toggleTheme() {
    this.theme = this.theme === 'light' ? 'dark' : 'light';
    localStorage.setItem(themeKey, this.theme);
  },
});
