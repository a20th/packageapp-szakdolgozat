
import { apiLocation } from '$lib/server/config';

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
                if (!Array.isArray(resp)){
                    return {found: false}
                }

                if ((resp as []).length == 0){
                    return {found: false}
                }
                return {found: true, status: resp, packageId: id }
            }
        }).catch((e) => {
            console.log(e)
        })
    }
    return {search: true}
}