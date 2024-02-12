import uploadFeature from '@adminjs/upload';

import { componentLoader } from '../admin/component-loader.js';
import Thumbnail from '../models/thumbnail.js';
import UploadProvider from '../admin/upload-provider.js';

export const thumbnails = {
  resource: Thumbnail,
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
      crop: {
        type: 'boolean',
      },
      filepath: {
        isVisible: {
          list: false,
          filter: false,
          show: true,
          edit: false,
        },
      },
    },
    actions: {
      EnableCrop: {
        actionType: 'bulk',
        component: false,
        handler: (request, response, context) => {
          const { records, currentAdmin } = context;
          // eslint-disable-next-line no-restricted-syntax, prefer-const
          for (let record of records) {
            const recordId = record.params._id as number;
            Thumbnail.findByIdAndUpdate(recordId, { crop: true }).exec();
          }
          return { records };
        },
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
