// @ts-nocheck
import type { LayoutServerLoad } from './$types';
import { apiLocation } from '$lib/server/config';
export const load = async ({ locals }: Parameters<LayoutServerLoad>[0]) => {
    if(locals.loggedIn){
        return { loggedIn: true, profile: locals.profile }
    }
    return { loggedIn: false }
}