import type { mentorList } from './model';

export class MentorManagerView {
	mentors = $state<mentorList[]>([]);
	page = $state<number>(1);
	limit = $state<number>(15);
	total_row = $state<number>(15);
	isDesktop = $state<boolean>(false);
	mentorsIsLoading = $state<boolean>(false);

	setMentors(newList: mentorList[]) {
		this.mentors = newList;
	}
	setPaginationData(page: number, limit: number, total_row: number) {
		this.page = page;
		this.limit = limit;
		this.total_row = total_row;
	}
	setIsDesktop(b: boolean) {
		this.isDesktop = b;
	}
}
