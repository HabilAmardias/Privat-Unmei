import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const [adminProfile, payments] = await Promise.all([
		controller.getProfile(fetch),
		controller.getPayments(fetch)
	]);
	if (!adminProfile.success) {
		throw error(adminProfile.status, { message: adminProfile.message });
	}
	if (!payments.success) {
		throw error(payments.status, { message: payments.message });
	}
	return {
		payments: payments.resBody.data,
		isVerified: adminProfile.resBody.data.status === 'verified'
	};
};

export const actions = {
	getPayments: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getPayments(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { payments: resBody.data };
	},
	deletePayment: async ({ fetch, request }) => {
		const { success, message, status } = await controller.deletePayment(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { message };
	},
	updatePayment: async ({ fetch, request }) => {
		const { success, message, status } = await controller.updatePayment(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { message };
	},
	createPayment: async ({ fetch, request }) => {
		const { success, message, status } = await controller.createPayment(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { message };
	}
} satisfies Actions;
