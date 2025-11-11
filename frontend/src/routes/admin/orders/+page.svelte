<script lang="ts">
    import * as m from "$lib/paraglide/messages";
    import {
        Button,
        Alert,
        FormGroup,
        Input,
        ButtonGroup,
        OffcanvasBackdrop,
        Accordion,
        AccordionItem,
    } from "@sveltestrap/sveltestrap";
    import { countries } from "../../countries";
    import { getLocale } from "$lib/paraglide/runtime.js";
    
    let { data } = $props();
    let id = $state();

    const translateCountry = (code: string) => {
        if (getLocale() == "en") {
            return countries.filter((value) => value.code == code)[0].en_name;
        }
        return countries.filter((value) => value.code == code)[0].hu_name;
    };
</script>

<main class="container content p-4">
    <h2 class="text-center mb-4">{m.orders()}</h2>
    <Accordion>
        {#each data.orders as order}
            <AccordionItem>
                <div class="m-0 w-100 d-block" slot="header">{m.order()} {order.Id} <span class="float-end m-0 me-4 {order.Active ? "green-dot" : "red-dot"}"></span></div>
                {m.client()}: {order.AccountName} {order.AccountEmail} {order.AccountPhone}
                <hr/>
                {m.invoice_address()}: {order.Name}
                {order.TaxNumber == undefined ? "" : order.TaxNumber}
                {translateCountry(order.Country)}, {order.ZIPCode}, {order.City},
                {order.Address}
                {order.Number},
                <a class="btn btn-secondary" href="/admin/order?id={order.Id}">{m.edit()}</a>
                <ul class="list-group mt-1 w-100">
                    {#each order.Packages as pack}
                        <li
                            class="list-group-item d-flex justify-content-between align-items-start p-3"
                        >
                            <div class="ms-2">
                                <div class="fw-bold">
                                    {m.package()}
                                    {pack.Id}
                                    <a class="btn btn-secondary" href="/admin/status?id={pack.Id}">{m.status()}</a>
                                    <a class="btn btn-secondary" href="/admin/package?id={pack.Id}">{m.edit()}</a>
                                </div>
                                {m.size_of_package()}: {pack.Height} x {pack.Width}
                                x {pack.Length}
                                cm
                                <hr />
                                {m.sender()}: {pack.FromName}
                                {pack.FromPhone}
                                {pack.FromEmail}
                                <br />
                                {translateCountry(pack.FromCountry)}, {pack.FromZIP},
                                {pack.FromCity},
                                {pack.FromAddress}
                                {pack.FromNumber}, {pack.FromOther}
                                <hr />
                                {m.recipient()}: {pack.ToName}
                                {pack.ToPhone}
                                {pack.ToEmail}
                                <br />
                                {translateCountry(pack.ToCountry)}, {pack.ToZIP},
                                {pack.ToCity}, {pack.ToAddress}
                                {pack.ToNumber}, {pack.ToOther}
                                
                            </div>
                        </li>
                    {/each}
                </ul>
            </AccordionItem>
        {/each}
    </Accordion>
</main>
