import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ cookies }) => {
	cookies.delete('auth_token', { path: '/' });
	cookies.delete('refresh_token', { path: '/' });
	cookies.delete('status', { path: '/' });
	cookies.delete('role', { path: '/' });
	throw redirect(303, '/login');
};
