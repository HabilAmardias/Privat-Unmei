import { fail, redirect, type Actions } from '@sveltejs/kit';
import { controller } from './controller';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ cookies }) => {
	if (cookies.get('auth_token')) {
		redirect(303, '/courses');
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
				maxAge: c.maxAge,
				sameSite: c.sameSite,
				secure: false
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
	googlelogin: async ({ fetch, cookies }) => {
		const { success, message, status, res, cookiesData } = await controller.googleLogin(fetch);
		if (!success) {
			return fail(status, { message });
		}
		const redirectUrl = res?.headers.get('Location');
		if (redirectUrl) {
			cookiesData.forEach((val) => {
				cookies.set(val.key, val.value, { path: val.path });
			});
			redirect(303, redirectUrl);
		}
	}
} satisfies Actions;
