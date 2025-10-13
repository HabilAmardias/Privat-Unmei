<script lang="ts">
	import { Button } from 'bits-ui';

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
				class="cursor-pointer rounded-tl-lg bg-[var(--tertiary-color)] p-1 text-[var(--secondary-color)] hover:text-[var(--primary-color)]"
				onclick={menu.onClick}>{menu.header}</Button.Root
			>
		{:else if i === menus.length - 1}
			<Button.Root
				class="cursor-pointer rounded-tr-lg bg-[var(--tertiary-color)] p-1 text-[var(--secondary-color)] hover:text-[var(--primary-color)]"
				onclick={menu.onClick}>{menu.header}</Button.Root
			>
		{:else}
			<Button.Root
				class="cursor-pointer bg-[var(--tertiary-color)] p-1 text-[var(--secondary-color)] hover:text-[var(--primary-color)]"
				onclick={menu.onClick}>{menu.header}</Button.Root
			>
		{/if}
	{/each}
</div>
