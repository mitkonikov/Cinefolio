<script>
	import MainSlides from './components/MainSlides.svelte';
	import FooterSnap from './components/FooterSnap.svelte';

	import './theme/common.scss';
	import './theme/_smui-theme.scss';
</script>

<main>
	<div class="bg-img" />
	<div id="name-tag" class="noselect">
		<h1>Mitko Nikov</h1>
	</div>

	<div id="homepage">
		<MainSlides />

		<div id="footer-container">
			<div id="footer-flex" class="center-vh">
				<FooterSnap title="PROJECTS" number="100+" content="Take a look..."></FooterSnap>
				<FooterSnap title="VIDEOS" number="100+" content="Let's see them!"></FooterSnap>
				<FooterSnap title="SATISFIED CLIENTS" number="100+" content="Meet my clients."></FooterSnap>
			</div>
		</div>
	</div>
</main>

<style>
	.bg-img {
		overflow: hidden;
		position: absolute;
		width: 100vw;
		height: 100vh;
		background-image: url('./img/02.jpg');
		background-repeat: no-repeat;
		background-size: cover;
		opacity: 0.3;
	}

	#homepage {
		position: relative;
		width: 100vw;
		height: 100vh;
	}

	#name-tag {
		position: absolute;
		right: 3em;
		top: 0.5em;
		font-size: 18px;
	}

	#footer-container {
		position: relative;
		height: 12em;
		width: 100%;
		background-color: rgba(0, 0, 0, 0.3);
	}

	#footer-flex {
		display: flex;
		justify-content: center;
	}
	
    @media screen and (max-width: 39em) {
		#footer-container {
			display: none;
		}
	}
</style>