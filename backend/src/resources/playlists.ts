import Playlist from '../models/playlist.js';

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
      },
      site: {
        type: 'reference',
        reference: 'Site',
      },
    },
  },
};
