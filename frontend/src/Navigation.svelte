<script>
    import { onMount } from 'svelte';

    let activeLink;
    export let links;

    onMount(() => {
        activeLink = links[0];
    })
</script>

<ul>
    {#each links as link}
    <li>
        {#if activeLink == link}
            {link}
        {:else}
            <a on:click={() => activeLink = link} href="#{link}">{link}</a>
        {/if}
    </li>
    {/each}
</ul>

{#if activeLink == "battles"}
<slot name="battles">
</slot>
{/if}

{#if activeLink == "feed"}
<slot name="feed">
</slot>
{/if}

{#if activeLink == "battleDetails"}
<slot name="battleDetails">
</slot>
{/if}