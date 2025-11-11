
import { apiLocation } from '$lib/server/config';
import type { Actions } from '@sveltejs/kit';
import type { Order, Package } from '../../app/order/+page.server.js';
import { fail, isActionFailure } from '@sveltejs/kit';
import type { ActionFailure } from '@sveltejs/kit';
import validator from 'validator';


declare interface OrderDTO {
    Id: string;
    Name: string;
    TaxNumber?: string;
    ZIPCode: string;
    City: string;
    Country: string;
    Address: string;
    Number: string;
    Active: boolean
}

export const load = async (event) => {

    const id = event.url.searchParams.get("id")
    if (id) {
        const url = apiLocation + "/admin/order?id=" + id
        return await event.fetch(url, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + event.cookies.get("access")
            },

        }).then(async (res) => {
            if (res.ok) {
                const resp = await res.json() as OrderDTO
                return { found: true, order: resp, orderId: id }
            }
        }).catch((e) => {
            console.log(e)
        })
    }
    return { search: true }
}

function createOrderFromFormData(formData: FormData): OrderDTO | ActionFailure<{ error: string }> | ActionFailure {
    const fromCity = formData.get("city")
    const fromZip = formData.get('zip');
    const fromCountry = formData.get('country');
    const fromAddress = formData.get('address');
    const fromAddressNumber = formData.get('address-number')
    const fromName = formData.get('name')
    const fromTaxNumber = formData.get('tax-number')
    const id = formData.get('id')
    const active = formData.get('active')
    if ([fromCity, fromAddress, fromAddressNumber, fromName, fromCountry, fromZip].includes(null)) {
        return fail(400, { error: "form fields missing" })
    }
    
    let order = {} as OrderDTO
    order.Active = active ? true : false
    order.Address = fromAddress!.toString()
    order.City = fromCity!.toString()
    order.Country = fromCountry!.toString()
    order.Number = fromAddressNumber!.toString()
    order.TaxNumber = fromTaxNumber ? fromTaxNumber.toString() : ""
    order.ZIPCode = fromZip!.toString()
    order.Name = fromName!.toString()
    order.Id = id!.toString()

    return order
}


export const actions: Actions = {
    default: async (event) => {
        const formData = await event.request.formData();

        const order = createOrderFromFormData(formData)
        if (isActionFailure(order)) {
            return order
        }
        const body = JSON.stringify(order)

        try {
            const response = await event.fetch(apiLocation + "/admin/order", {
                method: "POST",
                body: body,
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": "Bearer " + event.cookies.get("access")
                },
            })
            const data = await response.json()
            if (!response.ok) {
                if (data.error) {
                    return fail(response.status, { error: data.error })
                }
                return fail(response.status, { error: "Internal Error" })
            }
            return { success: true }
        }
        catch (e) {
            return fail(500, { error: "Internal error: " + e })
        }
    }

}