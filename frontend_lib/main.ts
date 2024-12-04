import Alpine from 'alpinejs';
import { theme } from './theme';

// On page load, apply the saved theme or default to light
document.addEventListener('alpine:init', () => {
  Alpine.data('theme', theme);
  Alpine.store('auth');
});

Alpine.start();
