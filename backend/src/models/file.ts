import mongoose from 'mongoose';

// Define File interface
export interface IFile extends Document {
  filepath: string;
  link: string;
  title: string;
  description: string;
  published: Date;
  playlist: mongoose.Schema.Types.ObjectId,
  thumbnail: mongoose.Schema.Types.ObjectId,
}

const schema = new mongoose.Schema<IFile>({
  filepath: String,
  link: String,
  title: String,
  description: String,
  published: Date,
  playlist: { type: mongoose.Schema.Types.ObjectId, ref: 'Playlist' },
  thumbnail: { type: mongoose.Schema.Types.ObjectId, ref: 'Thumbnail' },
});

const File = mongoose.model('File', schema);

export default File;
