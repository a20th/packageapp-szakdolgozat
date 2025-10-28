
import { apiLocation } from '$lib/server/config';

export const load = async (event) => {

    const id = event.url.searchParams.get("id")
    if (id) {
        const url = apiLocation + "/track"
        const body = JSON.stringify({ id: id })
        return await event.fetch(url, {
            method: "GET",
            body: body,
            headers: { "Content-Type": "application/json" }
        }).then(async (res) => {
            if (res.ok) {
                return {status: await res.json(), packageId: id }
            }
        }).catch((e) => {
            console.log(e)
        })
    }
}