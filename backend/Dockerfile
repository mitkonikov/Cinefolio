# Use an official Node.js runtime as a parent image
FROM node:20-alpine

# Install bash
RUN apk update
RUN apk add bash
RUN apk add curl

# Copy docker setup script
COPY ./start.sh /start.sh

# Copy package(-lock).json to the working directory
COPY ./package.json /package.json
COPY ./package-lock.json /package-lock.json
COPY ./tsconfig.json /tsconfig.json
COPY ./.env /.env
COPY ./.eslintrc.cjs /.eslintrc.cjs
COPY ./.prettierrc /.prettierrc

# Bundle your app source code into the container
COPY ./src /src

# Install app dependencies
WORKDIR /
RUN npm install .

# Copy custom feature
RUN mkdir /_features/
COPY feature_setup.sh /feature_setup.sh
COPY relations /_features/relations

# Define the command to run your app
CMD ["bash", "/start.sh"]
