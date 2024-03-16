import React from 'react';
import useSWR from 'swr';

import { ColumnDef } from '@tanstack/react-table';
import { fetchAdminAPI, fetchAPI } from '@/api.ts';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
} from '@/components/ui/breadcrumb.tsx';
import { Skeleton } from '@/components/ui/skeleton.tsx';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert.tsx';
import { ExclamationTriangleIcon, ReloadIcon } from '@radix-ui/react-icons';
import { DataTable } from '@/components/ui/data-table.tsx';
import { timeToNatural } from '@/utils.ts';
import { Link, useNavigate } from 'react-router-dom';
import { Worker, WorkersResponse } from '@/worker-model.ts';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog.tsx';
import { Button } from '@/components/ui/button.tsx';
import { Input } from '@/components/ui/input.tsx';

export const columns: ColumnDef<Worker>[] = [
  {
    accessorKey: 'workerId',
    header: 'ID',
    size: 750,
    cell: (entry) => {
      const name = ('/' + entry.row.original.name) as string;
      const value = entry.getValue() as string;

      return (
        <Link to={name} className="h-6 flex items-center">
          {value}
        </Link>
      );
    },
  },
  {
    accessorKey: 'createTime',
    header: 'Created',
    size: 225,
    cell: (entry) => {
      const value = entry.getValue() as string;
      return <span>{timeToNatural(value)}</span>;
    },
  },
  {
    accessorKey: 'lastCheckinTime',
    header: 'Last active',
    size: 225,
    cell: (entry) => {
      const value = entry.getValue() as string;
      if (!value) return <span>Never</span>;
      return <span>{timeToNatural(value)}</span>;
    },
  },
];

export default function Workers() {
  const [pageTokens, setPageTokens] = React.useState<string[]>(['']);
  const [currentPage, setCurrentPage] = React.useState(0);
  const { data, error, isLoading } = useSWR<WorkersResponse>(
    '/v1/workers?pageToken=' + pageTokens[currentPage] || '',
    fetchAdminAPI,
  );

  // Store the page token for the next page
  React.useEffect(() => {
    if (data?.nextPageToken) {
      // If there are more history than the current page then take everything until that point
      // and create a new array with the new page token
      setPageTokens((pageTokens) => [
        ...pageTokens.slice(0, currentPage + 1),
        data.nextPageToken,
      ]);
    }
  }, [data]);

  const canNextPage = pageTokens.length > currentPage + 1;
  const canPreviousPage = currentPage > 0;
  const nextPage = () => setCurrentPage((page) => page + 1);
  const previousPage = () => setCurrentPage((page) => page - 1);

  const navigate = useNavigate();
  const [submitting, setSubmitting] = React.useState(false);
  const [worker, setWorker] = React.useState<Worker | undefined>(undefined);

  const create = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const workerId = e.currentTarget.workerId.value;

    setSubmitting(true);
    const res: Worker = await fetchAdminAPI('/v1/workers', {
      method: 'POST',
      body: JSON.stringify({ workerId }),
    });

    setWorker(res);
  };

  const goTo = () => {
    navigate('/' + worker?.name);
  };

  return (
    <>
      <div className="flex items-center justify-between">
        <Breadcrumb className="pb-2">
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbPage>Workers</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
        <Dialog>
          <DialogTrigger asChild>
            <Button variant="secondary">Create new worker</Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[525px]">
            {worker ? (
              <>
                <DialogHeader>
                  <DialogTitle>{worker.workerId}</DialogTitle>
                  <DialogDescription>
                    This is the only time the secret will be shown, so make sure
                    to save it somewhere safe.
                  </DialogDescription>
                </DialogHeader>
                <div className="bg-slate-50 dark:bg-slate-950 p-4 rounded-lg">
                  <p className="font-mono text-sm break-all">
                    {worker.apiSecret}
                  </p>
                </div>
                <DialogFooter>
                  <form onSubmit={goTo}>
                    <Button type="submit">Go to worker</Button>
                  </form>
                </DialogFooter>
              </>
            ) : (
              <>
                <DialogHeader>
                  <DialogTitle>Create new worker</DialogTitle>
                  <DialogDescription>
                    This ID will be used to uniquely identify this worker.
                  </DialogDescription>
                </DialogHeader>
                <form onSubmit={create} className="flex flex-col space-y-4">
                  <Input name="workerId" placeholder="Worker ID" />
                  <DialogFooter>
                    <Button type="submit" disabled={submitting}>
                      {submitting && (
                        <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
                      )}
                      Create
                    </Button>
                  </DialogFooter>
                </form>
              </>
            )}
          </DialogContent>
        </Dialog>
      </div>
      {isLoading && (
        <>
          <Skeleton className="h-4 w-full" />
          {[...Array(3)].map((_, i) => (
            <div key={i} className="flex space-x-1">
              {columns.map((column) => (
                <Skeleton
                  className="w-1/3 h-4"
                  key={column.header?.toString()}
                />
              ))}
            </div>
          ))}
        </>
      )}
      {error && (
        <Alert variant="destructive">
          <ExclamationTriangleIcon className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>
            There was an error fetching workers.
          </AlertDescription>
        </Alert>
      )}
      {data && (
        <DataTable
          columns={columns}
          data={data.workers}
          nextPage={nextPage}
          previousPage={previousPage}
          canNextPage={canNextPage}
          canPreviousPage={canPreviousPage}
        />
      )}
    </>
  );
}
