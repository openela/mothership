/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

import React from 'react';
import OpenELA from '@/openela.tsx';
import GithubIcon from '@/github-icon.tsx';
import { Button } from '@/components/ui/button.tsx';
import { Navigate } from 'react-router-dom';

export default function AuthPage() {
  if (window.username) {
    return <Navigate to="/" />;
  }

  return (
    <div className="w-screen h-screen flex bg-gray-950">
      <div className="w-full items-center justify-center bg-grey hidden xl:flex pr-12">
        <OpenELA forceLight />
      </div>
      <div className="bg-white min-w-400 w-full xl:w-2/4 p-12 lg:p-48 xl:p-24 items-center justify-center flex flex-col relative">
        <h3 className="text-4xl font-bold text-black mb-8">Mothership</h3>
        <a href="/connect/github">
          <Button
            variant="outline"
            type="button"
            className="bg-gray-950 text-white border-0"
          >
            <GithubIcon />
            Sign in with GitHub
          </Button>
        </a>
        <div className="absolute" style={{ bottom: '100px' }}>
          <h4 className="text-black font-medium">
            &copy; 2024 Mustafa Gezen and Ctrl IQ, Inc.
          </h4>
        </div>
      </div>
    </div>
  );
}
