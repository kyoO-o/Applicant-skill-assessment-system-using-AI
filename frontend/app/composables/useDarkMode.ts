const STORAGE_KEY = 'mh-dark';

export const useDarkMode = () => {
  const isDark = useState<boolean>('dark-mode', () => false);

  const apply = (dark: boolean) => {
    if (import.meta.client) {
      document.documentElement.classList.toggle('dark', dark);
    }
  };

  const init = () => {
    if (!import.meta.client) return;
    const stored = localStorage.getItem(STORAGE_KEY);
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const dark = stored !== null ? stored === '1' : prefersDark;
    isDark.value = dark;
    apply(dark);
  };

  const toggle = () => {
    isDark.value = !isDark.value;
    apply(isDark.value);
    if (import.meta.client) {
      localStorage.setItem(STORAGE_KEY, isDark.value ? '1' : '0');
    }
  };

  return { isDark, toggle, init };
};
