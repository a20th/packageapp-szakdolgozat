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

<main class="container content p-4">
    <h2 class="text-center">Rendelések</h2>
    {#each data.orders as order}
        <h3 class="text-center">Rendelés {order.Id}</h3>
        <ul>
            {#each order.Packages as pack}
                <li
                    class="list-group-item d-flex justify-content-between align-items-start"
                >
                    <div class="ms-2 me-auto">
                        <div class="fw-bold">{status.status}</div>
                        {#if status.description}
                            {status.description}
                        {/if}
                    </div>
                    <span class="badge text-bg-primary rounded-pill"
                        >{status.date}</span
                    >
                </li>
            {/each}
        </ul>
        {/each}
</main>
