# Cinefolio
A Personal Portfolio Website made using Svelte & Firebase. Hosted at [javascripter.web.app](https://javascripter.web.app/). The Portfolio is mostly meant to showcase videos, but also in the future we are planning a better support for showcasing coding projects using short videos and/or embedded demos.

<br>
<br>

## Developing

---

### First time setup

Run the following commands if you are setting up for the first time on your machine:

```sh
    cd cinefolio
    npm install
    cd public
    mkdir build
```
<br>

### Building and uploading the files to the cloud

 - `npm run build` (inside ./cinefolio) - to build the server
 - `firebase deploy` - (./) To deploy the public (build) folder to the cloud

### Future plans & improvements

 - Converting the Javascript code to Typescript
 - Creating filters
 - Connecting real-life demos of the projects

### Experiencing problems with building Svelte?

A common problem with our command is that it expects a `build` directory in the `public` directory.