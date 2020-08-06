<script>
    export let thumbnail;
    export let alt;
    export let title;
    export let content;
</script>

<div class="slide">
    <div class="frame-container">
        <img src={thumbnail} class="thumbnail" alt={alt}/>
    </div>
    <div class="text-container noselect">
        <h3>{title}</h3><br>
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
    	height: 9.5em;
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

    @media screen and (max-width: 39em) {
		.slide {
            width: 17em;
        }

		.text-container {
			display: none;
		}
	}
</style>