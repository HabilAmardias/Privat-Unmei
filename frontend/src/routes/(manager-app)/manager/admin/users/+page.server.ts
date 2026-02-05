import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const [mentorsList, studentList, adminProfile] = await Promise.all([
		controller.getMentors(fetch),
		controller.getStudents(fetch),
		controller.getProfile(fetch)
	]);
	if (!mentorsList.success) {
		throw error(mentorsList.status, { message: mentorsList.message });
	}
	if (!adminProfile.success) {
		throw error(adminProfile.status, { message: adminProfile.message });
	}
	if (!studentList.success) {
		throw error(studentList.status, { message: studentList.message });
	}
	return {
		mentorsList: mentorsList.resBody.data,
		studentList: studentList.resBody.data,
		isVerified: adminProfile.resBody.data.status === 'verified'
	};
};

export const actions = {
	getStudents: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getStudents(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { studentList: resBody.data };
	},
	deleteStudent: async ({ fetch, request }) => {
		const { success, message, status } = await controller.deleteStudent(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { success, message, status };
	},
	getMentors: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getMentors(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { mentorsList: resBody.data };
	},
	deleteMentor: async ({ fetch, request }) => {
		const { success, message, status } = await controller.deleteMentor(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { success, message, status };
	}
} satisfies Actions;
