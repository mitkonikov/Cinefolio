import mongoose from 'mongoose';

const schema = new mongoose.Schema({
  bucket: String,
  filepath: String,
  title: String,
  description: String,
  published: Date,
  playlist: { type: mongoose.Schema.Types.ObjectId, ref: 'Playlist' },
});

const File = mongoose.model('File', schema);

export default File;
