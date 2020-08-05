<div id="main-container">
    <div id="main-container-slide" class="center-vh">
        <div id="frame-container">
            <img src='http://i3.ytimg.com/vi/uOulc9dy69I/maxresdefault.jpg' class="thumbnail" alt="volumetrics"/>
        </div>
        <div id="text-container" class="noselect">
            <h3>Ultimate Volumetrics Guide in Blender</h3><br>
            We dive deep into the volumetrics optimizations that can be done inside Blender.
        </div>
    </div>
</div>

<style>
    .thumbnail {
		width: 16em;
		height: 9em;
		padding: 0.2em;
		background-color: rgba(0, 0, 0, 0.5);
    }
    
    #main-container {
		position: relative;
		width: 100%;
		height: calc(100% - 12em);
	}

	#main-container-slide {
		width: 38em;
    	height: 9.5em;
		padding: 0.5em 1em 0.5em 0.7em;
	}

	#main-container-slide:hover {
		background-color: rgba(0, 0, 0, 0.4);
	}

	#frame-container {
		margin: 0.2em 1em 0 0.2em;
		display: inline-block;
	}

	#text-container {
		display: inline-block;
		position: absolute;
		top: 0%;
		margin-right: 0.6em;
    }
    
    @media screen and (max-width: 39em) {
        #main-container {
		    height: 100%;
        }

		#main-container-slide {
            width: 17em;
        }

		#text-container {
			display: none;
		}
	}
</style>