<script>
	import { Router } from "@sveltech/routify";
  	import { routes } from "@sveltech/routify/tmp/routes";

	import './theme/common.scss';
	import './theme/_smui-theme.scss';
</script>

<main>
	<div class="bg-img" />
	<div id="name-tag" class="noselect">
		<div id="name">Mitko Nikov</div>
		<div id="coder">CODER</div>
	</div>

	<Router {routes} />
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
		right: 3em;
		top: 0.5em;
		font-size: 18px;
		text-align: right;
	}

	#name {
		font-size: 2em;
		margin-top: 0.67em;
		font-weight: 700;
	}

	#coder {
		margin-top: 0.1em;
		font-weight: 300;
	}
</style>