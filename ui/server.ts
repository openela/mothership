/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import express from 'express';
import httpProxyMiddleware from 'http-proxy-middleware';
import helmet from 'helmet';
import session from 'express-session';
import grant from 'grant';
import bodyParser from 'body-parser';
import RedisStore from 'connect-redis';
import IORedis from 'ioredis';

declare module 'express-session' {
  interface SessionData {
    username?: string;
    access_token?: string;
  }
}

const prod = process.env.NODE_ENV === 'production';
const app = express();
const port = process.env.PORT || 9111;
const secret = process.env.SECRET || 'secret';
const self = process.env.SELF || `http://localhost:${port}`;
const clientID = process.env.CLIENT_ID || '';
const clientSecret = process.env.CLIENT_SECRET || '';
const githubTeam = process.env.GITHUB_TEAM || 'openela/teams/tsc';
const apiURI = process.env.API_URI || 'http://localhost:6678';
const adminApiURI = process.env.ADMIN_API_URI || 'http://localhost:6688';
const repoBaseURI =
  process.env.REPO_BASE_URI || 'https://github.com/openela-main';

// If we're in prod, then a secret has to be present
if (prod && (!secret || secret.length < 32)) {
  throw 'secret has to be at least 32 characters';
}

const { createProxyMiddleware } = httpProxyMiddleware;

const store =
  prod || process.env.USE_REDIS === 'true'
    ? (() => {
        console.log('Using Redis store');
        const client = new IORedis(
          parseInt(process.env.REDIS_PORT || '0') || 6379,
          process.env.REDIS_HOST || 'localhost',
          {
            password: process.env.REDIS_PASSWORD || undefined,
          },
        );

        return new RedisStore({ client, prefix: 'mship:' });
      })()
    : new session.MemoryStore();

app.use(bodyParser.json());
app.use(
  session({
    secret,
    store,
    resave: false,
    saveUninitialized: false,
    cookie: {
      secure: prod,
      sameSite: prod ? 'strict' : 'lax',
      maxAge: 24 * 60 * 60 * 1000 * 365,
      httpOnly: true,
    },
  }),
);
app.use(
  grant.express({
    defaults: {
      origin: self,
      transport: 'session',
      state: true,
    },
    github: {
      key: clientID,
      secret: clientSecret,
      callback: '/_auth',
      scope: ['read:user', 'read:org'],
    },
  }),
);

app.use(
  '/ui/api',
  createProxyMiddleware({
    target: apiURI,
    changeOrigin: true,
    headers: {
      host: apiURI,
    },
    pathRewrite: {
      '^/ui/api': '',
    },
  }),
);

app.use('/ui/admin-api', (req, res, next) => {
  if (!req.session.access_token) {
    return res.status(403).send('Unauthorized');
  }

  createProxyMiddleware({
    target: adminApiURI,
    changeOrigin: true,
    headers: {
      host: adminApiURI,
      Authorization: `Bearer ${req.session.access_token}`,
    },
    pathRewrite: {
      '^/ui/admin-api': '',
    },
  })(req, res, next);
});

app.get('/favicon.png', (req, res) => {
  res.sendFile(path.join(__dirname, 'favicon.png'));
});

app.get('/_auth', async (req, res) => {
  const session = req.session as any;
  if (!session.grant || !session.grant.response || session.username) {
    return res.redirect('/');
  }
  const access_token = session.grant.response.access_token;

  // Fetch username from GitHub
  const gh = await fetch('https://api.github.com/user', {
    headers: {
      Accept: 'application/vnd.github+json',
      Authorization: `Bearer ${access_token}`,
      'X-GitHub-Api-Version': '2022-11-28',
    },
  });
  const json = await gh.json();
  const username = json.login;

  // Verify that the user is in the correct team
  const team = await fetch(
    `https://api.github.com/orgs/${githubTeam}/memberships/${username}`,
    {
      headers: {
        Accept: 'application/vnd.github+json',
        Authorization: `Bearer ${access_token}`,
        'X-GitHub-Api-Version': '2022-11-28',
      },
    },
  );
  const teamJson = await team.json();
  if (teamJson.state !== 'active') {
    return res.redirect('/entries');
  }

  session.username = username;
  session.access_token = access_token;

  req.session.save(() => {
    res.redirect('/entries');
  });
});

app.get('/_logout', (req, res) => {
  req.session.destroy(() => {
    res.redirect('/');
  });
});

if (prod) {
  app.set('trust proxy', 1);
  // Enable security hardening in prod
  app.use(
    helmet({
      contentSecurityPolicy: false,
      hidePoweredBy: true,
    }),
  );
  app.set('etag', false);

  app.use('/assets', express.static('dist/assets'));

  app.get('*', (req, res) => {
    // Serve index.html
    let template = fs
      .readFileSync(path.resolve(__dirname, 'dist/index.html'), 'utf-8')
      .replace('REPO_BASE_URI', repoBaseURI);

    if (req.session.username) {
      template = template.replace(
        'username = null',
        `username = '${req.session.username}'`,
      );
    }

    res
      .status(200)
      .set({
        'Content-Type': 'text/html',
        'Cache-Control':
          'no-store, no-cache, must-revalidate, proxy-revalidate',
        Expires: '0',
        'Surrogate-Control': 'no-store',
      })
      .end(template);
  });
} else {
  const viteImport = await import('vite');
  const { createServer, createViteRuntime } = viteImport;

  const __dirname = path.dirname(fileURLToPath(import.meta.url));
  const viteConfigPaths = await import('vite-tsconfig-paths');
  const vite = await createServer({
    server: {
      middlewareMode: true,
    },
    appType: 'custom',
    plugins: [viteConfigPaths.default()],
    build: {
      sourcemap: true,
    },
  });

  app.use(vite.middlewares);

  app.get('*', async (req, res, next) => {
    if (req.session.access_token) {
      console.log(`Access token: ${req.session.access_token}`);
    }
    const url = req.originalUrl;

    try {
      let template = fs
        .readFileSync(path.resolve(__dirname, 'index.html'), 'utf-8')
        .replace('REPO_BASE_URI', repoBaseURI);

      if (req.session.username) {
        template = template.replace(
          'username = null',
          `username = '${req.session.username}'`,
        );
      }

      template = await vite.transformIndexHtml(url, template);

      res.status(200).set({ 'Content-Type': 'text/html' }).end(template);
    } catch (e: any) {
      vite.ssrFixStacktrace(e);
      next(e);
    }
  });
}

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
