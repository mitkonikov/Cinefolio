<script>
	import Slide from './Slide.svelte';
</script>

<div id="main-container">
	<div class="center-vh">
		<Slide 
			thumbnail='http://i3.ytimg.com/vi/uOulc9dy69I/maxresdefault.jpg'
			alt="volumetrics"
			title='Ultimate Volumetrics Guide in Blender'
			content='We dive deep into the volumetrics optimizations that can be done inside Blender.'
		/>
	</div>
</div>

<style>
    #main-container {
		position: relative;
		width: 100%;
		height: calc(100% - 12em);
	}

    @media screen and (max-width: 39em) {
        #main-container {
		    height: 100%;
        }
	}
</style>