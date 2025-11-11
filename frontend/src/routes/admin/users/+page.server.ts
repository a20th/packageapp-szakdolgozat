import { apiLocation } from '$lib/server/config'
import type { Actions } from '@sveltejs/kit'
import { redirect } from '@sveltejs/kit'

interface DTO {
    admins: Array<string>
}

export const load = async (event) => {
    const url = apiLocation + "/admin/getall"
    const resp = await event.fetch(url, {
        method: "GET",
        headers: { "Content-Type": "application/json", "Authorization": "Bearer " + event.cookies.get("access"), }
    }).then(async (res) => {
        if (res.ok) {
            const resp = await res.json()
            return { admins: resp }
        }
        return { admins: null }
    }).catch((e) => {
        console.log(e)
        return { admins: null }
    })
    if (resp.admins === null) {
        redirect(303, "/admin")
    }
    return resp
}

export const actions: Actions = {
    delete: async (event) => {
        const formData = await event.request.formData();
        const username = formData.get("username")
        const url = apiLocation + "/admin/user?id=" + username?.toString()
        return await event.fetch(url, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + event.cookies.get("access"),
            },
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
    create: async (event) => {
        const formData = await event.request.formData();
        const username = formData.get("username")
        const password = formData.get("password")

        const body = JSON.stringify({ username: username?.toString(), password: password?.toString(), })
        console.log(body)
        const url = apiLocation + "/admin/user"
        return await event.fetch(url, {
            method: "PUT",
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