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
        <h4 class="alert-heading text-capitalize">Error</h4>
        {form.error}
    </Alert>
{:else if form?.success}
    <Alert class="container content" dismissible color="success">
        <h4 class="alert-heading text-capitalize">Success</h4>
        Sikeresen leadtad a rendelést! Munkatársunk hamarosan felveszi veled a kapcsolatot.
    </Alert>
{/if}

<main class="container content">
    <div class="p-3">
        <h2 class="text-center">1. lépés: Adatok megadása</h2>
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
                    ><span class="text-danger">*</span> - Required fields to fill</span
                >
                <hr />
                <h3 class="mb-3">Feladó adatai:</h3>
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
                                name="from-address-number"
                                required
                            />
                            <div slot="label">
                                Number <span class="text-danger">*</span>
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
                            <div slot="label">Other</div>
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
                                Full Name <span class="text-danger">*</span>
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
                                Phone number <span class="text-danger">*</span>
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
                <h3 class="mb-3">Számlázási adatok:</h3>
                <Input
                    type="checkbox"
                    name="invoice"
                    class="mb-4"
                    style="min-width: 0px"
                    bind:checked={checkedState}
                    label="Számlázási cím megegyezik"
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
                                    name="invoice-zip"
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
                                    name="invoice-city"
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
                                    name="invoice-address"
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
                                    name="invoice-address-number"
                                    required
                                />
                                <div slot="label">
                                    Number <span class="text-danger">*</span>
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
                                    Full Name / Company <span
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
                                <div slot="label">Tax number</div>
                            </FormGroup>
                        </div>
                    </div>
                {/if}
                <hr />
                {#each packages as item, i}
                    <div class="row justify-content-between">
                        <div class="col"><h3>Csomag {i + 1}</h3></div>

                        <div class="col">
                            <Button
                                color="danger"
                                type="button"
                                on:click={() =>
                                    (packages = packages.filter(
                                        (_, index) => index != i,
                                    ))}
                                style="float: right;"
                                >Törlés <i class="bi bi-trash"></i></Button
                            >
                            <Button
                                color="primary"
                                type="button"
                                class="mx-2"
                                on:click={() => {
                                    let clone = { ...packages[i] };
                                    packages.push(clone);
                                }}
                                style="float: right;">Duplázás</Button
                            >
                        </div>
                    </div>
                    <p>Csomag méretei:</p>
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
                                        Length <span class="text-danger">*</span
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
                                        Width <span class="text-danger">*</span
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
                                        Height <span class="text-danger">*</span
                                        ></label
                                    >
                                </div>
                                <span class="input-group-text">cm</span>
                            </div>
                        </div>
                    </div>

                    <p>Címzett adatai:</p>
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
                                    name="zip-{i}"
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
                                    name="city-{i}"
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
                                    name="address-{i}"
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
                                    name="address-number-{i}"
                                    required
                                />
                                <div slot="label">
                                    Number <span class="text-danger">*</span>
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
                                <div slot="label">Other</div>
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
                                    Full Name <span class="text-danger">*</span>
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
                                    Phone number <span class="text-danger"
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
                    class="w-100 p-4 my-1">Csomag hozzáadása</Button
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
                    >Árkalkuláció</Button
                >
            </div>
            <div class={form?.prices ? "" : "d-none"}>
                <hr />
                <h2 class="text-center mb-3">
                    2. lépés: Rendelés véglegesítése
                </h2>
                <div class="row justify-content-between">
                    <div class="col"><h3>Árajánlat</h3></div>

                    <div class="col">
                        <Button
                            color="danger"
                            type="button"
                            on:click={() => {if(form?.prices) {form = null}}}
                            style="float: right;"
                            >Vissza <i class="bi bi-arrow-left"></i></Button
                        >
                    </div>
                </div>
                {#if form?.prices}
                    <Table>
                        {#each form.prices as price, i}
                            <tr><td>Csomag {i}</td><td>{price} HUF</td></tr>
                        {/each}
                        <tr
                            ><td class="fw-bold">Összesen</td><td
                                >{form.prices.reduce((acc, x) => acc + x)} HUF</td
                            ></tr
                        >
                    </Table>
                {/if}

                <Button
                    type="submit"
                    class="w-100 mt-2"
                    disabled={submitInProgress}
                    >Megerősítés</Button
                >
            </div>
        </form>
    </div>
</main>
