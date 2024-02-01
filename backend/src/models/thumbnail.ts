import mongoose from 'mongoose';

const schema = new mongoose.Schema({
  filepath: String,
  film: { type: mongoose.Schema.Types.ObjectId, ref: 'File' },
});

const Thumbnail = mongoose.model('Thumbnail', schema);

export default Thumbnail;
