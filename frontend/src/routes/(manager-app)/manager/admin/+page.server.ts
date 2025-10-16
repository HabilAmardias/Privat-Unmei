import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';
import type { ServerResponse } from '$lib/types';
import type { adminProfile } from './model';

export const load : PageServerLoad = async ({fetch}) => {
    const {success, message, status, res} = await controller.getProfile(fetch)
    if (!success){
        throw error(status, {message})
    }
    
    const resBody : ServerResponse<adminProfile> = await res?.json()
    if (resBody.data.status !== 'verified'){
        throw redirect(303, '/manager/admin/verify')
    }
    return {profile: resBody.data}
}