import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { CourseCategory, MentorCourse } from './model';

class CourseManagementController {
	async getMyCourses(fetch: Fetch, req?: Request) {
		let url = '/api/v1/mentors/me/courses?';
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
			if (args.length > 0) {
				url += args.join('&');
			}
		}
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<PaginatedResponse<MentorCourse>> = await res?.json();
		return { success, message, status, resBody };
	}
	async deleteCourse(fetch: Fetch, req: Request) {
		const formData = await req.formData();
		const id = formData.get('id');
		if (!id) {
			return { success: false, message: 'no course selected', status: 400 };
		}
		const url = `/api/v1/courses/${id}`;
		const { success, message, status } = await FetchData(fetch, url, 'DELETE');
		return { success, message, status };
	}
	async getCourseCategories(fetch: Fetch, req?: Request) {
		let url = '/api/v1/course-categories?';
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

export const controller = new CourseManagementController();
