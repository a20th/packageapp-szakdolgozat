import { apiLocation } from '$lib/server/config'
import type { Actions } from '@sveltejs/kit'
import { redirect } from '@sveltejs/kit'

interface DTO {
    kmprice: Number
    baseprice: Number
}

export const load = async (event) => {
    const url = apiLocation + "/admin/pricing"
    const resp = await event.fetch(url, {
        method: "GET",
        headers: { "Content-Type": "application/json", "Authorization": "Bearer " + event.cookies.get("access"), }
    }).then(async (res) => {
        if (res.ok) {
            const resp = await res.json() as DTO
            return { pricing: resp }
        }
        return { pricing: null }
    }).catch((e) => {
        console.log(e)
        return { pricing: null }
    })
    if (resp === null) {
        redirect(303, "/admin")
    }
    return resp
}

export const actions: Actions = {
    default: async (event) => {
        const formData = await event.request.formData();
        const baseprice = formData.get("baseprice")
        const kmprice = formData.get("kmprice")

        const dto = {} as DTO
        if(!baseprice || !kmprice){
            return {success: false}
        }
        dto.baseprice = Number.parseInt(baseprice?.toString())
        dto.kmprice = Number.parseInt(kmprice?.toString())
        const body = JSON.stringify(dto)
        const url = apiLocation + "/admin/pricing"
        return await event.fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + event.cookies.get("access"),
            },
            body: body
        }).then(async (res) => {
            if (res.ok) {
                return { success: true }
            }
            return { success: false }
        }).catch((e) => {
            console.log(e)
            return { success: false }
        })
    },
}