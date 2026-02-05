import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const [paymentMethodsRes, generatePasswordRes, adminProfileRes] = await Promise.all([
		controller.getPaymentMethods(fetch),
		controller.getRandomizedPassword(fetch),
		controller.getAdminProfile(fetch)
	]);
	if (!adminProfileRes.success) {
		throw error(adminProfileRes.status, { message: adminProfileRes.message });
	}
	if (!paymentMethodsRes.success) {
		throw error(paymentMethodsRes.status, { message: paymentMethodsRes.message });
	}
	if (!generatePasswordRes.success) {
		throw error(generatePasswordRes.status, { message: generatePasswordRes.message });
	}
	return {
		paymentMethods: paymentMethodsRes.resBody.data.entries,
		generatedPassword: generatePasswordRes.resBody.data.password,
		isVerified: adminProfileRes.resBody.data.status === 'verified'
	};
};

export const actions = {
	getPaymentMethods: async ({ fetch, request }) => {
		const { success, status, message, resBody } = await controller.getPaymentMethods(
			fetch,
			request
		);
		if (!success) {
			throw fail(status, { message });
		}
		return { paymentMethods: resBody.data.entries };
	},
	generatePassword: async ({ fetch }) => {
		const { success, status, message, resBody } = await controller.getRandomizedPassword(fetch);
		if (!success) {
			return fail(status, { message });
		}
		return { password: resBody.data.password };
	},
	createMentor: async ({ fetch, request }) => {
		const { success, status, message } = await controller.CreateNewMentor(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		throw redirect(303, '/manager/admin/users');
	}
} satisfies Actions;
