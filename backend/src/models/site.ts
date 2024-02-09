import { Document, Schema, model } from 'mongoose';

// Define Site interface
export interface ISite extends Document {
  name: string;
}

// Define site schema
export const siteSchema = new Schema<ISite>({
  name: {
    type: String,
    unique: true,
    required: true,
  },
});

const Site = model<ISite>('Site', siteSchema);
export default Site;
