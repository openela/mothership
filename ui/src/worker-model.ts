export interface Worker {
  // Output only. The resource name of the worker.
  // Format: `workers/{worker}`
  name: string;

  // Unique identifier selected during creation.
  // Cannot be changed. Must conform to RFC-1034.
  workerId: string;

  // When the worker was created.
  createTime: string;

  // Last check-in time of the worker.
  lastCheckinTime?: string;

  // API secret that the worker should use to authenticate itself.
  // This is only returned when creating a new worker.
  // Can not be retrieved or changed later.
  apiSecret?: string;
}

export interface WorkersResponse {
  // The list of workers.
  workers: Worker[];

  // The next page token.
  nextPageToken: string;
}
