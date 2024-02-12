<script lang="ts">
	import { onMount } from 'svelte';
    import Playlist from './../components/Playlist.svelte';
    import { type IPlaylist } from './../types/playlist';
	import { API_URL } from '$lib/constants';

    let playlists: IPlaylist[] = [];
    $: playlistsLoaded = false;

    onMount(async () => {
        fetch(`${API_URL}/query/playlists/Cinefolio`, {
            method: 'GET',
            mode: 'cors'
        }).then((response) => response.json())
        .then((value) => {
            playlists = value.playlists;
            playlistsLoaded = true;
        });
    });
</script>

<main>
    <div class="noselect" id="logo">
        <div id="main-logo">cinefolio</div>
        <div class="logo-dec"></div>
    </div>
	
    <div id="videos-container-center">
        {#each playlists as playlist}
            <Playlist playlist={playlist}/>
        {/each}
    </div>
</main>

<style>
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
    }

    @media only screen and (max-width: 600px) {
        #logo {
            font-size: 0.4em;
        }
    }
</style>