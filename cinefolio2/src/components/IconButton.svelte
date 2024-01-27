<script lang="ts">
	import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="btn-icon" on:click={(event) => dispatch('click', event.detail)}>
    <slot/>
</div>

<style>
    .btn-icon {
        display: flex;
        width: 1.5em;
        height: 1.5em;
        justify-content: center;
        align-items: center;
    }

    .btn-icon:before {
        position: absolute;
        content: "";
        width: 100%;
        height: 100%;
        border-radius: 50%;
    }

    .btn-icon:hover::before {
        background-color: rgba(255, 255, 255, 0.2);
        cursor: pointer;
    }
</style>
