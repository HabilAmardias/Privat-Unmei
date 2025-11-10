import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const [adminProfile, additionalCosts, discounts] = await Promise.all([
		controller.getProfile(fetch),
		controller.getCosts(fetch),
		controller.getDiscounts(fetch)
	]);
	if (!discounts.success) {
		throw error(discounts.status, { message: discounts.message });
	}
	if (!adminProfile.success) {
		throw error(adminProfile.status, { message: adminProfile.message });
	}
	if (!additionalCosts.success) {
		throw error(additionalCosts.status, { message: additionalCosts.message });
	}

	return {
		discounts: discounts.resBody.data,
		costs: additionalCosts.resBody.data,
		isVerified: adminProfile.resBody.data.status === 'verified'
	};
};

export const actions = {
	getCosts: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getCosts(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { costs: resBody.data };
	},
	deleteCost: async ({ fetch, request }) => {
		const { success, message, status } = await controller.deleteCost(fetch, request);
		if (!success) {
			return fail(status, message);
		}
		return { message };
	},
	addCost: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.createCost(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { newCost: resBody?.data, message };
	},
	updateCostAmount: async ({ fetch, request }) => {
		const { success, message, status } = await controller.updateCostAmount(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { message };
	},
	getDiscounts: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getDiscounts(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { discounts: resBody.data };
	},
	deleteDiscount: async ({ fetch, request }) => {
		const { success, message, status } = await controller.deleteDiscount(fetch, request);
		if (!success) {
			return fail(status, message);
		}
		return { message };
	},
	addDiscount: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.createDiscount(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { newDiscount: resBody?.data, message };
	},
	updateDiscountAmount: async ({ fetch, request }) => {
		const { success, message, status } = await controller.updateDiscountAmount(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { message };
	}
} satisfies Actions;
