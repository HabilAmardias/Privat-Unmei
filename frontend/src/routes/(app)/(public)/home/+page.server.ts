import { error } from '@sveltejs/kit';
import { controller } from './controller';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const [courseRes, mentorRes] = await Promise.all([
		controller.getCourses(fetch),
		controller.getMentors(fetch)
	]);
	if (!courseRes.success) {
		throw error(courseRes.status, { message: courseRes.message });
	}
	if (!mentorRes.success) {
		throw error(mentorRes.status, { message: mentorRes.message });
	}
	return {
		courses: courseRes.resBody.data,
		mentors: mentorRes.resBody.data
	};
};
