<script>
    import { collectionData } from 'rxfire/firestore';
    import { startWith } from 'rxjs/operators';
    import { db } from './../js/firebase';
    import Slide from './../components/Slide.svelte';

    const query = db.collection('Projects').orderBy('rank', 'desc');
    const projects = collectionData(query).pipe(startWith([]));
</script>

<div class="center-vh">
    {#each $projects as project}
        <Slide {...project}/>
    {/each}
</div>