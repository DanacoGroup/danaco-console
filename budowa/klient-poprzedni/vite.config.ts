import { defineConfig } from 'vite';

// Konfiguracja budowania interfejsu.
// `fs.allow` obejmuje katalog nadrzędny, ponieważ warstwa protokołu importuje
// kontrakt z `budowa/shared/` — jedynego źródła prawdy nazw.
export default defineConfig({
  root: '.',
  build: {
    outDir: 'dist',
    target: 'es2022',
    emptyOutDir: true,
    sourcemap: true,
  },
  server: {
    port: 5173,
    strictPort: true,
    fs: { allow: ['..'] },
  },
});
