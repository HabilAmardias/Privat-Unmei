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
	verifyAdmin: async ({ fetch, request, cookies }) => {
		const { success, message, status } = await controller.verifyAdmin(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		cookies.delete('auth_token', { path: '/', secure: PUBLIC_ENVIRONMENT_OPTION === Production });
		cookies.delete('refresh_token', {
			path: '/',
			secure: PUBLIC_ENVIRONMENT_OPTION === Production
		});
		cookies.delete('status', { path: '/', secure: PUBLIC_ENVIRONMENT_OPTION === Production });
		cookies.delete('role', { path: '/', secure: PUBLIC_ENVIRONMENT_OPTION === Production });
		return { status, message };
	}
} satisfies Actions;
