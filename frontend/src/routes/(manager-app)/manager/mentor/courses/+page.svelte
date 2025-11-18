<script lang="ts">
	import type { PageProps } from './$types';
	import { CourseManagementView } from './view.svelte';
	import AlertDialog from '$lib/components/dialog/AlertDialog.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Loading from '$lib/components/loader/Loading.svelte';
	import { enhance } from '$app/forms';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import Search from '$lib/components/search/Search.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import Button from '$lib/components/button/Button.svelte';

	let { data }: PageProps = $props();
	const View = new CourseManagementView(data.categories, data.courses);
</script>

<svelte:head>
	<title>Course Management - Privat Unmei</title>
	<meta name="description" content="Course Management - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

{#snippet deleteDialogTitle()}
	Delete Course Confirmation
{/snippet}

{#snippet deleteDialogDescription()}
	Irreversible action, are you sure want to proceed?
{/snippet}

<div class="flex h-full flex-col p-4">
	<div class="mb-4 flex items-center justify-between">
		<h3 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">Courses</h3>
		<div class="h-fit w-fit rounded-lg bg-[var(--tertiary-color)] p-2">
			<Link href="/manager/mentor/courses/create">Create Course</Link>
		</div>
	</div>
	<form
		bind:this={View.searchCategoryForm}
		use:enhance={View.onSearchCategory}
		action="?/getCourseCategories"
		method="POST"
	></form>
	<form
		use:enhance={View.onSearchCourse}
		action="?/getMyCourses"
		class="mb-4 flex gap-4"
		method="POST"
	>
		<Input bind:value={View.search} type="text" name="search" id="search" placeholder="Search" />
		<Search
			bind:value={View.selectedCategory}
			items={View.categories}
			label="Categories"
			keyword={View.searchCategory}
			onKeywordChange={View.onSearchCategoryChange}
		/>
		<Button type="submit">Search</Button>
	</form>
	<div class="flex flex-1">
		{#if View.isLoading}
			<Loading />
		{:else if !View.courses || View.courses.length === 0}
			<b class="mx-auto self-center text-[var(--tertiary-color)]">No courses found</b>
		{:else}
			<ScrollArea orientation="vertical" class="flex-1" viewportClasses="max-h-[500px]">
				<table class="w-full table-fixed border-separate border-spacing-4">
					<thead>
						<tr>
							<th class="text-[var(--tertiary-color)]">Name</th>
							<th class="text-[var(--tertiary-color)]">Price</th>
							<th class="text-[var(--tertiary-color)]">Method</th>
						</tr>
					</thead>
					<tbody>
						{#each View.courses as c, i (c.id)}
							<tr>
								<td class="text-center">
									{c.title}
								</td>
								<td class="text-center">
									{c.price}
								</td>
								<td class="text-center">
									{c.method}
								</td>
								<td>
									<div class="flex items-center justify-center">
										<div class="h-fit w-fit rounded-lg bg-[var(--tertiary-color)] p-2">
											<Link href={`/manager/mentor/courses/${c.id}`}>Detail</Link>
										</div>
									</div>
								</td>
								<td>
									<div class="flex items-center justify-center">
										<AlertDialog
											action="?/deleteCourse"
											bind:open={View.deleteCourseOpenDialogs[i]}
											enhancement={View.onDeleteCourse}
											title={deleteDialogTitle}
											onClick={() => {
												View.courseToDelete = c.id;
											}}
											description={deleteDialogDescription}>Delete</AlertDialog
										>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</ScrollArea>
		{/if}
	</div>
	<form
		use:enhance={View.onPageChangeForm}
		action="?/getMyCourses"
		class="flex w-full items-center justify-center"
		method="POST"
		bind:this={View.paginationForm}
	>
		<Pagination
			onPageChange={View.onPageChange}
			bind:pageNumber={View.page}
			perPage={View.limit}
			count={View.totalRow}
		/>
	</form>
</div>
