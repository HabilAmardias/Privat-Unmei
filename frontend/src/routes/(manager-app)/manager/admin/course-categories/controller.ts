import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { CourseCategory, NewCategory, adminProfile } from './model';

class CourseCategoryManagementController {
	async deleteCategory(fetch: Fetch, req: Request) {
		const formData = await req.formData();
		const id = formData.get('id');
		if (!id) {
			return { success: false, status: 400, message: 'no course category selected' };
		}
		const url = `http://160.19.167.63/api/v1/course-categories/${id}`;
		const { success, message, status } = await FetchData(fetch, url, 'DELETE');
		return { success, message, status };
	}
	async updateCategory(fetch: Fetch, req: Request) {
		const formData = await req.formData();
		const id = formData.get('id');
		const name = formData.get('name');

		if (!id) {
			return { success: false, status: 400, message: 'no course category selected' };
		}
		const reqBody = JSON.stringify({
			name
		});

		const url = `http://160.19.167.63/api/v1/course-categories/${id}`;

		const { success, message, status } = await FetchData(fetch, url, 'PATCH', reqBody);
		return { success, message, status };
	}
	async createCategory(fetch: Fetch, req: Request) {
		const formData = await req.formData();
		const name = formData.get('name');
		if (!name) {
			return { success: false, status: 400, message: 'no name given' };
		}
		const reqBody = JSON.stringify({
			name
		});
		const url = 'http://160.19.167.63/api/v1/course-categories';
		const { success, message, status, res } = await FetchData(fetch, url, 'POST', reqBody);
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<NewCategory> = await res?.json();
		return { success, message, status, resBody };
	}
	async getCategories(fetch: Fetch, req?: Request) {
		let url = 'http://160.19.167.63/api/v1/course-categories?';
		if (req) {
			const args: string[] = [];
			const formData = await req.formData();
			const search = formData.get('search');
			const page = formData.get('page');
			if (search) {
				args.push(`search=${search}`);
			}
			if (page) {
				args.push(`page=${page}`);
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
	async getProfile(fetch: Fetch) {
		const url = 'http://160.19.167.63/api/v1/admins/me';
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<adminProfile> = await res?.json();
		return { success, message, status, resBody };
	}
}

export const controller = new CourseCategoryManagementController();
