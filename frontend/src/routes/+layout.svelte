<script>
    import "../app.css";
    import "../css/fonts.css";
    /** @type {{children?: import('svelte').Snippet}} */
    let { children } = $props();
</script>

<svelte:head>
	<title>Cinefolio</title>
</svelte:head>

{@render children?.()}
