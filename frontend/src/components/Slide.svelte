<script lang="ts">
    import Hoverable from './Hoverable.svelte';
    import Modal from './Modal.svelte';
    import Play from 'svelte-material-icons/Play.svelte';
    import IconButton from './IconButton.svelte';
    import { thumbnailURL, type Film } from '$lib/portfolio';

    let { data }: { data: Film } = $props();
    let thumbnail = $derived(thumbnailURL(data));
    let imageStyle = $derived(thumbnail ? `background-image: url("${thumbnail}");` : '');
    let showModal = $state(false);
</script>

<Hoverable >
    {#snippet children({ hovering })}
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
            class={"slide"}
            onclick={() => {
                showModal = true;
            }}
        >
            <div class="frame-container">
                <div class="thumbnail-container" style={imageStyle} style:background-size={data.thumbnail?.crop ? '132%' : null}></div>
                {#if hovering}
                    <div class="slide-overlay"></div>
                    <div class="text-container noselect">
                        <div>
                            <div class="text-title">{data.title}</div>
                            <div class="text-content">{data.description ? data.description : ""}</div>
                            <div class="timestamp">  Published {data.published.slice(0, 4)}</div>
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
    {/snippet}
</Hoverable>

<Modal bind:showModal>
    <div class="thumbnail-container">
        <img src={thumbnail || undefined} class="thumbnail" alt={data.title} />
        <div class="play-btn">
            <IconButton onclick={() => {
                window.open(data.link || data.videoUrl, "_blank", "noopener");
            }}>
                <Play></Play>
            </IconButton>
        </div>
    </div>
    <div class="modal-info noselect">
        <div>
            <div class="text-title">{data.title}</div>
            <div class="text-content">{data.description ? data.description : ""}</div>
            <div class="timestamp">  Published {data.published.slice(0, 4)}</div>
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

    .slide-overlay {
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
