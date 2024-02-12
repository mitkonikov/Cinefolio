import Site from '../models/site.js';

export const sites = {
  resource: Site,
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
    },
  },
};
