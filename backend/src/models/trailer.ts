import mongoose from 'mongoose';

// Define Thumbnail interface
export interface ITrailer extends Document {
  filepath: string;
  visible: boolean;
  film: mongoose.Schema.Types.ObjectId;
}

const schema = new mongoose.Schema<ITrailer>({
  filepath: String,
  visible: Boolean,
  film: { type: mongoose.Schema.Types.ObjectId, ref: 'File' },
});

const Trailer = mongoose.model('Trailer', schema);

export default Trailer;
