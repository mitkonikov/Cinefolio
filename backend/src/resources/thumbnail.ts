import uploadFeature from '@adminjs/upload';

import { componentLoader } from '../admin/component-loader.js';
import Thumbnail from '../models/thumbnail.js';
import UploadProvider from '../admin/upload-provider.js';

export const thumbnails = {
  resource: Thumbnail,
  options: {
    properties: {
      crop: {
        type: 'boolean',
      },
    },
  },
  features: [
    uploadFeature({
      componentLoader,
      provider: new UploadProvider('public/thumbs', {
        baseUrl: '/thumbs',
      }),
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
