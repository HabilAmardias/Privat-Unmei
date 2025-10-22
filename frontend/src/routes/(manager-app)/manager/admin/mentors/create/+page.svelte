<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Select from '$lib/components/select/Select.svelte';
	import { onMount } from 'svelte';
	import { degreeOpts } from './constants';
	import { CreateMentorView } from './view.svelte';
	import type { PageProps } from './$types';
	import Search from '$lib/components/search/Search.svelte';
	import { enhance } from '$app/forms';
	const View = new CreateMentorView();

	const { data }: PageProps = $props();

	onMount(() => {
		View.setPaymentMethods(data.paymentMethods);
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
		<div class="flex gap-4">
			<Input type="text" placeholder="Input mentor name" name="name" id="name" />
			<Input type="email" placeholder="Input mentor email" name="email" id="email" />
		</div>
		<div class="flex items-center justify-between">
			<div class="flex gap-3">
				<p>Password:</p>
				<p>{View.generatedPassword ? View.generatedPassword : data.generatedPassword}</p>
			</div>
			<Button type="button" onClick={View.generatePassword}>Generate</Button>
		</div>
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
		<Button type="button" onClick={View.addMentorPaymentMethod}>Add</Button>
		<div class="flex gap-4">
			<Select
				options={degreeOpts}
				defaultLable="Mentor degree"
				name="degree"
				bind:value={View.degree}
			/>
			<Input type="text" name="major" id="major" placeholder="Insert mentor major" />
		</div>
	</form>
</div>
