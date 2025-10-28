import { fail, redirect, type Handle } from '@sveltejs/kit';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { apiLocation } from '$lib/server/config';

/*
export const handle: Handle = ({ event, resolve }) => paraglideMiddleware(event.request, async ({ request, locale }) => {
	event.request = request;

	return resolve(event, {
		transformPageChunk: ({ html }) => html.replace('%paraglide.lang%', locale)
	});
});
*/

export const handle: Handle = async ({ event, resolve }) => paraglideMiddleware(event.request, async ({ request, locale }) => {
	event.request = request
	const accessToken = event.cookies.get('access')
	const refreshToken = event.cookies.get('refresh')

	let loggedIn: boolean = accessToken && refreshToken ? true : false

	if (loggedIn) {
		event.locals.access = accessToken!
		event.locals.refresh = refreshToken!
		event.locals.loggedIn = true

		const url = apiLocation + "/get"
		try {
			const response = await fetch(url, {
				method: 'get',
				headers: {
					'Content-Type': "application/json",
					'Authorization': 'Bearer ' + accessToken,

				},
			})

			const result = await response.json()
			
			if (result.status == 403) {
				event.locals.access = null
				event.locals.refresh = null
				event.locals.loggedIn = false

				return redirect(303, "/app")
			}
			event.locals.profile = result
		}
		catch (e) {
			console.log(e)
			return fail(500)
		}
	}

	//Must be logged in
	const path = (path: string) => event.url.pathname.startsWith(path)
	if (path("/app/order")) {
		if (!loggedIn) return redirect(303, "/app")
	}

	return await resolve(event, {
		transformPageChunk: ({ html }) => html.replace('%paraglide.lang%', locale)
	});
});