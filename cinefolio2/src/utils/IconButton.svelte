<script>
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();
</script>

<div
	class="btn-icon"
	role="button"
	tabindex="0"
	on:click={(event) => dispatch('click', event.detail)}
	on:keydown|preventDefault={(event) => {
		if (event.key == "Enter") {
			dispatch('click', event.detail);
		}
	}}
>
	<div class="btn-icon-center">
		<slot />
	</div>
</div>

<style>
	.btn-icon {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: transparent;
		width: 2.5em;
		height: 2.5em;
		cursor: pointer;
	}

	.btn-icon:hover::before {
		opacity: 0.2;
	}

	.btn-icon::before {
		position: absolute;
		content: "";
		background-color: black;
		width: 2.2em;
		height: 2.2em;
		border-radius: 50%;
		opacity: 0;
		z-index: -2;
		transition: opacity 0.1s cubic-bezier(0.65, 0, 0.35, 1);
	}

	.btn-icon-center {
		width: 0.9em;
		height: 0.9em;
		font-size: 1.3em;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.btn-icon:hover {
		transform: scale(1.15);
		transition: transform 0.15s cubic-bezier(0.65, 0, 0.35, 1);
	}
</style>
