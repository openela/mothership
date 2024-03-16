import React from 'react';
import useSWR from 'swr';

import { ColumnDef } from '@tanstack/react-table';
import { EntriesResponse, Entry, EntryState } from '@/entry-model.ts';
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
import { Badge } from '@/components/ui/badge.tsx';
import { capitalizeFirstLetter, timeToNatural } from '@/utils.ts';
import { Link } from 'react-router-dom';
import { Input } from '@/components/ui/input.tsx';
import { Button } from '@/components/ui/button.tsx';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select.tsx';

export const columns: ColumnDef<Entry>[] = [
  {
    accessorKey: 'entryId',
    header: 'NVRA',
    size: 750,
    cell: (entry) => {
      const state = entry.row.original.state.toString();
      const name = ('/' + entry.row.original.name) as string;
      const value = entry.getValue() as string;

      if (state !== 'ARCHIVED') {
        return (
          <div className="h-6 flex items-center">
            <Link to={name}>{value}</Link>
            <Badge variant="outline" className="ml-2">
              {capitalizeFirstLetter(state)}
            </Badge>
          </div>
        );
      }

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
    accessorKey: 'osRelease',
    header: 'OS Release',
    size: 150,
  },
];

function isEmptyFilter(value: string, filterWrapper: string): string {
  if (
    !value.includes(':') &&
    !value.includes('>') &&
    !value.includes('<') &&
    !value.includes('=') &&
    !value.includes('AND') &&
    !value.includes('OR')
  ) {
    return filterWrapper.replace('{value}', value);
  }

  return value;
}

export default function Entries() {
  const [pageTokens, setPageTokens] = React.useState<string[]>(['']);
  const [currentPage, setCurrentPage] = React.useState(0);
  const [filter, setFilter] = React.useState('');
  const { data, error, isLoading } = useSWR<EntriesResponse>(
    '/v1/entries?pageToken=' +
      (pageTokens[currentPage] || '') +
      '&filter=' +
      filter,
    fetchAPI,
  );

  // Store the page token for the next page
  React.useEffect(() => {
    if (data?.nextPageToken) {
      // If there are more history than the current page then take everything until that point
      // and create a new array with the new page token
      setPageTokens([
        ...pageTokens.slice(0, currentPage + 1),
        data.nextPageToken,
      ]);
    }

    if (currentPage === 0 && !data?.nextPageToken) {
      setPageTokens(['']);
    }
  }, [data]);

  const canNextPage = pageTokens.length > currentPage + 1;
  const canPreviousPage = currentPage > 0;
  const nextPage = () => setCurrentPage((page) => page + 1);
  const previousPage = () => setCurrentPage((page) => page - 1);

  const filterRef = React.useRef<HTMLInputElement>(null);
  const [currentQuickFilter, setCurrentQuickFilter] = React.useState('');

  const search = (evt: React.FormEvent<HTMLFormElement>) => {
    evt.preventDefault();

    const value = new FormData(evt.currentTarget).get('filter') as string;
    setCurrentQuickFilter('');
    setFilter(isEmptyFilter(value, 'entryId:"{value}"'));
    setCurrentPage(0);
    setPageTokens(['']);
  };

  const quickFilter = (value: string) => {
    setCurrentQuickFilter(value);
    filterRef.current!.value = value;
    setFilter(value);
  };

  return (
    <>
      <Breadcrumb className="pb-2">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbPage>Entries</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <form onSubmit={search}>
        <div className="flex space-x-4 w-full">
          <Select onValueChange={quickFilter} value={currentQuickFilter}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Quick filters" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value={'state="ARCHIVED"'}>Archived</SelectItem>
              <SelectItem value={'state="ARCHIVING"'}>In progress</SelectItem>
              <SelectItem value={'state="ON_HOLD"'}>On hold</SelectItem>
              <SelectItem value={'state="RETRACTED" OR state="RETRACTING"'}>
                Retracted
              </SelectItem>
              <SelectItem value={'state="FAILED" OR state="CANCELLED"'}>
                Failure
              </SelectItem>
            </SelectContent>
          </Select>
          <div className="flex justify-between items-center space-x-4 w-full">
            <Input
              ref={filterRef}
              name="filter"
              className="shadow"
              placeholder="Filter"
            />
            <Button type="submit">Search</Button>
          </div>
        </div>
      </form>
      {isLoading && (
        <>
          <Skeleton className="h-4 w-full" />
          {[...Array(25)].map((_, i) => (
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
            There was an error fetching entries.
          </AlertDescription>
        </Alert>
      )}
      {data && (
        <DataTable
          columns={columns}
          data={data.entries}
          nextPage={nextPage}
          previousPage={previousPage}
          canNextPage={canNextPage}
          canPreviousPage={canPreviousPage}
        />
      )}
    </>
  );
}
