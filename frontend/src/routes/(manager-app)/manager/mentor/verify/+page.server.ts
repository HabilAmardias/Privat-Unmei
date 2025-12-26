import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';
import { Production } from '$lib/utils/constants';
import { PUBLIC_ENVIRONMENT_OPTION } from '$env/static/public';

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
		return { status, message };
	}
} satisfies Actions;
