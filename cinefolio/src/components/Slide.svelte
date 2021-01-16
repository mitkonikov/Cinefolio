<script>
    import Hoverable from "./../js/Hoverable.svelte";

    export let thumb;
    export let alt;
    export let title;
    export let content;
    export let link;
	export let timestamp;
	export let embedlink;
	export let textHeight = 6;

	// Calculating the width and height after hovering on the element
    let width, height;
    function getWidthAndHeight(node) {
        width = getComputedStyle(node).width;
		height = width.substr(0, width.length - 2) * 9 / 16;
	}
	
	// We check if data saver is enabled on the device
    let dataSaver = false;
    if ("connection" in navigator) {
        if (navigator.connection.saveData === true) {
            console.log("Data saver is enabled.");
            dataSaver = true;
        }
	}

	// Animations
	function easeInOutCubic(x) {
		return x < 0.5 ? 4 * x * x * x : 1 - Math.pow(-2 * x + 2, 3) / 2;
	}

	function scaleV(node, { duration }) {
		return {
			duration,
			css: t => {
                return `transform: translate(0, -${10 * easeInOutCubic(1 - t)}px) scale(1, ${easeInOutCubic(t)}); height: ${textHeight * easeInOutCubic(t)}em`
            }
		};
	}
</script>

<Hoverable let:hovering>
    <div
        class={hovering ? "slide slide-more" : "slide"}
        on:click={() => {
            window.open(link, "_black");
        }}
    >
        <div class="frame-container">
            {#if hovering && !dataSaver && typeof embedlink != "undefined"}
                <div class="thumbnail-video-wrapper" in:getWidthAndHeight>
					<!-- Embed the Youtube Link -->
					<iframe
						title="thumbnail-video"
						{width}
						{height}
                        src={embedlink}
                        frameborder="0"
						autoplay="1"
                        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                        allowfullscreen
                    />
                </div>
			{:else}
				<!-- Or just display the thumbnail -->
                <img src={thumb} class="thumbnail" {alt} />
            {/if}
        </div>
		{#if hovering}
			<div 
				class="text-container noselect"
				in:scaleV = "{{ duration: 150 }}"
				out:scaleV = "{{ duration: 150 }}"
				style="height: {textHeight}em"
			>
				<div>
					<div class="text-title">{title}</div>
					<div class="text-content">{content}</div>
					<div class="timestamp">  Published {timestamp.toDate().getFullYear()}</div>
				</div>
			</div>
        {/if}
    </div>
</Hoverable>

<style>
    .thumbnail {
        width: 16em;
        background-color: rgba(0, 0, 0, 0.5);
    }

    .thumbnail-video-wrapper {
        width: 100%;
    }

    .slide {
        -webkit-box-shadow: 4px 10px 13px -5px rgba(0, 0, 0, 0.64);
        -moz-box-shadow: 4px 10px 13px -5px rgba(0, 0, 0, 0.64);
        box-shadow: 4px 10px 13px -5px rgba(0, 0, 0, 0.64);
        width: 18em;
    	transition: 0.4s cubic-bezier(0.83, 0, 0.17, 1);
		overflow: hidden;
    }

	.slide-more {
		width: 27em;
	}

    .slide {
        margin: 0.5em;
        cursor: pointer;
		border-radius: 0.3em;
		background-color: rgba(0, 0, 0, 0.9);
    }

    .thumbnail,
    .frame-container {
        width: 100%;
    }

    .frame-container {
        background-color: rgba(0, 0, 0, 0.8);
    }

    .slide-more .frame-container {
        display: inline-block;
    }

    .text-container {
        padding: 0.5em 1em 0.9em 1em;
		position: relative;
    }

    .text-title {
        font-weight: 700;
        margin-top: 0.4em;
        margin-bottom: 0.4em;
    }

    .text-content {
        font-size: 0.8em;
    }

	.timestamp {
		position: absolute;
		font-weight: 100;
		font-size: 0.7em;
		bottom: 1em;
	}
</style>