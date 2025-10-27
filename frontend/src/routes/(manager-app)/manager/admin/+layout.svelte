<script lang="ts">
	import { onMount } from 'svelte';
	import type { LayoutProps } from './$types';
	import { adminRole } from '$lib/utils/constants';
	import { goto } from '$app/navigation';
	import Link from '$lib/components/button/Link.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import { Banknote, House, LogOut, Menu, Percent, User } from '@lucide/svelte';
	import Button from '$lib/components/button/Button.svelte';
	import Menubar from '$lib/components/menubar/Menubar.svelte';
	import MenuItem from '$lib/components/menubar/MenuItem.svelte';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';
	import { adminLayoutView } from './view.svelte';
	import { CreateToast } from '$lib/utils/helper';

	const { children, data }: LayoutProps = $props();
	const View = new adminLayoutView();

	onMount(() => {
		if (data.role !== adminRole) {
			goto('/courses', { replaceState: true });
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
		<Link href="/manager/admin">
			<div class="flex flex-col items-center gap-1">
				<House />
				<p class="duration-300 ease-in-out {View.menuOpen ? 'opacity-100' : 'opacity-0'}">Home</p>
			</div>
		</Link>
		<Link href="/manager/admin/mentors">
			<div class="flex flex-col items-center gap-1">
				<User />
				<p class="duration-300 ease-in-out {View.menuOpen ? 'opacity-100' : 'opacity-0'}">
					Mentors
				</p>
			</div>
		</Link>
		<Link href="/manager/admin/costs">
			<div class="flex flex-col items-center gap-1">
				<Percent />
				<p class="duration-300 ease-in-out {View.menuOpen ? 'opacity-100' : 'opacity-0'}">Costs</p>
			</div>
		</Link>
		<Link href="/manager/admin/payment-methods">
			<div class="flex flex-col items-center gap-1">
				<Banknote />
				<p class="duration-300 ease-in-out {View.menuOpen ? 'opacity-100' : 'opacity-0'}">
					Payments
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
		<ScrollArea class="flex-1" orientation="vertical" viewportClasses="h-full max-h-full">
			{@render children()}
		</ScrollArea>
	</main>
</div>
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
			<Link href="/manager/admin/costs">
				<div class="flex flex-col items-center gap-1">
					<Percent />
					<p>Costs</p>
				</div>
			</Link>
		</MenuItem>
		<MenuItem>
			<Link href="/manager/admin/payment">
				<div class="flex flex-col items-center gap-1">
					<Banknote />
					<p>Payments</p>
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
