<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Select from '$lib/components/select/Select.svelte';
	import { dayofWeeks, degreeOpts } from './constants';
	import { UpdateMentorProfileView } from './view.svelte';
	import type { PageProps } from './$types';
	import Search from '$lib/components/search/Search.svelte';
	import { enhance } from '$app/forms';
	import FileInput from '$lib/components/form/FileInput.svelte';
	import { CloudUpload, X } from '@lucide/svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import Textarea from '$lib/components/form/Textarea.svelte';
	import Image from '$lib/components/image/Image.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { Pencil } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import Link from '$lib/components/button/Link.svelte';

	const { data }: PageProps = $props();
	const View = new UpdateMentorProfileView(
		data.paymentMethods,
		data.mentorSchedules,
		data.mentorPayments,
		data.profile
	);

	onMount(() => {
		View.isDesktop = window.innerWidth >= 768;
		function setIsDesktop() {
			View.isDesktop = window.innerWidth >= 768;
		}
		window.addEventListener('resize', setIsDesktop);

		return () => {
			window.removeEventListener('resize', setIsDesktop);
		};
	});
</script>

<svelte:head>
	<title>Update Profile - Privat Unmei</title>
	<meta name="description" content="Update Profile - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex flex-col p-4">
	<h3 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">Update Profile</h3>

	<form
		bind:this={View.paymentMethodForm}
		use:enhance={View.onGetPaymentMethods}
		action="?/getPaymentMethods"
		method="POST"
	></form>
	<form
		use:enhance={View.onUpdateMentor}
		action="?/updateProfile"
		method="POST"
		class="flex flex-col gap-4"
		enctype="multipart/form-data"
	>
		<div class="flex gap-4">
			<FileInput
				accept="image/png"
				bind:files={View.profileImage}
				id="profile_image"
				name="profile_image"
			>
				<div class="group relative inline-block overflow-hidden rounded-full">
					{#if View.profileImage}
						<Image
							src={URL.createObjectURL(View.profileImage[0])}
							width={View.size}
							height={View.size}
							className="rounded-full shadow-2xl border-gray-400 brightness-60 md:brightness-100 md:border-none md:shadow-none md:hover:shadow-2xl md:group-hover:border-gray-400 md:transition-all md:duration-300 md:group-hover:brightness-60"
						/>
					{:else}
						<CldImage
							src={data.profile.profile_image}
							width={View.size}
							height={View.size}
							className="rounded-full shadow-2xl border-gray-400 brightness-60 md:brightness-100 md:border-none md:shadow-none md:hover:shadow-2xl md:group-hover:border-gray-400 md:transition-all md:duration-300 md:group-hover:brightness-60"
						/>
					{/if}
					<div
						class="absolute inset-0 flex items-center justify-center bg-opacity-0 transition-all duration-300 group-hover:bg-opacity-30"
					>
						<Pencil
							class="text-white md:scale-50 md:transform md:opacity-0 md:transition-all md:duration-300 md:group-hover:scale-100 md:group-hover:opacity-100"
						/>
					</div>
				</div>
			</FileInput>
			<div class="flex flex-col md:flex-row md:items-center md:gap-4">
				<p class="font-bold text-[var(--tertiary-color)]">Name:</p>
				<Input
					bind:value={View.name}
					type="text"
					placeholder="Input mentor name"
					name="name"
					id="name"
				/>
			</div>
		</div>
		<div class="flex flex-col gap-4">
			<Textarea bind:value={View.bio} name="bio" id="bio" placeholder="Insert new bio"
				><p class="font-bold text-[var(--tertiary-color)]">Bio:</p></Textarea
			>
		</div>
		<div class="flex items-center gap-4">
			<p class="font-bold text-[var(--tertiary-color)]">Experience:</p>
			<Input
				type="number"
				placeholder="Years of Experience"
				name="years_of_experience"
				id="years_of_experience"
				bind:value={View.yearsOfExperience}
				min={0}
				max={40}
				err={View.yoeErr}
			/>
		</div>

		<p class="font-bold text-[var(--tertiary-color)]">Education:</p>
		<Input bind:value={View.campus} type="text" name="campus" id="campus" placeholder="Campus" />
		<div class="flex gap-4">
			<Select options={degreeOpts} defaultLable="Degree" name="degree" bind:value={View.degree} />
			<Input bind:value={View.major} type="text" name="major" id="major" placeholder="Major" />
		</div>
		<p class="font-bold text-[var(--tertiary-color)]">Resume:</p>
		{#if View.resumeErr}
			<p class="text-red-500">{View.resumeErr.message}</p>
		{/if}
		<FileInput
			bind:files={View.resumeFile}
			accept="application/pdf"
			name="resume_file"
			id="resume_file"
		>
			<div
				class="border-1 flex w-full flex-col items-center justify-center rounded-lg border-dashed border-[var(--tertiary-color)] p-3 font-bold text-[var(--tertiary-color)] hover:text-[var(--primary-color)]"
			>
				<CloudUpload />
				{View.resumeFile ? View.resumeFile[0].name : 'Upload Resume'}
			</div>
		</FileInput>
		<p class="font-bold text-[var(--tertiary-color)]">Schedules:</p>
		<div class="grid grid-cols-2 place-items-center gap-4">
			<div class="flex w-full flex-col gap-4">
				<div class="flex gap-4">
					<Select
						options={dayofWeeks}
						defaultLable="Weekday"
						name="day_of_week"
						bind:value={View.selectedDayOfWeek}
					/>
					<Input bind:value={View.selectedStartTime} step={1} type="time" name="start" id="start" />
					<Input bind:value={View.selectedEndTime} step={1} type="time" name="end" id="end" />
				</div>
				{#if View.selectMentorScheduleErr}
					<p class="text-red-500">{View.selectMentorScheduleErr.message}</p>
				{/if}
				<Button
					disabled={View.disableAddMentorSchedule}
					full
					type="button"
					onClick={View.addMentorSchedule}>Add Mentor Schedule</Button
				>
			</div>
			<ScrollArea class="w-full" orientation="vertical" viewportClasses="max-h-[100px]">
				<ul>
					{#each View.mentorSchedules as sch, i}
						<li class="flex justify-between">
							<p class="text-[var(--tertiary-color)]">
								{`${sch.day_of_week_label}, ${View.TimeOnlyToString(sch.start_time)}-${View.TimeOnlyToString(sch.end_time)}`}
							</p>
							<Button
								type="button"
								withBg={false}
								withPadding={false}
								textColor="dark"
								onClick={() => {
									View.removeMentorSchedule(i);
								}}
							>
								<X />
							</Button>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		</div>
		<p class="font-bold text-[var(--tertiary-color)]">Payment Methods:</p>
		<div class="grid grid-cols-2 place-items-center gap-4">
			<div class="flex w-full flex-col gap-4">
				<div class="flex flex-col gap-4 md:flex-row">
					<Search
						bind:value={View.selectedPaymentMethod}
						items={View.paymentMethods}
						label="Payment Method"
						keyword={View.searchValue}
						onKeywordChange={View.onSearchPaymentMethodChange}
					/>
					<Input
						type="text"
						name="account_number"
						id="account_number"
						placeholder="Account Number"
						bind:value={View.accountNumber}
					/>
				</div>
				{#if View.selectPaymentMethodErr}
					<p class="text-red-500">{View.selectPaymentMethodErr.message}</p>
				{/if}
				<Button
					disabled={View.disableAddPaymentMethod}
					full
					type="button"
					onClick={View.addMentorPaymentMethod}>Add Payment Method</Button
				>
			</div>
			<ScrollArea class="w-full" orientation="vertical" viewportClasses="max-h-[100px]">
				<ul>
					{#each View.mentorPaymentMethods as pym, i}
						<li class="flex justify-between">
							<p class="text-[var(--tertiary-color)]">
								{`${pym.payment_method_name} - ${pym.account_number}`}
							</p>
							<Button
								type="button"
								withBg={false}
								withPadding={false}
								textColor="dark"
								onClick={() => {
									View.removeMentorPaymentMethod(i);
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
				<Link href="/manager/mentor">Cancel</Link>
			</div>
			<Button disabled={View.disableUpdateMentor} type="submit">Update</Button>
		</div>
	</form>
</div>
