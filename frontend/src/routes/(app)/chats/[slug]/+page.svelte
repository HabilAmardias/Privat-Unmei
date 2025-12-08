<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import type { MessageInfo } from './model';
	import { ChatroomView } from './view.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import { enhance } from '$app/forms';
	import Input from '$lib/components/form/Input.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { Send } from '@lucide/svelte';

	const { data, params }: PageProps = $props();

	const View = new ChatroomView(data.messages);

	onMount(() => {
		const url = `ws://localhost:8080/ws/v1/chatrooms/${params.slug}/messages`;
		const socket = new WebSocket(url);
		socket.onopen = () => {
			console.log('Connection estabilished');
		};
		socket.onmessage = (ev: MessageEvent<string>) => {
			const msg: MessageInfo = JSON.parse(ev.data);
			View.messages.push(msg);
			View.limit += 1;
			View.lastID = msg.id;
		};
		socket.onclose = () => {
			console.log('connection closed');
		};
		socket.onerror = () => {
			console.warn('error');
		};
		return () => {
			socket.close();
		};
	});
</script>

<div class="flex flex-col gap-4 p-4">
	<div class="flex w-full items-center gap-4 rounded-lg bg-[var(--tertiary-color)] p-2">
		<CldImage src={data.chatroom.profile_image} width={70} height={70} className="rounded-full" />
		<div>
			<p class="text-[var(--primary-color)]">{data.chatroom.username}</p>
			<p class="text-[var(--secondary-color)]">{data.chatroom.email}</p>
		</div>
	</div>
	<ScrollArea orientation="vertical" viewportClasses={`h-[550px] max-h-[550px]`}>
		<ul class="flex h-full w-full flex-col justify-end gap-2">
			{#each View.messages as msg (msg.id)}
				<li
					class={`w-fit rounded-lg bg-[var(--tertiary-color)] p-2 ${data.profile.id === msg.sender_id ? 'ml-auto' : ''}`}
				>
					<p class="text-[var(--secondary-color)]">{msg.content}</p>
				</li>
			{/each}
		</ul>
		<div bind:this={View.endRef}></div>
	</ScrollArea>
	<form class="flex gap-2" use:enhance={View.onSendMessage} action="?/sendMessage" method="POST">
		<Input
			bind:value={View.messageContent}
			type="text"
			name="message"
			id="message"
			placeholder="Write a message"
		/>
		<Button disabled={View.disableSendMessage} type="submit"><Send /></Button>
	</form>
</div>
