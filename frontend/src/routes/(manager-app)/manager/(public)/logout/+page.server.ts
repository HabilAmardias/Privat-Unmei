import { redirect, type Actions } from '@sveltejs/kit';

export const actions = {
	default: async ({ cookies }) => {
		cookies.delete('auth_token', { path: '/' });
		cookies.delete('refresh_token', { path: '/' });
		cookies.delete('status', { path: '/' });
		cookies.delete('role', { path: '/' });
		throw redirect(303, '/login');
	}
} satisfies Actions;
