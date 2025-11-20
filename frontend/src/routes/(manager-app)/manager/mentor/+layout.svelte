<script lang="ts">
	import { onMount } from 'svelte';
	import type { LayoutProps } from './$types';
	import { mentorRole } from '$lib/utils/constants';
	import { goto } from '$app/navigation';
	import { CreateToast } from '$lib/utils/helper';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';
	import { MentorLayoutView } from './view.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { Menu, House, List, LogOut, Folders, MessageCircle } from '@lucide/svelte';
	import Link from '$lib/components/button/Link.svelte';
	import Menubar from '$lib/components/menubar/Menubar.svelte';
	import MenuItem from '$lib/components/menubar/MenuItem.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';

	const { children, data }: LayoutProps = $props();

	const View = new MentorLayoutView();

	onMount(() => {
		if (data.role !== mentorRole) {
			goto('/login', { replaceState: true });
		}
		if (window.location.pathname !== '/manager/mentor/verify' && data.userStatus !== 'verified') {
			goto('/manager/mentor/verify', { replaceState: true });
		}
	});

	function onLogout() {
		const loadID = CreateToast('loading', 'logging out....');
		loadingStore.setLogOutLoadID(loadID);
		goto('/manager/logout', { replaceState: true });
	}
</script>

<div class="flex h-dvh flex-row">
	<nav
		class="
            hidden h-dvh flex-col items-center gap-4 bg-[var(--tertiary-color)] p-4 transition-all
            duration-300 ease-in-out md:flex
            {View.menuOpen ? 'w-[15%]' : 'w-20'} 
            "
	>
		<div class="mb-16 w-full justify-self-start">
			<Button onClick={() => View.handleMenu()} full>
				<Menu class="duration-300 ease-in-out {View.menuOpen ? 'rotate-90' : ''}" />
			</Button>
		</div>
		<Link href="/manager/mentor">
			<div class="flex flex-col items-center gap-1">
				<House />
				<p class="duration-300 ease-in-out {View.menuOpen ? 'opacity-100' : 'opacity-0'}">Home</p>
			</div>
		</Link>
		<Link href="/manager/mentor/courses">
			<div class="flex flex-col items-center gap-1">
				<Folders />
				<p class="duration-300 ease-in-out {View.menuOpen ? 'opacity-100' : 'opacity-0'}">
					Courses
				</p>
			</div>
		</Link>
		<Link href="/manager/mentor/requests">
			<div class="flex flex-col items-center gap-1">
				<List />
				<p class="duration-300 ease-in-out {View.menuOpen ? 'opacity-100' : 'opacity-0'}">
					Requests
				</p>
			</div>
		</Link>
		<Link href="/manager/mentor/messages">
			<div class="flex flex-col items-center gap-1">
				<MessageCircle />
				<p class="duration-300 ease-in-out {View.menuOpen ? 'opacity-100' : 'opacity-0'}">
					Messages
				</p>
			</div>
		</Link>
		<Button onClick={onLogout} withBg={false} textColor="light" withPadding={false}>
			<div class="flex flex-col items-center gap-1">
				<LogOut />
				<p class="duration-300 ease-in-out {View.menuOpen ? 'opacity-100' : 'opacity-0'}">Logout</p>
			</div>
		</Button>
	</nav>
	<main class="flex flex-1 pb-24 md:pb-0">
		<ScrollArea class="flex-1" orientation="vertical" viewportClasses="h-full max-h-[850px]">
			{@render children()}
		</ScrollArea>
	</main>
</div>
<div class="block md:hidden">
	<Menubar>
		<MenuItem>
			<Link href="/manager/mentor">
				<div class="flex flex-col items-center gap-1">
					<House />
					<p>Home</p>
				</div>
			</Link>
		</MenuItem>
		<MenuItem>
			<Link href="/manager/mentor/courses">
				<div class="flex flex-col items-center gap-1">
					<Folders />
					<p>Courses</p>
				</div>
			</Link>
		</MenuItem>
		<MenuItem>
			<Link href="/manager/mentor/requests">
				<div class="flex flex-col items-center gap-1">
					<List />
					<p>Requests</p>
				</div>
			</Link>
		</MenuItem>
		<MenuItem>
			<Link href="/manager/mentor/messages">
				<div class="flex flex-col items-center gap-1">
					<MessageCircle />
					<p>Messages</p>
				</div>
			</Link>
		</MenuItem>
		<MenuItem>
			<Button onClick={onLogout} withBg={false} textColor="light" withPadding={false}>
				<div class="flex flex-col items-center gap-1">
					<LogOut />
					<p>Logout</p>
				</div>
			</Button>
		</MenuItem>
	</Menubar>
</div>
