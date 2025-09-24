import { fail, type Actions } from '@sveltejs/kit';
import { controller } from './controller';

export const actions = {
	login: async ({ request, cookies }) => {
		const { data, success } = await controller.login(request);
		if (!success) {
			return fail(400);
		}
		cookies.set('status', data!.status, { path: '/' });
		cookies.set('token', data!.token, { path: '/' });
		return { success: true, message: 'succesfully login' };
	},
	register: async ({ request }) => {
		const success = await controller.register(request);
		if (!success) {
			return fail(400);
		}
		return { success: true, message: 'succesfully registered' };
	}
} satisfies Actions;
