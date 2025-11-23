import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const [paymentMethodRes, schedulesRes, mentorPaymentRes, mentorProfileRes] = await Promise.all([
		controller.getPaymentMethods(fetch),
		controller.getMentorSchedules(fetch),
		controller.getMentorPayments(fetch),
		controller.getMentorProfile(fetch)
	]);
	if (!paymentMethodRes.success) {
		throw error(paymentMethodRes.status, { message: paymentMethodRes.message });
	}
	if (!schedulesRes.success) {
		throw error(schedulesRes.status, { message: schedulesRes.message });
	}
	if (!mentorPaymentRes.success) {
		throw error(mentorPaymentRes.status, { message: mentorPaymentRes.message });
	}
	if (!mentorProfileRes.success) {
		throw error(mentorProfileRes.status, { message: mentorProfileRes.message });
	}
	return {
		paymentMethods: paymentMethodRes.resBody.data.entries,
		mentorSchedules: schedulesRes.resBody.data,
		mentorPayments: mentorPaymentRes.resBody.data,
		profile: mentorProfileRes.resBody.data
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
	updateProfile: async ({ fetch, request }) => {
		const { success, status, message } = await controller.updateProfile(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		throw redirect(303, '/manager/mentor');
	}
} satisfies Actions;
