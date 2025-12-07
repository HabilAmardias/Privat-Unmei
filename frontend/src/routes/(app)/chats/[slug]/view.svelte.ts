import type { SeekPaginatedResponse } from '$lib/types';
import type { MessageInfo } from './model';

export class ChatroomView {
	messages = $state<MessageInfo[]>([]);
	lastID = $state<number>(15);
	limit = $state<number>(15);

	constructor(m: SeekPaginatedResponse<MessageInfo>) {
		this.messages = m.entries;
		this.lastID = m.page_info.last_id;
		this.limit = m.page_info.limit;
	}
}
