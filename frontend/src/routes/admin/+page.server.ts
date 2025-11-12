
import { fail, redirect } from '@sveltejs/kit';
import validator from 'validator';
import type { Actions } from '../app/$types';
import * as m from "$lib/paraglide/messages"
import { refreshAll } from '$app/navigation';
import { apiLocation } from '$lib/server/config';

export const load = async (event) => {

    
}


export const actions = {
    login: async (event) => {
        const formData = await event.request.formData();
        const username = formData.get('username')
        const password = formData.get('password')
        if (username === null) {
            return fail(400, { error: m.form_login_error_no_email() })
        }
        if (password === null) {
            return fail(400, { error: m.form_login_error_no_password() })
        }

        const body = JSON.stringify({ "username": username, "password": password })
        try {
            const response = await event.fetch(apiLocation + "/admin/login", {
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

            if (!(data.access_token && data.refresh_token)) {
                console.log("Internal error, invalid login response" + data)
                return fail(response.status, { error: "Internal Error" })
            }

            const accessToken = data.access_token
            const refreshToken = data.refresh_token

            event.cookies.set("access", accessToken, { path: "/admin" })
            event.cookies.set("refresh", refreshToken, { path: "/admin" })
        }
        catch (e) {
            return fail(500, { error: "Internal error: " + e })
        }

        return { success: true }
    },
    logout: async (event) => {
        const url = apiLocation + "/admin/logout"
        event.fetch(url, {
            method: "POST",
            headers: {
                "Authorization": "Bearer " + event.locals.access
            }
        }).then((r) => {
        }).catch((e) => {
            console.log(e)
        })
        event.locals.admin!.access = null
        event.locals.admin!.refresh = null
        event.locals.admin!.loggedIn = false
        event.cookies.delete("access", { path: "/admin" })
        event.cookies.delete("refresh", { path: "/admin" })
        return { success: "Logged out successfully" }
    },

} satisfies Actions;