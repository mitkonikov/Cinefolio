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

	// We check if data saver is enabled on the device
    let dataSaver = false;
    if ("connection" in navigator) {
        if ((navigator.connection as any).saveData === true) {
            console.log("Data saver is enabled.");
            dataSaver = true;
        }
	}
</script>

<Hoverable let:hovering>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
        class={"slide"}
        on:click={() => {
            window.open(link, "_black");
        }}
    >
        <div class="frame-container">
            <img src={thumb} class="thumbnail" {alt} />
            {#if hovering}
                <div class="slide-overlay"></div>
                <div class="text-container noselect">
                    <div>
                        <div class="text-title">{title}</div>
                        <div class="text-content">{content}</div>
                        <div class="timestamp">  Published {timestamp.toDate().getFullYear()}</div>
                    </div>
                </div>
            {/if}
        </div>
    </div>
</Hoverable>

<style>
    .thumbnail {
        background-color: rgba(0, 0, 0, 0.5);
    }

    .slide {
        -webkit-box-shadow: 4px 10px 13px -5px rgba(0, 0, 0, 0.64);
        -moz-box-shadow: 4px 10px 13px -5px rgba(0, 0, 0, 0.64);
        box-shadow: 4px 10px 13px -5px rgba(0, 0, 0, 0.64);
    	transition: 0.4s cubic-bezier(0.83, 0, 0.17, 1);
		overflow: hidden;
    }

    .slide {
        margin: 0.2em;
        cursor: pointer;
		border-radius: 0.3em;
		background-color: rgba(0, 0, 0, 0.9);
    }

    .thumbnail,
    .frame-container {
        width: 100%;
        position: relative;
    }

    .frame-container {
        background-color: rgba(0, 0, 0, 0.8);
    }

    .text-container {
        padding: 0.5em 1em 0.9em 1em;
		position: absolute;
        z-index: 10000;
        top: 0;
        bottom: 0;
    }

    .slide-overlay {
        background-color: black;
        opacity: 0.5;
        position: absolute;
        height: 100%;
        width: 100%;
        top: 0;
        left: 0;
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