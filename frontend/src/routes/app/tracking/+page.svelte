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

    const translateStatus: Function = (value: string) => {
        switch (value) {
            case "confirmation":
                return m.status_confirmation();

            case "confirmed":
                return m.status_confirmed()

            case "shipping":
                return m.status_shipping()

            case "shipping-fail":
                return m.status_shipping_fail()

            case "shipped":
                return m.status_shipped()

            default:
                return value;
        }
    };

</script>

{#if !data.found && !data.search}
    <Alert class="container content" dismissible id="error" color="danger">
        <h4 class="alert-heading text-capitalize">{m.error()}</h4>
        {m.package_not_found()}
    </Alert>
{/if}

<main class="container content p-4">
    <h2 class="text-center mb-4">{m.app_navbar_tracking()}</h2>

    {#if data.status}
        <h3 class="text-center mb-3">{m.package()} {data.packageId}</h3>
        <ul class="list-group">
            {#each data.status as status}
                <li
                    class="list-group-item d-flex justify-content-between align-items-start"
                >
                    <div class="ms-2 me-auto">
                        <div class="fw-bold">{translateStatus(status.status)}</div>
                        {status.description}
                    </div>
                    <span class="badge text-bg-primary rounded-pill"
                        >{new Date(Date.parse(status.date)).toLocaleString()}</span
                    >
                </li>
            {/each}
        </ul>
    {:else}
        <FormGroup floating label={m.packageid()}>
            <Input bind:value={id} required />
        </FormGroup>
        <a href="/app/tracking?id={id}" class="btn btn-secondary d-block w-100 {id ? "" : "pe-none"}"
            >{m.track()}</a
        >
    {/if}
</main>
