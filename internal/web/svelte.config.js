import { vitePreprocess } from '@sveltejs/vite-plugin-svelte'

export default {
  preprocess: vitePreprocess(),
  
  // Hide evil a11y warnings (normally I never supress warnings, but
  // web dev is garbage and each warning is a bug-workaround)
  onwarn: (warning, handler) => {
    if (warning.code.startsWith('a11y-')) {
      return;
    }
    handler(warning);
  },
}
