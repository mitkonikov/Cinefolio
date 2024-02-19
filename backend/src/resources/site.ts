import { ResourceOptions } from 'adminjs';

import Site from '../models/site.js';
import { SuperUserAction } from '../admin/privileges.js';

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
        isRequired: true,
      },
    },
    actions: SuperUserAction,
  } as ResourceOptions,
};
