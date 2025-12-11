<script lang="ts">
	import Link from '$lib/components/button/Link.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import type { PageProps } from './$types';
	import { ChatListView } from './view.svelte';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import { enhance } from '$app/forms';
	import Loading from '$lib/components/loader/Loading.svelte';

	const { data }: PageProps = $props();

	const View = new ChatListView(data.chatrooms);
</script>

<svelte:head>
	<title>Chats - Privat Unmei</title>
	<meta name="description" content="My Chats - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-full w-full flex-col gap-4 p-4">
	<h1 class="text-2xl font-bold text-[var(--tertiary-color)]">Chats</h1>
	{#if View.isLoading}
		<Loading />
	{:else if View.chats.length === 0 || !View.chats}
		<b class="font-bold text-[var(--tertiary-color)]">No Chats Found</b>
	{:else}
		<ScrollArea
			orientation="vertical"
			viewportClasses="h-[550px] max-h-[550px] md:h-[500px] md:max-h-[500px]"
		>
			<ul class="flex gap-4">
				{#each View.chats as ch (ch.id)}
					<li class="w-full rounded-lg bg-[var(--tertiary-color)] p-2">
						<Link href={`/manager/mentor/messages/${ch.id}`}>
							<div class="flex items-center gap-4">
								<CldImage src={ch.profile_image} width={70} height={70} className="rounded-full" />
								<div>
									<p class="text-[var(--primary-color)]">{ch.username}</p>
									<p class="text-[var(--secondary-color)]">{ch.public_id}</p>
								</div>
							</div>
						</Link>
					</li>
				{/each}
			</ul>
		</ScrollArea>
	{/if}
	<form
		class="flex items-center justify-center"
		use:enhance={View.onPageChange}
		action="?/getChats"
		method="POST"
	>
		<Pagination bind:pageNumber={View.page} perPage={View.limit} count={View.totalRow} />
	</form>
</div>
