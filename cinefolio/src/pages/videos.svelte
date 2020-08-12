<script>
    import { collectionData } from 'rxfire/firestore';
    import { startWith } from 'rxjs/operators';
    import { db } from './../js/firebase';
    import Slide from './../components/Slide.svelte';

    const query = db.collection('Videos').orderBy('timestamp', 'desc').limit(5);
    const videos = collectionData(query).pipe(startWith([]));
</script>

<div id="videos-container-center">
    <div id="videos-container">
        {#each $videos as video}
            <Slide {...video}/>
        {/each}
    </div>
</div>

<style>
    #videos-container-center {
        position: relative;
        width: 100%;
        height: 100%;
        overflow-x: auto;
    }

    #videos-container {
        position: relative;
        margin: auto;
        width: fit-content;
        margin-top: 8em;
    }
</style>