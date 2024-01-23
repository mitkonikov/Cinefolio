<script lang="ts">
    // import { collectionData } from 'rxfire/firestore';
    // import { startWith } from 'rxjs/operators';
    import { db } from '$lib/firebase';
	import { collection, getDocs, limit, query } from 'firebase/firestore';
    import Playlist from './../components/Playlist.svelte';

    let metadata = {
        Projects: 0,
        Clients: 0,
        Videos: 0
    };

    const q = query(collection(db, 'Home'), limit(1))
    const querySnapshot = getDocs(q).then((snapshot) => {
        snapshot.forEach((doc) => {
            metadata = doc.data() as any;
        });
    });
</script>

<main>
	<div class="bg-img" />
	
    <div id="videos-container-center">
        <Playlist repositoryName={"Videos"} displayName={"Latest"}/>
        <Playlist repositoryName={"Strumica Open Festival"}/>
    </div>
</main>

<style>
	.bg-img {
		overflow: hidden;
		position: absolute;
		width: 100vw;
		height: 100vh;
		background-image: url('/img/Stars.jpg');
		background-repeat: no-repeat;
		background-size: cover;
		opacity: 0.3;
	}
    
    #videos-container-center {
        position: relative;
        padding-top: 2em;
        width: 100%;
        height: 100%;
    }
</style>