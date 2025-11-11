
import { apiLocation } from '$lib/server/config';
import type { Actions } from '@sveltejs/kit';
import type { Package } from '../../app/order/+page.server.js';
import { fail, isActionFailure } from '@sveltejs/kit';
import type { ActionFailure } from '@sveltejs/kit';
import validator from 'validator';

export const load = async (event) => {

    const id = event.url.searchParams.get("id")
    if (id) {
        const url = apiLocation + "/admin/package?id=" + id
        return await event.fetch(url, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + event.cookies.get("access")
            },

        }).then(async (res) => {
            if (res.ok) {
                const resp = await res.json() as Package
                return { found: true, package: resp, packageId: id }
            }
        }).catch((e) => {
            console.log(e)
        })
    }
    return { search: true }
}

function createPackageFromFormData(formData: FormData): Package | ActionFailure<{ error: string }> | ActionFailure {
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
        return fail(400, { error: "form fields missing" })
    }
    const _package = {} as Package
    const length = formData.get('length');
    const width = formData.get('width' );
    const height = formData.get('height' );
    const toCountry = formData.get('to-country' );
    const toZip = formData.get('to-zip' );
    const toCity = formData.get('to-city' );
    const toAddress = formData.get('to-address' );
    const toAddressNumber = formData.get('to-address-number' )
    const toName = formData.get('to-name' )
    const toPhone = formData.get('to-phone' )
    const toOther = formData.get('to-other' )
    const toEmail = formData.get('to-email' )

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

    return _package
}


export const actions: Actions = {
    default: async (event) => {
        const formData = await event.request.formData();

        const _package = createPackageFromFormData(formData)
        if (isActionFailure(_package)) {
            return _package
        }
        const data = _package as any
        data.id = formData.get("id")?.toString()
        const body = JSON.stringify(data)

        try {
            const response = await event.fetch(apiLocation + "/admin/package", {
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