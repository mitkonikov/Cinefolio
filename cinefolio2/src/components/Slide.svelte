<script lang="ts">
    import Hoverable from "./Hoverable.svelte";
	import Modal from './Modal.svelte';
	import type { IFilm } from "../types/film";
    import Play from "svelte-material-icons/Play.svelte";
	import IconButton from "./IconButton.svelte";

    export let data: IFilm;

	// We check if data saver is enabled on the device
    let dataSaver = false;
    if ("connection" in navigator) {
        if ((navigator.connection as any).saveData === true) {
            console.log("Data saver is enabled.");
            dataSaver = true;
        }
	}

    let showModal = false;
</script>

<Hoverable let:hovering>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
        class={"slide"}
        on:click={() => {
            showModal = true;
        }}
    >
        <div class="frame-container">
            <img src={data.thumb} class="thumbnail" alt={data.alt} />
            {#if hovering}
                <div class="slide-overlay"></div>
                <div class="text-container noselect">
                    <div>
                        <div class="text-title">{data.title}</div>
                        <div class="text-content">{data.content}</div>
                        <div class="timestamp">  Published {data.timestamp.toDate().getFullYear()}</div>
                    </div>
                </div>
            {/if}
        </div>
    </div>
</Hoverable>

<Modal bind:showModal>
    <div class="thumbnail-container">
        <img src={data.thumb} class="thumbnail" alt={data.alt} />
        <div class="play-btn">
            <IconButton on:click={() => {
                window.open(data.link, "_black");
            }}>
                <Play></Play>
            </IconButton>
        </div>
    </div>
    <div class="modal-info noselect">
        <div>
            <div class="text-title">{data.title}</div>
            <div class="text-content">{data.content}</div>
            <div class="timestamp">  Published {data.timestamp.toDate().getFullYear()}</div>
        </div>
    </div>
</Modal>

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

    .modal-info {
        padding: 1em 1em 2em 1em;
    }

    .modal-info .text-title {
        font-size: 1.5em;
    }

    .modal-info .text-content {
        margin-bottom: 1em;
        font-size: 1.1em;
        max-width: 600px;
    }

    .modal-info .timestamp {
        margin-bottom: 1em;
    }

    .thumbnail-container {
        position: relative;
    }

    .play-btn {
        position: absolute;
        bottom: 0.5em;
        right: 1em;
        font-size: 2em;
    }

    @media only screen and (max-width: 600px) {
        .modal-info .text-title {
            font-size: 1.2em;
        }

        .modal-info .text-content {
            font-size: 0.8em;
        }
    }
</style>
