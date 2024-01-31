import express from 'express';

import File from '../models/file.js';
import Trailer from '../models/trailer.js';
import Playlist from '../models/playlist.js';
import Site from '../models/site.js';

const router = express.Router();

router.get('/playlists/:site', async (req, res) => {
  try {
    const { site } = req.params;

    const siteDoc = await Site.findOne({ name: site });

    if (!siteDoc) {
      console.log('No site found by the name ', site);
      res.json({ error: 'No site found by the name.' });
      return;
    }

    const playlists = await Playlist.find({ site: siteDoc._id })
      .sort({ published: -1 }) // Sort in descending order based on the 'published' field
      .limit(10); // Limit the result to the last 10 files

    res.json({ playlists });
  } catch (exception) {
    console.log(exception);
    res.json({ error: 'Internal server error' });
  }
});

router.get('/films/:playlist', async (req, res) => {
  try {
    const { playlist } = req.params;

    const playlistDoc = await Playlist.findOne({ name: playlist });

    if (!playlistDoc) {
      console.log('No playlist found by the name ', playlist);
      res.json({ error: 'No playlist found by the name.' });
      return;
    }

    const files = await File.find({ playlist: playlistDoc._id }).populate('playlist').populate('thumbnail');

    res.json({ files });
  } catch (exception) {
    console.log(exception);
    res.json({ error: 'Internal server error' });
  }
});

router.get('/trailers', async (req, res) => {
  try {
    const trailers = await Trailer.find({ visible: true }).populate('film').populate('thumbnail');

    res.json({ trailers });
  } catch (exception) {
    console.log(exception);
    res.json({ error: 'Internal server error' });
  }
});

export default router;
