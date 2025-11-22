<script lang="ts">
	import type { PageProps } from './$types';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import { RequestDetailView } from './view.svelte';
	import AlertDialog from '$lib/components/dialog/AlertDialog.svelte';
	import Dialog from '$lib/components/dialog/Dialog.svelte';

	let { data }: PageProps = $props();
	const View = new RequestDetailView();
</script>

<svelte:head>
	<title>Request Detail - Privat Unmei</title>
	<meta name="description" content="Profile - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

{#snippet description()}
	Irreversible Action, do you want to proceed?
{/snippet}

{#snippet rejectTitle()}
	Reject Request Confirmation
{/snippet}

{#snippet acceptTitle()}
	Accept Request Confirmation
{/snippet}

{#snippet confirmTitle()}
	Confirm Payment Confirmation
{/snippet}

<div class="flex flex-col gap-4 p-4">
	<div class="flex items-center justify-between gap-4">
		<h1 class="text-xl font-bold text-[var(--tertiary-color)]">{data.detail.course_name}</h1>
		<div class="flex gap-4">
			{#if data.detail.status === 'reserved'}
				<AlertDialog
					action="?/rejectRequest"
					bind:open={View.rejectDialogOpen}
					title={rejectTitle}
					{description}>Reject</AlertDialog
				>
				<AlertDialog
					action="?/acceptRequest"
					bind:open={View.acceptDialogOpen}
					title={acceptTitle}
					{description}>Accept</AlertDialog
				>
			{:else if data.detail.status === 'pending payment'}
				<AlertDialog
					action="?/confirmPayment"
					bind:open={View.confirmDialogOpen}
					title={confirmTitle}
					{description}>Confirm Payment</AlertDialog
				>
			{/if}
		</div>
	</div>
	<div class="grid grid-cols-2 gap-4 text-center md:flex md:justify-between">
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Course Name:</p>
			<p>{data.detail.course_name}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Student Name:</p>
			<p>{data.detail.student_name}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Participant:</p>
			<p>{data.detail.number_of_participant}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Session:</p>
			<p>{data.detail.number_of_sessions}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Status:</p>
			<p>{data.detail.status}</p>
		</div>
		{#if data.detail.expired_at}
			<div>
				<p class="font-bold text-[var(--tertiary-color)]">Expired At:</p>
				<p>{data.detail.expired_at}</p>
			</div>
		{/if}
	</div>
	<h2 class="text-lg font-bold text-[var(--tertiary-color)]">Payment Info</h2>
	<div class="grid grid-cols-2 gap-4 text-center md:flex md:justify-between">
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Payment Method:</p>
			<p>{data.detail.payment_method}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Account Number:</p>
			<p>{data.detail.account_number}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Subtotal:</p>
			<p>{data.detail.subtotal}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Operational Cost:</p>
			<p>{data.detail.operational_cost}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Total Price:</p>
			<p>{data.detail.total_price}</p>
		</div>
	</div>
	<div class="flex flex-col gap-4">
		<p class="font-bold text-[var(--tertiary-color)]">Schedules:</p>
		{#if data.detail.schedules}
			<ScrollArea orientation="vertical" viewportClasses="max-h-[300px]">
				<ul class="flex flex-col gap-4">
					{#each data.detail.schedules as sch (sch.date)}
						<li class="rounded-lg bg-[var(--tertiary-color)] p-4 text-[var(--secondary-color)]">
							<p class="font-bold text-[var(--primary-color)]">{sch.date}</p>
							<p>{sch.start_time} - {sch.end_time}</p>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		{:else}
			<b class="text-[var(--tertiary-color)]">No topic found</b>
		{/if}
	</div>
</div>
