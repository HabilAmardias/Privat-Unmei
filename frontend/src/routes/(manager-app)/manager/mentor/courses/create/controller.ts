import type { Fetch, ServerResponse, PaginatedResponse } from '$lib/types';
import type { CourseCategory } from './model';
import { FetchData } from '$lib/utils';

class CreateCourseController {
	async getCourseCategories(fetch: Fetch, req?: Request) {
		let url = 'http://habilog.xyz/api/v1/course-categories?';
		if (req) {
			const formData = await req.formData();
			const limit = formData.get('limit');
			const search = formData.get('search');
			const args: string[] = [];
			if (limit) {
				args.push(`limit=${limit}`);
			}
			if (search) {
				args.push(`search=${search}`);
			}
			url += args.join('&');
		}
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<PaginatedResponse<CourseCategory>> = await res?.json();
		return { success, message, status, resBody };
	}
	async createCourse(fetch: Fetch, req: Request) {
		const url = 'http://habilog.xyz/api/v1/courses';
		const formData = await req.formData();
		const title = formData.get('title');
		const description = formData.get('description');
		const domicile = formData.get('domicile');
		const price = formData.get('price');
		const method = formData.get('method');
		const sessionDuration = formData.get('session_duration');
		const maxSession = formData.get('max_session');
		const categories = formData.get('categories');
		const topics = formData.get('topics');

		const reqBody = JSON.stringify({
			title,
			description,
			domicile,
			price: parseFloat(price as string),
			method,
			session_duration_minutes: parseInt(sessionDuration as string),
			max_total_session: parseInt(maxSession as string),
			course_categories: (categories as string).split(',').map<number>((item) => parseInt(item)),
			course_topics: JSON.parse(topics as string)
		});

		const { success, message, status } = await FetchData(fetch, url, 'POST', reqBody);
		return { success, message, status };
	}
}

export const controller = new CreateCourseController();
