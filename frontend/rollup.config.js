import commonjs from '@rollup/plugin-commonjs';
import resolve from '@rollup/plugin-node-resolve';
import livereload from 'rollup-plugin-livereload';
import replace from 'rollup-plugin-replace';
import svelte from 'rollup-plugin-svelte';
import {terser} from 'rollup-plugin-terser';
import autoPreprocess from 'svelte-preprocess';

const production = !process.env.ROLLUP_WATCH;

export default {
  input: 'src/main.js',
  output:
      {sourcemap: true, format: 'iife', name: 'app', file: 'public/bundle.js'},
  plugins: [
    replace({
      ENDPOINT: JSON.stringify(
          process.env['STHUB_ENDPOINT'] || 'http://localhost:1323'),
    }),
    svelte({
      preprocess: autoPreprocess({
        postcss: {plugins: [require('tailwindcss')], extract: true},
      }),
      // enable run-time checks when not in production
      dev: !production,
      // we'll extract any component CSS out into
      // a separate file  better for performance
      css: css => {
        css.write('public/bundle.css');
      }
    }),

    // If you have external dependencies installed from
    // npm, you'll most likely need these plugins. In
    // some cases you'll need additional configuration 
    // consult the documentation for details:
    // https://github.com/rollup/rollup-plugin-commonjs
    resolve({browser: true}), commonjs(),

    // Watch the `public` directory and refresh the
    // browser on changes when not in production
    !production && livereload('public'),

    // If we're building for production (npm run build
    // instead of npm run dev), minify
    production && terser()
  ],
  watch: {
    clearScreen: false,
    chokidar: {
      include: 'src/**',
      usePolling: true,
    },
  }
};
