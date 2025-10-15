<script lang="ts">
	import { Button, Separator } from 'bits-ui';

	type NavMenu = {
		header: string;
		onClick: () => void;
	};

	type NavigationButtonsProps = {
		menus: NavMenu[];
	};

	let { menus }: NavigationButtonsProps = $props();
</script>

<div class="flex">
	{#each menus as menu, i}
		{#if i === 0}
			<Button.Root
				class="cursor-pointer rounded-tl-lg bg-[var(--tertiary-color)] p-2 text-[var(--secondary-color)] hover:text-[var(--primary-color)]"
				onclick={menu.onClick}>{menu.header}</Button.Root
			>
			<Separator.Root
				orientation="vertical"
				class="bg-border my-4 shrink-0 data-[orientation=horizontal]:h-px data-[orientation=vertical]:h-full data-[orientation=horizontal]:w-full data-[orientation=vertical]:w-[1px]"
			/>
		{:else if i === menus.length - 1}
			<Button.Root
				class="cursor-pointer rounded-tr-lg bg-[var(--tertiary-color)] p-2 text-[var(--secondary-color)] hover:text-[var(--primary-color)]"
				onclick={menu.onClick}>{menu.header}</Button.Root
			>
		{:else}
			<Button.Root
				class="cursor-pointer bg-[var(--tertiary-color)] p-2 text-[var(--secondary-color)] hover:text-[var(--primary-color)]"
				onclick={menu.onClick}>{menu.header}</Button.Root
			>
			<Separator.Root
				orientation="vertical"
				class="bg-border my-4 shrink-0 data-[orientation=horizontal]:h-px data-[orientation=vertical]:h-full data-[orientation=horizontal]:w-full data-[orientation=vertical]:w-[1px]"
			/>
		{/if}
	{/each}
</div>
