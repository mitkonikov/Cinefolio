echo " > Running docker compose down..."
docker compose down
echo " > docker compose down finished."

echo " > Copy .env"
cp ~/env/backend.env ~/cine/backend/.env

# Start the Docker containers
echo " > Running docker compose up..."
docker compose up --build -d
