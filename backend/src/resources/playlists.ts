import Playlist from '../models/playlist.js';
import { FilePlaylist } from '../features/mtm.js';
import { componentLoader } from '../admin/component-loader.js';

export const playlists = {
  resource: Playlist,
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
      name: {
        type: 'string',
        isRequired: true,
      },
      site: {
        type: 'reference',
        reference: 'Site',
      },
      visible: {
        type: 'boolean',
      },
      order: {
        type: 'number',
        isRequired: true,
      },
    },
    sort: {
      direction: 'asc',
      sortBy: 'order',
    },
  },
  features: [
    FilePlaylist(componentLoader),
  ],
};
