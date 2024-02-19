import { Document, Schema, model } from 'mongoose';

// Define Playlist interface
export interface IPlaylist extends Document {
  name: string;
  site: Schema.Types.ObjectId;
  visible: boolean;
  order: number;
}

// Define playlist schema
export const playlistSchema = new Schema<IPlaylist>({
  name: {
    type: String,
    unique: true,
    required: true,
  },
  site: {
    type: Schema.Types.ObjectId,
    ref: 'Site',
  },
  visible: {
    type: Boolean,
    default: true,
  },
  order: {
    type: Number,
    required: true,
  },
});

const Playlist = model<IPlaylist>('Playlist', playlistSchema);
export default Playlist;
