import { error } from '@sveltejs/kit';
import { controller } from './controller';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, cookies }) => {
	console.log(cookies.getAll());
	const [courseRes, mentorRes, mostBoughtRes] = await Promise.all([
		controller.getCourses(fetch),
		controller.getMentors(fetch),
		controller.getMostBought(fetch)
	]);
	if (!courseRes.success) {
		throw error(courseRes.status, { message: courseRes.message });
	}
	if (!mentorRes.success) {
		throw error(mentorRes.status, { message: mentorRes.message });
	}
	if (!mostBoughtRes.success) {
		throw error(mostBoughtRes.status, { message: mostBoughtRes.message });
	}
	return {
		courses: courseRes.resBody.data,
		mentors: mentorRes.resBody.data,
		mostBought: mostBoughtRes.resBody.data
	};
};
