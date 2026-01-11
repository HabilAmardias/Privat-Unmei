import { redirect } from '@sveltejs/kit';
import { Production } from '$lib/utils/constants';
import { PUBLIC_ENVIRONMENT_OPTION, PUBLIC_COOKIE_DOMAIN } from '$env/static/public';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const cookiesOption = {
		path: '/',
		secure: PUBLIC_ENVIRONMENT_OPTION === Production,
		httpOnly: true,
		domain: PUBLIC_COOKIE_DOMAIN
	};
	cookies.delete('auth_token', cookiesOption);
	cookies.delete('refresh_token', cookiesOption);
	cookies.delete('status', cookiesOption);
	cookies.delete('role', cookiesOption);
	throw redirect(303, '/login');
};
