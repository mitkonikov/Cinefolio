<script>
    import { collectionData } from 'rxfire/firestore';
    import { startWith } from 'rxjs/operators';
    import { db } from './../js/firebase';
    import Slide from './../components/Slide.svelte';

    const query = db.collection('Videos').orderBy('timestamp', 'desc').limit(1);
    const videos = collectionData(query).pipe(startWith([]));
</script>

<div id="main-container">
	<div class="center-vh">
		{#each $videos as video}
			<Slide {...video} />
		{/each}
	</div>
</div>

<style>
    #main-container {
		position: relative;
		width: 100%;
		height: calc(100% - 12em);
	}

    @media screen and (max-width: 39em) {
        #main-container {
		    height: 100%;
        }
	}
</style>