import mongoose from 'mongoose';

// Define Thumbnail interface
export interface IThumbnail extends Document {
  title: string;
  filepath: string;
}

const schema = new mongoose.Schema<IThumbnail>({
  title: String,
  filepath: String,
});

const Thumbnail = mongoose.model('Thumbnail', schema);

export default Thumbnail;
