<script lang="ts">
	import { createEventDispatcher } from 'svelte';

    export let size = "1.5em";
    export let border = false;

    const dispatch = createEventDispatcher();
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
    class={"btn-icon" + (border ? " btn-icon-border" : "")}
    style={`width: ${size}; height: ${size};`}
    on:click={(event) => dispatch('click', event.detail)}
    >
        <slot/>
</div>

<style>
    .btn-icon {
        display: flex;
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

    .btn-icon-border {
        border: 0.03em rgba(255, 255, 255, 0.6) solid;
        border-radius: 50%;
    }
</style>
