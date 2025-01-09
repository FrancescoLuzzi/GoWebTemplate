import Alpine from 'alpinejs';
import { theme } from './theme';
import { password } from './password';

// On page load, apply the saved theme or default to light
document.addEventListener('alpine:init', () => {
  Alpine.data('theme', theme);
  Alpine.data('password', password);
  Alpine.store('auth');
});

Alpine.start();
