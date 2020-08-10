<script>
    import FooterSnap from './../components/FooterSnap.svelte';
    import MainSlides from '../components/MainSlides.svelte';

    import { collectionData } from 'rxfire/firestore';
    import { startWith } from 'rxjs/operators';
    import { db } from './../js/firebase';

    const query = db.collection('Home').limit(1);
    const info = collectionData(query).pipe(startWith([]));
    $: projects = typeof $info[0] == 'undefined' ? 0 : $info[0].Projects;
    $: videos = typeof $info[0] == 'undefined' ? 0 : $info[0].Videos;
    $: clients = typeof $info[0] == 'undefined' ? 0 : $info[0].Clients;
</script>

<div id="homepage">
    <MainSlides />

    <div id="footer-container">
        <div id="footer-flex" class="center-vh">
            <FooterSnap title="PROJECTS" number={projects + '+'} content="Take a look..." where='projects'></FooterSnap>
            <FooterSnap title="VIDEOS" number={videos + '+'} content="Let's see them!" where='videos'></FooterSnap>
            <FooterSnap title="SATISFIED CLIENTS" number={clients + '+'} content="Meet my clients." where='clients'></FooterSnap>
        </div>
    </div>
</div>

<style>
    #homepage {
		position: relative;
		width: 100vw;
		height: 100vh;
	}

    #footer-container {
        position: absolute;
        bottom: 0;

        overflow: hidden;
        height: 12em;
		width: 100%;
        
		background-color: rgba(0, 0, 0, 0.3);

        animation: 0.5s ease-out 0s 1 footerAnimation;
	}

	#footer-flex {
		display: flex;
		justify-content: center;
	}

    @keyframes footerAnimation {
        0% {
            height: 0em;
        }
        100% {
            height: 12em;
        }
    }
	
    @media screen and (max-width: 39em) {
		#footer-container {
			display: none;
		}
	}
</style>