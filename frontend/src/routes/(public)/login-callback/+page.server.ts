import { fail, type Actions } from '@sveltejs/kit';
import { controller } from './controller';
import { Production } from '$lib/utils/constants';
import { PUBLIC_ENVIRONMENT_OPTION, PUBLIC_COOKIE_DOMAIN } from '$env/static/public';

export const actions = {
	login: async ({ fetch, request, cookies }) => {
		const { success, message, status, cookiesData } = await controller.login(fetch, request);
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
		const cookiesOption = {
			path: '/',
			secure: PUBLIC_ENVIRONMENT_OPTION === Production,
			httpOnly: true,
			domain: PUBLIC_COOKIE_DOMAIN
		};
		cookies.delete('login_token', cookiesOption);
		return { success };
	},
	resendOTP: async ({ fetch, cookies }) => {
		const { success, message, status, cookiesData } = await controller.resendOTP(fetch);
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
	}
} satisfies Actions;
