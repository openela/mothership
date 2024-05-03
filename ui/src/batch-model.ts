/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

export interface Batch {
  // Output only. Unique ID of the batch.
  name: string;

  // Custom ID of the batch. Optional
  batchId: string;

  // Worker ID that created the batch.
  workerId: string;

  // Output only. Timestamp when the batch was created.
  createTime: string;

  // Output only. Timestamp when the batch was last updated.
  updateTime?: string;

  // Output only. Timestamp when the batch was sealed.
  // Batches are automatically sealed after an hour of inactivity.
  sealTime?: string;

  // Output only. Bugtracker URI of the batch.
  bugtrackerUri?: string;

  // Output only. Entry count of the batch.
  entry_count: number;
}

export interface BatchesResponse {
  batches: Batch[];
  nextPageToken: string;
}
