import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

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
		cookies.delete('auth_token', { path: '/', secure: false });
		cookies.delete('refresh_token', { path: '/', secure: false });
		cookies.delete('status', { path: '/', secure: false });
		cookies.delete('role', { path: '/', secure: false });
		return { status, message };
	}
} satisfies Actions;
