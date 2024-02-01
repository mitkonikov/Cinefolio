import mongoose from 'mongoose';

const schema = new mongoose.Schema({
  filepath: String,
  film: { type: mongoose.Schema.Types.ObjectId, ref: 'File' },
});

const Trailer = mongoose.model('Trailer', schema);

export default Trailer;
