# Cinefolio

Personal portfolio website made using Svelte & Firebase. Hosted at [cinefolio.web.app](https://cinefolio.web.app/). The Portfolio is meant to showcase short films and videos.

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
To run the backend and frontend together you need a MongoDB instance
running in the background with the proper `./backend/.env` config to connect to the instance.
Afterwards, you can run docker compose to start the backend server:

```sh
docker compose up --build
```

It is assumed that the database users are already created and that the `.env` file
in the `./backend` folder contains the proper credentials.

### Future plans & improvements

 - Box for more information
 - Refactor code
