import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { type Actions, fail } from '@sveltejs/kit';
import { controller } from './controller';

export const load: PageServerLoad = ({ cookies }) => {
	if (!cookies.get('auth_token')) {
		redirect(303, '/login');
	}
};

export const actions = {
	refresh: async ({ cookies, fetch }) => {
		if (!cookies.get('auth_token')) {
			return fail(401, { message: 'unauthorized' });
		}
		const { success, cookiesData, message, status } = await controller.refresh(fetch);
		if (!success) {
			return fail(status, { message });
		}
		cookiesData?.forEach((val) => {
			cookies.set(val.key, val.value, {
				path: val.path,
				domain: val.domain,
				httpOnly: val.httpOnly,
				maxAge: val.maxAge,
				sameSite: val.sameSite
			});
		});
		return { success: true, message };
	}
} satisfies Actions;
