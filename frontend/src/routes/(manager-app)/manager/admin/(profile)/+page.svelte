<script lang="ts">
	import type { PageProps } from './$types';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { adminProfileView } from './view.svelte';
	import Card from '$lib/components/card/Card.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { enhance } from '$app/forms';
	import InputSecret from '$lib/components/form/InputSecret.svelte';
	import { onMount } from 'svelte';

	let { data }: PageProps = $props();
	const View = new adminProfileView();
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

<svelte:head>
	<title>Profile - Privat Unmei</title>
	<meta name="description" content="Profile - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

{#if View.isEdit}
	<div class="flex h-dvh w-full items-center justify-center">
		<Card>
			<h2 class="mb-3 text-2xl font-bold text-[var(--tertiary-color)]">Change Password</h2>
			<form
				use:enhance={View.onChangePasswordSubmit}
				action="?/changePassword"
				method="post"
				class="flex flex-col gap-4"
			>
				<InputSecret
					err={View.passwordError}
					onBlur={() => View.passwordOnBlur()}
					id="password"
					placeholder="New Password"
					bind:value={View.password}
					name="password"
				/>
				<InputSecret
					err={View.repeatPasswordError}
					onBlur={() => View.repeatPasswordOnBlur()}
					id="password"
					placeholder="Repeat New Password"
					bind:value={View.repeatPassword}
					name="repeat-password"
				/>
				<div class="flex gap-1">
					<Button onClick={() => View.switchForm()}>Cancel</Button>
					<Button disabled={View.isDisabled} type="submit">Submit</Button>
				</div>
			</form>
		</Card>
	</div>
{:else}
	<div class="flex h-full flex-col gap-4 p-4">
		<div class="flex items-center gap-4">
			<CldImage
				width={View.size}
				height={View.size}
				src={data.profile.profile_image}
				className="rounded-full"
			/>
			<div class="flex flex-col gap-1">
				<div class="flex gap-1">
					<b class="text-xl text-[var(--tertiary-color)]">{data.profile.name}</b>
				</div>
				<p class="text-md">{data.profile.email}</p>
				<Button onClick={() => View.switchForm()}>Change Password</Button>
			</div>
		</div>
		<div class="flex flex-col gap-2">
			<b class="text-xl text-[var(--tertiary-color)]">Bio:</b>
			<p class="text-justify">{data.profile.bio}</p>
		</div>
	</div>
{/if}
