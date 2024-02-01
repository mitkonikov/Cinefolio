import React from 'react';

import ThumbnailCreator from './thumbnailCreator.js';

function Dashboard() {
  return (
    <div style={{ padding: '2em' }}>
      <h1 style={{ fontSize: '2em' }}>
        Welcome to the Admin Panel
      </h1>

      <div style={{ padding: '2em' }}>
        <ThumbnailCreator />
      </div>
    </div>
  );
}

export default Dashboard;
