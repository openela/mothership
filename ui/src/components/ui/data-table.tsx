/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

'use client';

import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table';

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

import { Button } from '@/components/ui/button';

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  nextPage?: () => void;
  previousPage?: () => void;
  canNextPage?: boolean;
  canPreviousPage?: boolean;
}

export function DataTable<TData, TValue>({
  columns,
  data,
  nextPage,
  previousPage,
  canNextPage,
  canPreviousPage,
}: DataTableProps<TData, TValue>) {
  const table = useReactTable({
    data: data || [],
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div>
      <div className="rounded-md border">
        <Table>
          <TableHeader className="bg-slate-100 dark:bg-slate-900">
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead
                      key={header.id}
                      colSpan={header.colSpan}
                      style={{ width: `${header.getSize()}px` }}
                    >
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext(),
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody className="bg-slate-50 dark:bg-slate-950">
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && 'selected'}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext(),
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      {nextPage && (
        <div className="flex items-center justify-end space-x-2 py-4">
          <Button
            variant="outline"
            size="sm"
            onClick={previousPage}
            disabled={!canPreviousPage}
          >
            Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={nextPage}
            disabled={!canNextPage}
          >
            Next
          </Button>
        </div>
      )}
    </div>
  );
}
