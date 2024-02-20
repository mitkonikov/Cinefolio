import { Document, Schema, model } from 'mongoose';

// Define File to Playlist interface
export interface IFileToPlaylist extends Document {
  fileId: Schema.Types.ObjectId;
  playlistId: Schema.Types.ObjectId;
}

// Define file to playlist schema
export const filePlaylistSchema = new Schema<IFileToPlaylist>({
  fileId: {
    type: Schema.Types.ObjectId,
    ref: 'File',
  },
  playlistId: {
    type: Schema.Types.ObjectId,
    ref: 'Playlist',
  },
});

const FileToPlaylist = model<IFileToPlaylist>('FileToPlaylist', filePlaylistSchema);
export default FileToPlaylist;
