<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$lib/components/button/Button.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import Menubar from '$lib/components/menubar/Menubar.svelte';
	import MenuItem from '$lib/components/menubar/MenuItem.svelte';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';
	import { List, LogIn, LogOut, MessageCircleMore, User } from '@lucide/svelte';
	import { ScrollArea } from 'bits-ui';
	import toast from 'svelte-french-toast';
	import type { LayoutProps } from './$types';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { onMount } from 'svelte';
	import { AppLayoutView } from './view.svelte';

	function onLogout() {
		const loadID = toast.loading('logging out....', { position: 'top-right' });
		loadingStore.setLogOutLoadID(loadID);
		goto('/logout', { replaceState: true });
	}
	const View = new AppLayoutView();

	onMount(() => {
		View.setIsDesktop(window.innerWidth >= 768);
		function setIsDesktop() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		window.addEventListener('resize', setIsDesktop);
		return () => {
			window.removeEventListener('resize', setIsDesktop);
		};
	});
	let { data, children }: LayoutProps = $props();
</script>

<main class="h-screen pb-24 md:pb-0 md:pt-24">
	<ScrollArea.Root class="h-full">
		<ScrollArea.Viewport class="h-full">
			{@render children()}
		</ScrollArea.Viewport>
	</ScrollArea.Root>
</main>
<Menubar>
	<MenuItem>
		<Link href="/courses">
			<div class="flex h-full items-center">
				<CldImage src={View.logoSrc} width={View.logoWidth} height={View.logoHeight} />
			</div>
		</Link>
	</MenuItem>
	<MenuItem>
		<Link href="/courses">
			<div class="flex flex-col items-center">
				<List />
				Courses
			</div>
		</Link>
	</MenuItem>
	{#if data.isLoggedIn}
		<MenuItem>
			<Link href="/chats">
				<div class="flex flex-col items-center">
					<MessageCircleMore />
					Chats
				</div>
			</Link>
		</MenuItem>
		<MenuItem>
			<Link href="/profile">
				<div class="flex flex-col items-center">
					<User />
					Profile
				</div>
			</Link>
		</MenuItem>
		<MenuItem>
			<Button onClick={onLogout} withBg={false} textColor="light" withPadding={false}>
				<div class="flex flex-col items-center">
					<LogOut />
					Logout
				</div>
			</Button>
		</MenuItem>
	{:else}
		<MenuItem>
			<Link href="/login">
				<div class="flex flex-col items-center">
					<LogIn />
					Login
				</div>
			</Link>
		</MenuItem>
	{/if}
</Menubar>
