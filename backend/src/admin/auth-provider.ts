import { DefaultAuthProvider } from 'adminjs';

import { findUser } from '../models/user.js';

import componentLoader from './component-loader.js';

const provider = new DefaultAuthProvider({
  componentLoader,
  authenticate: async ({ email, password }) => {
    const user = await findUser(email, password);
    return user;
  },
});

export default provider;
