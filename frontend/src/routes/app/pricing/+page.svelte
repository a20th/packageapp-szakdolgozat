<script lang="ts">
    import {
        Button,
        Alert,
        FormGroup,
        Input,
        ButtonGroup,
        InputGroupText,
        InputGroup,
    } from "@sveltestrap/sveltestrap";
    import * as m from "$lib/paraglide/messages";
    import type { PageProps } from "./$types";
    import { enhance } from "$app/forms";
    let { data, form } = $props();
    import { setLocale, getLocale } from "$lib/paraglide/runtime.js";
    const countrySort = (a: any, b: any) => {
        if (getLocale() === "en") {
            return a.en_name.localeCompare(b.en_name);
        }
        return a.hu_name.localeCompare(b.hu_name);
    };

    let submitInProgress = $state(false);
</script>

<main class="container content">
    {#if form?.error}
            <Alert dismissible color="danger">
                <h4 class="alert-heading text-capitalize">Error</h4>
                {form.error}
            </Alert>
    {/if}
    <div class="p-3">
        <form
            action=""
            method="post"
            use:enhance={() => {
                submitInProgress = true;

                return async ({ update }) => {
                    await update({reset: false, invalidateAll: false});
                    submitInProgress = false;
                };
            }}
        >
            <h2 class="text-center mb-3">{m.quote()}</h2>
            <p>{m.size_of_package()}:</p>
            <div class="row row-cols-2 row-cols-md-3">
                <div class="col">
                    <div class="input-group mb-3">
                        <div class="form-floating">
                            <input
                                type="number"
                                class="form-control"
                                placeholder="Length"
                                name="length"
                                disabled={submitInProgress}
                                required
                            />
                            <label for="floatingInputGroup1">
                                {m.length()} <span class="text-danger">*</span></label
                            >
                        </div>
                        <span class="input-group-text">cm</span>
                    </div>
                </div>
                <div class="col">
                    <div class="input-group mb-3">
                        <div class="form-floating">
                            <input
                                type="number"
                                class="form-control"
                                placeholder="Width"
                                name="width"
                                disabled={submitInProgress}
                                required
                            />
                            <label for="floatingInputGroup1">
                                {m.width()} <span class="text-danger">*</span></label
                            >
                        </div>
                        <span class="input-group-text">cm</span>
                    </div>
                </div>
                <div class="col">
                    <div class="input-group mb-3">
                        <div class="form-floating">
                            <input
                                type="number"
                                class="form-control"
                                placeholder="Height"
                                name="height"
                                disabled={submitInProgress}
                                required
                            />
                            <label for="floatingInputGroup1">
                                {m.width()} <span class="text-danger">*</span></label
                            >
                        </div>
                        <span class="input-group-text">cm</span>
                    </div>
                </div>
            </div>
            <p>{m.sender()}:</p>
            <div class="row">
                <div class="col col-sm-6">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            type="select"
                            name="from-country"
                            required
                        >
                            {#each data.countries.sort(countrySort) as country}
                                <option
                                    value={country.code}
                                    selected={country.code == "HU"}
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
                            name="from-zip"
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
                            name="from-city"
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
                            name="from-address"
                            required
                        />
                        <div slot="label">
                            {m.address()} <span
                                class="text-danger">*</span
                            >
                        </div>
                    </FormGroup>
                </div>
                <div class="col">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="from-address-number"
                            required
                        />
                        <div slot="label">
                            {m.number()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>
            <p>{m.recipient()}:</p>
            <div class="row">
                <div class="col col-sm-6">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            type="select"
                            name="to-country"
                            required
                        >
                            {#each data.countries.sort(countrySort) as country}
                                <option
                                    value={country.code}
                                    selected={country.code == "HU"}
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
                            name="to-zip"
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
                            name="to-city"
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
                            name="to-address"
                            required
                        />
                        <div slot="label">
                            {m.address()} <span
                                class="text-danger">*</span
                            >
                        </div>
                    </FormGroup>
                </div>
                <div class="col">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="to-address-number"
                            required
                        />
                        <div slot="label">
                            {m.number()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>
            <Button type="submit" disabled={submitInProgress} class="w-100">{m.calculate_prices()}</Button>
        </form>
        {#if form?.success}
            <p>{m.price_of_delivery()}: {form.success} HUF</p>
        {/if}
    </div>
</main>
