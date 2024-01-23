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
    });
</script>

<div class="playlist-wrapper">
    <div class="title">{displayName}</div>
    <div class="videos-slider">
        <Splide aria-label="Playlist" options={splideOptions}>
            {#each videos as video}
                <SplideSlide>
                    <Slide {...video} first={videos.at(0) == video}/>
                </SplideSlide>
            {/each}
        </Splide>
    </div>
</div>

<style>
    .title {
        font-size: 1.6em;
        font-weight: 300;
        margin-bottom: 0.5em;
        padding: 2em 2em 0 1em;
    }

    .playlist-wrapper {
        position: relative;
    }

    .videos-slider {
        position: relative;
        overflow-x: hidden;
        padding: 6em 0;
        margin: -6em 0;
    }
</style>