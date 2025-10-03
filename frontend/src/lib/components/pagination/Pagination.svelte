<script lang="ts">
	import { Pagination } from 'bits-ui';
	import { ChevronLeft, ChevronRight } from '@lucide/svelte';

	type paginationProps = {
		count: number;
		perPage: number;
		onPageChange?: (num: number) => void;
	};

	let { count, perPage, onPageChange }: paginationProps = $props();
</script>

<Pagination.Root {onPageChange} {count} {perPage}>
	{#snippet children({ pages })}
		<div class="my-4 flex items-center">
			<Pagination.PrevButton
				class="hover:bg-dark-10 disabled:text-muted-foreground mr-[25px] inline-flex size-10 cursor-pointer items-center justify-center rounded-[9px] bg-transparent hover:bg-[var(--tertiary-color)] hover:text-[var(--secondary-color)] active:scale-[0.98] disabled:cursor-not-allowed hover:disabled:bg-transparent hover:disabled:text-[var(--dark-color)]"
			>
				<ChevronLeft class="size-6" />
			</Pagination.PrevButton>
			<div class="flex items-center gap-2.5">
				{#each pages as page (page.key)}
					{#if page.type === 'ellipsis'}
						<div class="text-foreground-alt cursor-not-allowed select-none text-[15px] font-medium">
							...
						</div>
					{:else}
						<Pagination.Page
							{page}
							class="hover:bg-dark-10 data-selected:bg-[var(--tertiary-color)] data-selected:text-[var(--secondary-color)] data-selected:text-background inline-flex size-10 cursor-pointer select-none items-center justify-center rounded-[9px] bg-transparent text-[15px] font-medium active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-50 hover:disabled:bg-transparent"
						>
							{page.value}
						</Pagination.Page>
					{/if}
				{/each}
			</div>
			<Pagination.NextButton
				class="disabled:text-muted-foreground ml-[29px] inline-flex size-10 cursor-pointer items-center justify-center rounded-[9px] bg-transparent hover:bg-[var(--tertiary-color)] hover:text-[var(--secondary-color)] active:scale-[0.98] disabled:cursor-not-allowed disabled:text-[var(--dark-color)] hover:disabled:bg-transparent"
			>
				<ChevronRight class="size-6" />
			</Pagination.NextButton>
		</div>
	{/snippet}
</Pagination.Root>
