import type { LayoutServerLoad } from './$types';
import { apiLocation } from '$lib/server/config';
export const load: LayoutServerLoad = async ({ locals }) => {
    if(locals.loggedIn){
        return { loggedIn: true, profile: locals.profile }
    }
    return { loggedIn: false }
}