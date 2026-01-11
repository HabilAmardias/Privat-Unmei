import { type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from '../controller';
import { fail } from '@sveltejs/kit';
import { Production } from '$lib/utils/constants';
import { PUBLIC_ENVIRONMENT_OPTION, PUBLIC_COOKIE_DOMAIN } from '$env/static/public';

export const load: PageServerLoad = ({ cookies, params }) => {
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
	cookies.set('auth_token', params.slug, cookiesOption);
	return { returnHome: false };
};

export const actions = {
	reset: async ({ fetch, cookies, request }) => {
		const { success, message, status } = await controller.resetPassword(request, fetch);
		if (!success) {
			return fail(status, { message });
		}
		const cookiesOption = {
			path: '/',
			secure: PUBLIC_ENVIRONMENT_OPTION === Production,
			httpOnly: true,
			domain: PUBLIC_COOKIE_DOMAIN
		};
		cookies.delete('auth_token', cookiesOption);
		return { success: true };
	}
} satisfies Actions;
