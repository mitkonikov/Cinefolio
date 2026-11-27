<script lang="ts">
    import IconButton from './IconButton.svelte';
    import Play from 'svelte-material-icons/Play.svelte';
    import { playbackURL, thumbnailURL, type Trailer } from '#lib/portfolio.ts';

    let { trailer }: { trailer: Trailer | null } = $props();
    let poster = $derived(trailer ? thumbnailURL(trailer.film) : '');
    let filmURL = $derived(trailer ? playbackURL(trailer.film) : '');
    let endedURL = $state('');
</script>

{#if filmURL}
    <div class="trailer-frozen trailer-play-container">
        <div class="trailer-play">
            <IconButton border={true} onclick={() => window.open(filmURL, '_blank', 'noopener')}>
                <Play />
            </IconButton>
        </div>
    </div>
{/if}

<div class="trailer">
    {#if trailer && (endedURL !== trailer.url || !poster)}
        <!-- svelte-ignore a11y_media_has_caption -->
        <video width="100%" autoplay muted playsinline poster={poster || undefined} onended={() => endedURL = trailer?.url ?? ''}>
            <source src={trailer.url} type="video/mp4">
        </video>
    {:else if poster}
        <div class="trailer-frozen" style:background-image={`url("${poster}")`}></div>
    {/if}
</div>

<div class="fade-decoration"></div>

<div class="noselect" id="logo">
    <div id="main-logo">cinefolio</div>
    <div class="logo-dec"></div>
    <div id="info-logo">Mitko Nikov's cinematography portfolio</div>
</div>

<style>
    .fade-decoration {
        position: absolute;
        width: 100%;
        margin-top: calc(100vw / 16 * 9 / 1.4);
        height: calc(100vw / 16 * 9 - 100vw / 16 * 9 / 1.4);
        background: -moz-linear-gradient(top, rgba(20,20,20,0), rgba(20,20,20, 1));
        background: -webkit-linear-gradient(top, rgba(20,20,20,0), rgba(20,20,20, 1));
    }

    #logo {
        position: absolute;
        width: 100%;
        padding-top: calc(100vw / 16 * 9 / 4);
        padding-left: 4em;
    }

    #main-logo {
        display: inline-block;
        padding: 0.5em 0 0 0;
        font-size: 10vw;
        font-weight: 100;
        color: #eee;
        text-shadow: 0.03em 0.03em 0 black;
    }

    #info-logo {
        display: inline-block;
        padding: 0.5em 0 0 0;
        font-size: 1vw;
        font-weight: 100;
        color: #eee;
        text-shadow: 0.03em 0.03em 0 black;
    }

    #main-logo:hover {
        padding: 0.6em 0 0.1em 0;
        cursor: pointer;
        color: white;
    }

    .trailer {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 0;
    }

    .logo-dec {
        width: calc(100vw / 16 * 9 / 2);
        height: 0.15em;
        background-color: rgba(239, 239, 239, 0.895);
    }

    .trailer-frozen {
        width: 100%;
        height: calc(100vw / 16 * 9);
        background-repeat: no-repeat;
        background-size: cover;
    }
    
    .trailer-play-container {
        position: absolute;
        z-index: 1;
        top: 0;
    }

    .trailer-play {
        position: absolute;
        bottom: 10vw;
        right: 10vw;
        font-size: 7em;
    }

    @media only screen and (max-width: 600px) {
        #logo {
            font-size: 0.4em;
        }

        #info-logo {
            font-size: 2em;
        }

        .trailer-play {
            bottom: 7vw;
            right: 7vw;
            font-size: 3em;
        }
    }
</style>
