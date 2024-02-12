import mongoose from 'mongoose';

// Define Thumbnail interface
export interface IThumbnail extends Document {
  title: string;
  filepath: string;
  crop: boolean;
}

const schema = new mongoose.Schema<IThumbnail>({
  title: String,
  filepath: String,
  crop: Boolean,
});

const Thumbnail = mongoose.model('Thumbnail', schema);

export default Thumbnail;
