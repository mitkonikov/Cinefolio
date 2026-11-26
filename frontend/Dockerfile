# Use a lightweight Nginx image
FROM nginx:alpine

RUN apk update
RUN apk add curl

# Copy your static HTML files to the Nginx document root
COPY ./build /usr/share/nginx/html

# Expose port 80 to the outside world
EXPOSE 80

# Start Nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
