<div id="name-tag" class="noselect" on:click={() => { window.location.href = window.location.origin; }}>
    <div id="name">Mitko Nikov</div>
    <div id="coder">CODER</div>
</div>

<style>
    #name-tag {
		position: absolute;
		right: 2.5em;
		top: 1.17em;
		font-size: 18px;
		text-align: right;
		padding: 0.3em;
	}

	#name-tag:hover {
		transform: scale(1.03);
	}

	#name {
		font-size: 2em;
		font-weight: 700;
	}

	#coder {
		margin-top: 0.1em;
		font-weight: 300;
	}

	@media screen and (max-width: 20em) {
		#name-tag {
			right: 1.6em;
		}
	}
</style>