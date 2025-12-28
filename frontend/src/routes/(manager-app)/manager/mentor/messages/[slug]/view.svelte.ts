import type { EnhancementArgs, EnhancementReturn, SeekPaginatedResponse } from '$lib/types';
import { CreateToast } from '$lib/utils/helper';
import type { MessageInfo } from './model';
import { type BeforeNavigate } from '@sveltejs/kit';
import { goto } from '$app/navigation';
import { resolve } from '$app/paths';

export class ChatroomView {
	messages = $state<MessageInfo[]>([]);
	lastID = $state<number>(15);
	limit = $state<number>(15);
	totalRow = $state<number>(15);
	messageContent = $state<string>('');
	isNavigatingAfterSubmit = $state<boolean>(false);
	disableSendMessage = $derived<boolean>(
		this.messageContent.length > 180 || this.messageContent.length === 0
	);
	updateLastReadForm = $state<HTMLFormElement>();
	endRef = $state<HTMLDivElement>();
	getMessageForm = $state<HTMLFormElement>();
	isInitialLoad = $state<boolean>(true);
	isLoading = $state<boolean>(false);
	viewPortRef = $state<HTMLDivElement | null>(null);
	prevScrollHeight = $state<number>(0);

	constructor(m: SeekPaginatedResponse<MessageInfo>) {
		this.messages = m.entries;
		this.lastID = m.page_info.last_id;
		this.limit = m.page_info.limit;
		this.totalRow = m.page_info.total_row;
	}
	onNavigate = (n: BeforeNavigate) => {
		if (this.isNavigatingAfterSubmit) {
			this.isNavigatingAfterSubmit = false;
			return;
		}
		if (n.to?.route.id === '/(manager-app)/manager/mentor/messages') {
			n.cancel();
			this.updateLastReadForm?.requestSubmit();
			return;
		}
		const formData = new FormData();
		navigator.sendBeacon('?/updateLastRead', formData);
	};
	scrollToBottom = (smooth: boolean = true) => {
		if (this.endRef) {
			this.endRef.scrollIntoView({ behavior: smooth ? 'smooth' : 'auto' });
		}
	};
	handleInitialScroll = () => {
		if (this.isInitialLoad && this.messages.length > 0) {
			this.scrollToBottom(false);
			this.isInitialLoad = false;
		}
	};
	onUpdateLastRead = () => {
		return async ({ update }: EnhancementReturn) => {
			await update({ reset: false });
			this.isNavigatingAfterSubmit = true;
			goto(resolve('/(manager-app)/manager/mentor/messages'));
		};
	};
	onSendMessage = () => {
		return async ({ result }: EnhancementReturn) => {
			if (result.type === 'success') {
				this.scrollToBottom(true);
				this.messageContent = '';
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	restoreScrollPosition = () => {
		if (this.viewPortRef && this.prevScrollHeight > 0) {
			const newScrollHeight = this.viewPortRef.scrollHeight;
			const heightDiff = newScrollHeight - this.prevScrollHeight;
			this.viewPortRef.scrollTop += heightDiff;
			this.prevScrollHeight = 0;
		}
	};
	onIntersect = () => {
		if (this.messages.length < this.totalRow) {
			if (this.viewPortRef) {
				this.prevScrollHeight = this.viewPortRef.scrollHeight;
			}
			this.getMessageForm?.requestSubmit();
		}
	};
	onGetMessage = (args: EnhancementArgs) => {
		this.isLoading = true;
		args.formData.append('last_id', `${this.lastID}`);
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'success') {
				const newMessages: MessageInfo[] = result.data?.messages.entries;
				this.messages = [...newMessages, ...this.messages];
				this.lastID = result.data?.messages.page_info.last_id;
				this.limit = result.data?.messages.page_info.limit;
				this.totalRow = result.data?.messages.page_info.total_row;
				this.restoreScrollPosition();
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
