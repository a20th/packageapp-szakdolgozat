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
    {#if data.package}
        <h3 class="text-center mb-3">{m.package()} {data.packageId}</h3>
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
            <input type="hidden" name="id" value={data.packageId} />
            <p>{m.size_of_package()}:</p>
            <div class="row row-cols-2 row-cols-md-3">
                <div class="col">
                    <div class="input-group mb-3">
                        <div class="form-floating">
                            <input
                                type="number"
                                class="form-control"
                                placeholder="Length"
                                bind:value={data.package.Length}
                                name="length"
                                disabled={submitInProgress}
                                required
                            />
                            <label for="floatingInputGroup1">
                                {m.length()}
                                <span class="text-danger">*</span></label
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
                                bind:value={data.package.Width}
                                disabled={submitInProgress}
                                required
                            />
                            <label for="floatingInputGroup1">
                                {m.width()}
                                <span class="text-danger">*</span></label
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
                                bind:value={data.package.Height}
                                disabled={submitInProgress}
                                required
                            />
                            <label for="floatingInputGroup1">
                                {m.height()}
                                <span class="text-danger">*</span></label
                            >
                        </div>
                        <span class="input-group-text">cm</span>
                    </div>
                </div>
            </div>
            <h3 class="mb-3">{m.sender()}:</h3>
            <div class="row">
                <div class="col col-sm-6">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            type="select"
                            name="from-country"
                            required
                        >
                            {#each countries.sort(countrySort) as country}
                                <option
                                    value={country.code}
                                    selected={country.code ==
                                        data.package.FromCountry}
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
                            bind:value={data.package.FromZIP}
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
                            bind:value={data.package.FromCity}
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
                            bind:value={data.package.FromAddress}
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
                            name="from-address-number"
                            bind:value={data.package.FromNumber}
                            required
                        />
                        <div slot="label">
                            {m.number()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>
            <div class="row">
                <div class="col">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            bind:value={data.package.FromOther}
                            name="from-other"
                        />
                        <div slot="label">{m.other()}</div>
                    </FormGroup>
                </div>
            </div>

            <div class="row">
                <div class="col-12 col-sm-12">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="from-name"
                            bind:value={data.package.FromName}
                            required
                        />
                        <div slot="label">
                            {m.full_name()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>
            <div class="row">
                <div class="col-12 col-sm-6">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="from-phone"
                            bind:value={data.package.FromPhone}
                            required
                        />
                        <div slot="label">
                            {m.phone_number()}
                            <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
                <div class="col">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="from-email"
                            bind:value={data.package.FromEmail}
                            type="email"
                        />
                        <div slot="label">Email</div>
                    </FormGroup>
                </div>
            </div>
            <h3 class="mb-3">{m.recipient()}:</h3>
            <div class="row">
                <div class="col col-sm-6">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            type="select"
                            name="to-country"
                            required
                        >
                            {#each countries.sort(countrySort) as country}
                                <option
                                    value={country.code}
                                    selected={country.code ==
                                        data.package.ToCountry}
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
                            bind:value={data.package.ToZIP}
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
                            bind:value={data.package.ToCity}
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
                            bind:value={data.package.ToAddress}
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
                            name="to-address-number"
                            bind:value={data.package.ToNumber}
                            required
                        />
                        <div slot="label">
                            {m.number()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>
            <div class="row">
                <div class="col">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            bind:value={data.package.ToOther}
                            name="to-other"
                        />
                        <div slot="label">{m.other()}</div>
                    </FormGroup>
                </div>
            </div>

            <div class="row">
                <div class="col-12 col-sm-12">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="to-name"
                            bind:value={data.package.ToName}
                            required
                        />
                        <div slot="label">
                            {m.full_name()} <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
            </div>
            <div class="row">
                <div class="col-12 col-sm-6">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="to-phone"
                            bind:value={data.package.ToPhone}
                            required
                        />
                        <div slot="label">
                            {m.phone_number()}
                            <span class="text-danger">*</span>
                        </div>
                    </FormGroup>
                </div>
                <div class="col">
                    <FormGroup floating>
                        <Input
                            disabled={submitInProgress}
                            name="to-email"
                            bind:value={data.package.ToEmail}
                            type="email"
                        />
                        <div slot="label">Email</div>
                    </FormGroup>
                </div>
            </div>
            <Button type="submit" class="w-100 mt-2" disabled={submitInProgress}
                >{m.update()}</Button
            >
        </form>
    {:else}
        <h2 class="text-center mb-4">{m.package()}</h2>
        <FormGroup floating label={m.packageid()}>
            <Input bind:value={id} required />
        </FormGroup>
        <a
            href="/admin/package?id={id}"
            class="btn btn-secondary d-block w-100 {id ? '' : 'pe-none'}"
            >{m.search()}</a
        >
    {/if}
</main>
