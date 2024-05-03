/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

export {};

declare global {
  interface Window {
    username: string;
    repoBaseURI: string;
  }
}
