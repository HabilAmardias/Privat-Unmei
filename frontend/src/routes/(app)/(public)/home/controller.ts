import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { CourseList, MentorList } from './model';

class HomeController {
	async getCourses(fetch: Fetch, req?: Request) {
		let url = 'http://habilog.xyz/api/v1/courses?';
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
			const category = formData.get('category');
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
	async getMentors(fetch: Fetch) {
		const url = 'http://habilog.xyz/api/v1/mentors';
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<PaginatedResponse<MentorList>> = await res?.json();
		return { success, message, status, resBody };
	}
	async getMostBought(fetch: Fetch) {
		const url = 'http://habilog.xyz/api/v1/courses/most-bought';
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<CourseList[]> = await res?.json();
		return { success, resBody, status, message };
	}
}
export const controller = new HomeController();
