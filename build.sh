echo " > Running docker compose down..."
docker compose down
echo " > docker compose down finished."

echo " > Copy .env"
cp ~/env/backend.env ~/cine/backend/.env

echo " > Copy additional scripts..."
cp ~/env/feature_setup.sh ~/cine/backend/feature_setup.sh
cp -R ~/env/relations ~/cine/backend/relations
cp -R ~/env/src/features ~/cine/backend/src/features

# Start the Docker containers
echo " > Running docker compose up..."
docker compose up --build -d
