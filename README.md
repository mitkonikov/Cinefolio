# Cinefolio
Portfolio website using Firebase

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

 ### Experiencing problems with building Svelte?

 A common problem with our command is that it expects a `build` directory in the `public` directory.