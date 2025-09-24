<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Card from '$lib/components/card/Card.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import InputSecret from '$lib/components/form/InputSecret.svelte';
	import { View } from './view.svelte';
	import { PUBLIC_RECAPTCHA_SITE_KEY } from '$env/static/public';
	import { enhance } from '$app/forms';
	import type { PageProps } from './$types';

	function switchForm() {
		View.switchForm();
	}

	let { data, form }: PageProps = $props();
</script>

<svelte:head>
	<title>Login</title>
	<meta name="description" content="Login - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<script src="https://www.google.com/recaptcha/api.js?render={PUBLIC_RECAPTCHA_SITE_KEY}"></script>
</svelte:head>

<div class="flex h-full w-full items-center justify-center">
	{#if !View.login}
		<Card>
			<h2 class="mb-3 text-2xl font-bold text-[var(--tertiary-color)]">Register</h2>
			<form
				action="?/register"
				method="post"
				class="flex flex-col gap-4"
				use:enhance={async ({ formData }) => {
					await new Promise((resolve) => {
						grecaptcha.ready(resolve);
					});
					const token = await grecaptcha.execute(PUBLIC_RECAPTCHA_SITE_KEY, { action: 'submit' });
					formData.append('token', token);
				}}
			>
				<Input type="text" name="name" placeholder="Username" id="name" bind:value={View.name} />
				<Input type="email" name="email" placeholder="Email" id="email" bind:value={View.email} />
				<InputSecret
					id="password"
					placeholder="Password"
					name="password"
					bind:value={View.password}
				/>
				<InputSecret
					id="repeat-password"
					placeholder="Repeat Password"
					name="repeat-password"
					bind:value={View.repeatPassword}
				/>
				<Button full={true} type="submit" formAction="?/register">Register</Button>
			</form>
			<Button withBg={false} onClick={switchForm}>Already have an account?</Button>
		</Card>
	{:else}
		<Card>
			<h2 class="mb-3 text-2xl font-bold text-[var(--tertiary-color)]">Login</h2>
			<form action="?/login" method="post" class="flex flex-col gap-4">
				<Input type="email" name="email" placeholder="Email" id="email" bind:value={View.email} />
				<InputSecret
					id="password"
					placeholder="Password"
					name="password"
					bind:value={View.password}
				/>
				<Button full={true} type="submit" formAction="?/login">Login</Button>
			</form>
			<Button withBg={false} onClick={switchForm}>Want to create account?</Button>
		</Card>
	{/if}
	{#if form?.success}
		<p>{form.message}</p>
	{/if}
</div>
