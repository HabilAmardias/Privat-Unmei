<script lang="ts">
	import type { PageProps } from './$types';
	import Select from '$lib/components/select/Select.svelte';
	import { enhance } from '$app/forms';
	import Input from '$lib/components/form/Input.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { CourseListView } from './view.svelte';
	import { methodOpts } from './constants';
	import Search from '$lib/components/search/Search.svelte';
	import Loading from '$lib/components/loader/Loading.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import Pagination from '$lib/components/pagination/Pagination.svelte';

	let { data }: PageProps = $props();
	const View = new CourseListView(data.courses);
</script>

<svelte:head>
	<title>Courses - Privat Unmei</title>
	<meta name="description" content="Courses - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>
<form
	use:enhance={View.onGetCategory}
	bind:this={View.categoryForm}
	action="?/getCategories"
	method="POST"
></form>
<div class="flex flex-col gap-4 p-4">
	<h1 class="text-2xl font-bold text-[var(--tertiary-color)]">Courses</h1>
	<form
		use:enhance={View.onSearchCourse}
		class="flex flex-col gap-4 md:flex-row"
		action="?/getCourses"
		method="POST"
	>
		<div class="grid grid-cols-2 gap-4">
			<Input
				bind:value={View.search}
				width="full"
				placeholder="Search"
				id="search"
				name="search"
				type="text"
			/>
			<Select defaultLable="Method" name="method" options={methodOpts} bind:value={View.method} />
		</div>
		<Search
			bind:value={View.selectedCategory}
			items={View.categories}
			label="Category"
			keyword={View.searchCategory}
			onKeywordChange={View.onSearchCategory}
		/>
		<Button disabled={View.isLoading} type="submit" full>Search</Button>
	</form>
	<div class="flex-1">
		{#if View.isLoading}
			<Loading />
		{:else if !View.courses || View.courses.length === 0}
			<b class="mx-auto self-center text-[var(--tertiary-color)]">No courses found</b>
		{:else}
			<ScrollArea orientation="vertical" viewportClasses="h-[400px] max-h-[400px]">
				<ul class="flex flex-col gap-4 md:grid md:grid-cols-3">
					{#each View.courses as c (c.id)}
						<li class="transition-transform hover:-translate-y-1">
							<Link href={`/courses/${c.id}`}>
								<div
									class="flex h-[100px] flex-col justify-between rounded-lg bg-[var(--tertiary-color)] p-4"
								>
									<div>
										<p class="font-bold text-[var(--primary-color)]">{c.title}</p>
										<p class="text-[var(--secondary-color)]">{c.mentor_name}</p>
									</div>
									<div>
										<p class="text-[var(--secondary-color)]">
											{new Intl.NumberFormat('id-ID', {
												style: 'currency',
												currency: 'IDR'
											}).format(c.price)}
										</p>
									</div>
								</div>
							</Link>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		{/if}
	</div>
	<form
		use:enhance={View.onPageChangeEnhance}
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
