import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ cookies, params, fetch }) => {
	if (cookies.get('auth_token')) {
		cookies.delete('auth_token', { path: '/' });
		cookies.delete('refresh_token', { path: '/' });
	}
	cookies.set('auth_token', params.slug, { path: '/' });
	const { success, message, status } = await controller.verify(fetch);
	if (!success) {
		cookies.delete('auth_token', { path: '/' });
		throw error(status, { message });
	}
	cookies.delete('auth_token', { path: '/' });
	return { success, message, status };
};
