<script lang="ts">
	import { onMount } from 'svelte';
	import type { LayoutProps } from './$types';
	import { adminRole } from '$lib/utils/constants';
	import { goto } from '$app/navigation';
	import Link from '$lib/components/button/Link.svelte';
	import { ScrollArea } from 'bits-ui';
	import { House, Menu } from '@lucide/svelte';
	import Button from '$lib/components/button/Button.svelte';

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
</script>

<main class="flex">
	<nav
		class="
            flex h-dvh flex-col items-center gap-4 bg-[var(--tertiary-color)] p-4
            transition-all duration-300 ease-in-out
            {menuOpen ? 'w-[30%] md:w-[15%]' : 'w-20'} 
            "
	>
		<div class="mb-16 w-full justify-self-start">
			<Button onClick={handleMenu} full={true}>
				<Menu class="duration-300 ease-in-out {menuOpen ? 'rotate-90' : ''}" />
			</Button>
		</div>

		<Link href="/manager/admin">
			<div class="flex flex-col items-center gap-1">
				<House />
				<p class="duration-300 ease-in-out {menuOpen ? 'opacity-100' : 'opacity-0'}">Home</p>
			</div>
		</Link>

		<Link href="/manager/admin">
			<div class="flex flex-col items-center gap-1">
				<House />
				<p class="duration-300 ease-in-out {menuOpen ? 'opacity-100' : 'opacity-0'}">Home</p>
			</div>
		</Link>
	</nav>
	<ScrollArea.Root class="h-full flex-1">
		<ScrollArea.Viewport class="h-full flex-1">
			{@render children()}
		</ScrollArea.Viewport>
	</ScrollArea.Root>
</main>
