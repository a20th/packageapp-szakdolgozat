import type { ActionFailure, Actions, Load } from '@sveltejs/kit'
import { fail, isActionFailure } from '@sveltejs/kit'
import * as m from "$lib/paraglide/messages"
import validator from 'validator';
import { apiLocation } from '$lib/server/config';
import { countries } from '$lib/server/countries'
export interface Order {
    Name: string
    TaxNumber: string
    ZIPCode: string
    Country: string
    City: string
    Address: string
    Number: string
    Packages: Array<Package>
}

export interface Package {
    Length: number
    Width: number
    Height: number
    ToName: string
    ToPhone: string
    ToEmail: string
    ToCountry: string
    ToZIP: string
    ToCity: string
    ToAddress: string
    ToNumber: string
    ToOther: string
    FromName: string
    FromPhone: string
    FromEmail: string
    FromCountry: string
    FromZIP: string
    FromCity: string
    FromAddress: string
    FromNumber: string
    FromOther: string
}

export const load: Load = () => {
    return { countries: countries }
}

function createOrderFromFormData(formData: FormData): Order | ActionFailure<{ error: string }> | ActionFailure {
    const count = formData.get("count")
    const checked = formData.get("invoice")
    if (!count) {
        return fail(400)
    }


    const numberOfPackages = parseInt(count.toString())
    let order = {} as Order
    order.Packages = new Array<Package>(numberOfPackages)

    const fromCity = formData.get("from-city")
    const fromZip = formData.get('from-zip');
    const fromCountry = formData.get('from-country');
    const fromAddress = formData.get('from-address');
    const fromAddressNumber = formData.get('from-address-number')
    const fromName = formData.get('from-name')
    const fromOther = formData.get('from-other')
    const fromEmail = formData.get('from-email')
    const fromPhone = formData.get('from-phone')

    if ([fromCity, fromAddress, fromAddressNumber, fromName, fromCountry, fromZip].includes(null)) {
        return fail(400, { error: "from fields missing" })
    }

    order.City = fromCity!.toString()
    order.Address = fromAddress!.toString()
    order.Country = fromCountry!.toString()
    order.Name = fromName!.toString()
    order.Number = fromAddressNumber!.toString()
    order.ZIPCode = fromZip!.toString()

    if (!checked) {
        const invoiceCity = formData.get("invoice-city")
        const invoiceZip = formData.get('invoice-zip');
        const invoiceCountry = formData.get('invoice-country');
        const invoiceAddress = formData.get('invoice-address');
        const invoiceAddressNumber = formData.get('invoice-address-number')
        const invoiceName = formData.get('invoice-name')
        const invoiceTaxNumber = formData.get('invoice-tax-number')

        if ([invoiceCity, invoiceAddress, invoiceAddressNumber, invoiceName, invoiceCountry, invoiceZip].includes(null)) {
            return fail(400, { error: "invoice fields missing" })
        }

        order.City = invoiceCity!.toString()
        order.Address = invoiceAddress!.toString()
        order.Country = invoiceCountry!.toString()
        order.Name = invoiceName!.toString()
        order.Number = invoiceAddressNumber!.toString()
        order.ZIPCode = invoiceZip!.toString()
        order.TaxNumber = invoiceTaxNumber ? invoiceTaxNumber.toString() : ""
    }




    for (let i = 0; i < numberOfPackages; i++) {
        const _package = {} as Package
        const length = formData.get('length-' + i);
        const width = formData.get('width-' + i);
        const height = formData.get('height-' + i);
        const toCountry = formData.get('country-' + i);
        const toZip = formData.get('zip-' + i);
        const toCity = formData.get('city-' + i);
        const toAddress = formData.get('address-' + i);
        const toAddressNumber = formData.get('address-number-' + i)
        const toName = formData.get('name-' + i)
        const toPhone = formData.get('phone-' + i)

        const toOther = formData.get('other-' + i)
        const toEmail = formData.get('email-' + i)

        if ([length, width, height, toCountry, toZip, toCity, toAddress, toAddressNumber, toName, toPhone].includes(null)) {
            return fail(400, { error: "to fields missing" })
        }

        if (!validator.isMobilePhone(toPhone!.toString())) {
            return fail(400, { error: "invalid mobile" })
        }

        if (toEmail && !validator.isEmail(toEmail.toString())) {
            return fail(400, { error: "invalid email entered" })
        }

        if (!validator.isPostalCode(toZip!.toString(), 'any')) {
            return fail(400, { error: "invalid zip" })
        }

        if (!validator.isInt(height!.toString())) {
            return fail(400, { error: "height not a number" })
        }
        if (!validator.isInt(length!.toString())) {
            return fail(400, { error: "length not a number" })
        }
        if (!validator.isInt(width!.toString())) {
            return fail(400, { error: "width not a number" })
        }
        if (!validator.isInt(toAddressNumber!.toString())) {
            return fail(400, { error: "addr number not a number" })
        }

        _package.Height = parseInt(height!.toString())
        _package.Length = parseInt(length!.toString())
        _package.Width = parseInt(width!.toString())
        _package.ToCountry = toCountry!.toString()
        _package.ToZIP = toZip!.toString()
        _package.ToCity = toCity!.toString()
        _package.ToAddress = toAddress!.toString()
        _package.ToNumber = toAddressNumber!.toString()
        _package.ToOther = toOther ? toOther.toString() : ""
        _package.ToEmail = toEmail ? toEmail.toString() : ""
        _package.ToPhone = toPhone!.toString()
        _package.ToName = toName!.toString()

        _package.FromName = fromName!.toString()
        _package.FromCountry = fromCountry!.toString()
        _package.FromZIP = fromZip!.toString()
        _package.FromCity = fromCity!.toString()
        _package.FromAddress = fromAddress!.toString()
        _package.FromNumber = fromAddressNumber!.toString()
        _package.FromOther = fromOther ? fromOther.toString() : ""
        _package.FromEmail = fromEmail ? fromEmail.toString() : ""
        _package.FromPhone = fromPhone!.toString()

        order.Packages[i] = _package
    }

    return order
}

function fromOneLine(element: Package): string {
    return element.FromNumber + "," + element.FromAddress + "," + element.FromCity + "," + element.FromZIP + "," + element.FromCountry
}

function toOneLine(element: Package): string {
    return element.ToNumber + "," + element.ToAddress + "," + element.ToCity + "," + element.ToZIP + "," + element.ToCountry
}

export const actions: Actions = {
    calculate: async (event) => {
        const formData = await event.request.formData();
        const order = createOrderFromFormData(formData)
        if (isActionFailure(order)) {
            return order
        }

        let prices = new Array<number>((order as Order).Packages.length)
        for (let index = 0; index < (order as Order).Packages.length; index++) {
            const element = (order as Order).Packages[index]
            const parsedHeight = element.Height
            const parsedWidth = element.Width
            const parsedLength = element.Length

            const to = toOneLine(element)
            const from = fromOneLine(element)
            const size = parsedHeight + parsedWidth + parsedLength - Math.min(parsedHeight, parsedWidth, parsedLength)
            const body = JSON.stringify({ from: from, to: to, size: size })

            try {
                const response = await event.fetch(apiLocation + "/price", {
                    method: "POST",
                    body: body,
                    headers: {
                        "Content-Type": "application/json",
                    },
                })
                const data = await response.json()

                if (!response.ok) {
                    if (data.error) {
                        return fail(response.status, { error: data.error })
                    }
                    return fail(response.status, { error: "Internal Error" })
                }

                if (!data.price) {
                    console.log("Internal error, invalid price response" + data)
                    return fail(response.status, { error: "Internal Error" })
                }

                prices[index] = data.price
            }
            catch (e) {
                return fail(500, { error: "Internal error: " + e })
            }
            return {prices: prices}
        }
    },
    confirm: async (event) => {
        const formData = await event.request.formData();

        const order = createOrderFromFormData(formData)
        if (isActionFailure(order)) {
            return order
        }

        const body = JSON.stringify(order)

        try {
            const response = await event.fetch(apiLocation + "/order/create", {
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