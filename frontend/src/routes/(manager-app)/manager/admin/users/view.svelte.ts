import type { EnhancementArgs, EnhancementReturn, PaginatedResponse } from '$lib/types';
import type { mentorList, studentList } from './model';
import { CreateToast } from '$lib/utils/helper';
import { debounce } from '$lib/utils/helper';

export class UsersManagerView {
	isDesktop = $state<boolean>(false);
	menuState = $state<'mentors' | 'students'>('mentors');

	mentors = $state<mentorList[]>([]);
	mentorPage = $state<number>(1);
	mentorLimit = $state<number>(15);
	mentorTotalRow = $state<number>(15);
	mentorsIsLoading = $state<boolean>(false);
	mentorToDelete = $state<string>();
	sortByYears = $state<boolean | null>(null);
	mentorSearch = $state<string>('');
	mentorsAlertOpen = $state<boolean[]>([]);
	mentorsSearchForm = $state<HTMLFormElement | null>(null);
	#mentorSearchSubmit = debounce(() => {
		this.mentorsSearchForm?.requestSubmit();
	}, 500);

	students = $state<studentList[]>([]);
	studentPage = $state<number>(1);
	studentLimit = $state<number>(15);
	studentTotalRow = $state<number>(15);
	studentIsLoading = $state<boolean>(false);
	studentToDelete = $state<string>();
	studentSearch = $state<string>('');
	studentsAlertOpen = $state<boolean[]>([]);
	studentsSearchForm = $state<HTMLFormElement | null>(null);
	#studentSearchSubmit = debounce(() => {
		this.studentsSearchForm?.requestSubmit();
	}, 500);

	constructor(m: PaginatedResponse<mentorList>, s: PaginatedResponse<studentList>) {
		this.mentors = m.entries;
		this.students = s.entries;
		this.mentorsAlertOpen = new Array<boolean>(this.mentors.length).fill(false);
		this.studentsAlertOpen = new Array<boolean>(this.students.length).fill(false);

		this.mentorPage = m.page_info.page;
		this.mentorLimit = m.page_info.limit;
		this.mentorTotalRow = m.page_info.total_row;

		this.studentPage = s.page_info.page;
		this.studentLimit = s.page_info.limit;
		this.studentTotalRow = s.page_info.total_row;
	}

	onMentorSearchInput = (e: Event & { currentTarget: EventTarget & HTMLInputElement }) => {
		this.mentorSearch = e.currentTarget.value;
		this.#mentorSearchSubmit();
	};
	onStudentSearchInput = (e: Event & { currentTarget: EventTarget & HTMLInputElement }) => {
		this.studentSearch = e.currentTarget.value;
		this.#studentSearchSubmit();
	};
	removeMentors(id: string) {
		this.mentors = this.mentors.filter((m) => m.id !== id);
	}
	removeStudents(id: string) {
		this.students = this.students.filter((s) => s.id !== id);
	}
	setIsDesktop(b: boolean) {
		this.isDesktop = b;
	}
	onSort() {
		switch (this.sortByYears) {
			case null:
				this.sortByYears = true;
				break;
			case true:
				this.sortByYears = false;
				break;
			case false:
				this.sortByYears = null;
				break;
		}
	}
	resetFilterForm() {
		if (this.menuState === 'mentors') {
			this.mentorSearch = '';
			this.sortByYears = null;
		} else {
			this.studentSearch = '';
		}
	}
	onDeleteMentor = (args: EnhancementArgs) => {
		this.mentorsIsLoading = true;
		if (this.mentorToDelete) {
			args.formData.append('id', this.mentorToDelete);
		}
		return async ({ result }: EnhancementReturn) => {
			this.mentorsIsLoading = false;
			if (result.type === 'success') {
				if (this.mentorToDelete) {
					this.removeMentors(this.mentorToDelete);
					this.mentorsAlertOpen = new Array<boolean>(this.mentors.length).fill(false);
				}
				this.mentorToDelete = undefined;
				this.mentorTotalRow -= 1;
				CreateToast('success', 'successfully delete mentor');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onSearchMentors = (args: EnhancementArgs) => {
		this.mentorsIsLoading = true;
		if (this.sortByYears !== null) {
			args.formData.append('sort_year_of_experience', `${this.sortByYears}`);
		}
		this.mentorPage = 1;
		args.formData.append('page', `${this.mentorPage}`);
		return async ({ result }: EnhancementReturn) => {
			this.mentorsIsLoading = false;
			if (result.type === 'success') {
				this.mentors = result.data?.mentorsList.entries;
				this.mentorsAlertOpen = new Array<boolean>(this.mentors.length).fill(false);
				this.mentorPage = result.data?.mentorsList.page_info.page;
				this.mentorLimit = result.data?.mentorsList.page_info.limit;
				this.mentorTotalRow = result.data?.mentorsList.page_info.total_row;
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onPageChangeMentors = (args: EnhancementArgs) => {
		this.mentorsIsLoading = true;
		if (this.sortByYears !== null) {
			args.formData.append('sort_year_of_experience', `${this.sortByYears}`);
		}
		args.formData.append('search', this.mentorSearch);
		args.formData.append('page', `${this.mentorPage}`);
		return async ({ result }: EnhancementReturn) => {
			this.mentorsIsLoading = false;
			if (result.type === 'success') {
				this.mentors = result.data?.mentorsList.entries;
				this.mentorsAlertOpen = new Array<boolean>(this.mentors.length).fill(false);
				this.mentorPage = result.data?.mentorsList.page_info.page;
				this.mentorLimit = result.data?.mentorsList.page_info.limit;
				this.mentorTotalRow = result.data?.mentorsList.page_info.total_row;
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onDeleteStudent = (args: EnhancementArgs) => {
		this.studentIsLoading = true;
		if (this.studentToDelete) {
			args.formData.append('id', this.studentToDelete);
		}
		return async ({ result }: EnhancementReturn) => {
			this.studentIsLoading = false;
			if (result.type === 'success') {
				if (this.studentToDelete) {
					this.removeStudents(this.studentToDelete);
					this.studentsAlertOpen = new Array<boolean>(this.students.length).fill(false);
				}
				this.studentToDelete = undefined;
				this.studentTotalRow -= 1;
				CreateToast('success', 'successfully delete student');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onSearchStudent = (args: EnhancementArgs) => {
		this.studentIsLoading = true;
		this.studentPage = 1;
		args.formData.append('page', `${this.studentPage}`);
		return async ({ result }: EnhancementReturn) => {
			this.studentIsLoading = false;
			if (result.type === 'success') {
				this.students = result.data?.studentList.entries;
				this.studentsAlertOpen = new Array<boolean>(this.students.length).fill(false);
				this.studentPage = result.data?.studentList.page_info.page;
				this.studentLimit = result.data?.studentList.page_info.limit;
				this.studentTotalRow = result.data?.studentList.page_info.total_row;
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onPageChangeStudent = (args: EnhancementArgs) => {
		this.studentIsLoading = true;
		args.formData.append('search', this.studentSearch);
		args.formData.append('page', `${this.studentPage}`);
		return async ({ result }: EnhancementReturn) => {
			this.studentIsLoading = false;
			if (result.type === 'success') {
				this.students = result.data?.studentList.entries;
				this.studentsAlertOpen = new Array<boolean>(this.students.length).fill(false);
				this.studentPage = result.data?.studentList.page_info.page;
				this.studentLimit = result.data?.studentList.page_info.limit;
				this.studentTotalRow = result.data?.studentList.page_info.total_row;
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
