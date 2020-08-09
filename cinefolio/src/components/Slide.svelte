<script>
    export let thumb;
    export let alt;
    export let title;
	export let content;
	export let link;
</script>

<div class="slide" on:click={() => { window.open(link, '_black'); }}>
    <div class="frame-container">
        <img src={thumb} class="thumbnail" alt={alt}/>
    </div>
    <div class="text-container noselect">
        <div class="text-title">{title}</div>
        {content}
    </div>
</div>

<style>
     .thumbnail {
		width: 16em;
		height: 9em;
		padding: 0.2em;
		background-color: rgba(0, 0, 0, 0.5);
    }

    .slide {
		width: 38em;
		padding: 0.5em 1em 0.5em 0.7em;
	}

	.slide:hover {
		background-color: rgba(0, 0, 0, 0.4);
	}

	.frame-container {
		margin: 0.2em 1em 0 0.2em;
		display: inline-block;
	}

	.text-container {
		display: inline-block;
		position: absolute;
		top: 0%;
		margin-right: 0.6em;
    }

	.text-title {
		font-size: 1.17em;
		margin-top: 1em;
		margin-bottom: 1em;
		font-weight: 700;
	}

    @media screen and (max-width: 39em) {
		.slide {
            width: 17em;
        }

		.text-container {
			display: block;
    		position: relative;
			text-align: center;
		}

		.text-title {
			margin-top: 0.4em;
			margin-bottom: 0.4em;
			text-align: center;
		}
	}
</style>