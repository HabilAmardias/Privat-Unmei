import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const [profile, orders] = await Promise.all([
		controller.getProfile(fetch),
		controller.getOrders(fetch)
	]);
	if (!profile.success) {
		throw error(profile.status, { message: profile.message });
	}
	return {
		profile: profile.resBody.data,
		orders: orders.resBody,
		userStatus: profile.resBody.data.status
	};
};

export const actions = {
	updateProfile: async ({ request, fetch }) => {
		const { success, message, status } = await controller.updateProfile(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { message };
	},
	sendVerification: async ({ fetch }) => {
		const { success, message, status } = await controller.sendVerificationLink(fetch);
		if (!success) {
			return fail(status, { message });
		}
		return { success, message, status };
	},
	myOrders: async ({ request, fetch }) => {
		const { success, message, status, resBody } = await controller.getOrders(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return {
			orders: resBody.entries,
			page: resBody.page_info.page,
			limit: resBody.page_info.limit,
			totalRow: resBody.page_info.total_row,
			message
		};
	}
} satisfies Actions;
