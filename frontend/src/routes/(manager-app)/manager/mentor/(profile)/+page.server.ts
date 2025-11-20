import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const [mentorProfile, mentorPayments, mentorSchedules] = await Promise.all([
		controller.getMentorProfile(fetch),
		controller.getMentorPayments(fetch),
		controller.getMentorSchedules(fetch)
	]);
	if (!mentorPayments.success) {
		if (mentorPayments.status !== 404) {
			throw error(mentorPayments.status, { message: mentorPayments.message });
		}
	}
	if (!mentorSchedules.success) {
		if (mentorSchedules.status !== 404) {
			throw error(mentorSchedules.status, { message: mentorSchedules.message });
		}
	}
	if (!mentorProfile.success) {
		throw error(mentorProfile.status, { message: mentorProfile.message });
	}
	return {
		profile: mentorProfile.resBody.data,
		schedules: mentorSchedules.resBody?.data,
		payments: mentorPayments.resBody?.data
	};
};

// export const actions = {
// 	deleteMentor: async ({ fetch, params }) => {
// 		const { success, message, status } = await controller.deleteMentor(fetch, params.slug);
// 		if (!success) {
// 			return fail(status, { message });
// 		}
// 		throw redirect(303, '/manager/admin/mentors');
// 	}
// } satisfies Actions;
