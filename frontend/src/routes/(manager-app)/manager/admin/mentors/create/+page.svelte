<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Select from '$lib/components/select/Select.svelte';
	import { onMount } from 'svelte';
	import { dayofWeeks, degreeOpts } from './constants';
	import { CreateMentorView } from './view.svelte';
	import type { PageProps } from './$types';
	import Search from '$lib/components/search/Search.svelte';
	import { enhance } from '$app/forms';
	const View = new CreateMentorView();

	const { data }: PageProps = $props();

	onMount(() => {
		View.setPaymentMethods(data.paymentMethods);
		View.setGeneratedPassword(data.generatedPassword);
	});
</script>

<div class="flex h-full flex-col p-4">
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
	<form action="?/createMentor" method="POST" class="flex flex-col gap-4">
		<Input type="text" placeholder="Input mentor name" name="name" id="name" />
		<Input
			type="number"
			placeholder="Years of Experience"
			name="years_of_experience"
			id="years_of_experience"
			min={0}
			max={100}
		/>
		<Input type="email" placeholder="Input mentor email" name="email" id="email" />
		<div class="flex items-center gap-4">
			<div
				class="flex flex-1 gap-1 rounded-lg bg-[var(--tertiary-color)] px-4 py-2 text-[var(--secondary-color)]"
			>
				<p>Password:</p>
				<p>{View.generatedPassword}</p>
			</div>
			<Button type="button" onClick={View.generatePassword}>Generate</Button>
		</div>
		<Input type="text" name="campus" id="campus" placeholder="Mentor Campus" />
		<div class="flex gap-4">
			<Select
				options={degreeOpts}
				defaultLable="Mentor degree"
				name="degree"
				bind:value={View.degree}
			/>
			<Input type="text" name="major" id="major" placeholder="Insert mentor major" />
		</div>
		<div class="flex flex-col gap-4">
			<div class="flex gap-4">
				<Select
					options={dayofWeeks}
					defaultLable="Weekday"
					name="day_of_week"
					bind:value={View.degree}
				/>
				<Input step={1} type="time" name="start" id="start" />
				<Input step={1} type="time" name="end" id="end" />
			</div>
			<Button
				disabled={View.disableAddPaymentMethod}
				full
				type="button"
				onClick={View.addMentorPaymentMethod}>Add Mentor Schedule</Button
			>
		</div>
		<div class="flex flex-col gap-4">
			<div class="flex gap-4">
				<Search
					bind:value={View.selectedPaymentMethod}
					items={View.paymentMethods}
					label="Payment Method"
					onKeywordChange={View.onKeyWordChange}
				/>
				<Input
					type="text"
					name="account_number"
					id="account_number"
					placeholder="Account Number"
					bind:value={View.accountNumber}
				/>
			</div>
			<Button
				disabled={View.disableAddPaymentMethod}
				full
				type="button"
				onClick={View.addMentorPaymentMethod}>Add Payment Method</Button
			>
		</div>
		<div class="flex gap-4">
			<Button type="button">Cancel</Button>
			<Button type="submit">Create</Button>
		</div>
	</form>
</div>
