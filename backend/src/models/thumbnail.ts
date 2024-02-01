import mongoose from 'mongoose';

const schema = new mongoose.Schema({
  title: String,
  filepath: String,
});

const Thumbnail = mongoose.model('Thumbnail', schema);

export default Thumbnail;
