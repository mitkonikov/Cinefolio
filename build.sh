echo " > Running docker-compose down..."
docker-compose down
echo " > docker-compose down finished."

echo " > Copy .env"
cp ~/env/backend.env ~/cine/backend/.env
cp ~/env/frontend.env ~/cine/cinefolio2/.env

echo " > Building Svelte..."
cd cinefolio2 && bash build.sh

# Start the Docker containers
echo " > Running docker-compose up..."
docker-compose up --build -d
