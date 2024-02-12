<script>
	let hovering = false;

	function enter() {
        hovering = true;
	}

	function leave() {
		hovering = false;
	}
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="hoverable" on:mouseenter={enter} on:mouseleave={leave}>
	<slot hovering={hovering}></slot>
</div>

<style>
	.hoverable {
		display: inline-block;
		width: 100%;
	}
</style>