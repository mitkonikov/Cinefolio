import mongoose from 'mongoose';

const schema = new mongoose.Schema({
  bucket: String,
  filepath: String,
  title: String,
  description: String,
  published: Date,
});

const File = mongoose.model('File', schema);

export default File;
