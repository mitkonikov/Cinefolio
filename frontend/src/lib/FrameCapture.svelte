<script lang="ts">
    let { onCapture }: { onCapture: (file: File) => void } = $props();
    let source = $state<File>();
    let video = $state<HTMLVideoElement>();
    let url = $state('');
    let ready = $state(false);
    let error = $state('');

    $effect(() => {
        if (!source) return;
        const next = URL.createObjectURL(source);
        url = next;
        return () => URL.revokeObjectURL(next);
    });

    function capture() {
        if (!video || !video.videoWidth) return;
        video.pause();
        const canvas = document.createElement('canvas');
        canvas.width = Math.min(960, video.videoWidth);
        canvas.height = Math.round(canvas.width * video.videoHeight / video.videoWidth);
        const context = canvas.getContext('2d');
        if (!context) {
            error = 'Frame capture is unavailable in this browser.';
            return;
        }
        context.drawImage(video, 0, 0, canvas.width, canvas.height);
        canvas.toBlob((blob) => {
            if (blob && blob.size <= 512 * 1024) {
                onCapture(new File([blob], 'video-frame.jpg', { type: 'image/jpeg' }));
            } else {
                error = 'The selected frame could not be saved as a JPEG under 512 KiB.';
            }
        }, 'image/jpeg', 0.8);
    }
</script>

<details>
    <summary>Capture a frame from a video</summary>
    <label>Choose a local video
        <input type="file" accept="video/*" onchange={(event) => {
            source = event.currentTarget.files?.[0];
            ready = false;
            error = '';
        }} />
    </label>
    {#if source}
        <!-- svelte-ignore a11y_media_has_caption -->
        <video bind:this={video} src={url || undefined} controls playsinline preload="metadata"
            onloadeddata={() => ready = true}
            onseeking={() => ready = false}
            onseeked={() => ready = true}
            onerror={() => { ready = false; error = 'This video could not be opened in your browser.'; }}></video>
        <button type="button" onclick={capture} disabled={!ready}>Use current frame</button>
    {/if}
    {#if error}<p role="alert">{error}</p>{/if}
</details>

<style>
    details { padding: 0.65rem 0; }
    label { display: grid; gap: 0.3rem; margin: 0.5rem 0; }
    video { display: block; max-width: 100%; max-height: 20rem; margin-bottom: 0.5rem; }
</style>
