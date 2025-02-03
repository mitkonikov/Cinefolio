import mongoose from 'mongoose';
import { Database, Resource } from '@adminjs/mongoose';
import AdminJS from 'adminjs';

AdminJS.registerAdapter({ Database, Resource });

const initialize = async () => {
  try {
    const db = await mongoose.connect(process.env.DATABASE_URL as string, {
      user: process.env.DATABASE_USER,
      pass: process.env.DATABASE_PASSWORD,
      authSource: process.env.DATABASE_AUTH_SOURCE,
    });

    // Check connection status and report it
    if (db.connection.readyState === 1) {
      console.log('Database connected successfully.');
    } else {
      console.error('Error connecting to database: ', db.connection.readyState);
    }

    return { db };
  } catch (error) {
    console.error('Failed to connect to the database:', error);
    throw new Error('Database connection failed');
  }
};

export default initialize;
