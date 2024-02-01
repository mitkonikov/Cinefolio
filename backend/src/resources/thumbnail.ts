import uploadFeature, { LocalUploadOptions } from '@adminjs/upload';

import { componentLoader } from '../admin/component-loader.js';
import Thumbnail from '../models/thumbnail.js';

const localProvider: LocalUploadOptions = {
  bucket: 'public/thumbs',
  opts: {
    baseUrl: '/thumbs',
  },
};

export const thumbnails = {
  resource: Thumbnail,
  options: {},
  features: [
    uploadFeature({
      componentLoader,
      provider: { local: localProvider },
      validation: {
        mimeTypes: [
          'image/jpeg',
          'image/jpg',
        ],
        maxSize: 524288,
      },
      properties: {
        key: 'filepath',
      },
    }),
  ],
};
