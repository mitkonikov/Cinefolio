import Site from '../models/site.js';

export const sites = {
  resource: Site,
  options: {
    properties: {
      name: {
        type: 'string',
      },
    },
  },
};
