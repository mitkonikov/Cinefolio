<script>
    import Card, { PrimaryAction } from '@smui/card';
    import { goto } from '@sveltech/routify'
    
    export let title = '';
    export let number = '';
    export let where = '';
</script>

<div class="footer-snap">
    <Card>
        <PrimaryAction on:click={$goto('/' + where)}>
            <div class="footer-snap-content noselect">
                <div class="snap-number">{number}</div>
                <div class="snap-title">{title}</div>
                <slot></slot>
            </div>
        </PrimaryAction>
    </Card>
</div>

<style>
    .footer-snap {
        padding: 0.5em 6%;
    }
    
    .footer-snap-content {
        position: relative;
		height: 2.7em;
        padding: 0.2em;
        text-align: center;
    }
    
    .footer-snap-content:hover {
		background-color: rgba(255, 255, 255, 0.05);
    }

    .snap-title {
        display: inline-block;
        padding: 0.8em 0.3em 0.3em 0.3em;
        text-align: center;
        font-size: 1.1em;
        font-weight: 300;
    }

    .snap-number {
        display: inline-block;
        text-align: center;
        font-size: 1.4em;
        font-weight: 700;
        color: rgb(255, 196, 0);
    }
</style>