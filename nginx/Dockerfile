FROM nginx:alpine

# Remove the default Nginx configuration file
RUN rm -f /etc/nginx/conf.d/default.conf

# Copy your static HTML files to the Nginx document root
COPY ./nginx.conf /etc/nginx/nginx.conf

# Expose port 80 to the outside world
EXPOSE 80

# Start Nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
