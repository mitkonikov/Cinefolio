<script lang="ts">
    // import { collectionData } from 'rxfire/firestore';
    // import { startWith } from 'rxjs/operators';
    // import { db } from '$lib/firebase';
	// import { collection, getDocs, limit, query } from 'firebase/firestore';
    import Playlist from './../components/Playlist.svelte';
</script>

<main>
	<div class="bg-img" />

    <div class="noselect" id="logo">
        <div id="main-logo">cinefolio</div>
        <div class="logo-dec"></div>
    </div>
	
    <div id="videos-container-center">
        <Playlist repositoryName={"Videos"} displayName={"Latest"}/>
        <Playlist repositoryName={"Strumica Open Festival"}/>
    </div>
</main>

<style>
	.bg-img {
		overflow: hidden;
		position: absolute;
		width: 100vw;
		height: 100vh;
		background-image: url('/img/Stars.jpg');
		background-repeat: no-repeat;
		background-size: cover;
		opacity: 0.3;
        z-index: -1;
	}

    #logo {
        width: 100%;
        text-align: center;
        margin: 1em 0;
    }

    #main-logo {
        display: inline-block;
        padding: 0.5em 0 0 0;
        font-size: 6em;
        font-weight: 100;
        transition: 0.2s;
        color: #eee;
    }

    #main-logo:hover {
        padding: 0.6em 0 0.1em 0;
        cursor: pointer;
        color: white;
    }

    .logo-dec {
        width: 20em;
        height: 0.15em;
        background-color: rgba(239, 239, 239, 0.517);
        margin: auto;
    }
    
    #videos-container-center {
        position: relative;
        padding-top: 2em;
        width: 100%;
        height: 100%;
    }

    @media only screen and (max-width: 600px) {
        #logo {
            font-size: 0.4em;
        }
    }
</style>