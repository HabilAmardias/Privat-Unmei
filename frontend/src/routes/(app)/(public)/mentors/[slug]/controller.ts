import type { Fetch, ServerResponse, PaginatedResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { MentorProfile, MentorScheduleInfo, MentorCourse } from './model';

class MentorProfileController {
	async getMentorProfile(fetch: Fetch, id: string) {
		const url = `http://localhost:8080/api/v1/mentors/${id}`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorProfile> = await res?.json();
		return { success, status, message, resBody };
	}
	async getMentorSchedules(fetch: Fetch, id: string) {
		const url = `http://localhost:8080/api/v1/mentors/${id}/availability`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorScheduleInfo[]> = await res?.json();
		return { success, status, message, resBody };
	}
	async getMentorCourses(fetch: Fetch, id: string, req?: Request) {
		let url = `http://localhost:8080/api/v1/mentors/${id}/courses?`;
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
}

export const controller = new MentorProfileController();
