import { error } from '@sveltejs/kit';
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
		dows: dowRes.resBody.data
	};
};
