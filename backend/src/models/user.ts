import mongoose, { Document, Schema } from 'mongoose';
import bcrypt from 'bcrypt';

import { ISite } from './site.js';

// Define User interface
interface IUser extends Document {
  email: string;
  password: string;
  role: 'SUPER_ADMIN' | 'SITE_ADMIN'
  site?: ISite;
}

// Define User schema
const userSchema = new Schema<IUser>({
  email: {
    type: String,
    unique: true,
    required: true,
  },
  password: {
    type: String,
    required: true,
  },
  role: {
    type: String,
    required: true,
  },
  site: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Site',
  },
});

// Hash the user's password before saving
// eslint-disable-next-line @typescript-eslint/no-explicit-any, func-names
userSchema.pre<IUser>('save', async function (next) {
  // eslint-disable-next-line @typescript-eslint/no-this-alias
  const user: IUser = this;

  if (user && (user.isModified('password') || user.isNew)) {
    const salt = await bcrypt.genSalt(10);
    const hash = await bcrypt.hash(user.password, salt);
    user.password = hash;
  }
  return next();
});

// Create a model
const User = mongoose.model<IUser>('User', userSchema);

// Function to search for a user by email and password
export async function findUser(email: string, password: string): Promise<IUser | null> {
  const user = await User.findOne({ email });
  if (!user) {
    return null;
    // throw new Error('User not found');
  }

  const isMatch = await bcrypt.compare(password, user.password);
  if (!isMatch) {
    return null;
    // throw new Error('Incorrect password');
  }

  return user;
}

// Function to add a new user
export async function addUser(email: string, password: string): Promise<IUser> {
  const newUser = new User({ email, password });
  await newUser.save();
  return newUser;
}
