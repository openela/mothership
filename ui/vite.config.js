/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

import { defineConfig } from 'vite';
import viteConfigPaths from 'vite-tsconfig-paths';

/** @type {import('vite').UserConfig} */
export default defineConfig({
  appType: 'spa',
  plugins: [viteConfigPaths()],
  build: {
    sourcemap: 'hidden',
    minify: 'esbuild',
  },
});
