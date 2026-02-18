import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const { success, message, status, resBody } = await controller.getProfile(fetch);
	if (!success) {
		throw error(status, { message });
	}
	if (resBody.data.status !== 'verified') {
		throw redirect(303, '/manager/admin/verify');
	}
	return { profile: resBody.data, isVerified: resBody.data.status === 'verified' };
};

export const actions = {
	changePassword: async ({ fetch, request }) => {
		const { success, message, status } = await controller.changePassword(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { success };
	}
} satisfies Actions;
