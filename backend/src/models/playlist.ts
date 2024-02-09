import { Document, Schema, model } from 'mongoose';

// Define Playlist interface
export interface IPlaylist extends Document {
  name: string;
  site: Schema.Types.ObjectId;
}

// Define playlist schema
export const playlistSchema = new Schema<IPlaylist>({
  name: {
    type: String,
    unique: true,
    required: true,
  },
  site: { type: Schema.Types.ObjectId, ref: 'Site' },
});

const Playlist = model<IPlaylist>('Playlist', playlistSchema);
export default Playlist;
