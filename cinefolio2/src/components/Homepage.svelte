<script>
    import FooterSnap from './../components/FooterSnap.svelte';
    // import FooterSnapMobile from './../components/FooterSnapMobile.svelte';
    import MainSlides from '../components/MainSlides.svelte';
    import NameTag from '../components/NameTag.svelte';
</script>

<div id="homepage">
    <MainSlides />
    <NameTag />

    <div id="footer-container">
        <div id="footer-flex" class="center-vh">
            <!-- <FooterSnap title="PROJECTS" number={metadata.Projects + '+'} content="Take a look..." where='projects'></FooterSnap>
            <FooterSnap title="VIDEOS" number={metadata.Videos + '+'} content="Let's see them!" where='videos'></FooterSnap>
            <FooterSnap title="SATISFIED CLIENTS" number={metadata.Clients + '+'} content="Meet my clients." where='clients'></FooterSnap> -->
        </div>
        <div id="footer-flex-mobile">
            <!-- <FooterSnapMobile title="PROJECTS" number={projects + '+'} where='projects'></FooterSnapMobile>
            <FooterSnapMobile title="VIDEOS" number={videos + '+'}  where='videos'></FooterSnapMobile>
            <FooterSnapMobile    title="SATISFIED CLIENTS" number={clients + '+'} where='clients'></FooterSnapMobile> -->
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

    #footer-flex-mobile {
        display: none;
        bottom: 0;

        padding: 0.5em 0em;

		background-color: rgba(0, 0, 0, 0.3);
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
		#footer-flex {
            display: none;
		}

        #footer-container {
            height: initial;
        }

        #footer-flex-mobile {
            display: block;
        }
	}
</style>