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
		const { cookiesData, success, message, status, userStatus } = await controller.login(
			request,
			fetch
		);
		if (!success) {
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
		return { success, userStatus };
	},
	register: async ({ request, fetch }) => {
		const { success, message, status } = await controller.register(request, fetch);
		if (!success) {
			return fail(status, { message });
		}
		return { success, message: 'succesfully registered' };
	},
	googlelogin: async ({ fetch }) => {
		const { success, message, status, res } = await controller.googleLogin(fetch);
		if (!success) {
			return fail(status, { message });
		}
		if (res?.redirected) {
			redirect(303, res.url);
		}
	}
} satisfies Actions;
