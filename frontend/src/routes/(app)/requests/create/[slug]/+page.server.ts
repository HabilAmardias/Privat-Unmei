import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const [courseRes, dowRes] = await Promise.all([
		controller.getCourseDetail(fetch, params.slug),
		controller.getAvailableDayOfWeek(fetch, params.slug)
	]);
	if (!courseRes.success) {
		throw error(courseRes.status, { message: courseRes.message });
	}
	if (!dowRes.success) {
		throw error(dowRes.status, { message: dowRes.message });
	}
	const [paymentRes, scheduleRes] = await Promise.all([
		controller.getMentorPayments(fetch, courseRes.resBody.data.mentor_id),
		controller.getMentorSchedules(fetch, courseRes.resBody.data.mentor_id)
	]);
	if (!paymentRes.success) {
		throw error(paymentRes.status, { message: paymentRes.message });
	}
	if (!scheduleRes.success) {
		throw error(scheduleRes.status, { message: scheduleRes.message });
	}

	return {
		detail: courseRes.resBody.data,
		schedules: scheduleRes.resBody.data,
		payments: paymentRes.resBody.data,
		dows: dowRes.resBody.data.day_of_weeks
	};
};

export const actions = {
	createRequest: async ({ fetch, request, params }) => {
		if (!params.slug) {
			return fail(400, { message: 'no course selected' });
		}
		const { success, status, message } = await controller.createRequest(
			fetch,
			params.slug,
			request
		);
		if (!success) {
			return fail(status, { message });
		}
		throw redirect(303, '/profile');
	}
} satisfies Actions;
