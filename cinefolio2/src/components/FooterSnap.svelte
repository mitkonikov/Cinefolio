<script>
    
    export let title = '';
    export let number = '';
    export let content = '';
    export let where = '';
</script>

<div class="footer-snap">
    <div class="footer-snap-content noselect">
        <div class="snap-title">{title}</div>
        <div class="snap-number center-vh">{number}</div>
        <div class="snap-content">{content}</div>
        <slot></slot>
    </div>
</div>

<style>
    .footer-snap {
		margin: 0em 1em;
    }
    
    .footer-snap-content {
        position: relative;
		height: 10em;
		width: 12em;
        cursor: pointer;
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