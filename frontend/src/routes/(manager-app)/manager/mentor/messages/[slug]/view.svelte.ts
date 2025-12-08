import type { EnhancementReturn, SeekPaginatedResponse } from '$lib/types';
import { CreateToast } from '$lib/utils/helper';
import type { MessageInfo } from './model';

export class ChatroomView {
	messages = $state<MessageInfo[]>([]);
	lastID = $state<number>(15);
	limit = $state<number>(15);
	messageContent = $state<string>('');
	disableSendMessage = $derived<boolean>(
		this.messageContent.length > 180 || this.messageContent.length === 0
	);
	endRef = $state<HTMLDivElement>();
	constructor(m: SeekPaginatedResponse<MessageInfo>) {
		this.messages = m.entries;
		this.lastID = m.page_info.last_id;
		this.limit = m.page_info.limit;
	}
	onSendMessage = () => {
		return async ({ result }: EnhancementReturn) => {
			if (result.type === 'success') {
				this.endRef?.scrollIntoView({ behavior: 'smooth' });
				this.messageContent = '';
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
