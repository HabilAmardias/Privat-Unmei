import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { Production } from '$lib/utils/constants';
import { PUBLIC_ENVIRONMENT_OPTION, PUBLIC_COOKIE_DOMAIN } from '$env/static/public';
import { controller } from './controller';

export const load: PageServerLoad = async ({ cookies }) => {
	const status = cookies.get('status');
	const authToken = cookies.get('auth_token');
	const refreshToken = cookies.get('refresh_token');
	const verifyToken = cookies.get('verify_token');

	if ((!authToken || !refreshToken || !verifyToken) && !status) {
		throw error(401, { message: 'Login failed' });
	}

	const cookiesOption = {
		path: '/',
		secure: PUBLIC_ENVIRONMENT_OPTION === Production,
		httpOnly: true,
		domain: PUBLIC_COOKIE_DOMAIN
	};
	cookies.delete('oauthstate', cookiesOption);
	return { success: true, userStatus: status };
};

export const actions = {
	verify: async ({ fetch, cookies }) => {
		const { success, message, status, cookiesData } = await controller.verify(fetch);
		if (!success) {
			return fail(status, { message });
		}
		const cookiesOption = {
			path: '/',
			secure: PUBLIC_ENVIRONMENT_OPTION === Production,
			httpOnly: true,
			domain: PUBLIC_COOKIE_DOMAIN
		};
		cookies.delete('verify_token', cookiesOption);
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
