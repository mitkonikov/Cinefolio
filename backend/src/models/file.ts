import mongoose from 'mongoose';

const schema = new mongoose.Schema({
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
