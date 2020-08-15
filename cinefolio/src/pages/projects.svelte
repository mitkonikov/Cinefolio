<script>
    import { collectionData } from 'rxfire/firestore';
    import { startWith } from 'rxjs/operators';
    import { db } from './../js/firebase';
    import Slide from './../components/Slide.svelte';
    import NameTag from '../components/NameTag.svelte';

    const query = db.collection('Projects').orderBy('rank', 'desc').limit(5);
    const projects = collectionData(query).pipe(startWith([]));
</script>

<div id="project-container-center">
    <div id="name-tag-container">
        <NameTag />
    </div>
    
    <div id="project-container">
        {#each $projects as project}
            <Slide {...project}/>
        {/each}
    </div>
</div>

<style>
    #name-tag-container {
        height: 8em;
    }

    #project-container-center {
        position: relative;
        width: 100%;
        height: 100%;
        overflow-x: auto;
    }

    #project-container {
        position: relative;
        margin-left: auto;
        margin-right: auto;
        width: fit-content;
    }
</style>