import { type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from '../controller';
import { fail } from '@sveltejs/kit';
import { Production } from '$lib/utils/constants';
import { PUBLIC_ENVIRONMENT_OPTION } from '$env/static/public';

export const load: PageServerLoad = ({ cookies, params }) => {
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
		httpOnly: true,
		secure: PUBLIC_ENVIRONMENT_OPTION === Production
	});
	return { returnHome: false };
};

export const actions = {
	reset: async ({ fetch, cookies, request }) => {
		const { success, message, status } = await controller.resetPassword(request, fetch);
		if (!success) {
			return fail(status, { message });
		}
		cookies.delete('auth_token', {
			path: '/',
			secure: PUBLIC_ENVIRONMENT_OPTION === Production,
			httpOnly: true
		});
		return { success: true };
	}
} satisfies Actions;
