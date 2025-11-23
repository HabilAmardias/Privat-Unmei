<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Select from '$lib/components/select/Select.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import { dayofWeeks, degreeOpts } from './constants';
	import { CreateMentorView } from './view.svelte';
	import type { PageProps } from './$types';
	import Search from '$lib/components/search/Search.svelte';
	import { enhance } from '$app/forms';
	import FileInput from '$lib/components/form/FileInput.svelte';
	import { CloudUpload, X } from '@lucide/svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';

	const { data }: PageProps = $props();
	const View = new CreateMentorView(data.paymentMethods, data.generatedPassword);
</script>

<svelte:head>
	<title>Add Mentor - Privat Unmei</title>
	<meta name="description" content="Add Mentor - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex flex-col p-4">
	<h3 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">Create New Mentor</h3>
	<form
		bind:this={View.paymentMethodForm}
		use:enhance={View.onGetPaymentMethods}
		action="?/getPaymentMethods"
		method="POST"
	></form>
	<form
		action="?/generatePassword"
		method="POST"
		bind:this={View.generatePasswordForm}
		use:enhance={View.onGetPassword}
	></form>
	<form
		use:enhance={View.onCreateMentor}
		action="?/createMentor"
		method="POST"
		class="flex flex-col gap-4"
		enctype="multipart/form-data"
	>
		<div class="grid grid-cols-[30%_70%] md:grid-cols-[10%_90%]">
			<p class="font-bold text-[var(--tertiary-color)]">Name:</p>
			<Input
				bind:value={View.name}
				type="text"
				placeholder="Input mentor name"
				name="name"
				id="name"
			/>
		</div>
		<div class="grid grid-cols-[30%_70%] md:grid-cols-[10%_90%]">
			<p class="font-bold text-[var(--tertiary-color)]">Email:</p>
			<Input
				err={View.emailErr}
				bind:value={View.email}
				type="email"
				placeholder="Input mentor email"
				name="email"
				id="email"
			/>
		</div>
		<div class="grid grid-cols-[30%_40%_30%] md:grid-cols-[10%_80%_10%]">
			<p class="font-bold text-[var(--tertiary-color)]">Password:</p>
			<div
				class="flex flex-1 gap-1 rounded-lg bg-[var(--tertiary-color)] px-4 py-2 text-[var(--secondary-color)]"
			>
				<p>{View.generatedPassword}</p>
			</div>
			<div class="flex justify-end">
				<Button type="button" onClick={View.generatePassword}>Generate</Button>
			</div>
		</div>
		<div class="grid grid-cols-[30%_70%] md:grid-cols-[10%_90%]">
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
				<Link href="/manager/admin">Cancel</Link>
			</div>
			<Button disabled={View.disableCreateMentor} type="submit">Create</Button>
		</div>
	</form>
</div>
