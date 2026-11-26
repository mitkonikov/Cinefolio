<script lang="ts">
    import { onMount } from 'svelte';
    import Playlist from '../components/Playlist.svelte';
    import Trailer from '../components/Trailer.svelte';
    import { request } from '$lib/api';
    import type { Portfolio } from '$lib/portfolio';

    let portfolio: Portfolio | undefined = $state();
    let failed = $state(false);

    onMount(() => {
        request<Portfolio>('/api/public/portfolio')
            .then((data) => { portfolio = data; })
            .catch(() => { failed = true; });
    });
</script>

<main>
    <Trailer trailer={portfolio?.trailer ?? null} />

    <div id="videos-container-center">
        {#each portfolio?.playlists ?? [] as playlist (playlist.id)}
            <Playlist {playlist} />
        {/each}
        {#if failed}
            <p class="status">The portfolio is temporarily unavailable.</p>
        {:else if portfolio && !portfolio.playlists.some((playlist) => playlist.films.length)}
            <p class="status">There are no films to display yet.</p>
        {/if}
    </div>
</main>

<style>
    .status { padding: 1rem; }
    main {
        position: relative;
    }

    #videos-container-center {
        position: relative;
        width: 100%;
        padding-top: calc(100vw / 16 * 9 - 100vw / 16 * 9 / 10);
    }
</style>
