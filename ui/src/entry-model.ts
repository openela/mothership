/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

export enum EntryState {
  STATE_UNSPECIFIED = 0,
  ARCHIVING = 1,
  ARCHIVED = 2,
  ON_HOLD = 3,
  CANCELLED = 4,
  FAILED = 5,
  RETRACTING = 6,
  RETRACTED = 7,
}

export interface Entry {
  // Unique ID of the entry. Format: `entries/{entry_id}`
  name: string;

  // Package NEVRA (name-epoch:version-release.arch) of the package being archived.
  entryId: string;

  // When the package was archived.
  createTime: string;

  // OS release value the package was pulled from.
  osRelease: string;

  // SHA256 of the package.
  sha256Sum: string;

  // Repository name of the package as in which repository the package was archived from.
  repository: string;

  // Worker ID of the worker that archived the package.
  // If not set, the package was archived by a user instead of a worker.
  workerId: string | null;

  // Name of the batch the package was archived in.
  // If not set, the package was not archived in a batch.
  batch: string | null;

  // User email of the user that archived the package.
  // If not set, the package was archived by a worker instead of a user.
  userEmail: string | null;

  // URI to view commit
  commitUri: string;

  // Commit hash of the resulting import
  commitHash: string;

  // Commit branch of the resulting import
  commitBranch: string;

  // Commit tag of the resulting import
  commitTag: string;

  // State of the entry.
  state: EntryState;

  // Name of the package being archived.
  pkg: string;

  // Error message if on hold
  errorMessage: string;
}

export interface EntriesResponse {
  entries: Entry[];
  nextPageToken: string;
}
