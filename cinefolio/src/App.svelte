<script>
	import { Router ,goto } from "@sveltech/routify";
  	import { routes } from "@sveltech/routify/tmp/routes";

	import './theme/common.scss';
	import './theme/_smui-theme.scss';
</script>

<main>
	<div class="bg-img" />
	
	<Router {routes} />

	<div id="name-tag" class="noselect" on:click={() => { window.location.href = window.location.origin; }}>
		<div id="name">Mitko Nikov</div>
		<div id="coder">CODER</div>
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

	#name-tag {
		position: absolute;
		right: 2.5em;
		top: 1.17em;
		font-size: 18px;
		text-align: right;
		padding: 0.3em;
	}

	#name-tag:hover {
		transform: scale(1.03);
	}

	#name {
		font-size: 2em;
		font-weight: 700;
	}

	#coder {
		margin-top: 0.1em;
		font-weight: 300;
	}

	@media screen and (max-width: 20em) {
		#name-tag {
			right: 1.6em;
		}
	}
</style>