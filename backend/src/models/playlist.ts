import { Document, Schema, model } from 'mongoose';

// Define Playlist interface
export interface IPlaylist extends Document {
  name: string;
}

// Define playlist schema
export const playlistSchema = new Schema<IPlaylist>({
  name: {
    type: String,
    unique: true,
    required: true,
  },
});

const Playlist = model<IPlaylist>('Playlist', playlistSchema);
export default Playlist;
