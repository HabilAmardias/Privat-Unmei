import { fail, redirect, type Actions } from '@sveltejs/kit';
import { controller } from './controller';
import type { PageServerLoad } from './$types';
import { Production } from '$lib/utils/constants';
import { PUBLIC_ENVIRONMENT_OPTION } from '$env/static/public';

export const load: PageServerLoad = ({ cookies }) => {
	if (cookies.get('auth_token') || cookies.get('refresh_token')) {
		throw redirect(303, '/courses');
	}
};

export const actions = {
	login: async ({ request, cookies, fetch }) => {
		const { cookiesData, success, message, status } = await controller.login(request, fetch);
		if (!success) {
			return fail(status, { message });
		}
		cookiesData?.forEach((c) => {
			cookies.set(c.key, c.value, {
				path: c.path,
				httpOnly: c.httpOnly,
				maxAge: c.maxAge,
				domain: c.domain,
				sameSite: c.sameSite,
				secure: PUBLIC_ENVIRONMENT_OPTION === Production
			});
		});
		return { success };
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
				cookies.set(val.key, val.value, {
					path: val.path,
					httpOnly: val.httpOnly,
					maxAge: val.maxAge,
					domain: val.domain,
					sameSite: val.sameSite,
					secure: PUBLIC_ENVIRONMENT_OPTION === Production
				});
			});
			redirect(303, redirectUrl);
		}
	}
} satisfies Actions;
