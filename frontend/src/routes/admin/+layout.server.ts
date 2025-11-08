import type { LayoutServerLoad } from '../app/$types';
import { apiLocation } from '$lib/server/config';
export const load: LayoutServerLoad = async ({ locals }) => {
    if(locals.admin?.loggedIn){
        return { loggedIn: true }
    }
    return { loggedIn: false }
}