<script>
    export let repositoryName;
    export let displayName = repositoryName;

    import { collectionData } from 'rxfire/firestore';
    import { startWith } from 'rxjs/operators';
    import { db } from './../js/firebase';
    import Slide from './../components/Slide.svelte';

    const query = db.collection(repositoryName).orderBy('timestamp', 'desc').limit(5);
    const videos = collectionData(query).pipe(startWith([]));
</script>

<div class="wrapper">
    <div class="wrapper-absolute">
        <div class="playlist-container">
            <div class="title">{displayName}</div>
            <div class="videos-container">
                <!-- <div class="right-arrow"></div> -->
                {#each $videos as video}
                    <Slide {...video} first={$videos.at(0) == video}/>
                {/each}
            </div>
        </div>
    </div>
</div>

<style>
    .title {
        position: absolute;
        top: 3em;
        left: 1em;
        font-size: 1.6em;
    }

    .wrapper {
        position: relative;
        height: 20em;
    }

    .wrapper-absolute {
        width: 98vw;
        position: absolute;
        transform: translate(0, -5em);
        padding: 0 2em;
    }

    .wrapper-absolute:hover {
        z-index: 100000;
    }
    
    .playlist-container {
        height: 30em;
        display: flex;
        align-items: center;
    }

    .videos-container {
        height: 13em;
        position: relative;
        display: flex;
        align-items: center;
    }

    .right-arrow {
        position: absolute;
        width: 5em;
        height: 100%;
        background-color: black;
        opacity: 0.5;
        z-index: 100001;
    }
</style>