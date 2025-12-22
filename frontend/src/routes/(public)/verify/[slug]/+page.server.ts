import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ cookies, params, fetch }) => {
	if (cookies.get('auth_token')) {
		cookies.delete('auth_token', { path: '/', secure: false });
		cookies.delete('refresh_token', { path: '/', secure: false });
		cookies.delete('status', { path: '/', secure: false });
		cookies.delete('role', { path: '/', secure: false });
	}
	cookies.set('auth_token', params.slug, { path: '/', secure: false });
	const { success, message, status } = await controller.verify(fetch);
	cookies.delete('auth_token', { path: '/', secure: false });
	if (!success) {
		throw error(status, { message });
	}
	return { success, message, status };
};
