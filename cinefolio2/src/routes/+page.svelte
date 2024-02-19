<script lang="ts">
	import { onMount } from 'svelte';
    import Playlist from './../components/Playlist.svelte';
    import { type IPlaylist } from './../types/playlist';
	import { API_URL } from '$lib/constants';
	import Trailer from '../components/Trailer.svelte';

    let playlists: IPlaylist[] = [];

    onMount(async () => {
        fetch(`${API_URL}/query/playlists/Cinefolio`, {
            method: 'GET',
            mode: 'cors'
        }).then((response) => response.json())
        .then((value) => {
            playlists = value.playlists;
            playlists.sort((a, b) => {
                return a.order - b.order;
            });
        });
    });
</script>

<main>
    <Trailer/>
	
    <div id="videos-container-center">
        {#each playlists as playlist}
            {#if playlist.visible}
                <Playlist playlist={playlist}/>
            {/if}
        {/each}
    </div>
</main>

<style>
    main {
        position: relative;
    }
    
    #videos-container-center {
        position: relative;
        width: 100%;
        padding-top: calc(100vw / 16 * 9 - 100vw / 16 * 9 / 10);
    }
</style>