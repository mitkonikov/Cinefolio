<script>
    import IconButton from './../utils/IconButton.svelte';

    import YouTube from 'svelte-material-icons/Youtube.svelte';
    import GitHub from 'svelte-material-icons/GitHub.svelte';
    import Gmail from 'svelte-material-icons/Gmail.svelte'
    import Facebook from 'svelte-material-icons/Facebook.svelte';
</script>

<div id="links">
    <IconButton on:click={() => { window.open('https://www.youtube.com/channel/UCEejMIwMqDOXxl05ZrN-E4g', '_blank'); }}>
        <YouTube/>
    </IconButton>
    <IconButton on:click={() => { window.open('https://github.com/mitkonikov', '_blank'); }}>
        <GitHub/>
    </IconButton>
    <IconButton on:click={() => { window.open('mailto:mitkonikov01@gmail.com', '_blank'); }}>
        <Gmail/>
    </IconButton>
    <IconButton on:click={() => { window.open('https://www.facebook.com/mitko.nikov.3', '_blank'); }}>
        <Facebook/>
    </IconButton>
</div>

<style>
    #links {
        position: relative;
        width: 100%;
        height: 3em;    
        display: flex;
        justify-content: center;
        align-content: center;
        margin-top: 0.3em;
    }
</style>