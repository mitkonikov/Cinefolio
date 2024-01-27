<script lang="ts">
    import Close from "svelte-material-icons/Close.svelte";
    import IconButton from "./IconButton.svelte";

	export let showModal: boolean;

	let dialog: HTMLDialogElement;

	$: if (dialog && showModal) dialog.showModal();
</script>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-noninteractive-element-interactions -->
<dialog
	bind:this={dialog}
	on:close={() => (showModal = false)}
	on:click|self={() => dialog.close()}
>
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div on:click|stopPropagation>
		<slot />
		<hr class="decoration" />
		<!-- svelte-ignore a11y-autofocus -->
		<div class="close-btn" autofocus on:click={() => dialog.close()}>
            <IconButton>
                <Close></Close>
            </IconButton>
        </div>
	</div>
</dialog>

<style>
	dialog {
        background-color: rgb(0, 12, 21);
        color: white;
        border: none !important;
        padding: 0;
        outline: 0px solid !important;
	}
	dialog::backdrop {
		background: rgba(0, 0, 0, 0.814);
	}
	dialog[open] {
		animation: zoom 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
	}
	@keyframes zoom {
		from {
			transform: scale(0.95);
		}
		to {
			transform: scale(1);
		}
	}
	dialog[open]::backdrop {
		animation: fade 0.2s ease-out;
	}
	@keyframes fade {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}
    .close-btn {
        position: absolute;
        top: 1em;
        right: 1em;
    }
    .decoration {
        margin: 0.5em 0 0 0;
    }
</style>
