<script lang="ts">
    import Hoverable from "./Hoverable.svelte";
	import Modal from './Modal.svelte';
	import type { IFilm } from "../types/film";
    import Play from "svelte-material-icons/Play.svelte";
	import IconButton from "./IconButton.svelte";
	import { STATIC_URL } from "$lib/constants";
	import { onMount } from "svelte";

    export let data: IFilm;

    $: imgStyle = '';

	// We check if data saver is enabled on the device
    let dataSaver = false;
    if ("connection" in navigator) {
        if ((navigator.connection as any).saveData === true) {
            console.log("Data saver is enabled.");
            dataSaver = true;
        }
	}

    function getYouTubeID(url: string) {
        // Regular expression to find YouTube video ID
        var regExp = /^(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})/;
        var match = url.match(regExp);
        if (match && match[1]) {
            return match[1];
        } else {
            return null;
        }
    }

    function getVimeoID(url: string) {
        let s = url.split('/');
        return s[s.length - 1];
    }

    function getThumbnailPath() {
        if (data.thumbnail) {
            return `${STATIC_URL}/thumbs/${data.thumbnail.filepath}`;
        } else if (data.link) {
            if (data.link.includes('vimeo')) {
                const VimeoID = getVimeoID(data.link);
                return `https://vumbnail.com/${VimeoID}.jpg`;
            } else {
                const YoutubeID = getYouTubeID(data.link);
                return `https://img.youtube.com/vi/${YoutubeID}/sddefault.jpg`;
            }
        }
        return "";
    }

    onMount(() => {
        let thumbnailLink = getThumbnailPath();
        let backgroundSize = (data.thumbnail && data.thumbnail.crop ? ' background-size: 132%' : '');
        imgStyle = `background-image: url(${thumbnailLink});` + backgroundSize;
    });

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
            <div class="thumbnail-container" style={imgStyle}></div>
            {#if hovering}
                <div class="slide-overlay"></div>
                <div class="text-container noselect">
                    <div>
                        <div class="text-title">{data.title}</div>
                        <div class="text-content">{data.description ? data.description : ""}</div>
                        <div class="timestamp">  Published {new Date(data.published).getFullYear()}</div>
                    </div>
                </div>
            {:else}
                <div class="text-thumb-container noselect">
                    <div class="thumbnail-title-container">
                        {#if data.title.split('-').length == 1}
                            <div class="text-thumb-title">{data.title}</div>
                        {:else}
                            <div class="text-thumb-title">{data.title.split('-')[0]}</div>
                            <div class="text-thumb-title">{data.title.split('-')[1]}</div>
                        {/if}
                    </div>
                </div>
            {/if}
        </div>
    </div>
</Hoverable>

<Modal bind:showModal>
    <div class="thumbnail-container">
        <img src={getThumbnailPath()} class="thumbnail" alt={data.title} />
        <div class="play-btn">
            <IconButton on:click={() => {
                if (data.link) {
                    window.open(data.link, "_black");
                } else {
                    window.open(`${STATIC_URL}/files/${data.filepath}`);
                }
            }}>
                <Play></Play>
            </IconButton>
        </div>
    </div>
    <div class="modal-info noselect">
        <div>
            <div class="text-title">{data.title}</div>
            <div class="text-content">{data.description ? data.description : ""}</div>
            <div class="timestamp">  Published {new Date(data.published).getFullYear()}</div>
        </div>
    </div>
</Modal>

<style>
    .thumbnail-container {
        width: 100%;
        aspect-ratio: 16 / 9;
        background-size: cover;
        background-position: center;
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
		border-radius: 0.8em;
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

    .text-thumb-container {
        padding: 0.5em 1em 0.9em 1em;
		position: absolute;
        display: flex;
        align-items: left;
        justify-content: center;
        flex-direction: column;
        z-index: 10000;
        top: 0;
        bottom: 0;
    }

    .thumbnail-title-container {
        padding: 0.5em 0.5em;
        border-radius: 0.2em;
        background-color: rgba(0, 0, 0, 0.8);
    }

    .slide-overlay, .slide-thumb-overlay {
        background-color: black;
        opacity: 0.6;
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

    .text-thumb-title {
        font-weight: 700;
        font-size: 1em;
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

        .text-thumb-container {
            display: none;
        }
    }
</style>
