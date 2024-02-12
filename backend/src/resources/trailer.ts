import uploadFeature from '@adminjs/upload';

import { componentLoader } from '../admin/component-loader.js';
import Trailer from '../models/trailer.js';
import UploadProvider from '../admin/upload-provider.js';

export const trailers = {
  resource: Trailer,
  options: {
    properties: {
      _id: {
        isVisible: {
          list: false,
          filter: false,
          show: true,
          edit: false,
        },
      },
      visible: {
        type: 'boolean',
      },
      film: {
        type: 'reference',
        reference: 'File',
      },
    },
  },
  features: [
    uploadFeature({
      componentLoader,
      provider: new UploadProvider('public/trailers', {
        baseUrl: '/trailers',
      }),
      validation: {
        mimeTypes: [
          'video/mp4',
          'video/x-msvideo',
          'video/mpeg',
          'video/ogg',
          'video/mp2t',
          'video/webm',
          'video/3gpp',
          'video/3gpp2',
        ],
        maxSize: 15728640,
      },
      properties: {
        key: 'filepath',
      },
    }),
  ],
};
