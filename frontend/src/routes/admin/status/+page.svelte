<script lang="ts">
    import { enhance } from "$app/forms";
    import * as m from "$lib/paraglide/messages";
    import {
        Button,
        Alert,
        FormGroup,
        Input,
        ButtonGroup,
    } from "@sveltestrap/sveltestrap";
    let { data, form } = $props();
    let id = $state();
    let submitInProgress = $state(false);

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

{#if form?.success}
    <Alert class="container content" dismissible id="error" color="success">
        <h4 class="alert-heading text-capitalize">{m.success()}</h4>
    </Alert>

{:else if form?.success === false}
    <Alert class="container content" dismissible id="error" color="danger">
        <h4 class="alert-heading text-capitalize">{m.error()}</h4>
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
                        <div class="fw-bold">
                            {translateStatus(status.status)}
                        </div>
                        {status.description}
                    </div>
                    <span class="badge text-bg-primary rounded-pill"
                        >{new Date(
                            Date.parse(status.date),
                        ).toLocaleString()}</span
                    >
                </li>
            {/each}
        </ul>

        <!-- svelte-ignore component_name_lowercase -->
        <form
            action=""
            method="post"
            class="mt-3"
            use:enhance={() => {
                submitInProgress = true;

                return async ({ update }) => {
                    submitInProgress = false;
                    await update({ reset: false, invalidateAll: true });
                };
            }}
        >
            <span
                ><span class="text-danger">*</span> - {m.required_fields()}</span
            >
            <hr />

            <div class="row">
                <div class="col">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="status"
                            type="select"
                            required
                        >
                        <option value="confirmed">{translateStatus("confirmed")}</option>
                        <option value="shipping">{translateStatus("shipping")}</option>
                        <option value="shipping-fail">{translateStatus("shipping-fail")}</option>
                        <option value="shipped">{translateStatus("shipped")}</option>
                    </Input>
                    <div slot="label">
                            {m.status()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
                <div class="col-12 col-sm-8">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="description"
                        />
                        <div slot="label">
                            {m.description()}
                        </div>
                    </FormGroup>
                </div>
                <input type="hidden" name="id" value={data.packageId}/>
                <Button
                    type="submit"
                    class="w-100 mt-2"
                    disabled={submitInProgress}
                    >{m.update()}</Button
                >
            </div>
        </form>
    {:else}
        <FormGroup floating label="Package ID">
            <Input bind:value={id} required />
        </FormGroup>
        <a
            href="/admin/status?id={id}"
            class="btn btn-secondary d-block w-100 {id ? '' : 'pe-none'}"
            >{m.track()}</a
        >
    {/if}
</main>
