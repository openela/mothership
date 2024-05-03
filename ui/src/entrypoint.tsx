/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

import React from 'react';
import ReactDOM from 'react-dom/client';

import './globals.css';

import { ThemeProvider } from '@/components/theme-provider';
import Root from '@/root.tsx';

import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
  RouterProvider,
} from 'react-router-dom';
import AuthPage from '@/auth-page.tsx';

const root = ReactDOM.createRoot(
  document.getElementById('root') || document.body,
);

const router = createBrowserRouter(
  createRoutesFromElements(
    <>
      <Route path="/auth" element={<AuthPage />} />
      <Route path="*" element={<Root />} />
    </>,
  ),
);

root.render(
  <React.StrictMode>
    <ThemeProvider defaultTheme="system" storageKey="mothership-theme">
      <RouterProvider router={router} />
    </ThemeProvider>
  </React.StrictMode>,
);
