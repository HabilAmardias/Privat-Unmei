import type { Fetch, ServerResponse } from '$lib/types';
import type { RequestDetail, ChatroomID, CreateReview, IsReviewed } from './model';
import { FetchData } from '$lib/utils';

class RequestDetailController {
	async getRequestDetail(fetch: Fetch, id: string) {
		const url = `/api/v1/me/course-requests/${id}`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<RequestDetail> = await res?.json();
		return { success, message, status, resBody };
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
	async createReview(fetch: Fetch, req: Request) {
		const formData = await req.formData();
		const rating = formData.get('rating');
		const feedback = formData.get('feedback');
		const id = formData.get('id');

		const url = `/api/v1/courses/${id}/reviews`;

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
	async isCourseReviewed(fetch: Fetch, courseID: number) {
		const url = `/api/v1/courses/${courseID}/reviews/me`;
		const { success, status, message, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<IsReviewed> = await res?.json();
		return { success, status, message, resBody };
	}
}

export const controller = new RequestDetailController();
