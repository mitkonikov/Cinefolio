<script lang="ts">
    interface Props {
        size?: string;
        border?: boolean;
        onclick?: (event: MouseEvent) => void;
        children?: import('svelte').Snippet;
    }

    let { size = "1.5em", border = false, onclick, children }: Props = $props();
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
    class={"btn-icon" + (border ? " btn-icon-border" : "")}
    style={`width: ${size}; height: ${size};`}
    onclick={onclick}
    >
        {@render children?.()}
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
