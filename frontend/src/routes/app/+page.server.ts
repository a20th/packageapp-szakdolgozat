
import { fail, redirect } from '@sveltejs/kit';
import validator from 'validator';
import type { Actions } from './$types';
import * as m from "$lib/paraglide/messages"
import { refreshAll } from '$app/navigation';
import { apiLocation } from '$lib/server/config';

export const load = async (event) => {

    const verify = event.url.searchParams.get("verify")
    const email = event.url.searchParams.get("email")
    if (verify && email) {
        const url = apiLocation + "/verify"
        const body = JSON.stringify({ email: email, code: verify })
        const result = await event.fetch(url, {
            method: "POST",
            body: body,
            headers: { "Content-Type": "application/json" }
        }).then((res) => {
            if (res.ok) {
                return { success: true, message: "Account verified successfully" }
            }
            return {success: false, message: "Account verification failed"}
        }).catch((e) => {
            console.log(e)
            return {
                success: false,
                message: "Internal error"
            }
        })
        return result
    }
}


export const actions = {
    login: async (event) => {
        const formData = await event.request.formData();
        const email = formData.get('email')
        const password = formData.get('password')
        if (email === null) {
            return fail(400, { error: m.form_login_error_no_email() })
        }
        if (password === null) {
            return fail(400, { error: m.form_login_error_no_password() })
        }
        if (!validator.isEmail(email.toString())) {
            return fail(400, { error: m.form_login_error_invalid_email() })
        }

        const body = JSON.stringify({ "email": email, "password": password })
        try {
            const response = await event.fetch(apiLocation + "/login", {
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

            event.cookies.set("access", accessToken, { path: "/" })
            event.cookies.set("refresh", refreshToken, { path: "/" })
        }
        catch (e) {
            return fail(500, { error: "Internal error: " + e })
        }

        return { success: true }
    },
    register: async (event) => {
        const formData = await event.request.formData()
        const email = formData.get('email')
        const password = formData.get('password')
        const password_confirm = formData.get('password-confirm')
        const phone = formData.get('phone')
        const name = formData.get('fullname')
        const lang = formData.get('lang')

        if(lang === null){
            return fail(400, { error: m.form_login_error_no_email() })
        }
        if (email === null) {
            return fail(400, { error: m.form_login_error_no_email() })
        }
        if (password === null) {
            return fail(400, { error: m.form_login_error_no_password() })
        }
        if (name === null) {
            return fail(400, { error: m.form_login_error_no_password() })
        }
        if (phone === null) {
            return fail(400, { error: m.form_login_error_no_password() })
        }
        if (password !== password_confirm) {
            return fail(400, { error: m.form_login_error_no_password() })
        }
        if (!validator.isMobilePhone(phone.toString(), 'hu-HU')) {
            return fail(400, { error: m.form_login_error_invalid_email() })
        }
        if (!validator.isEmail(email.toString())) {
            return fail(400, { error: m.form_login_error_invalid_email() })
        }

        if(lang != "hu" && lang != "en"){
            return fail(400)
        }

        const body = JSON.stringify({ "email": email, "password": password, "name": name, "phone_number": phone, "lang" : lang})

        try {
            const response = await event.fetch(apiLocation + "/register", {
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
        }
        catch (e) {
            return fail(500, { error: "Internal error: " + e })
        }
        return {success: "Check your email to activate your account!"}
    },
    logout: async (event) => {
        const url = apiLocation + "/logout"
        event.fetch(url, {
            method: "POST",
            headers: {
                "Authorization": "Bearer " + event.locals.access
            }
        }).then((r) => {
        }).catch((e) => {
            console.log(e)
        })
        event.locals.access = null
        event.locals.refresh = null
        event.locals.loggedIn = false
        event.locals.profile = null
        event.cookies.delete("access", { path: "/" })
        event.cookies.delete("refresh", { path: "/" })
        return { success: "Logged out successfully" }
    },

} satisfies Actions;