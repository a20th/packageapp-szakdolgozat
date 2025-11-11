import { redirect, type Handle } from '@sveltejs/kit';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { apiLocation } from '$lib/server/config';
import { sequence } from '@sveltejs/kit/hooks';

const adminHandle: Handle = async ({ event, resolve }) => {
	const path = (path: string) => event.url.pathname.startsWith(path)
	if (!path("/admin")) {
		return await resolve(event)
	}

	const accessToken = event.cookies.get('access')
	const refreshToken = event.cookies.get('refresh')

	let loggedIn: boolean = accessToken && refreshToken ? true : false

	try {
		const response = await fetch(apiLocation + "/admin/get", {
			method: 'GET',
			headers: {
				'Content-Type': "application/json",
				'Authorization': 'Bearer ' + accessToken,
			},
		})

		if (response.status == 401) {
			const response = await fetch(apiLocation + "/admin/refresh", {
				method: "GET",
				headers: {
					"Content-Type": "application/json",
					"Authorization": "Bearer " + event.cookies.get('refresh')
				},
			})
			const data = await response.json()
			if (!response.ok) {
				loggedIn = false
			}
			else if (!(data.access_token && data.refresh_token)) {
				loggedIn = false
			}
			else {
				const accessToken = data.access_token
				const refreshToken = data.refresh_token
				event.cookies.set("access", accessToken, { path: "/admin" })
				event.cookies.set("refresh", refreshToken, { path: "/admin" })
				loggedIn = true
			}
		}
	}
	catch (e) {
		console.log(e)
		loggedIn = false
	}

	if (loggedIn) {
		event.locals.admin = {} as App.AdminData
		event.locals.admin.access = event.cookies.get('access')!
		event.locals.admin.refresh = event.cookies.get('refresh')!
		event.locals.admin.loggedIn = true
	}
	else {
		event.locals.admin = null
		event.cookies.delete("access", { path: "/admin" })
		event.cookies.delete("refresh", { path: "/admin" })
	}

	if (path("/admin/")) {
		if (!loggedIn) return redirect(303, "/admin")
	}

	return await resolve(event);
};


const authHandle: Handle = async ({ event, resolve }) => {
	const path = (path: string) => event.url.pathname.startsWith(path)
	if (!path("/app")) {
		return await resolve(event)
	}

	const accessToken = event.cookies.get('access')
	const refreshToken = event.cookies.get('refresh')

	let loggedIn: boolean = accessToken && refreshToken ? true : false

	try {
		const response = await fetch(apiLocation + "/get", {
			method: 'get',
			headers: {
				'Content-Type': "application/json",
				'Authorization': 'Bearer ' + accessToken,
			},
		})

		const result = await response.json()
		event.locals.profile = result

		if (!response.ok) {
			const response = await event.fetch(apiLocation + "/refresh", {
				method: "GET",
				headers: {
					"Content-Type": "application/json",
					"Authorization": "Bearer " + event.cookies.get('refresh')
				},
			})
			const data = await response.json()
			if (!response.ok) {
				loggedIn = false
			}
			else if (!(data.access_token && data.refresh_token)) {
				loggedIn = false
			}
			else {
				const accessToken = data.access_token
				const refreshToken = data.refresh_token
				event.cookies.set("access", accessToken, { path: "/app" })
				event.cookies.set("refresh", refreshToken, { path: "/app" })
				loggedIn = true
				const resp = await fetch(apiLocation + "/get", {
					method: 'get',
					headers: {
						'Content-Type': "application/json",
						'Authorization': 'Bearer ' + accessToken,
					},
				})

				const result = await response.json()
				event.locals.profile = result
			}


		}
		else {
			event.locals.profile = result
		}
	}
	catch (e) {
		console.log(e)
		loggedIn = false
	}

	if (loggedIn) {
		event.locals.access = event.cookies.get('access')!
		event.locals.refresh = event.cookies.get('refresh')!
		event.locals.loggedIn = true
	}
	else {
		event.locals.access = null
		event.locals.refresh = null
		event.locals.loggedIn = false
		event.locals.profile = null
		event.cookies.delete("access", { path: "/app" })
		event.cookies.delete("refresh", { path: "/app" })
	}

	//Must be logged in
	if (path("/app/order") || path("/app/orders")) {
		if (!loggedIn) return redirect(303, "/app")
	}

	return await resolve(event);
};

const paraglideHandle: Handle = ({ event, resolve }) => paraglideMiddleware(event.request, async ({ request, locale }) => {
	event.request = request;

	return resolve(event, {
		transformPageChunk: ({ html }) => html.replace('%paraglide.lang%', locale)
	});
});

export const handle: Handle = sequence(authHandle, adminHandle, paraglideHandle)