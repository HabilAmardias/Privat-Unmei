<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$lib/components/button/Button.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import Menubar from '$lib/components/menubar/Menubar.svelte';
	import MenuItem from '$lib/components/menubar/MenuItem.svelte';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';
	import { House, Info, List, LogOut, MessageCircleMore, ShoppingCart, User } from '@lucide/svelte';
	import toast from 'svelte-french-toast';

	function onLogout() {
		const loadID = toast.loading('logging out....', { position: 'top-right' });
		loadingStore.setLogOutLoadID(loadID);
		goto('/logout', { replaceState: true });
	}

	let { children } = $props();
</script>

<main class="flex-1 pb-24 md:pb-0 md:pt-24">
	{@render children()}
</main>
<Menubar>
	<MenuItem>
		<Link href="/home">
			<div class="flex flex-col items-center justify-center">
				<House />
				Home
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
	<MenuItem>
		<Link href="/orders">
			<div class="flex flex-col items-center">
				<ShoppingCart />
				Orders
			</div>
		</Link>
	</MenuItem>
	<MenuItem>
		<Link href="/chats">
			<div class="flex flex-col items-center">
				<MessageCircleMore />
				Chats
			</div>
		</Link>
	</MenuItem>
	<MenuItem>
		<Link href="/about">
			<div class="flex flex-col items-center">
				<Info />
				About
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
</Menubar>
