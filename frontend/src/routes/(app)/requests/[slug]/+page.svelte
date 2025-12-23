<script lang="ts">
	import type { PageProps } from './$types';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import { RequestDetailView } from './view.svelte';
	import { enhance } from '$app/forms';
	import Button from '$lib/components/button/Button.svelte';
	import NavigationButton from '$lib/components/button/NavigationButton.svelte';
	import { onMount } from 'svelte';
	import { SvelteDate } from 'svelte/reactivity';

	let { data }: PageProps = $props();
	const View = new RequestDetailView(data.detail);

	onMount(() => {
		const intv = setInterval(() => {
			View.now = new SvelteDate();
		}, 1000);
		return () => {
			clearInterval(intv);
		};
	});
</script>

<svelte:head>
	<title>Request Detail - Privat Unmei</title>
	<meta name="description" content="Request Detail - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex flex-col gap-4 p-4">
	<div class="flex justify-between">
		<h1 class="text-xl font-bold text-[var(--tertiary-color)]">{data.detail.course_name}</h1>
		<form method="POST" action="?/messageMentor" use:enhance={View.onMessageMentor}>
			<Button type="submit">Message</Button>
		</form>
	</div>
	<p class="font-bold text-[var(--tertiary-color)]">Request Status:</p>
	<div class="grid grid-cols-5 gap-4 rounded-lg bg-[var(--tertiary-color)] p-4 text-center">
		{#each View.statuses as st}
			<div class="flex flex-col items-center gap-2">
				<st.icon
					class={st.id === View.status
						? 'text-[var(--primary-color)]'
						: 'text-[var(--secondary-color)]'}
				/>
				<p
					class={st.id === View.status
						? 'text-xs text-[var(--primary-color)]'
						: 'text-xs text-[var(--secondary-color)]'}
				>
					{st.label}
				</p>
			</div>
		{/each}
	</div>
	<p class="font-bold text-[var(--tertiary-color)]">Request Detail:</p>
	<div>
		<NavigationButton
			menus={[
				{
					header: 'Detail',
					onClick: () => (View.detailState = 'detail')
				},
				{
					header: 'Payment',
					onClick: () => (View.detailState = 'payment')
				}
			]}
		/>
		<div
			class="grid h-[250px] grid-cols-2 gap-4 rounded-lg rounded-tl-none bg-[var(--tertiary-color)] p-2 text-center"
		>
			{#if View.detailState === 'detail'}
				<div>
					<p class="font-bold text-[var(--primary-color)]">Mentor Name:</p>
					<p class="text-[var(--secondary-color)]">{data.detail.mentor_name}</p>
				</div>
				<div>
					<p class="font-bold text-[var(--primary-color)]">Mentor Public ID:</p>
					<p class="text-[var(--secondary-color)]">{data.detail.mentor_public_id}</p>
				</div>
				<div>
					<p class="font-bold text-[var(--primary-color)]">Participant:</p>
					<p class="text-[var(--secondary-color)]">{data.detail.number_of_participant}</p>
				</div>
				<div>
					<p class="font-bold text-[var(--primary-color)]">Session:</p>
					<p class="text-[var(--secondary-color)]">{data.detail.number_of_sessions}</p>
				</div>
				{#if data.detail.expired_at}
					<div>
						<p class="font-bold text-[var(--primary-color)]">Expired In:</p>
						<p class="text-[var(--secondary-color)]">
							{View.expiredIn}
						</p>
					</div>
				{/if}
			{:else}
				<div>
					<p class="font-bold text-[var(--primary-color)]">Payment Method:</p>
					<p class="text-[var(--secondary-color)]">{data.detail.payment_method}</p>
				</div>
				<div>
					<p class="font-bold text-[var(--primary-color)]">Account Number:</p>
					<p class="text-[var(--secondary-color)]">{data.detail.account_number}</p>
				</div>
				<div>
					<p class="font-bold text-[var(--primary-color)]">Subtotal:</p>
					<p class="text-[var(--secondary-color)]">
						{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
							data.detail.subtotal
						)}
					</p>
				</div>
				<div>
					<p class="font-bold text-[var(--primary-color)]">Operational Cost:</p>
					<p class="text-[var(--secondary-color)]">
						{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
							data.detail.operational_cost
						)}
					</p>
				</div>
				<div>
					<p class="font-bold text-[var(--primary-color)]">Total Price:</p>
					<p class="text-[var(--secondary-color)]">
						{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
							data.detail.total_price
						)}
					</p>
				</div>
			{/if}
		</div>
	</div>
	<div class="flex flex-col gap-4">
		<p class="font-bold text-[var(--tertiary-color)]">Schedules:</p>
		{#if data.detail.schedules}
			<ScrollArea orientation="vertical" viewportClasses="h-[500px] max-h-[500px]">
				<ul class="flex flex-col gap-4">
					{#each data.detail.schedules as sch (sch.date)}
						<li class="rounded-lg bg-[var(--tertiary-color)] p-4 text-[var(--secondary-color)]">
							<p class="font-bold text-[var(--primary-color)]">{View.convertToDate(sch.date)}</p>
							<p>{sch.start_time} - {sch.end_time}</p>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		{:else}
			<b class="text-[var(--tertiary-color)]">No schedules found</b>
		{/if}
	</div>
</div>
