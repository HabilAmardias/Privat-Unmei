import { redirect, type Actions } from '@sveltejs/kit';
import { Production } from '$lib/utils/constants';
import { PUBLIC_ENVIRONMENT_OPTION } from '$env/static/public';

export const actions = {
	default: async ({ cookies }) => {
		cookies.delete('auth_token', {
			path: '/',
			secure: PUBLIC_ENVIRONMENT_OPTION === Production,
			httpOnly: true
		});
		cookies.delete('refresh_token', {
			path: '/',
			secure: PUBLIC_ENVIRONMENT_OPTION === Production
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
		throw redirect(303, '/manager/login');
	}
} satisfies Actions;
