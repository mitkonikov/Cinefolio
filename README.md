# Cinefolio

Personal portfolio website made using Svelte & Firebase. Hosted at [cinefolio.web.app](https://cinefolio.web.app/). The Portfolio is mostly meant to showcase videos, but also in the future we are planning a better support for showcasing coding projects using short videos and/or embedded demos.

## Developing

### First time setup

Run the following commands if you are setting up for the first time on your machine:

```sh
    cd cinefolio2
    npm install
```

### Development

Use `npm run dev` for development inside the `/cinefolio2` folder.

There is automatic GitHub actions for preview of pull request and deployment after merging.

### Backend

The backend is based on AdminJS connected with MongoDB.
To run the backend and frontend together using docker compose by running:

```sh
docker compose up --build
```

It is assumed that the database users are already created and that the `.env` file
in the `./backend` folder contains the proper credentials.

### Future plans & improvements

 - Box for more information
 - Refactor code
