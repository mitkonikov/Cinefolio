import uploadFeature, { LocalUploadOptions } from '@adminjs/upload';

import { componentLoader } from '../admin/component-loader.js';
import Trailer from '../models/trailer.js';

const localProvider: LocalUploadOptions = {
  bucket: 'public/trailers',
  opts: {
    baseUrl: '/trailers',
  },
};

export const trailers = {
  resource: Trailer,
  options: {
    properties: {
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
      provider: { local: localProvider },
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
