<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Select from '$lib/components/select/Select.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import { UpdateCourseView } from './view.svelte';
	import type { PageProps } from './$types';
	import Search from '$lib/components/search/Search.svelte';
	import { enhance } from '$app/forms';
	import { X } from '@lucide/svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import Textarea from '$lib/components/form/Textarea.svelte';
	import { methodOpts } from './constants';

	const { data }: PageProps = $props();
	const View = new UpdateCourseView(
		data.categories,
		data.topics,
		data.courseCategories,
		data.detail
	);
</script>

<svelte:head>
	<title>Update Course - Privat Unmei</title>
	<meta name="description" content="Update Course - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex flex-col p-4">
	<h3 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">Update Course</h3>
	<form
		use:enhance={View.onGetCategory}
		method="POST"
		action="?/getCategories"
		bind:this={View.searchCategoryForm}
	></form>
	<form
		use:enhance={View.onUpdateCourse}
		action="?/updateCourse"
		method="POST"
		class="flex flex-col gap-4"
	>
		<div class="grid grid-cols-[30%_70%] md:grid-cols-[10%_90%]">
			<p class="font-bold text-[var(--tertiary-color)]">Title:</p>
			<Input
				bind:value={View.title}
				type="text"
				placeholder="Input course title"
				name="title"
				id="title"
			/>
		</div>
		<Textarea
			err={View.descriptionErr}
			bind:value={View.description}
			placeholder="Input course description"
			id="description"
			name="description"
			><p class="font-bold text-[var(--tertiary-color)]">Description:</p>
		</Textarea>
		<div class="grid grid-cols-[30%_70%] md:grid-cols-[10%_90%]">
			<p class="font-bold text-[var(--tertiary-color)]">Domicile:</p>
			<Input
				bind:value={View.domicile}
				type="text"
				placeholder="Input domicile"
				name="domicile"
				id="domicile"
			/>
		</div>
		<div class="grid grid-cols-[30%_70%] md:grid-cols-[10%_90%]">
			<p class="font-bold text-[var(--tertiary-color)]">Price:</p>
			<Input
				err={View.priceErr}
				bind:value={View.price}
				type="number"
				placeholder="Input course price"
				name="price"
				id="price"
				min={1}
			/>
		</div>
		<div class="grid grid-cols-[30%_70%] md:grid-cols-[10%_90%]">
			<p class="font-bold text-[var(--tertiary-color)]">Method:</p>
			<Select
				bind:value={View.method}
				options={methodOpts}
				defaultLable="Choose Method"
				name="method"
			/>
		</div>

		<div class="grid grid-cols-[30%_70%] md:grid-cols-[10%_90%]">
			<p class="font-bold text-[var(--tertiary-color)]">Session Duration (minutes):</p>
			<Input
				bind:value={View.sessionDuration}
				err={View.sessionDurationErr}
				type="number"
				placeholder="Input per session duration"
				name="session_duration"
				id="session_duration"
				min={1}
			/>
		</div>
		<div class="grid grid-cols-[30%_70%] md:grid-cols-[10%_90%]">
			<p class="font-bold text-[var(--tertiary-color)]">Max Session:</p>
			<Input
				err={View.maxSessionErr}
				bind:value={View.maxSession}
				type="number"
				placeholder="Input maximum number of session"
				name="max_session"
				id="max_session"
				min={1}
			/>
		</div>
		<p class="font-bold text-[var(--tertiary-color)]">Topics:</p>
		<div class="grid grid-cols-2 place-items-center gap-4">
			<div class="flex w-full flex-col gap-4">
				<div class="flex flex-col gap-4">
					<div class="flex flex-col">
						<p class="font-bold text-[var(--tertiary-color)]">Topic Title:</p>
						<Input
							bind:value={View.topicTitle}
							type="text"
							name="title"
							id="title"
							placeholder="Insert topic title"
						/>
					</div>

					<Textarea
						bind:value={View.topicDescription}
						placeholder="Input topic description"
						id="topic_description"
						name="topic_description"
						><p class="font-bold text-[var(--tertiary-color)]">Topic Description:</p>
					</Textarea>
				</div>
				<Button disabled={View.disableAddTopic} full type="button" onClick={View.addCourseTopic}
					>Add</Button
				>
			</div>
			<ScrollArea class="w-full" orientation="vertical" viewportClasses="max-h-[300px]">
				<ul>
					{#each View.addedTopic as t, i}
						<li class="flex justify-between">
							<p class="text-[var(--tertiary-color)]">
								{`${t.title} - ${t.description}`}
							</p>
							<Button
								type="button"
								withBg={false}
								withPadding={false}
								textColor="dark"
								onClick={() => {
									View.removeAddedTopic(i);
								}}
							>
								<X />
							</Button>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		</div>
		<p class="font-bold text-[var(--tertiary-color)]">Course Category:</p>
		<div class="grid grid-cols-2 place-items-center gap-4">
			<div class="flex w-full flex-col gap-4">
				<div class="flex gap-4">
					<Search
						bind:value={View.selectedCategory}
						items={View.categories}
						label="Course Category"
						keyword={View.searchCategory}
						onKeywordChange={View.onSearchCategory}
					/>
				</div>
				{#if View.addedCategoryErr}
					<p class="font-bold text-[var(--tertiary-color)]">{View.addedCategoryErr.message}</p>
				{/if}
				<Button
					disabled={View.disableAddCategory}
					full
					type="button"
					onClick={View.addCourseCategory}>Add</Button
				>
			</div>
			<ScrollArea class="w-full" orientation="vertical" viewportClasses="max-h-[300px]">
				<ul>
					{#each View.addedCategories as c, i}
						<li class="flex justify-between">
							<p class="text-[var(--tertiary-color)]">
								{c.name}
							</p>
							<Button
								type="button"
								withBg={false}
								withPadding={false}
								textColor="dark"
								onClick={() => {
									View.removeAddedCategories(i);
								}}
							>
								<X />
							</Button>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		</div>
		<div class="flex gap-4">
			<div class="w-fit rounded-lg bg-[var(--tertiary-color)] p-3">
				<Link href="/manager/mentor/courses">Cancel</Link>
			</div>
			<Button disabled={View.disableUpdateCourse} type="submit">Update</Button>
		</div>
	</form>
</div>
