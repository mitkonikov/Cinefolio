<script lang="ts">
	import type { Timestamp } from "firebase/firestore";
    import Hoverable from "./Hoverable.svelte";

    export let thumb : string;
    export let alt : string;
    export let title : string;
    export let content : string;
    export let link : string;
	export let timestamp : Timestamp;
	export let embedlink : string;
	export let textHeight = 6;
    export let first = false;

	// Calculating the width and height after hovering on the element
    let width = "0", height = "0";
    function getWidthAndHeight(node: Element) : __sveltets_2_SvelteTransitionReturnType {
        width = getComputedStyle(node).width;
		height = (parseInt(width.substring(0, width.length - 2)) * 9 / 16).toString();
        return {};
	}
	
	// We check if data saver is enabled on the device
    let dataSaver = false;
    if ("connection" in navigator) {
        if ((navigator.connection as any).saveData === true) {
            console.log("Data saver is enabled.");
            dataSaver = true;
        }
	}

	// Animations
	function easeInOutCubic(x: number) {
		return x < 0.5 ? 4 * x * x * x : 1 - Math.pow(-2 * x + 2, 3) / 2;
	}

	function scaleV(node: Element, info : { duration: number }) {
		return {
			duration: info.duration,
			css: (t: number) => {
                return `transform: translate(0, -${10 * easeInOutCubic(1 - t)}px) scale(1, ${easeInOutCubic(t)}); height: ${textHeight * easeInOutCubic(t)}em`
            }
		};
	}
</script>

<Hoverable let:hovering>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
        class={"slide " + (hovering ? (first ? "slide-more-first" : "slide-more") : "slide-non-hover")}
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
				out:scaleV = "{{ duration: 100 }}"
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
        font-size: 0.4em;
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
		transform: scale(1.4);
        z-index: 100000;
	}

    .slide-more-first {
		transform: scale(1.4) translate(3em, 0);
        z-index: 100000;
    }

    .slide {
        margin: 0.5em;
        cursor: pointer;
		border-radius: 0.3em;
		background-color: rgba(0, 0, 0, 0.9);
    }

    .slide-non-hover {
    	transition: 0.2s;
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
        font-size: 0.8em;
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