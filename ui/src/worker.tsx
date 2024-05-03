/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

import React from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import useSWR from 'swr';

import { Worker } from '@/worker-model.ts';
import { fetchAdminAPI, fetchAPI } from '@/api.ts';
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
import { Button } from '@/components/ui/button.tsx';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog.tsx';
import { useToast } from '@/components/ui/use-toast.ts';

export default function Worker() {
  const { toast } = useToast();
  const { workerId } = useParams();
  const navigate = useNavigate();
  const { data, error, isLoading } = useSWR<Worker>(
    '/v1/workers/' + workerId,
    fetchAdminAPI,
  );

  const [submitting, setSubmitting] = React.useState(false);
  const doDelete = async () => {
    setSubmitting(true);

    const res: { code?: number } = await fetchAdminAPI(
      '/v1/workers/' + workerId,
      {
        method: 'DELETE',
      },
    );
    if (res.code === 13) {
      setSubmitting(false);
      toast({
        title: 'Worker has content',
        description: 'This worker can not be deleted because it has been used.',
        variant: 'destructive',
      });
      return;
    }

    navigate('/workers');
  };

  return (
    <>
      {data && (
        <>
          <Breadcrumb className="pb-2">
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbLink asChild>
                  <Link to="/workers">Workers</Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbPage>
                  {data.name.substring('workers/'.length)}
                </BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
          <div className="flex items-center space-x-4">
            <h1 className="text-3xl font-light">{data.name}</h1>
            <div className="w-full flex justify-end">
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button
                    className="ml-auto"
                    variant="destructive"
                    disabled={submitting}
                  >
                    {submitting && (
                      <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Delete
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>
                      Are you sure you want to delete this worker?
                    </AlertDialogTitle>
                    <AlertDialogDescription>
                      This action cannot be undone. Any client that uses the
                      worker secret for{' '}
                      <span className="font-bold">{data.workerId}</span> will no
                      longer be able to communicate with the server.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction onClick={doDelete}>
                      Delete
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          </div>
          <div className="pb-4 flex space-x-4">
            <div>
              <span className="font-bold">Created </span>
              <span>{timeToNatural(data.createTime)}</span>
            </div>
            {data.lastCheckinTime && (
              <div>
                <span className="font-bold">Last active </span>
                <span>{timeToNatural(data.lastCheckinTime)}</span>
              </div>
            )}
          </div>
        </>
      )}
      {isLoading && <Skeleton className="h-4 w-full" />}
      {error && (
        <Alert variant="destructive">
          <ExclamationTriangleIcon className="h-4 w-4" />
          <AlertTitle>{error.code === 5 ? 'Not Found' : 'Error'}</AlertTitle>
          <AlertDescription>
            {error.code === 5
              ? 'The worker you are looking for does not exist.'
              : error.message}
          </AlertDescription>
        </Alert>
      )}
    </>
  );
}
