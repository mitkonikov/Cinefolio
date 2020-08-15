<script>
	import { Router ,goto } from "@sveltech/routify";
  	import { routes } from "@sveltech/routify/tmp/routes";

	import './theme/common.scss';
	import './theme/_smui-theme.scss';
</script>

<main>
	<div class="bg-img" />
	
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
</style>