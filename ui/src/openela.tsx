/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

import React from 'react';
import { cn } from '@/utils.ts';

export default function OpenELA(props: { forceLight?: boolean }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="193"
      height="193"
      fill="none"
      viewBox="0 0 193 193"
      className={cn(
        props.forceLight ? 'text-white' : 'text-black dark:text-white',
      )}
    >
      <path
        fill="currentColor"
        d="M34 73h35.005v8.28h-25.57v10.513h16.85v8.018h-16.85v10.91h25.832V119H34V73z"
      />
      <path
        fill="currentColor"
        fillRule="evenodd"
        d="M43.434 110.721v-10.91h16.851v-8.018h-16.85V81.279h25.57V73H34v46h35.267v-8.279H43.435zm24.833 1H42.435v-12.91h16.85v-6.018h-16.85V80.279h25.57V74H35v44h33.267v-6.279z"
        clipRule="evenodd"
      />
      <path fill="currentColor" d="M75.4 73h9.264v37.589h24.469V119H75.4V73z" />
      <path
        fill="currentColor"
        fillRule="evenodd"
        d="M84.664 110.589V73H75.4v46h33.733v-8.411h-24.47zm23.469 1h-24.47V74H76.4v44h31.733v-6.411z"
        clipRule="evenodd"
      />
      <path
        fill="currentColor"
        d="M130.699 73h10.79l18.244 46h-10.085l-4.561-12.092h-18.565L121.898 119H112.2l18.499-46zm11.562 26.286l-6.359-16.954-6.487 16.954h12.846z"
      />
      <path
        fill="currentColor"
        fillRule="evenodd"
        d="M145.087 106.908L149.648 119h10.085l-18.244-46h-10.79L112.2 119h9.698l4.624-12.092h18.565zM121.21 118l4.624-12.092h19.945L150.34 118h7.921L140.81 74h-9.435l-17.695 44h7.53zm22.494-17.714h-15.743l7.951-20.777 7.792 20.777zm-7.802-17.954l-6.487 16.954h12.846l-6.359-16.954z"
        clipRule="evenodd"
      />
      <path
        fill="currentColor"
        fillRule="evenodd"
        d="M96.5 174c42.802 0 77.5-34.698 77.5-77.5S139.302 19 96.5 19 19 53.698 19 96.5 53.698 174 96.5 174zm0 19c53.295 0 96.5-43.205 96.5-96.5S149.795 0 96.5 0 0 43.205 0 96.5 43.205 193 96.5 193z"
        clipRule="evenodd"
      />
    </svg>
  );
}
