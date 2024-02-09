import uploadFeature, { LocalUploadOptions } from '@adminjs/upload';

import { componentLoader } from '../admin/component-loader.js';
import File from '../models/file.js';

const localProvider: LocalUploadOptions = {
  bucket: 'public/files',
  opts: {
    baseUrl: '/files',
  },
};

export const files = {
  resource: File,
  options: {
    properties: {
      link: {
        type: 'string',
      },
      title: {
        type: 'string',
      },
      description: {
        type: 'textarea',
      },
      published: {
        type: 'date',
      },
      playlist: {
        type: 'reference',
        reference: 'Playlist',
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
        maxSize: 524288000,
      },
      properties: {
        key: 'filepath',
      },
    }),
  ],
};
