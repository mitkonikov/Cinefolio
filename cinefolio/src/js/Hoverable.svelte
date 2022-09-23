<script>
	let hovering;

	function enter() {
        hovering = true;
	}

	function leave() {
		hovering = false;
	}
</script>

<div class="hoverable" on:mouseenter={enter} on:mouseleave={leave}>
	<slot hovering={hovering}></slot>
</div>

<style>
	.hoverable {
		display: inline-block;
	}
</style>