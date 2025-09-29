import { type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from '../controller';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = ({ cookies, params }) => {
	if (cookies.get('auth_token')) {
		return { returnHome: true };
	}
	cookies.set('auth_token', params.slug, { path: '/' });
	return { returnHome: false };
};

export const actions = {
	reset: async ({ fetch, cookies, request }) => {
		const { success, message, status } = await controller.resetPassword(request, fetch);
		if (!success) {
			return fail(status, { message });
		}
		cookies.delete('auth_token', { path: '/' });
		return { success: true };
	}
} satisfies Actions;
