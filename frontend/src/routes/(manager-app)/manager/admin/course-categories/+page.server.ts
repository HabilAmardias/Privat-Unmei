import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const [adminProfile, categories] = await Promise.all([
		controller.getProfile(fetch),
		controller.getCategories(fetch)
	]);
	if (!adminProfile.success) {
		throw error(adminProfile.status, { message: adminProfile.message });
	}
	if (!categories.success) {
		throw error(categories.status, { message: categories.message });
	}
	return {
		categories: categories.resBody.data,
		isVerified: adminProfile.resBody.data.status === 'verified'
	};
};

export const actions = {
	getCategories: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getCategories(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { categories: resBody.data };
	},
	deleteCategory: async ({ fetch, request }) => {
		const { success, message, status } = await controller.deleteCategory(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { message };
	},
	updateCategory: async ({ fetch, request }) => {
		const { success, message, status } = await controller.updateCategory(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { message };
	},
	createCategory: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.createCategory(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { newCategory: resBody?.data };
	}
} satisfies Actions;
