/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

import React from 'react';
import useSWR from 'swr';

import { ColumnDef } from '@tanstack/react-table';
import { fetchAPI } from '@/api.ts';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
} from '@/components/ui/breadcrumb.tsx';
import { Skeleton } from '@/components/ui/skeleton.tsx';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert.tsx';
import { ExclamationTriangleIcon } from '@radix-ui/react-icons';
import { DataTable } from '@/components/ui/data-table.tsx';
import { timeToNatural } from '@/utils.ts';
import { Link } from 'react-router-dom';
import { Batch, BatchesResponse } from '@/batch-model.ts';

export const columns: ColumnDef<Batch>[] = [
  {
    accessorKey: 'name',
    header: 'Name',
    size: 750,
    cell: (entry) => {
      const value = entry.getValue() as string;

      return (
        <Link to={'/' + value} className="h-6 flex items-center">
          {value.substring('batches/'.length)}
        </Link>
      );
    },
  },
  {
    accessorKey: 'workerId',
    header: 'Worker',
    size: 225,
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
    accessorKey: 'sealTime',
    header: 'Sealed',
    size: 225,
    cell: (entry) => {
      const value = entry.getValue() as string;
      if (!value) return <span>Not sealed</span>;
      return <span>{timeToNatural(value)}</span>;
    },
  },
  {
    accessorKey: 'entryCount',
    header: 'Entries',
    size: 50,
  },
];

export default function Batches() {
  const [pageTokens, setPageTokens] = React.useState<string[]>(['']);
  const [currentPage, setCurrentPage] = React.useState(0);
  const { data, error, isLoading } = useSWR<BatchesResponse>(
    '/v1/batches?filter=entryCount>0&pageToken=' + pageTokens[currentPage] ||
      '',
    fetchAPI,
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

  return (
    <>
      <Breadcrumb className="pb-2">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbPage>Batches</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      {isLoading && (
        <>
          <Skeleton className="h-4 w-full" />
          {[...Array(25)].map((_, i) => (
            <div key={i} className="flex space-x-1">
              {columns.map((column) => (
                <Skeleton
                  className="w-1/5 h-4"
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
            There was an error fetching batches.
          </AlertDescription>
        </Alert>
      )}
      {data && (
        <DataTable
          columns={columns}
          data={data.batches}
          nextPage={nextPage}
          previousPage={previousPage}
          canNextPage={canNextPage}
          canPreviousPage={canPreviousPage}
        />
      )}
    </>
  );
}
