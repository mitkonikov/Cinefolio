<script>
	/** @type {{children?: import('svelte').Snippet<[any]>}} */
	let { children } = $props();
	let hovering = $state(false);

	function enter() {
        hovering = true;
	}

	function leave() {
		hovering = false;
	}
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="hoverable" onmouseenter={enter} onmouseleave={leave}>
	{@render children?.({ hovering, })}
</div>

<style>
	.hoverable {
		display: inline-block;
		width: 100%;
	}
</style>