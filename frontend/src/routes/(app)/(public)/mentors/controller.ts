import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { MentorList } from './model';

class MentorsController {
	async getMentors(fetch: Fetch, req?: Request) {
		let url = 'http://habilog.xyz/api/v1/mentors?';
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
		const resBody: ServerResponse<PaginatedResponse<MentorList>> = await res?.json();
		return { success, message, status, resBody };
	}
}
export const controller = new MentorsController();
