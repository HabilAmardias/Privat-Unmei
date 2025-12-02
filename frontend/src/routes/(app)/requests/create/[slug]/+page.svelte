<script lang="ts">
	import type { PageProps } from './$types';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import { CreateRequestView } from './view.svelte';
	import { enhance } from '$app/forms';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { dowMap } from './constants';
	import Input from '$lib/components/form/Input.svelte';
	import Select from '$lib/components/select/Select.svelte';
	import Datepicker from '$lib/components/calendar/Datepicker.svelte';
	import Button from '$lib/components/button/Button.svelte';

	let { data }: PageProps = $props();
	const View = new CreateRequestView(data.detail, data.payments);
</script>

<svelte:head>
	<title>{data.detail.title} - Privat Unmei</title>
	<meta name="description" content="Profile - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex flex-col gap-4 p-4">
	<div class="flex flex-col gap-2 md:flex-row md:justify-between">
		<h1 class="text-2xl font-bold text-[var(--tertiary-color)]">{data.detail.title}</h1>
		<div class="w-fit rounded-lg bg-[var(--tertiary-color)] p-2">
			<p class="font-bold text-[var(--secondary-color)]">
				{new Intl.NumberFormat('id-ID', { currency: 'IDR', style: 'currency' }).format(
					data.detail.price
				)} / session
			</p>
		</div>
	</div>
	<div class="grid grid-cols-2">
		<div class="flex flex-col gap-2">
			<div class="flex gap-2">
				<p class="font-bold text-[var(--tertiary-color)]">Method:</p>
				<p>{View.capitalizeFirstLetter(data.detail.method)}</p>
			</div>
			<div class="flex gap-2">
				<p class="font-bold text-[var(--tertiary-color)]">Domicile:</p>
				<p>{data.detail.domicile}</p>
			</div>
			<div class="flex items-center gap-2">
				<p class="font-bold text-[var(--tertiary-color)]">Per Session Duration (minutes):</p>
				<p>{data.detail.session_duration_minutes}</p>
			</div>
		</div>
	</div>
	<div class="flex flex-col">
		<p class="font-bold text-[var(--tertiary-color)]">Description</p>
		<p>{data.detail.description}</p>
	</div>
	<div>
		<p class="font-bold text-[var(--tertiary-color)]">Mentor</p>
		<div
			class="hover:-translate-y-0.25 flex transform items-center gap-4 rounded-md bg-[var(--tertiary-color)] p-2 transition-transform"
		>
			<CldImage
				src={data.detail.mentor_profile_image}
				width={50}
				height={50}
				className="rounded-full"
			/>
			<div>
				<p class="font-bold text-[var(--primary-color)]">{data.detail.mentor_name}</p>
				<p class="text-[var(--secondary-color)]">{data.detail.mentor_email}</p>
			</div>
		</div>
	</div>
	<div class="flex flex-col gap-4">
		<p class="font-bold text-[var(--tertiary-color)]">Schedules:</p>
		{#if data.schedules}
			<ScrollArea orientation="vertical" viewportClasses="max-h-[300px]">
				<ul class="flex flex-col gap-4">
					{#each data.schedules as sch, i (i)}
						<li class="rounded-lg bg-[var(--tertiary-color)] p-4 text-[var(--secondary-color)]">
							<p>{dowMap.get(sch.day_of_week)}, {sch.start_time} - {sch.end_time}</p>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		{:else}
			<b class="mx-auto self-center text-[var(--tertiary-color)]">No schedules found</b>
		{/if}
	</div>
	<h2 class="text-2xl font-bold text-[var(--tertiary-color)]">Create Request</h2>
	<form
		class="flex flex-col gap-4"
		use:enhance={View.onCreateRequest}
		action="?/createRequest"
		method="POST"
	>
		<div class="grid grid-cols-2 gap-4">
			<div>
				<p class="font-bold text-[var(--tertiary-color)]">Participant</p>
				<Input
					bind:value={View.participant}
					type="number"
					min={1}
					name="participant"
					id="participant"
				/>
			</div>
			<div>
				<p class="font-bold text-[var(--tertiary-color)]">Payment Method</p>
				<Select
					bind:value={View.selectedPayment}
					options={View.paymentOpts}
					name="payment"
					defaultLable="Please choose payment method"
				/>
			</div>
		</div>
		<div class="flex gap-4">
			<Datepicker dows={data.dows} onChange={View.onCalendarValueChange} label="Pick a date" />
			<div class="flex flex-col gap-4">
				<Input type="time" bind:value={View.selectedStartTime} name="start_time" id="start_time" />
				<Button type="button" disabled={View.disableAddSchedule}>Add</Button>
			</div>
		</div>
		<Button disabled={View.disableSubmit} type="submit">Create Request</Button>
	</form>
</div>
