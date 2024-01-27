<script lang="ts">
    export let repositoryName;
    export let displayName = repositoryName;

    import { db } from '$lib/firebase';
	import { collection, getDocs, limit, orderBy, query } from 'firebase/firestore';
    import Slide from './../components/Slide.svelte';

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

    $: splideInit = false;

    // @ts-ignore
    /**
	 * @type {any[]}
	 */
    let videos: string | any[] = [];

    const q = query(collection(db, repositoryName), orderBy('timestamp', 'desc'), limit(5));
    // @ts-ignore
    const querySnapshot = getDocs(q).then((snapshot) => {
        snapshot.forEach((doc) => {
            videos = [...videos, (doc.data())];
        });
        while (videos.length < 12) {
            videos = [...videos, ...videos];
        }
        while (videos.length > 12) {
            (videos as []).pop();
        }
        videos = [...videos];
        splideInit = true;
    });
</script>

<div class="playlist-wrapper">
    <div class="title">{displayName}</div>
    <div class="videos-slider">
        {#if splideInit}
            <Splide aria-label="Playlist" options={splideOptions}>
                {#each videos as video}
                    <SplideSlide>
                        <Slide data={video}/>
                    </SplideSlide>
                {/each}
            </Splide>
        {/if}
    </div>
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