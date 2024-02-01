import path from 'path';
import { fileURLToPath } from 'url';

import { ComponentLoader } from 'adminjs';

const componentLoader = new ComponentLoader();

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const components = {
  Thumbnail: componentLoader.add('Thumbnail', path.resolve(__dirname, './../components/thumbnailCreator')),
};

export { components, componentLoader };
