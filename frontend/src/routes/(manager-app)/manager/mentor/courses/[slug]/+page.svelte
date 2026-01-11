<script lang="ts">
	import type { PageProps } from './$types';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import { Trash } from '@lucide/svelte';
	import { CourseDetailView } from './view.svelte';
	import AlertDialog from '$lib/components/dialog/AlertDialog.svelte';
	import NavigationButton from '$lib/components/button/NavigationButton.svelte';

	let { data }: PageProps = $props();
	const View = new CourseDetailView();
</script>

<svelte:head>
	<title>{data.detail.title} - Privat Unmei</title>
	<meta name="description" content="Profile - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

{#snippet deleteDialogTitle()}
	Delete Course Confirmation
{/snippet}

{#snippet deleteDialogDescription()}
	Irreversible action, are you sure want to proceed?
{/snippet}

<div class="flex flex-col gap-4 p-4">
	<div class="flex items-center justify-between gap-4">
		<h1 class="text-xl font-bold text-[var(--tertiary-color)]">{data.detail.title}</h1>
		<div class="flex gap-4">
			<div class="h-fit w-fit rounded-lg bg-[var(--tertiary-color)] p-2">
				<Link href={`/manager/mentor/courses/${data.detail.id}/update`}>Update</Link>
			</div>
			<AlertDialog
				action="?/deleteCourse"
				bind:open={View.deleteDialogOpen}
				enhancement={View.onDeleteCourse}
				title={deleteDialogTitle}
				description={deleteDialogDescription}><Trash /></AlertDialog
			>
		</div>
	</div>
	<div>
		<NavigationButton
			menus={[
				{
					header: 'Description',
					onClick: () => (View.detailState = 'description')
				},
				{
					header: 'Detail',
					onClick: () => (View.detailState = 'detail')
				}
			]}
		/>
		<div
			class="flex h-[210px] flex-col gap-2 rounded-lg rounded-tl-none bg-[var(--tertiary-color)] p-4"
		>
			{#if View.detailState === 'detail'}
				<div class="flex gap-2">
					<p class="font-bold text-[var(--primary-color)]">Method:</p>
					<p class="text-[var(--secondary-color)]">
						{View.capitalizeFirstLetter(data.detail.method)}
					</p>
				</div>
				<div class="flex gap-2">
					<p class="font-bold text-[var(--primary-color)]">Domicile:</p>
					<p class="text-[var(--secondary-color)]">{data.detail.domicile}</p>
				</div>
				<div class="flex items-center gap-2">
					<p class="font-bold text-[var(--primary-color)]">Per Session Duration (minutes):</p>
					<p class="text-[var(--secondary-color)]">{data.detail.session_duration_minutes}</p>
				</div>
				<div class="flex gap-2">
					<p class="font-bold text-[var(--primary-color)]">Max Session Number:</p>
					<p class="text-[var(--secondary-color)]">{data.detail.max_total_session}</p>
				</div>
				<div class="flex gap-2">
					<p class="font-bold text-[var(--primary-color)]">Price:</p>
					<p class="text-[var(--secondary-color)]">
						{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
							data.detail.price
						)}
					</p>
				</div>
				<div class="flex gap-2">
					<p class="font-bold text-[var(--primary-color)]">Categories:</p>
					<ScrollArea orientation="horizontal" viewportClasses="max-w-[200px]">
						<ul class="flex gap-4">
							{#each data.courseCategories as cc, i (cc.id)}
								<li
									class="w-fit rounded-lg bg-[var(--tertiary-color)] text-[var(--secondary-color)]"
								>
									<p>{cc.name}</p>
								</li>
							{/each}
						</ul>
					</ScrollArea>
				</div>
			{:else}
				<ScrollArea orientation="vertical" viewportClasses="h-[150px] max-h-[150px]">
					<p class="text-justify text-[var(--secondary-color)]">{data.detail.description}</p>
				</ScrollArea>
			{/if}
		</div>
	</div>
	<div class="flex flex-col gap-4">
		<p class="font-bold text-[var(--tertiary-color)]">Topics:</p>
		{#if data.topics}
			<ScrollArea orientation="vertical" viewportClasses="max-h-[300px]">
				<ul class="flex flex-col gap-4">
					{#each data.topics as t, i (i)}
						<li class="rounded-lg bg-[var(--tertiary-color)] p-4 text-[var(--secondary-color)]">
							<p class="font-bold text-[var(--primary-color)]">{t.title}</p>
							<p>{t.description}</p>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		{:else}
			<b class="text-[var(--tertiary-color)]">No topic found</b>
		{/if}
	</div>
</div>
