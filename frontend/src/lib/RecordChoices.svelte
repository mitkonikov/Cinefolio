<script lang="ts">
    interface Item { id: number; text: string }
    let { title, items, selected = $bindable<number[]>([]) }: {
        title: string; items: Item[]; selected: number[];
    } = $props();
    let query = $state('');
    let matches = $derived(items.filter((item) => item.text.toLowerCase().includes(query.trim().toLowerCase())));
</script>

<fieldset>
    <legend>{title}</legend>
    {#if items.length > 10}<input aria-label="Filter available records" placeholder="Filter…" bind:value={query} />{/if}
    <div class="choices">
        {#each matches as item (item.id)}
            <label><input type="checkbox" checked={selected.includes(item.id)} onchange={(event) => {
                selected = event.currentTarget.checked
                    ? [...selected, item.id]
                    : selected.filter((id) => id !== item.id);
            }} /> {item.text}</label>
        {:else}
            <p>No matching records.</p>
        {/each}
    </div>
</fieldset>

<style>
    fieldset { border: 1px solid #46484c; border-radius: 0.45rem; margin: 0; }
    .choices { display: grid; gap: 0.35rem; max-height: 12rem; overflow: auto; }
    label { cursor: pointer; }
    input:not([type='checkbox']) { width: 100%; margin: 0.2rem 0 0.5rem; }
</style>
