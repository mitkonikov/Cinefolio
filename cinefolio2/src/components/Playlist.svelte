<script lang="ts">
    export let repositoryName;
    export let displayName = repositoryName;

    import { db } from '$lib/firebase';
	import { collection, getDocs, limit, orderBy, query } from 'firebase/firestore';
    import Slide from './../components/Slide.svelte';

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
    });
</script>

<div class="playlist-wrapper">
    <div class="title">{displayName}</div>
    <div class="videos-slider">
        <!-- <div class="left-arrow arrow-container">
            <div class="arrow">
                L
            </div>
        </div> -->
        <div class="videos-container">
            {#each videos as video}
                <Slide {...video} first={videos.at(0) == video}/>
            {/each}
        </div>
        <!-- <div class="right-arrow arrow-container">
            <div class="arrow">
                R
            </div>
        </div> -->
    </div>
</div>

<style>
    .title {
        font-size: 1.6em;
        font-weight: 300;
        margin-bottom: 0.5em;
    }

    .playlist-wrapper {
        padding: 2em 2em 0 2em;
        position: relative;
    }

    .videos-slider {
        position: relative;
        overflow-x: hidden;
        padding: 6em 0;
        margin: -6em 0;
    }

    .videos-container {
        height: 13em;
        position: relative;
        display: flex;
        align-items: center;
    }

    .arrow-container {
        display: flex;
        align-items: center;
        justify-content: center;
        background-color: black;
    }

    .left-arrow {
        width: 2em;
        height: 13em;
        position: absolute;
        z-index: 100000;
    }

    .right-arrow {
        width: 2em;
        height: 13em;
        position: absolute;
        right: 0;
        top: 6em;
    }

    .arrow {

    }
</style>