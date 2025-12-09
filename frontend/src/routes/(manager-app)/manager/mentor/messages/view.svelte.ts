import type { EnhancementArgs, EnhancementReturn, PaginatedResponse } from '$lib/types';
import { CreateToast } from '$lib/utils/helper';
import type { Chatroom } from './model';

export class ChatListView {
	chats = $state<Chatroom[]>([]);
	page = $state<number>(1);
	totalRow = $state<number>(15);
	limit = $state<number>(15);
	isLoading = $state<boolean>(false);

	constructor(d: PaginatedResponse<Chatroom>) {
		this.chats = d.entries;
		this.page = d.page_info.page;
		this.totalRow = d.page_info.total_row;
		this.limit = d.page_info.limit;
	}

	onPageChange = async (args: EnhancementArgs) => {
		this.isLoading = true;
		args.formData.append('page', `${this.page}`);
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'success') {
				this.chats = result.data?.chatrooms.entries;
				this.page = result.data?.chatrooms.page_info.page;
				this.totalRow = result.data?.chatrooms.page_info.total_row;
				this.limit = result.data?.chatrooms.page_info.limit;
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
