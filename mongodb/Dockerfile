FROM mongo:latest

# Setup env variables
ENV MONGO_INITDB_ROOT_USERNAME=${MONGO_INITDB_ROOT_USERNAME}
ENV MONGO_INITDB_ROOT_PASSWORD=${MONGO_INITDB_ROOT_PASSWORD}

# Copy the MongoDB configuration file into the container
COPY mongod.conf /etc/mongod.conf

# Copy and set permissions for the entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Use the entrypoint script to handle repairs before starting MongoDB
ENTRYPOINT ["/entrypoint.sh"]
