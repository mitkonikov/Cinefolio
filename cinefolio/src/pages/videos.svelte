<script>
    import NameTag from '../components/NameTag.svelte';
    import Playlist from '../components/Playlist.svelte';
</script>

<div id="videos-container-center">
    <div id="name-tag-container">
        <NameTag />
    </div>
    
    <Playlist repositoryName={"Videos"} displayName={"Latest"}/>
    <Playlist repositoryName={"Strumica Open Festival"}/>
</div>

<style>
    #name-tag-container {
        height: 8em;
    }

    #videos-container-center {
        position: relative;
        width: 100%;
        height: 100%;
    }
</style>