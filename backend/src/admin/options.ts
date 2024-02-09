import path from 'path';
import { fileURLToPath } from 'url';

import { AdminJSOptions } from 'adminjs';

import { files } from '../resources/files.js';
import { playlists } from '../resources/playlists.js';
import { thumbnails } from '../resources/thumbnail.js';
import { trailers } from '../resources/trailer.js';
import { sites } from '../resources/site.js';

import { componentLoader } from './component-loader.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const options: AdminJSOptions = {
  componentLoader,
  rootPath: '/cinefolio/admin',
  loginPath: '/cinefolio/admin/login',
  logoutPath: '/cinefolio/admin/logout',
  resources: [
    files,
    playlists,
    thumbnails,
    trailers,
    sites,
  ],
  databases: [],
  dashboard: {
    component: componentLoader.add('Dashboard', path.resolve(__dirname, './../components/dashboard')),
  },
  branding: {
    companyName: 'Cinefolio',
  },
};

export default options;
