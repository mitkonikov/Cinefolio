<script>
    import { collectionData } from 'rxfire/firestore';
    import { startWith } from 'rxjs/operators';
    import { db } from './../js/firebase';
    import Slide from './../components/Slide.svelte';
    import NameTag from '../components/NameTag.svelte';

    const query = db.collection('Videos').orderBy('timestamp', 'desc').limit(5);
    const videos = collectionData(query).pipe(startWith([]));
</script>

<div id="videos-container-center">
    <div id="name-tag-container">
        <NameTag />
    </div>
    <div id="videos-container">
        {#each $videos as video}
            <Slide {...video}/>
        {/each}
    </div>
</div>

<style>
    #name-tag-container {
        height: 8em;
    }

    #videos-container-center {
        position: relative;
        width: 100%;
        height: 100%;
        overflow-x: auto;
    }

    #videos-container {
        position: relative;
        margin-left: auto;
        margin-right: auto;
        width: fit-content;
    }
</style>