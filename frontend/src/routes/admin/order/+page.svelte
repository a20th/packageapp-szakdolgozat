<script lang="ts">
    import * as m from "$lib/paraglide/messages";
    import { getLocale } from "$lib/paraglide/runtime.js";
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
    import { countries } from "../../countries";
    import { enhance } from "$app/forms";
    const countrySort = (a: any, b: any) => {
        if (getLocale() === "en") {
            return a.en_name.localeCompare(b.en_name);
        }
        return a.hu_name.localeCompare(b.hu_name);
    };
</script>

{#if !data.found && !data.search}
    <Alert class="container content" dismissible id="error" color="danger">
        <h4 class="alert-heading text-capitalize">{m.error()}</h4>
        {m.package_not_found()}
    </Alert>
{/if}

{#if form?.error}
    <Alert class="container content" dismissible id="error" color="danger">
        <h4 class="alert-heading text-capitalize">{m.error()}</h4>
        {form.error}
    </Alert>
{:else if form?.success}
    <Alert class="container content" dismissible color="success">
        <h4 class="alert-heading text-capitalize">{m.success()}</h4>
    </Alert>
{/if}

<main class="container content p-4">
    {#if data.order}
        <h3 class="text-center mb-3">{m.order()} {data.orderId}</h3>
        <span><span class="text-danger">*</span> - {m.required_fields()}</span>
        <hr />
        <form
            action=""
            method="post"
            class=""
            use:enhance={() => {
                submitInProgress = true;

                return async ({ update }) => {
                    submitInProgress = false;
                    await update({ reset: false, invalidateAll: true });
                };
            }}
        >
            <input type="hidden" name="id" value={data.orderId} />
            <h3 class="mb-3">{m.invoice_address()}:</h3>
            <Input label={m.active()} class="mb-2" type="switch" name="active" bind:checked={data.order.Active}/>
            <div class="row">
                <div class="col col-sm-6">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            type="select"
                            name="country"
                            required
                        >
                            {#each countries.sort(countrySort) as country}
                                <option
                                    value={country.code}
                                    selected={country.code ==
                                        data.order.Country}
                                    >{getLocale() === "en"
                                        ? country.en_name
                                        : country.hu_name}</option
                                >
                            {/each}
                        </Input>
                        <div slot="label">
                            {m.country()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>
            <div class="row">
                <div class="col">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="zip"
                            bind:value={data.order.ZIPCode}
                            required
                        />
                        <div slot="label">
                            {m.postcode()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
                <div class="col-12 col-sm-8">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="city"
                            bind:value={data.order.City}
                            required
                        />
                        <div slot="label">
                            {m.city()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>
            <div class="row">
                <div class="col-12 col-sm-8">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="address"
                            bind:value={data.order.Address}
                            required
                        />
                        <div slot="label">
                            {m.address()}
                            <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
                <div class="col">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="address-number"
                            bind:value={data.order.Number}
                            required
                        />
                        <div slot="label">
                            {m.number()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>

            <div class="row">
                <div class="col-12 col-sm-6">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="name"
                            bind:value={data.order.Name}
                            required
                        />
                        <div slot="label">
                            {m.full_name()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
                <div class="col-12 col-sm-6">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="tax-number"
                            bind:value={data.order.TaxNumber}
                            
                        />
                        <div slot="label">
                            {m.tax_number()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>

        
            <Button type="submit" class="w-100 mt-2" disabled={submitInProgress}
                >{m.update()}</Button
            >
        </form>
    {:else}
        <h2 class="text-center mb-4">{m.order()}</h2>
        <FormGroup floating label={m.orderid()}>
            <Input bind:value={id} required />
        </FormGroup>
        <a
            href="/admin/order?id={id}"
            class="btn btn-secondary d-block w-100 {id ? '' : 'pe-none'}"
            >{m.search()}</a
        >
    {/if}
</main>
