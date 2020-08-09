<script>
    import { collectionData } from 'rxfire/firestore';
    import { startWith } from 'rxjs/operators';
    import { db } from './../js/firebase';
    import Slide from './../components/Slide.svelte';

    const query = db.collection('Videos').orderBy('timestamp', 'desc');
    const videos = collectionData(query).pipe(startWith([]));
</script>

<div class="center-vh">
    {#each $videos as video}
        <Slide {...video}/>
    {/each}
</div>