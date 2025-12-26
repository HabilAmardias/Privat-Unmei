import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { Production } from '$lib/utils/constants';
import { PUBLIC_ENVIRONMENT_OPTION } from '$env/static/public';

export const load: PageServerLoad = async ({ cookies }) => {
	const authToken = cookies.get('auth_token');
	const refreshToken = cookies.get('refresh_token');
	const status = cookies.get('status');
	if (!authToken || !refreshToken || !status) {
		throw error(401, { message: 'Google login failed' });
	}
	cookies.delete('oauthstate', {
		path: '/',
		secure: PUBLIC_ENVIRONMENT_OPTION === Production,
		httpOnly: true
	});
	return { success: true, status };
};
