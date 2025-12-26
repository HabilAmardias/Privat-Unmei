import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';
import { Production } from '$lib/utils/constants';
import { PUBLIC_ENVIRONMENT_OPTION } from '$env/static/public';

export const load: PageServerLoad = async ({ cookies, params, fetch }) => {
	if (cookies.get('auth_token') || cookies.get('refresh_token')) {
		cookies.delete('auth_token', {
			path: '/',
			secure: PUBLIC_ENVIRONMENT_OPTION === Production,
			httpOnly: true
		});
		cookies.delete('refresh_token', {
			path: '/',
			secure: PUBLIC_ENVIRONMENT_OPTION === Production,
			httpOnly: true
		});
		cookies.delete('status', {
			path: '/',
			secure: PUBLIC_ENVIRONMENT_OPTION === Production,
			httpOnly: true
		});
		cookies.delete('role', {
			path: '/',
			secure: PUBLIC_ENVIRONMENT_OPTION === Production,
			httpOnly: true
		});
	}
	cookies.set('auth_token', params.slug, {
		path: '/',
		secure: PUBLIC_ENVIRONMENT_OPTION === Production,
		httpOnly: true
	});
	const { success, message, status } = await controller.verify(fetch);
	cookies.delete('auth_token', {
		path: '/',
		secure: PUBLIC_ENVIRONMENT_OPTION === Production,
		httpOnly: true
	});
	if (!success) {
		throw error(status, { message });
	}
	return { success, message, status };
};
