
import { apiLocation } from '$lib/server/config';
import type { Actions } from '@sveltejs/kit';

export const load = async (event) => {

    const id = event.url.searchParams.get("id")
    if (id) {
        const url = apiLocation + "/track?id=" + id
        return await event.fetch(url, {
            method: "GET",
            headers: { "Content-Type": "application/json" }
        }).then(async (res) => {
            if (res.ok) {
                const resp = await res.json()
                if (!Array.isArray(resp)) {
                    return { found: false }
                }

                if ((resp as []).length == 0) {
                    return { found: false }
                }
                return { found: true, status: resp, packageId: id }
            }
        }).catch((e) => {
            console.log(e)
        })
    }
    return { search: true }
}

export const actions: Actions = {
    default: async (event) => {
        const formData = await event.request.formData();
        const status = formData.get("status")
        const id = formData.get("id")
        const description = formData.get("description")

        const body = JSON.stringify({id:id?.toString(), status:status?.toString(), description:description?.toString()})
        const url = apiLocation + "/admin/status"
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
        return {success: false}
    }).catch((e) => {
        console.log(e)
        return {success: false}
    })
    }
}