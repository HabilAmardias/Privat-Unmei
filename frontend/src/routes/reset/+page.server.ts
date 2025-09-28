import { type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from '../$types';
import { fail } from '@sveltejs/kit';
import { controller } from './controller';

export const load: PageServerLoad = ({ cookies }) => {
	if (cookies.get('auth_token')) {
		return { returnHome: true };
	}
	return { returnHome: false };
};

export const actions = {
	send: async ({ request, fetch }) => {
		const { success, status, message } = await controller.sendEmail(request, fetch);
		if (!success) {
			return fail(status, { message });
		}
		return { success, status, message };
	}
} satisfies Actions;
