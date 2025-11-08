<script lang="ts">
    import {
        Button,
        Alert,
        FormGroup,
        Input,
        Table,
        ButtonGroup,
        InputGroupText,
        InputGroup,
    } from "@sveltestrap/sveltestrap";
    import * as m from "$lib/paraglide/messages";
    import type { PageProps } from "./$types";
    import { enhance } from "$app/forms";
    import type { Package } from "./+page.server";
    import { getLocale } from "$lib/paraglide/runtime.js";
    let { data, form } = $props();

    let packages = $state(new Array<Package>());
    (() => packages.push({} as Package))();
    let submitInProgress = $state(false);
    let checkedState = $state(true);

    const countrySort = (a: any, b: any) => {
        if (getLocale() === "en") {
            return a.en_name.localeCompare(b.en_name);
        }
        return a.hu_name.localeCompare(b.hu_name);
    };

    function scrollIntoView() {
        const el = document.querySelector("#error");
        if (!el) return;
        el.scrollIntoView({
            behavior: "smooth",
        });
    }
</script>

{#if form?.error}
    <Alert class="container content" dismissible id="error" color="danger">
        <h4 class="alert-heading text-capitalize">{m.error()}</h4>
        {form.error}
    </Alert>
{:else if form?.success}
    <Alert class="container content" dismissible color="success">
        <h4 class="alert-heading text-capitalize">{m.success()}</h4>
        {m.successful_order_message()}
    </Alert>
{/if}

<main class="container content">
    <div class="p-3">
        <h2 class="text-center">{m.order_step1()}</h2>
        <form
            action="{form?.prices ? "/app/order?/confirm" : "/app/order?/calculate"}"
            method="post"
            class="mt-3"
            use:enhance={() => {
                submitInProgress = true;

                return async ({ update, formData }) => {
                    formData.append("countTest", packages.length.toString());
                    submitInProgress = false;
                    if (form?.error) {
                        scrollIntoView();
                    }
                    await update({ reset: false, invalidateAll: false });                
                };
            }}
        >
            <div class={form?.prices ? "d-none" : ""}>
                <span
                    ><span class="text-danger">*</span> - {m.required_fields()}</span
                >
                <hr />
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
                            <div slot="label">{m.address()}
                                 <span
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
                <div class="row">
                    <div class="col">
                        <FormGroup floating>
                            <Input
                                disabled={submitInProgress}
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
                                required
                            />
                            <div slot="label">
                                {m.phone_number()} <span class="text-danger">*</span>
                            </div>
                        </FormGroup>
                    </div>
                    <div class="col">
                        <FormGroup floating>
                            <Input
                                disabled={submitInProgress}
                                name="from-email"
                                type="email"
                            />
                            <div slot="label">Email</div>
                        </FormGroup>
                    </div>
                </div>
                <hr />
                <h3 class="mb-3">{m.invoice_address()}:</h3>
                <Input
                    type="checkbox"
                    name="invoice"
                    class="mb-4"
                    style="min-width: 0px"
                    bind:checked={checkedState}
                    label={m.invoice_address_same()}
                />
                {#if !checkedState}
                    <div class="row">
                        <div class="col col-sm-6">
                            <FormGroup floating>
                                <Input
                                    disabled={submitInProgress}
                                    name="invoice-country"
                                    type="select"
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
                                    name="invoice-zip"
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
                                    name="invoice-city"
                                    required
                                />
                                <div slot="label">
                                    {m.city()}  <span class="text-danger">*</span>
                                </div>
                            </FormGroup>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-12 col-sm-8">
                            <FormGroup floating>
                                <Input
                                    disabled={submitInProgress}
                                    name="invoice-address"
                                    required
                                />
                                <div slot="label">
                                    {m.address()}  <span
                                        class="text-danger">*</span
                                    >
                                </div>
                            </FormGroup>
                        </div>
                        <div class="col">
                            <FormGroup floating>
                                <Input
                                    disabled={submitInProgress}
                                    name="invoice-address-number"
                                    required
                                />
                                <div slot="label">
                                    {m.number()}  <span class="text-danger">*</span>
                                </div>
                            </FormGroup>
                        </div>
                    </div>

                    <div class="row">
                        <div class="col-12 col-sm-6">
                            <FormGroup floating>
                                <Input
                                    disabled={submitInProgress}
                                    name="invoice-name"
                                    required
                                />
                                <div slot="label">
                                    {m.full_name_company()} <span
                                        class="text-danger">*</span
                                    >
                                </div>
                            </FormGroup>
                        </div>
                        <div class="col-12 col-sm-6">
                            <FormGroup floating>
                                <Input
                                    disabled={submitInProgress}
                                    name="invoice-tax-number"
                                    required
                                />
                                <div slot="label">{m.tax_number()}</div>
                            </FormGroup>
                        </div>
                    </div>
                {/if}
                <hr />
                {#each packages as item, i}
                    <div class="row justify-content-between">
                        <div class="col"><h3>{m.package()} {i + 1}</h3></div>

                        <div class="col">
                            <Button
                                color="danger"
                                type="button"
                                on:click={() =>
                                    (packages = packages.filter(
                                        (_, index) => index != i,
                                    ))}
                                style="float: right;"
                                >{m.delete()} <i class="bi bi-trash"></i></Button
                            >
                            <Button
                                color="primary"
                                type="button"
                                class="mx-2"
                                on:click={() => {
                                    let clone = { ...packages[i] };
                                    packages.push(clone);
                                }}
                                style="float: right;">{m.duplicate()}</Button
                            >
                        </div>
                    </div>
                    <p>{m.size_of_package()}:</p>
                    <div class="row row-cols-2 row-cols-md-3">
                        <div class="col">
                            <div class="input-group mb-3">
                                <div class="form-floating">
                                    <input
                                        type="number"
                                        class="form-control"
                                        placeholder="Length"
                                        bind:value={item.Length}
                                        name="length-{i}"
                                        disabled={submitInProgress}
                                        required
                                    />
                                    <label for="floatingInputGroup1">
                                        {m.length()} <span class="text-danger">*</span
                                        ></label
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
                                        name="width-{i}"
                                        bind:value={item.Width}
                                        disabled={submitInProgress}
                                        required
                                    />
                                    <label for="floatingInputGroup1">
                                        {m.width()} <span class="text-danger">*</span
                                        ></label
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
                                        name="height-{i}"
                                        bind:value={item.Height}
                                        disabled={submitInProgress}
                                        required
                                    />
                                    <label for="floatingInputGroup1">
                                        {m.height()} <span class="text-danger">*</span
                                        ></label
                                    >
                                </div>
                                <span class="input-group-text">cm</span>
                            </div>
                        </div>
                    </div>

                    <p>{m.recipient()}:</p>
                    <div class="row">
                        <div class="col col-sm-6">
                            <FormGroup floating>
                                <Input
                                    disabled={submitInProgress}
                                    type="select"
                                    name="country-{i}"
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
                                    name="zip-{i}"
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
                                    name="city-{i}"
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
                                    name="address-{i}"
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
                                    name="address-number-{i}"
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
                                    name="other-{i}"
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
                                    name="name-{i}"
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
                                    name="phone-{i}"
                                    required
                                />
                                <div slot="label">
                                    {m.phone_number()} <span class="text-danger"
                                        >*</span
                                    >
                                </div>
                            </FormGroup>
                        </div>
                        <div class="col">
                            <FormGroup floating>
                                <Input
                                    disabled={submitInProgress}
                                    name="email-{i}"
                                    type="email"
                                />
                                <div slot="label">Email</div>
                            </FormGroup>
                        </div>
                    </div>
                    <hr />
                {/each}
                <Button
                    type="button"
                    disabled={submitInProgress}
                    on:click={() => packages.push({} as Package)}
                    class="w-100 p-4 my-1">{m.add_package()}</Button
                >
                <input
                    type="hidden"
                    name="count"
                    bind:value={packages.length}
                />
                <Button
                    type="submit"
                    class="w-100 mb-1"
                    disabled={submitInProgress || packages.length == 0}
                    >{m.quote()}</Button
                >
            </div>
            <div class={form?.prices ? "opacity-100" : "d-none opacity-0"}>
                <hr />
                <h2 class="text-center mb-3">
                    {m.order_step2()}
                </h2>
                <div class="row justify-content-between">
                    <div class="col"><h3>{m.quote()}</h3></div>

                    <div class="col">
                        <Button
                            color="danger"
                            type="button"
                            on:click={() => {if(form?.prices) {form = null}}}
                            style="float: right;"
                            >{m.back()} <i class="bi bi-arrow-left"></i></Button
                        >
                    </div>
                </div>
                {#if form?.prices}
                    <Table>
                        {#each form.prices as price, i}
                            <tr><td>{m.package()} {i + 1}</td><td>{price} HUF</td></tr>
                        {/each}
                        <tr
                            ><td class="fw-bold">{m.total()}</td><td
                                >{form.prices.reduce((acc, x) => acc + x)} HUF</td
                            ></tr
                        >
                    </Table>
                {/if}

                <Button
                    type="submit"
                    class="w-100 mt-2"
                    disabled={submitInProgress}
                    >{m.confirm_order()}</Button
                >
            </div>
        </form>
    </div>
</main>
