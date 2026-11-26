<script lang="ts">
    import Slide from './Slide.svelte';
    import type { Playlist } from '$lib/portfolio';
    import { Splide, SplideSlide } from '@splidejs/svelte-splide';
    import '@splidejs/svelte-splide/css';
    import '@splidejs/svelte-splide/css/sea-green';
    import '../css/splide.css';

    let { playlist }: { playlist: Playlist } = $props();
    const options = {
        type: 'loop',
        perPage: 6,
        drag: false,
        breakpoints: {
            500: { perPage: 2 },
            800: { perPage: 3 },
            1080: { perPage: 4 },
            1400: { perPage: 5 }
        }
    } as const;
</script>

<div class="playlist-wrapper">
    {#if playlist.films.length}
        <div class="title">{playlist.name}</div>
        <div class="videos-slider">
            <Splide aria-label="Playlist" {options}>
                {#each playlist.films as film (film.id)}
                    <SplideSlide>
                        <Slide data={film}/>
                    </SplideSlide>
                {/each}
            </Splide>
        </div>
    {/if}
</div>

<style>
    .title {
        font-size: 1.6em;
        font-weight: 300;
        margin-bottom: 0.5em;
        padding: 2em 2em 0 1em;
    }

    @media only screen and (max-width: 600px) {
        .title {
            font-size: 1em;
            margin-bottom: 0em;
        }
    }
</style>
