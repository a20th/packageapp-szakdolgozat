// @ts-nocheck
import type { Actions, Load } from '@sveltejs/kit'
import {fail} from '@sveltejs/kit'
import * as m from "$lib/paraglide/messages"
import validator from 'validator';
import { apiLocation } from '$lib/server/config';

import {countries }from '$lib/server/countries'
export const load = () => {
    return {countries: countries}
}

export const actions = {
    default: async (event: import('./$types').RequestEvent) => {
        const formData = await event.request.formData();
		const length = formData.get('length');
		const width = formData.get('width');
        const height = formData.get('height');

        const toCountry = formData.get('to-country');
        const toZip = formData.get('to-zip');
        const toCity = formData.get('to-city');
        const toAddress = formData.get('to-address');
        const toAddressNumber = formData.get('to-address-number')

        const fromCountry = formData.get('from-country');
        const fromZip = formData.get('from-zip');
        const fromCity = formData.get('from-city');
        const fromAddress = formData.get('from-address');
        const fromAddressNumber = formData.get('from-address-number')
        
        if(length === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }

        if(width === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }

        if(height === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }


        if(toCountry === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }
        if(toZip === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }
        if(toCity === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }
        if(toAddress === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }
        if(toAddressNumber === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }

        if(fromCountry === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }
        if(fromZip === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }
        if(fromCity === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }
        if(fromAddress === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }
        if(fromAddressNumber === null){
            return fail(400, {error: m.form_login_error_no_email()})
        }

        if(!validator.isPostalCode(toZip.toString(), "any")){
            return fail(400, {error: m.form_login_error_no_email()})
        }
        if(!validator.isPostalCode(fromZip.toString(), "any")){
            return fail(400, {error: m.form_login_error_no_email()})
        }

        if(!validator.isInt(toAddressNumber.toString())){
            return fail(400, {error: m.form_login_error_no_email()})
        }
        if(!validator.isInt(fromAddressNumber.toString())){
            return fail(400, {error: m.form_login_error_no_email()})
        }

        if(!validator.isInt(height.toString())){
            return fail(400, {error: m.form_login_error_no_email()})
        }
        if(!validator.isInt(length.toString())){
            return fail(400, {error: m.form_login_error_no_email()})
        }
        if(!validator.isInt(width.toString())){
            return fail(400, {error: m.form_login_error_no_email()})
        }

        const parsedHeight = parseInt(height.toString())
        const parsedWidth = parseInt(width.toString())
        const parsedLength = parseInt(length.toString())
        
        const to = toAddressNumber + "," + toAddress + "," + toCity + "," + toZip + "," + toCountry
        const from = fromAddressNumber + "," + fromAddress + "," + fromCity + "," + fromZip + "," + fromCountry
        const size = parsedHeight + parsedWidth + parsedLength - Math.min(parsedHeight, parsedWidth, parsedLength)
        const body = JSON.stringify({from: from, to: to, size: size})

        try{
            const response = await event.fetch(apiLocation + "/price", {
                method: "POST",
                body: body,
                headers: {
                    "Content-Type": "application/json",
                },
            })
            const data = await response.json()

            if(!response.ok){
                if(data.error){
                    return fail(response.status, {error: data.error})
                }
                return fail(response.status, {error: "Internal Error"})
            }

            if(!data.price){
                console.log("Internal error, invalid price response" + data)
                return fail(response.status, {error: "Internal Error"})
            }

            return {success: data.price}
        }
        catch (e)
        {
            return fail(500, {error: "Internal error: " + e})
        }


    }

};null as any as Load;;null as any as Actions;