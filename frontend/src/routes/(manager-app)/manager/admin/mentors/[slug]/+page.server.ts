import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const [mentorProfile, mentorPayments, mentorSchedules, adminProfile] = await Promise.all([
		controller.getMentorProfile(fetch, params.slug),
		controller.getMentorPayments(fetch, params.slug),
		controller.getMentorSchedules(fetch, params.slug),
		controller.getAdminProfile(fetch)
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
	if (!adminProfile.success) {
		throw error(adminProfile.status, { message: adminProfile.message });
	}
	return {
		isVerified: adminProfile.resBody.data.status === 'verified',
		profile: mentorProfile.resBody.data,
		schedules: mentorSchedules.resBody?.data,
		payments: mentorPayments.resBody?.data
	};
};

export const actions = {
	deleteMentor: async ({ fetch, params }) => {
		const { success, message, status } = await controller.deleteMentor(fetch, params.slug);
		if (!success) {
			return fail(status, { message });
		}
		throw redirect(303, '/manager/admin/mentors');
	}
} satisfies Actions;
