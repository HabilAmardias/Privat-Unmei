import { redirect, type Actions } from '@sveltejs/kit';

export const actions = {
	default: async ({ cookies }) => {
		cookies.delete('auth_token', { path: '/', secure: false });
		cookies.delete('refresh_token', { path: '/', secure: false });
		cookies.delete('status', { path: '/', secure: false });
		cookies.delete('role', { path: '/', secure: false });
		throw redirect(303, '/login');
	}
} satisfies Actions;
