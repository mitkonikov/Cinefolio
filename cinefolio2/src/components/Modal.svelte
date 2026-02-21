<script lang="ts">
	import Close from "svelte-material-icons/Close.svelte";
	import IconButton from "./IconButton.svelte";

	interface Props {
		showModal: boolean;
		children?: import('svelte').Snippet;
	}

	let { showModal = $bindable(), children }: Props = $props();

	let dialog: HTMLDialogElement = $state();

	$effect(() => {
		if (!dialog) return;

		if (showModal) {
			if (!dialog.open) dialog.showModal();
			return;
		}

		if (dialog.open) dialog.close();
	});
</script>

<!-- svelte-ignore a11y_click_events_have_key_events, a11y_no_noninteractive_element_interactions -->
<dialog
	bind:this={dialog}
	onclose={() => (showModal = false)}
	onclick={(event) => {
		if (event.target === event.currentTarget) dialog.close();
	}}
>
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div onclick={(event) => event.stopPropagation()}>
		{@render children?.()}
		<hr class="decoration" />
		<!-- svelte-ignore a11y_autofocus -->
		<div class="close-btn" autofocus onclick={() => dialog.close()}>
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
