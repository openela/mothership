import React from 'react';
import { Link, Navigate, NavLink, Route, Routes } from 'react-router-dom';
import { cn } from '@/utils.ts';
import Entries from '@/entries.tsx';
import Entry from '@/entry.tsx';
import Batches from '@/batches.tsx';
import Batch from '@/batch.tsx';
import Workers from '@/workers.tsx';
import Worker from '@/worker.tsx';
import GithubIcon from '@/github-icon.tsx';
import { Toaster } from '@/components/ui/toaster.tsx';

interface Route {
  component?: any;
  redirect?: string;
  text?: string;
  children?: Routes;
}

interface Routes {
  [key: string]: Route;
}

const adminRoutes: Routes = {
  '/workers': {
    component: <Workers />,
    text: 'Workers',
    children: {
      '/workers/:workerId': {
        component: <Worker />,
      },
    },
  },
};

const routes: Routes = Object.assign(
  {
    '/': {
      component: <Navigate to="/entries" />,
    },
    '/entries': {
      component: <Entries />,
      text: 'Entries',
      children: {
        '/entries/:entryId': {
          component: <Entry />,
        },
      },
    },
    '/batches': {
      component: <Batches />,
      text: 'Batches',
      children: {
        '/batches/:batchId': {
          component: <Batch />,
        },
      },
    },
  },
  window.username ? adminRoutes : {},
);

function generateRoutes(routes: Routes) {
  return Object.keys(routes)
    .map((routeKey) => {
      const route = routes[routeKey];
      let ret = [
        <Route key={routeKey} path={routeKey} element={route.component} />,
      ];
      if (route.children) {
        ret.push(...generateRoutes(route.children));
      }
      return ret;
    })
    .flat();
}

export default function Root() {
  const linkClasses =
    'text-sm font-medium transition-colors hover:text-primary';

  const jsxRoutes = generateRoutes(routes);

  return (
    <>
      <div className="border-b">
        <div className="flex h-16 items-center px-4">
          <Link
            to="/"
            className="text-xl font-bold pr-4 mr-4 border-r dark:text-[#69E190]"
          >
            Mothership
          </Link>
          <div className="flex items-center space-x-4 lg:space-x-6 w-full">
            {Object.keys(routes).map(
              (route) =>
                routes[route].text && (
                  <NavLink
                    to={route}
                    key={route}
                    className={({ isActive }) =>
                      cn(linkClasses, isActive || 'text-muted-foreground')
                    }
                  >
                    {routes[route].text}
                  </NavLink>
                ),
            )}
            <div className="w-full flex flex-start items-center">
              <a
                target="_blank"
                href={window.repoBaseURI}
                className="ml-auto mr-4"
              >
                <GithubIcon />
              </a>
              {window.username ? (
                <>
                  <span
                    className={cn(
                      linkClasses,
                      'text-muted-foreground mr-4 lg:mr-6',
                    )}
                  >
                    {window.username}
                  </span>
                  <a href="/_logout" className={linkClasses}>
                    Logout
                  </a>
                </>
              ) : (
                <Link
                  to="/auth"
                  className={cn(linkClasses, 'text-muted-foreground')}
                >
                  Login
                </Link>
              )}
            </div>
          </div>
        </div>
      </div>
      <div className="flex-1 space-y-4 p-4 pt-6">
        <Routes>{jsxRoutes}</Routes>
      </div>
      <Toaster />
    </>
  );
}
