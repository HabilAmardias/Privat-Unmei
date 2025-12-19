import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { CourseList, CourseCategory } from './model';

class CoursesController {
	async getCourses(fetch: Fetch, req?: Request) {
		let url = 'http://160.19.167.63/api/v1/courses?';
		if (req) {
			const args: string[] = [];
			const formData = await req.formData();
			const page = formData.get('page');
			if (page) {
				args.push(`page=${page}`);
			}
			const search = formData.get('search');
			if (search) {
				args.push(`search=${search}`);
			}
			const category = formData.get('course_category');
			if (category) {
				args.push(`course_category=${category}`);
			}
			const method = formData.get('method');
			if (method) {
				args.push(`method=${method}`);
			}
			if (args.length > 0) {
				url += args.join('&');
			}
		}
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<PaginatedResponse<CourseList>> = await res?.json();
		return { success, resBody, status, message };
	}
	async getCourseCategories(fetch: Fetch, req?: Request) {
		let url = 'http://160.19.167.63/api/v1/course-categories?';
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
}
export const controller = new CoursesController();
