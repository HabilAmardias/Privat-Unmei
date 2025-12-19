import type { Fetch, ServerResponse } from '$lib/types';
import type { CourseCategory, CourseDetail, CourseTopic } from './model';
import { FetchData } from '$lib/utils';

class CourseDetailController {
	async getCourseDetail(fetch: Fetch, id: string) {
		const url = `http://localhost/api/v1/courses/${id}`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<CourseDetail> = await res?.json();
		return { success, message, status, resBody };
	}
	async getCourseTopics(fetch: Fetch, id: string) {
		const url = `http://localhost/api/v1/courses/${id}/topics`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<CourseTopic[]> = await res?.json();
		return { success, message, status, resBody };
	}
	async getCourseDetailCategories(fetch: Fetch, id: string) {
		const url = `http://localhost/api/v1/courses/${id}/categories`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<CourseCategory[]> = await res?.json();
		return { success, message, status, resBody };
	}
	async deleteCourse(fetch: Fetch, id: string) {
		if (!id) {
			return { success: false, message: 'no course selected', status: 400 };
		}
		const url = `http://localhost/api/v1/courses/${id}`;
		const { success, message, status } = await FetchData(fetch, url, 'DELETE');
		return { success, message, status };
	}
}

export const controller = new CourseDetailController();
