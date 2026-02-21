<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { UsersManagerView } from './view.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { enhance } from '$app/forms';
	import Loading from '$lib/components/loader/Loading.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import { ArrowDown, ArrowUp, UserPlus } from '@lucide/svelte';
	import Link from '$lib/components/button/Link.svelte';
	import AlertDialog from '$lib/components/dialog/AlertDialog.svelte';
	import NavigationButton from '$lib/components/button/NavigationButton.svelte';

	let { data }: PageProps = $props();
	const View = new UsersManagerView(data.mentorsList, data.studentList);

	onMount(() => {
		View.setIsDesktop(window.innerWidth >= 768);
		function setIsDesktop() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		window.addEventListener('resize', setIsDesktop);

		return () => {
			window.removeEventListener('resize', setIsDesktop);
		};
	});
</script>

{#snippet dialogStudentTitle()}
	Delete Student Confirmation
{/snippet}

{#snippet dialogMentorTitle()}
	Delete Mentor Confirmation
{/snippet}

{#snippet dialogDescription()}
	Irreversible action, are you sure want to proceed?
{/snippet}

<svelte:head>
	<title>Users - Privat Unmei</title>
	<meta name="description" content="Users - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-full flex-col p-4">
	<div class="mb-4">
		<NavigationButton
			menus={[
				{
					header: 'Mentors',
					onClick: () => (View.menuState = 'mentors')
				},
				{
					header: 'Students',
					onClick: () => (View.menuState = 'students')
				}
			]}
		/>
	</div>
	{#if View.menuState === "students"}
		<div class="flex h-full flex-col gap-4 p-4">
		<div class="flex justify-between gap-4">
			<h3 class="text-xl font-bold text-[var(--tertiary-color)]">Students</h3>
		</div>
		<form
			use:enhance={View.onSearchStudent}
			bind:this={View.studentsSearchForm}
			action="?/getStudents"
			method="POST"
		>
			<Input
				onInput={View.onStudentSearchInput}
				width="full"
				placeholder="Search"
				id="search"
				name="search"
				type="text"
			/>
		</form>
		<div class="flex flex-1">
			{#if View.studentIsLoading}
				<Loading />
			{:else if !View.students || View.students.length === 0}
				<b class="mx-auto self-center text-[var(--tertiary-color)]">No students found</b>
			{:else}
				<ScrollArea orientation="both" viewportClasses="w-[80dvw] max-w-[80dvw] h-[50dvh] max-h-[50dvh]">
					<table class="w-full border-separate border-spacing-4">
						<thead>
							<tr>
								<th class="text-[var(--tertiary-color)]">Public ID</th>
								<th class="text-[var(--tertiary-color)]">Name</th>
								<th class="text-[var(--tertiary-color)]">Status</th>
							</tr>
						</thead>
						<tbody>
							{#each View.students as student, i (student.id)}
								<tr>
									<td class="overflow-x-auto text-center">{student.public_id}</td>
									<td class="overflow-x-auto text-center">{student.name}</td>
									<td class="overflow-x-auto text-center">{student.status}</td>
									<td>
										<AlertDialog
											action="?/deleteStudent"
											bind:open={View.studentsAlertOpen[i]}
											enhancement={View.onDeleteStudent}
											title={dialogStudentTitle}
											onClick={() => {
												View.studentToDelete = student.id;
											}}
											description={dialogDescription}>Delete</AlertDialog
										>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</ScrollArea>
			{/if}
		</div>
		<form
			use:enhance={View.onPageChangeStudent}
			action="?/getStudents"
			class="flex w-full items-center justify-center"
			method="POST"
		>
			<Pagination bind:pageNumber={View.studentPage} perPage={View.studentLimit} count={View.studentTotalRow} />
		</form>
	</div>
	{:else}
		<div class="flex h-full flex-col gap-4 p-4">
	<div class="flex justify-between gap-4">
		<h3 class="text-xl font-bold text-[var(--tertiary-color)]">Mentors</h3>
		<Link href="/manager/admin/users/mentors/create" theme="light">
			<div class="flex gap-3 rounded-lg bg-[var(--tertiary-color)] p-2">
				<p>Create Mentor</p>
				<UserPlus /> 
			</div>
		</Link>
	</div>
	<form
		use:enhance={View.onSearchMentors}
		bind:this={View.mentorsSearchForm}
		class="grid grid-cols-2 gap-4"
		action="?/getMentors"
		method="POST"
	>
		<Input
			onInput={View.onMentorSearchInput}
			width="full"
			placeholder="Search"
			id="search"
			name="search"
			type="text"
		/>
		<Button onClick={() => View.onSort()} full type="submit">
			<div class="flex w-full justify-between px-1">
				<p>Sort</p>
				{#if View.sortByYears === true}
					<ArrowUp />
				{:else if View.sortByYears === false}
					<ArrowDown />
				{/if}
			</div>
		</Button>
	</form>
	<div class="flex flex-1">
		{#if View.mentorsIsLoading}
			<Loading />
		{:else if !View.mentors || View.mentors.length === 0}
			<b class="mx-auto self-center text-[var(--tertiary-color)]">No mentors found</b>
		{:else}
			<ScrollArea orientation="both" viewportClasses="w-[80dvw] max-w-[80dvw] h-[50dvh] max-h-[50dvh]">
				<table class="w-full border-separate border-spacing-4">
					<thead>
						<tr>
							<th class="text-[var(--tertiary-color)]">Public ID</th>
							<th class="text-[var(--tertiary-color)]">Name</th>
							<th class="text-[var(--tertiary-color)]">YoE</th>
						</tr>
					</thead>
					<tbody>
						{#each View.mentors as mentor, i (mentor.id)}
							<tr>
								<td class="overflow-x-auto text-center">{mentor.public_id}</td>
								<td class="overflow-x-auto text-center">{mentor.name}</td>
								<td class="overflow-x-auto text-center">{mentor.years_of_experience}</td>
								<td class="text-center">
									<Link theme="dark" href={`/manager/admin/users/mentors/${mentor.id}`}>Detail</Link>
								</td>
								<td>
									<AlertDialog
										action="?/deleteMentor"
										bind:open={View.mentorsAlertOpen[i]}
										enhancement={View.onDeleteMentor}
										title={dialogMentorTitle}
										onClick={() => {
											View.mentorToDelete = mentor.id;
										}}
										description={dialogDescription}>Delete</AlertDialog
									>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</ScrollArea>
		{/if}
	</div>
	<form
		use:enhance={View.onPageChangeMentors}
		action="?/getMentors"
		class="flex w-full items-center justify-center"
		method="POST"
	>
		<Pagination bind:pageNumber={View.mentorPage} perPage={View.mentorLimit} count={View.mentorTotalRow} />
	</form>
</div>
	{/if}
</div>




