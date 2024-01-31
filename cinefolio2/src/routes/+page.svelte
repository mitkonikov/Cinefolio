<script lang="ts">
	import { onMount } from 'svelte';
    import Playlist from './../components/Playlist.svelte';
    import { type IPlaylist } from './../types/playlist';
    import { type ITrailer } from './../types/film';
	import { API_URL, STATIC_URL } from '$lib/constants';

    let playlists: IPlaylist[] = [];
    let trailers: ITrailer[] = [];
    $: playlistsLoaded = false;
    let paused = true;

    onMount(async () => {
        fetch(`${API_URL}/query/playlists/Cinefolio`, {
            method: 'GET',
            mode: 'cors'
        }).then((response) => response.json())
        .then((value) => {
            playlists = value.playlists;
            playlistsLoaded = true;
        });

        fetch(`${API_URL}/query/trailers`, {
            method: 'GET',
            mode: 'cors'
        }).then((response) => response.json())
        .then((value) => {
            trailers = value.trailers;
            // paused = false;
        });
    });

    function getTrailerPath() {
        return `${STATIC_URL}/trailers/${trailers[0].filepath}`;
    }
</script>

<main>
    <div class="trailer">
        {#if !paused}
            <!-- svelte-ignore a11y-media-has-caption -->
            <video width="100%" autoplay muted bind:paused>
                <source src={getTrailerPath()} type="video/mp4">
                Your browser does not support the video tag.
            </video>
        {:else}
            <div class="frozen-trailer"></div>
        {/if}
    </div>

    <div class="fade-decoration"></div>

    <div class="noselect" id="logo">
        <div id="main-logo">cinefolio</div>
        <div class="logo-dec"></div>
        <div id="info-logo">Mitko Nikov's cinematography portfolio</div>
    </div>
	
    <div id="videos-container-center">
        {#each playlists as playlist}
            <Playlist playlist={playlist}/>
        {/each}
    </div>
</main>

<style>
    main {
        position: relative;
    }

    .fade-decoration {
        position: absolute;
        width: 100%;
        margin-top: calc(100vw / 16 * 9 / 1.4);
        height: calc(100vw / 16 * 9 - 100vw / 16 * 9 / 1.4);
        background: -moz-linear-gradient(top, rgba(20,20,20,0), rgba(20,20,20, 1));
        background: -webkit-linear-gradient(top, rgba(20,20,20,0), rgba(20,20,20, 1));
    }

    #logo {
        position: absolute;
        width: 100%;
        padding-top: calc(100vw / 16 * 9 / 4);
        padding-left: 4em;
    }

    #main-logo {
        display: inline-block;
        padding: 0.5em 0 0 0;
        font-size: 10vw;
        font-weight: 100;
        color: #eee;
        text-shadow: 0.03em 0.03em 0 black;
    }

    #info-logo {
        display: inline-block;
        padding: 0.5em 0 0 0;
        font-size: 1vw;
        font-weight: 100;
        color: #eee;
        text-shadow: 0.03em 0.03em 0 black;
    }

    #main-logo:hover {
        padding: 0.6em 0 0.1em 0;
        cursor: pointer;
        color: white;
    }

    .trailer {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 0;
    }

    .logo-dec {
        width: calc(100vw / 16 * 9 / 2);
        height: 0.15em;
        background-color: rgba(239, 239, 239, 0.895);
    }

    .frozen-trailer {
        width: 100%;
        height: calc(100vw / 16 * 9);
        background-image: url('/img/Stars.jpg');
        background-repeat: no-repeat;
        background-size: cover;
    }
    
    #videos-container-center {
        position: relative;
        width: 100%;
        padding-top: calc(100vw / 16 * 9 - 100vw / 16 * 9 / 10);
    }

    @media only screen and (max-width: 600px) {
        #logo {
            font-size: 0.4em;
        }

        #info-logo {
            font-size: 2em;
        }
    }
</style>