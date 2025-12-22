import { type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from '../controller';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = ({ cookies, params }) => {
	if (cookies.get('auth_token') || cookies.get('refresh_token')) {
		cookies.delete('auth_token', { path: '/', secure: false });
		cookies.delete('refresh_token', { path: '/', secure: false });
		cookies.delete('status', { path: '/', secure: false });
		cookies.delete('role', { path: '/', secure: false });
	}
	cookies.set('auth_token', params.slug, { path: '/', secure: false });
	return { returnHome: false };
};

export const actions = {
	reset: async ({ fetch, cookies, request }) => {
		const { success, message, status } = await controller.resetPassword(request, fetch);
		if (!success) {
			return fail(status, { message });
		}
		cookies.delete('auth_token', { path: '/', secure: false });
		return { success: true };
	}
} satisfies Actions;
