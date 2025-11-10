import type { PaginatedResponse } from '$lib/types';
import { MAX_BIO_LENGTH } from '$lib/utils/constants';
import { IsAlphaOnly } from '$lib/utils/helper';
import type { StudentOrders, StudentProfile } from './model';

export class profileView {
	verifyIsLoading = $state<boolean>(false);
	ordersIsLoading = $state<boolean>(false);
	profileIsLoading = $state<boolean>(false);
	isDesktop = $state<boolean>(false);
	isEdit = $state<boolean>(false);
	name = $state<string>('');
	bio = $state<string>('');
	profileImage = $state<FileList>();
	status = $state<string>('');
	totalRow = $state<number>(1); // temporary
	limit = $state<number>(15);
	pageNumber = $state<number>(1);
	paginationForm = $state<HTMLFormElement | undefined>();
	nameError = $state<Error | undefined>();
	bioError = $state<Error | undefined>();
	orders = $state<StudentOrders[]>();

	updateProfileDisable = $derived.by<boolean>(() => {
		if (this.nameError || this.profileIsLoading || this.bioError) {
			return true;
		}
		return false;
	});

	size = $derived.by<number>(() => {
		if (this.isDesktop) {
			return 150;
		}
		return 100;
	});
	constructor(s: StudentProfile, o?: PaginatedResponse<StudentOrders>) {
		this.setBio(s.bio);
		this.setName(s.name);
		if (o) {
			this.setOrders(o.entries);
			this.setTotalRow(o.page_info.total_row);
			this.setPageNumber(o.page_info.page);
		}
	}
	setPageNumber(num: number) {
		this.pageNumber = num;
	}
	onPageChange(num: number) {
		this.pageNumber = num;
		this.paginationForm?.requestSubmit();
	}
	setTotalRow(row: number) {
		this.totalRow = row;
	}
	setOrdersIsLoading(b: boolean) {
		this.ordersIsLoading = b;
	}
	setVerifyIsLoading(b: boolean) {
		this.verifyIsLoading = b;
	}
	setProfileIsLoading(b: boolean) {
		this.profileIsLoading = b;
	}
	setBio(newBio: string) {
		this.bio = newBio;
	}
	setName(newName: string) {
		this.name = newName;
	}
	setIsDesktop(b: boolean) {
		this.isDesktop = b;
	}
	nameOnBlur() {
		if (this.name && !IsAlphaOnly(this.name)) {
			this.nameError = new Error('name must only contain alphabets');
		} else {
			this.nameError = undefined;
		}
	}
	bioOnBlur() {
		if (this.bio.length > MAX_BIO_LENGTH) {
			this.bioError = new Error(`bio must not more than ${MAX_BIO_LENGTH} characters`);
		} else {
			this.bioError = undefined;
		}
	}
	setBioError(e: Error | undefined) {
		this.bioError = e;
	}
	setNameError(e: Error | undefined) {
		this.nameError = e;
	}
	setIsEdit() {
		this.isEdit = !this.isEdit;
	}
	setOrders(newOrders: StudentOrders[]) {
		this.orders = newOrders;
	}
	setProfileImage(f: FileList | undefined) {
		this.profileImage = f;
	}
}
