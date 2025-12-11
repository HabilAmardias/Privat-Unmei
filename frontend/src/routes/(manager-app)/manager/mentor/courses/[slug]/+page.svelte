<script lang="ts">
	import type { PageProps } from './$types';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import { Trash } from '@lucide/svelte';
	import { CourseDetailView } from './view.svelte';
	import AlertDialog from '$lib/components/dialog/AlertDialog.svelte';

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
	<div class="grid grid-cols-2 gap-4 text-center md:flex md:justify-between">
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Method:</p>
			<p>{data.detail.method}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Domicile:</p>
			<p>{data.detail.domicile}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Max Session:</p>
			<p>{data.detail.max_total_session}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Session Duration:</p>
			<p>{data.detail.session_duration_minutes}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Price:</p>
			<p>
				{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
					data.detail.price
				)}
			</p>
		</div>
	</div>
	<div class="flex flex-col">
		<p class="font-bold text-[var(--tertiary-color)]">Description:</p>
		<p>{data.detail.description}</p>
	</div>
	<div class="flex flex-col gap-4 md:grid md:grid-cols-2">
		<div class="flex flex-col gap-4">
			<p class="font-bold text-[var(--tertiary-color)]">Categories:</p>
			{#if data.courseCategories}
				<ScrollArea orientation="vertical" viewportClasses="max-h-[300px]">
					<ul class="flex flex-col gap-4">
						{#each data.courseCategories as cc, i (cc.id)}
							<li class="rounded-lg bg-[var(--tertiary-color)] p-4 text-[var(--secondary-color)]">
								<p>{cc.name}</p>
							</li>
						{/each}
					</ul>
				</ScrollArea>
			{:else}
				<b class="mx-auto self-center text-[var(--tertiary-color)]">No categories found</b>
			{/if}
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
</div>
