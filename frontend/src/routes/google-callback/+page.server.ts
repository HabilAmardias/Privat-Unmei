import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
    const authToken = cookies.get('auth_token')
    const refreshToken = cookies.get('refresh_token')
    const status = cookies.get('status')
	if (!authToken || !refreshToken || !status) {
		throw error(401, { message: 'Google login failed' });
	}
	cookies.delete('oauthstate', { path: '/' });
	return { success: true , status};
};
