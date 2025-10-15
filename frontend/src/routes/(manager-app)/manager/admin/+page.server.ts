import type { PageServerLoad } from './$types';

export const load : PageServerLoad = async ({cookies}) => {
    const status = cookies.get('status')
    return {isVerified: status === 'verified'}
}