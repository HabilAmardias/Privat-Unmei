import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';
import { Production } from '$lib/utils/constants';
import { PUBLIC_ENVIRONMENT_OPTION, PUBLIC_COOKIE_DOMAIN } from '$env/static/public';

export const load: PageServerLoad = async ({ cookies }) => {
	const status = cookies.get('status');
	return { isVerified: status === 'verified' };
};

export const actions = {
	verifyMentor: async ({ fetch, request, cookies }) => {
		const { success, message, status } = await controller.verifyMentor(fetch, request);
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
		cookies.delete('refresh_token', cookiesOption);
		cookies.delete('status', cookiesOption);
		cookies.delete('role', cookiesOption);
		return { status, message };
	}
} satisfies Actions;
