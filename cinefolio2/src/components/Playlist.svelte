<script lang="ts">
    export let playlist: IPlaylist;

	import { onMount } from 'svelte';
    import Slide from './../components/Slide.svelte';
	import type { IFilm } from '../types/film';
	import type { IPlaylist } from '../types/playlist';
	import { API_URL } from '$lib/constants';

    import { Splide, SplideSlide } from '@splidejs/svelte-splide';
    import '@splidejs/svelte-splide/css';
    import '@splidejs/svelte-splide/css/sea-green';
    import './../css/splide.css';

    $: splideOptions = {
        type: 'loop',
        perPage: 6,
        drag: false,
        breakpoints: {
            500: {
                perPage: 2,
            },
            800: {
                perPage: 3,
            },
            1080: {
                perPage: 4,
            },
            1400: {
                perPage: 5,
            },
        }
    };

    let videos: IFilm[] = [];
    $: splideInit = false;

    onMount(async () => {
        fetch(`${API_URL}/query/films/${playlist.name}`, {
            method: 'GET',
            mode: 'cors'
        }).then((response) => response.json())
        .then((value) => {
            videos = value.files;
            if (videos.length > 0) {
                splideInit = true;
            }

            videos.sort((a: IFilm, b: IFilm) => {
                let dateA = new Date(a.published);
                let dateB = new Date(b.published);
                return Number(dateB) - Number(dateA);
            });
        });
    });
</script>

<div class="playlist-wrapper">
    {#if splideInit}
        <div class="title">{playlist.name}</div>
        <div class="videos-slider">
            <Splide aria-label="Playlist" options={splideOptions}>
                {#each videos as video}
                    <SplideSlide>
                        <Slide data={video}/>
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
