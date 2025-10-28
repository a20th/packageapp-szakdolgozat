<script lang="ts">
    import * as m from "$lib/paraglide/messages";
    import {
        Button,
        Alert,
        FormGroup,
        Input,
        ButtonGroup,
    } from "@sveltestrap/sveltestrap";
    let { data } = $props();
    let id = $state();
</script>
<main class="container">
<h2 class="text-center">{m.app_navbar_tracking()}</h2>

{#if data.status}
    <h3>Package {data.packageId}</h3>
    {#each data.status as status}
        <h4>{status.status} - {status.date}</h4>
        {#if status.description}
            <p>{status.description}</p>
        {/if}
    {/each}
{:else}
    <FormGroup floating label="Package ID">
        <Input
            bind:value={id}
            required
        />
    </FormGroup>
    <a href="/app/tracking?id={id}" class="btn btn-secondary d-block w-100">Track</a>
{/if}
</main>
