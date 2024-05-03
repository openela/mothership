/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

import React from 'react';
import { Link, useParams } from 'react-router-dom';
import useSWR from 'swr';

import { Entry } from '@/entry-model.ts';
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
import { Skeleton } from '@/components/ui/skeleton.tsx';
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
import { Button } from '@/components/ui/button.tsx';

const nvrRegex =
  /^(\S+)-([\w~%.+^]+)-(\w+(?:\.[\w~%+]+)+?)(?:\.(\w+))?(?:\.rpm)?$/;

export default function Entry() {
  const { entryId } = useParams();
  const { data, error, isLoading } = useSWR<Entry>(
    '/v1/entries/' + entryId,
    fetchAPI,
  );

  const dataView = {
    sha256Sum: 'Source RPM SHA256',
    workerId: 'Fetched by worker',
    commitUri: 'Commit URI',
    commitHash: 'Commit hash',
    batch: 'Batch',
  };

  const pkgName = data?.entryId.match(nvrRegex)?.[1];

  const [submitting, setSubmitting] = React.useState(false);
  const doRetract = async () => {
    setSubmitting(true);
    await fetchAdminAPI('/v1/' + data?.name + ':retract', {
      method: 'POST',
    });

    window.location.reload();
  };

  const doRescue = async () => {
    setSubmitting(true);
    await fetchAdminAPI('/v1/' + data?.name + ':rescueImport', {
      method: 'POST',
    });

    window.location.reload();
  };

  return (
    <>
      {data && (
        <>
          <Breadcrumb className="pb-2">
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbLink asChild>
                  <Link to="/entries">Entries</Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbPage>
                  {data.name.substring('entries/'.length)}
                </BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
          <div className="flex justify-between items-center">
            <div className="flex items-center space-x-4">
              <h1 className="text-3xl font-light flex-grow">{data.entryId}</h1>
              <Badge
                variant={
                  data.state.toString() === 'FAILED' ||
                  data.state.toString() === 'RETRACTED' ||
                  data.state.toString() === 'ON_HOLD'
                    ? 'destructive'
                    : data.state.toString() === 'ARCHIVING' ||
                        data.state.toString() === 'RETRACTING'
                      ? 'default'
                      : 'outline'
                }
              >
                {capitalizeFirstLetter(data.state.toString())}
              </Badge>
              <Badge>{data.osRelease}</Badge>
            </div>
            {window.username && (
              <div className="flex items-center space-x-4">
                {data.state.toString() === 'ARCHIVED' && (
                  <AlertDialog>
                    <AlertDialogTrigger asChild>
                      <Button variant="destructive" disabled={submitting}>
                        {submitting && (
                          <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
                        )}
                        Retract
                      </Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>
                          Are you sure you want to retract this entry?
                        </AlertDialogTitle>
                        <AlertDialogDescription>
                          This action cannot be undone. The SRPM must be
                          re-imported to make it available again.
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction
                          onClick={doRetract}
                          disabled={submitting}
                        >
                          Retract
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                )}
                {data.state.toString() === 'ON_HOLD' && (
                  <Button disabled={submitting} onClick={doRescue}>
                    {submitting && (
                      <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Rescue
                  </Button>
                )}
              </div>
            )}
          </div>
          <div className="pb-4">
            <span className="font-bold">Imported </span>
            <span>{timeToNatural(data.createTime)}</span>
          </div>
          <div className="bg-slate-100 dark:bg-slate-900 p-4">
            {Object.entries(dataView).map(
              ([key, value]) =>
                (data as any)[key] && (
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
                      ) : key === 'batch' ? (
                        <Link
                          className="text-blue-400"
                          to={'/' + (data as any)[key]}
                        >
                          {(data as any)[key]}
                        </Link>
                      ) : (
                        ((data as any)[key] as string)
                      )}
                    </div>
                  </div>
                ),
            )}
            {!data.commitUri && (
              <div className="flex items-center space-x-4 w-[1300px] py-3">
                <div className="w-80 font-bold">Repository</div>
                <div>
                  <a
                    target="_blank"
                    className="text-blue-400"
                    href={window.repoBaseURI + '/' + pkgName}
                  >
                    {window.repoBaseURI + '/' + pkgName}
                  </a>
                </div>
              </div>
            )}
          </div>
          {data.errorMessage && (
            <pre>
              <code>{data.errorMessage}</code>
            </pre>
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
              ? 'The entry you are looking for does not exist.'
              : error.message}
          </AlertDescription>
        </Alert>
      )}
    </>
  );
}
