import * as url from 'url';
import path from 'path';

import express from 'express';
import AdminJS from 'adminjs';
import { buildAuthenticatedRouter } from '@adminjs/express';
import bodyParser from 'body-parser';

import provider from './admin/auth-provider.js';
import options from './admin/options.js';
import initializeDb from './db/index.js';
import QueryAPI from './routers/query.js';

const port = process.env.PORT || 3000;

const jsonParser = bodyParser.json();

const start = async () => {
  const app = express();

  await initializeDb();

  const admin = new AdminJS(options);

  if (process.env.NODE_ENV === 'production') {
    await admin.initialize();
  } else {
    admin.watch();
  }

  const router = buildAuthenticatedRouter(
    admin,
    {
      cookiePassword: process.env.COOKIE_SECRET,
      cookieName: 'adminjs',
      provider,
    },
    null,
    {
      secret: process.env.COOKIE_SECRET,
      saveUninitialized: true,
      resave: true,
    },
  );

  app.use(admin.options.rootPath, router);
  app.use('/cinefolio/api/query', jsonParser, QueryAPI);

  const __dirname = url.fileURLToPath(new URL('.', import.meta.url));
  const publicDir = path.join(__dirname, '../public');
  console.log('Public Directory: ', publicDir);
  app.use('/cinefolio/public', express.static(publicDir));

  app.listen(port, () => {
    console.log(`AdminJS available at http://localhost:${port}${admin.options.rootPath}`);
  });
};

start();
