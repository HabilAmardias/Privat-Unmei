import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { adminProfile, mentorList } from './model';

class MentorManagerController {
	async deleteMentor(fetch: Fetch, req: Request) {
		const formData = await req.formData();
		const mentorID = formData.get('id');
		if (!mentorID) {
			return { success: false, message: 'no mentor selected', status: 400 };
		}
		const url = `http://160.19.167.63/api/v1/mentors/${mentorID}`;
		const { success, message, status } = await FetchData(fetch, url, 'DELETE');
		return { success, message, status };
	}
	async getMentors(fetch: Fetch, req?: Request) {
		let url = 'http://160.19.167.63/api/v1/mentors?';
		const queries: string[] = [];
		if (req) {
			const formData = await req.formData();
			const search = formData.get('search');
			const sortYearOfExperience = formData.get('sort_year_of_experience');
			const page = formData.get('page');
			if (search) {
				queries.push(`search=${search}`);
			}
			if (sortYearOfExperience) {
				queries.push(`sort_year_of_experience=${sortYearOfExperience}`);
			}
			if (page) {
				queries.push(`page=${page}`);
			}
			url += queries.join('&');
		}
		const { success, message, res, status } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<PaginatedResponse<mentorList>> = await res?.json();
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

export const controller = new MentorManagerController();
