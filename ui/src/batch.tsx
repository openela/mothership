/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

import React from 'react';
import { Link, useParams } from 'react-router-dom';
import useSWR from 'swr';

import { Batch } from '@/batch-model.ts';
import { fetchAPI } from '@/api.ts';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert.tsx';
import { ExclamationTriangleIcon, ReloadIcon } from '@radix-ui/react-icons';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb.tsx';
import { Badge } from '@/components/ui/badge.tsx';
import { capitalizeFirstLetter, timeToNatural } from '@/utils.ts';
import { EntriesResponse } from '@/entry-model.ts';
import { Skeleton } from '@/components/ui/skeleton.tsx';
import { DataTable } from '@/components/ui/data-table.tsx';
import { columns } from '@/entries.tsx';

const dataView = {
  workerId: 'Created by worker',
  bugtrackerUri: 'Bugtracker URI',
  entryCount: 'Entry count',
};

export default function Batch() {
  const { batchId } = useParams();
  const { data, error, isLoading } = useSWR<Batch>(
    '/v1/batches/' + batchId,
    fetchAPI,
  );

  const [pageTokens, setPageTokens] = React.useState<string[]>(['']);
  const [currentPage, setCurrentPage] = React.useState(0);
  const {
    data: entriesData,
    error: entriesError,
    isLoading: entriesIsLoading,
  } = useSWR<EntriesResponse>(
    '/v1/entries?pageToken=' +
      (pageTokens[currentPage] || '') +
      '&filter=batch="batches/' +
      batchId +
      '"',
    fetchAPI,
  );

  // Store the page token for the next page
  React.useEffect(() => {
    if (entriesData?.nextPageToken) {
      // If there are more history than the current page then take everything until that point
      // and create a new array with the new page token
      setPageTokens((pageTokens) => [
        ...pageTokens.slice(0, currentPage + 1),
        entriesData.nextPageToken,
      ]);
    }
  }, [entriesData]);

  const canNextPage = pageTokens.length > currentPage + 1;
  const canPreviousPage = currentPage > 0;
  const nextPage = () => setCurrentPage((page) => page + 1);
  const previousPage = () => setCurrentPage((page) => page - 1);

  return (
    <>
      {data && (
        <>
          <Breadcrumb className="pb-2">
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbLink asChild>
                  <Link to="/batches">Batches</Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbPage>
                  {data.name.substring('batches/'.length)}
                </BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
          <div className="flex items-center space-x-4">
            <h1 className="text-3xl font-light">{data.name}</h1>
            <Badge variant={data.sealTime ? 'default' : 'destructive'}>
              {data.sealTime ? 'Sealed' : 'Not sealed'}
            </Badge>
          </div>
          <div className="pb-4 flex space-x-4">
            <div>
              <span className="font-bold">Created </span>
              <span>{timeToNatural(data.createTime)}</span>
            </div>
            {data.sealTime && (
              <div>
                <span className="font-bold">Sealed </span>
                <span>{timeToNatural(data.sealTime)}</span>
              </div>
            )}
          </div>
          <div className="bg-slate-100 dark:bg-slate-900 p-4">
            {Object.entries(dataView).map(([key, value]) => (
              <div
                key={key}
                className="flex items-center space-x-4 w-[1300px] py-3"
              >
                <div className="w-80 font-bold">{value}</div>
                <div>
                  {key.indexOf('Uri') !== -1 ? (
                    <a
                      target="_blank"
                      className="text-blue-400"
                      href={(data as any)[key]}
                    >
                      {(data as any)[key]}
                    </a>
                  ) : (
                    ((data as any)[key] as string)
                  )}
                </div>
              </div>
            ))}
          </div>
          {entriesIsLoading && (
            <>
              <Skeleton className="h-4 w-full" />
              {[...Array(3)].map((_, i) => (
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
          {entriesError && (
            <Alert variant="destructive">
              <ExclamationTriangleIcon className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>
                There was an error fetching entries.
              </AlertDescription>
            </Alert>
          )}
          {entriesData && (
            <DataTable
              columns={columns}
              data={entriesData.entries}
              nextPage={nextPage}
              previousPage={previousPage}
              canNextPage={canNextPage}
              canPreviousPage={canPreviousPage}
            />
          )}
        </>
      )}
      {isLoading && <Skeleton className="h-4 w-full" />}
      {error && (
        <Alert variant="destructive">
          <ExclamationTriangleIcon className="h-4 w-4" />
          <AlertTitle>{error.code === 5 ? 'Not Found' : 'Error'}</AlertTitle>
          <AlertDescription>
            {error.code === 5
              ? 'The batch you are looking for does not exist.'
              : error.message}
          </AlertDescription>
        </Alert>
      )}
    </>
  );
}
