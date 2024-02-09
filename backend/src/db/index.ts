import mongoose from 'mongoose';
import { Database, Resource } from '@adminjs/mongoose';
import AdminJS from 'adminjs';

AdminJS.registerAdapter({ Database, Resource });

const initialize = async () => {
  const db = await mongoose.connect(process.env.DATABASE_URL as string, {
    user: process.env.DATABASE_USER,
    pass: process.env.DATABASE_PASSWORD,
    authSource: process.env.DATABASE_AUTH_SOURCE,
  });

  return { db };
};

export default initialize;
