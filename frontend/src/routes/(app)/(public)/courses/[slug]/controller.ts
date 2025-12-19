import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import type {
	CourseCategory,
	CourseDetail,
	CourseReview,
	CourseTopic,
	CreateReview,
	StudentProfile,
	ChatroomID
} from './model';
import { FetchData } from '$lib/utils';

class CourseDetailController {
	async getCourseDetail(fetch: Fetch, id: string) {
		const url = `/api/v1/courses/${id}`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<CourseDetail> = await res?.json();
		return { success, message, status, resBody };
	}
	async getCourseTopics(fetch: Fetch, id: string) {
		const url = `/api/v1/courses/${id}/topics`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<CourseTopic[]> = await res?.json();
		return { success, message, status, resBody };
	}
	async getCourseDetailCategories(fetch: Fetch, id: string) {
		const url = `/api/v1/courses/${id}/categories`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<CourseCategory[]> = await res?.json();
		return { success, message, status, resBody };
	}
	async getCourseReviews(fetch: Fetch, id: string, req?: Request) {
		let url = `/api/v1/courses/${id}/reviews?`;
		if (req) {
			const formData = await req.formData();
			const page = formData.get('page');
			if (page) {
				url += `page=${page}`;
			}
		}
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<PaginatedResponse<CourseReview>> = await res?.json();
		return { success, message, status, resBody };
	}
	async createReview(fetch: Fetch, req: Request, id: string) {
		const url = `/api/v1/courses/${id}/reviews`;
		const formData = await req.formData();
		const rating = formData.get('rating');
		const feedback = formData.get('feedback');

		const reqBody = JSON.stringify({
			rating: parseInt(rating as string),
			feedback
		});

		const { success, status, message, res } = await FetchData(fetch, url, 'POST', reqBody);
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<CreateReview> = await res?.json();
		return { success, status, message, resBody };
	}
	async getProfile(fetch: Fetch) {
		const url = '/api/v1/me';
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<StudentProfile> = await res?.json();
		return { resBody, status, success, message };
	}
	async messageMentor(fetch: Fetch, req: Request) {
		const formData = await req.formData();
		const id = formData.get('id');
		const url = `/api/v1/chatrooms/users/${id}`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<ChatroomID> = await res?.json();
		return { resBody, status, success, message };
	}
}

export const controller = new CourseDetailController();
