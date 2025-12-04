import type { PaginatedResponse } from '$lib/types';
import { MAX_BIO_LENGTH } from '$lib/utils/constants';
import { IsAlphaOnly } from '$lib/utils/helper';
import type { StudentOrders, StudentProfile } from './model';
import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';
export class profileView {
	verifyIsLoading = $state<boolean>(false);
	ordersIsLoading = $state<boolean>(false);
	profileIsLoading = $state<boolean>(false);
	search = $state<string>('');
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
	nameError = $derived.by<Error | undefined>(() => {
		if (this.name && !IsAlphaOnly(this.name)) {
			return new Error('name must only contain alphabets');
		}
		return undefined;
	});
	bioError = $derived.by<Error | undefined>(() => {
		if (this.bio.length > MAX_BIO_LENGTH) {
			return new Error('name must only contain alphabets');
		}
		return undefined;
	});
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
	onPageChange = () => {
		this.paginationForm?.requestSubmit();
	};
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
	onUpdateProfile = () => {
		const loadID = CreateToast('loading', 'updating....');
		this.setProfileIsLoading(true);
		return async ({ result, update }: EnhancementReturn) => {
			this.setIsEdit();
			this.setProfileIsLoading(false);
			DismissToast(loadID);
			if (result.type === 'success') {
				CreateToast('success', 'update profile success');
				await update();
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onSetPage = (args: EnhancementArgs) => {
		this.setOrdersIsLoading(true);
		args.formData.append('page', `${this.pageNumber}`);
		args.formData.append('status', `${this.status}`);
		args.formData.append('search', this.search);
		return async ({ result }: EnhancementReturn) => {
			this.setOrdersIsLoading(false);
			if (result.type === 'success') {
				this.setOrders(result.data?.orders);
				this.setTotalRow(result.data?.totalRow);
				CreateToast('success', result.data?.message);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};

	onSearchOrders = (args: EnhancementArgs) => {
		this.setOrdersIsLoading(true);
		this.pageNumber = 1;
		args.formData.append('page', `${this.pageNumber}`);
		return async ({ result }: EnhancementReturn) => {
			this.setOrdersIsLoading(false);
			if (result.type === 'success') {
				this.setOrders(result.data?.orders);
				this.setTotalRow(result.data?.totalRow);
				CreateToast('success', result.data?.message);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};

	onVerifySubmit = () => {
		this.setVerifyIsLoading(true);
		const loadID = CreateToast('loading', 'sending....');
		return async ({ result }: EnhancementReturn) => {
			this.setVerifyIsLoading(false);
			DismissToast(loadID);
			if (result.type === 'success') {
				CreateToast('success', result.data?.message);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	capitalizeFirstLetter(s: string) {
		if (s.length === 0) {
			return s;
		}
		if (s.length === 1) {
			return s.toUpperCase();
		}
		return s.charAt(0).toUpperCase() + s.slice(1);
	}
}
