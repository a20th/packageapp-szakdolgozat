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
            <h2 class="text-center mb-3">Árkalkulátor</h2>
            <p>Csomag méretei:</p>
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
                                Length <span class="text-danger">*</span></label
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
                                Width <span class="text-danger">*</span></label
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
                                Height <span class="text-danger">*</span></label
                            >
                        </div>
                        <span class="input-group-text">cm</span>
                    </div>
                </div>
            </div>
            <p>Feladó:</p>
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
                            Country <span class="text-danger">*</span>
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
                            Postcode <span class="text-danger">*</span>
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
                            City <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>
            <div class="row">
                <div class="col-12 col-sm-8">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            autocomplete="email"
                            name="from-address"
                            required
                        />
                        <div slot="label">
                            Name and type of public space <span
                                class="text-danger">*</span
                            >
                        </div>
                    </FormGroup>
                </div>
                <div class="col">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            autocomplete="email"
                            name="from-address-number"
                            required
                        />
                        <div slot="label">
                            Number <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>
            <p>Címzett:</p>
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
                            Country <span class="text-danger">*</span>
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
                            Postcode <span class="text-danger">*</span>
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
                            City <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>
            <div class="row">
                <div class="col-12 col-sm-8">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            autocomplete="email"
                            name="to-address"
                            required
                        />
                        <div slot="label">
                            Name and type of public space <span
                                class="text-danger">*</span
                            >
                        </div>
                    </FormGroup>
                </div>
                <div class="col">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            autocomplete="email"
                            name="to-address-number"
                            required
                        />
                        <div slot="label">
                            Number <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>
            <Button type="submit" class="w-100">Calculate price</Button>
        </form>
        {#if form?.success}
            <p>Price of delivery: {form.success} HUF</p>
        {/if}
    </div>
</main>
