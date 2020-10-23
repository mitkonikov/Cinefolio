<script>
	import Hoverable from './../js/Hoverable.svelte';

    export let thumb;
    export let alt;
    export let title;
	export let content;
	export let link;
	export let timestamp;
</script>

<Hoverable let:hovering={hovering}>
		<div class="{hovering ? 'slide-more' : 'slide'}" on:click={() => { window.open(link, '_black'); }}>
			<div class="frame-container">
				<img src={thumb} class="thumbnail" alt={alt}/>
			</div>
			{#if hovering}
				<div class="text-container noselect">
					<div class="text-title">{title}</div>
					<div class="text-content">{content}</div>
				</div>
			{/if}
		</div>
</Hoverable>

<style>
     .thumbnail {
		width: 16em;
		background-color: rgba(0, 0, 0, 0.5);
    }

    .slide {
		-webkit-box-shadow: 4px 10px 13px -5px rgba(0,0,0,0.64);
		-moz-box-shadow: 4px 10px 13px -5px rgba(0,0,0,0.64);
		box-shadow: 4px 10px 13px -5px rgba(0,0,0,0.64);
	}

	.slide-more {
		width: 17em;
	}

	.slide, .slide-more {
		margin: 0.5em;
	}

	.slide-more .thumbnail {
		width: 100%;
	}

	.slide:hover, .slide-more:hover {
		background-color: rgba(0, 0, 0, 0.2);
	}

	.frame-container {
		background-color: rgba(0, 0, 0, 0.5);
	}

	.slide-more .frame-container {
		display: inline-block;
	}

	.text-container {
		text-align: center;
		padding: 0.5em 1em 0.9em 1em;
    }

	.text-title {
		font-weight: 700;
		margin-top: 0.4em;
		margin-bottom: 0.4em;
		text-align: center;
	}

	.text-content {
		font-size: 0.8em;
	}
</style>