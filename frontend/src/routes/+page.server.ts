import { fail, redirect, type Actions } from '@sveltejs/kit';
import { controller } from './controller';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ cookies }) => {
	if (cookies.get('auth_token')) {
		redirect(303, '/home');
	}
};

export const actions = {
	login: async ({ request, cookies, fetch }) => {
		const { responseBody, cookiesData, success, message, status } = await controller.login(
			request,
			fetch
		);
		if (!success || !responseBody) {
			return fail(status, { message });
		}
		cookiesData?.forEach((c) => {
			cookies.set(c.key, c.value, {
				path: c.path,
				domain: c.domain,
				httpOnly: c.httpOnly,
				maxAge: c.maxAge
			});
		});
		cookies.set('status', responseBody.data.status, { path: '/', httpOnly: false });
		redirect(303, '/home');
	},
	register: async ({ request, fetch }) => {
		const { success, message, status } = await controller.register(request, fetch);
		if (!success) {
			return fail(status, { message });
		}
		return { success, message: 'succesfully registered' };
	}
} satisfies Actions;
