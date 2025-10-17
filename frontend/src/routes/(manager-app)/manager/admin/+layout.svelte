<script lang="ts">
	import { onMount } from 'svelte';
	import type { LayoutProps } from './$types';
	import { adminRole } from '$lib/utils/constants';
	import { goto } from '$app/navigation';
	import Link from '$lib/components/button/Link.svelte';
	import { ScrollArea } from 'bits-ui';
	import { House, LogOut, Menu, Percent, User } from '@lucide/svelte';
	import Button from '$lib/components/button/Button.svelte';
	import Menubar from '$lib/components/menubar/Menubar.svelte';
	import MenuItem from '$lib/components/menubar/MenuItem.svelte';
	import toast from 'svelte-french-toast';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';

	const { children, data }: LayoutProps = $props();

	onMount(() => {
		if (data.role !== adminRole) {
			goto('/courses', { replaceState: true });
		}
	});
	let menuOpen = $state<boolean>(false);
	function handleMenu() {
		menuOpen = !menuOpen;
	}
	function onLogout() {
		const loadID = toast.loading('logging out....', { position: 'top-right' });
		loadingStore.setLogOutLoadID(loadID);
		goto('/manager/logout', { replaceState: true });
	}
</script>

<main class="flex h-dvh">
	<nav
		class="
            hidden h-dvh flex-col items-center gap-4 bg-[var(--tertiary-color)] p-4 transition-all
            duration-300 ease-in-out md:flex
            {menuOpen ? 'w-[15%]' : 'w-20'} 
            "
	>
		<div class="mb-16 w-full justify-self-start">
			<Button onClick={handleMenu} full>
				<Menu class="duration-300 ease-in-out {menuOpen ? 'rotate-90' : ''}" />
			</Button>
		</div>
		<Link href="/manager/admin">
			<div class="flex flex-col items-center gap-1">
				<House />
				<p class="duration-300 ease-in-out {menuOpen ? 'opacity-100' : 'opacity-0'}">Home</p>
			</div>
		</Link>
		<Link href="/manager/admin/mentors">
			<div class="flex flex-col items-center gap-1">
				<User />
				<p class="duration-300 ease-in-out {menuOpen ? 'opacity-100' : 'opacity-0'}">Mentors</p>
			</div>
		</Link>
		<Link href="/manager/admin/maintenance">
			<div class="flex flex-col items-center gap-1">
				<Percent />
				<p class="duration-300 ease-in-out {menuOpen ? 'opacity-100' : 'opacity-0'}">Maintenance</p>
			</div>
		</Link>
		<Button onClick={onLogout} withBg={false} textColor="light" withPadding={false}>
			<div class="flex flex-col items-center gap-1">
				<LogOut />
				<p class="duration-300 ease-in-out {menuOpen ? 'opacity-100' : 'opacity-0'}">Logout</p>
			</div>
		</Button>
	</nav>
	<ScrollArea.Root class="flex-1">
		<ScrollArea.Viewport class="h-full">
			{@render children()}
		</ScrollArea.Viewport>
	</ScrollArea.Root>
	<div class="block md:hidden">
		<Menubar>
			<MenuItem>
				<Link href="/manager/admin">
					<div class="flex flex-col items-center gap-1">
						<House />
						<p>Home</p>
					</div>
				</Link>
			</MenuItem>
			<MenuItem>
				<Link href="/manager/admin/mentors">
					<div class="flex flex-col items-center gap-1">
						<User />
						<p>Mentors</p>
					</div>
				</Link>
			</MenuItem>
			<MenuItem>
				<Link href="/manager/admin/maintenance">
					<div class="flex flex-col items-center gap-1">
						<Percent />
						<p>Maintenance</p>
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
</main>
