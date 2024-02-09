FROM mongo:latest

# Setup env variables
ENV MONGO_INITDB_ROOT_USERNAME=${MONGO_INITDB_ROOT_USERNAME}
ENV MONGO_INITDB_ROOT_PASSWORD=${MONGO_INITDB_ROOT_PASSWORD}

# Copy the MongoDB configuration file into the container
COPY mongod.conf /etc/mongod.conf

# Start MongoDB with authentication enabled
CMD ["mongod", "--config", "/etc/mongod.conf"]
