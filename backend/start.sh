if [ ! -d "public/files" ]; then
  echo "Files folder doesn't exist. Creating it..."
  mkdir -p public/files
fi

if [ ! -d "public/thumbs" ]; then
  echo "Thumbs folder doesn't exist. Creating it..."
  mkdir -p public/thumbs
fi

if [ ! -d "public/trailers" ]; then
  echo "Trailers folder doesn't exist. Creating it..."
  mkdir -p public/trailers
fi

echo "Finished setup."

npm run start
