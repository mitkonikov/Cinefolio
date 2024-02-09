import Playlist from '../models/playlist.js';

export const playlists = {
  resource: Playlist,
  options: {
    properties: {
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
