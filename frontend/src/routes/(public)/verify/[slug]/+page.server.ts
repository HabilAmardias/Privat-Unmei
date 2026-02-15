import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';
import { Production } from '$lib/utils/constants';
import { PUBLIC_ENVIRONMENT_OPTION, PUBLIC_COOKIE_DOMAIN } from '$env/static/public';

export const load: PageServerLoad = async ({ cookies, params, fetch }) => {
	if (cookies.get('auth_token') || cookies.get('refresh_token')) {
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
	}
	const cookiesOption = {
		path: '/',
		secure: PUBLIC_ENVIRONMENT_OPTION === Production,
		httpOnly: true,
		domain: PUBLIC_COOKIE_DOMAIN
	};
	cookies.set('verify_token', params.slug, cookiesOption);
	const { success, message, status } = await controller.verify(fetch);
	cookies.delete('verify_token', cookiesOption);
	if (!success) {
		throw error(status, { message });
	}
	return { success, message, status };
};
