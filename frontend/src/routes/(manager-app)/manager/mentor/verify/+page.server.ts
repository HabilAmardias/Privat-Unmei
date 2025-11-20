import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

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
		cookies.delete('auth_token', { path: '/' });
		cookies.delete('refresh_token', { path: '/' });
		cookies.delete('status', { path: '/' });
		cookies.delete('role', { path: '/' });
		return { status, message };
	}
} satisfies Actions;
