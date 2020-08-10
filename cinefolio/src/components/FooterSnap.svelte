<script>
    import Card, { PrimaryAction } from '@smui/card';
    import { goto } from '@sveltech/routify'
    
    export let title = '';
    export let number = '';
    export let content = '';
    export let where = '';
</script>

<div class="footer-snap">
    <Card>
        <PrimaryAction on:click={$goto('/' + where)}>
            <div class="footer-snap-content noselect">
                <div class="snap-title">{title}</div>
                <div class="snap-number center-vh">{number}</div>
                <div class="snap-content">{content}</div>
                <slot></slot>
            </div>
        </PrimaryAction>
    </Card>
</div>

<style>
    .footer-snap {
		margin: 0em 1em;
    }
    
    .footer-snap-content {
        position: relative;
		height: 10em;
		width: 12em;
    }
    
    .footer-snap-content:hover {
		background-color: rgba(255, 255, 255, 0.05);
    }

    .snap-title {
        padding: 0.8em 0.3em 0.3em 0.3em;
        text-align: center;
        font-size: 1.1em;
        font-weight: 300;
    }

    .snap-number {
        position: absolute;
        text-align: center;
        font-size: 2.5em;
        font-weight: 700;
        color: rgb(255, 196, 0);
    }

    .snap-content {
        position: absolute;
        bottom: 0.8em;
        width: 100%;
        text-align: center;
        font-weight: 300;
    }
</style>