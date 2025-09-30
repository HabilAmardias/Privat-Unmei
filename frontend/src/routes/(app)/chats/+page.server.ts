import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ cookies }) => {
	if (!cookies.get('auth_token')) {
		redirect(303, '/');
	}
};
