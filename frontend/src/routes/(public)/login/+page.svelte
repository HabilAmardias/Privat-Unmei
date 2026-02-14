<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Card from '$lib/components/card/Card.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import InputSecret from '$lib/components/form/InputSecret.svelte';
	import { AuthView } from './view.svelte';
	import { PUBLIC_RECAPTCHA_SITE_KEY } from '$env/static/public';
	import { enhance } from '$app/forms';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import Link from '$lib/components/button/Link.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';
	import Google from '$lib/components/icons/Google.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { PrivatUnmeiLogo } from '$lib/utils/constants';
	import { CreateToast, DismissToast } from '$lib/utils/helper';
	import Dialog from '$lib/components/dialog/Dialog.svelte';
	import TAC from './TAC.svelte';

	const View = new AuthView();

	onMount(() => {
		View.setIsDesktop(window.innerWidth >= 768);
		function isDesktop() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		if (loadingStore.logOutLoadID) {
			DismissToast(loadingStore.logOutLoadID);
			loadingStore.removeLogOutLoadID();
			CreateToast('success', 'logout success');
		}
		window.addEventListener('resize', isDesktop);
		return () => {
			window.removeEventListener('resize', isDesktop);
			View.removeGoogleScript();
		};
	});

	async function onRegisterSubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = CreateToast('loading', 'creating account....');
		await new Promise((resolve) => {
			grecaptcha.ready(resolve);
		});
		const token = await grecaptcha.execute(PUBLIC_RECAPTCHA_SITE_KEY, { action: 'submit' });
		args.formData.append('token', token);

		return async ({ result, update }: EnhancementReturn) => {
			View.setIsLoading(false);
			DismissToast(loadID);
			if (result.type === 'success') {
				CreateToast('success', 'successfully registered');
				View.switchForm();
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			update();
		};
	}

	
</script>

{#snippet TACDialogTitle()}
	Terms and Condition
{/snippet}

{#snippet TACDialogDescription()}
	<p class="text-justify text-[var(--secondary-color)]"></p>
{/snippet}

<svelte:head>
	<title>Login - Privat Unmei</title>
	<meta name="description" content="Login - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<script
		bind:this={View.googleScript}
		src="https://www.google.com/recaptcha/api.js?render={PUBLIC_RECAPTCHA_SITE_KEY}"
	></script>
</svelte:head>

<div class="flex h-screen w-full flex-col items-center justify-center gap-4 md:flex-row md:gap-0">
	<div class="hidden md:flex md:flex-1">
		<CldImage src={PrivatUnmeiLogo} width={400} height={125} />
	</div>
	<div class="block md:hidden">
		<CldImage src={PrivatUnmeiLogo} width={200} height={60} />
	</div>
	{#if !View.login}
		<Card>
			<h2 class="mb-3 text-2xl font-bold text-[var(--tertiary-color)]">Register</h2>
			<form
				action="?/register"
				method="post"
				class="flex flex-col gap-4"
				use:enhance={onRegisterSubmit}
			>
				<Input
					err={View.nameError}
					type="text"
					name="name"
					placeholder="Username"
					id="name"
					bind:value={View.name}
				/>
				<Input
					err={View.emailError}
					type="email"
					name="email"
					placeholder="Email"
					id="email"
					bind:value={View.email}
				/>
				<InputSecret
					err={View.passwordError}
					id="password"
					placeholder="Password"
					name="password"
					bind:value={View.password}
				/>
				<InputSecret
					err={View.repeatPasswordError}
					id="repeat-password"
					placeholder="Repeat Password"
					name="repeat-password"
					bind:value={View.repeatPassword}
				/>
				<Dialog
					bind:open={View.openDialog}
					buttonText="Terms and Condition"
					buttonClass="text-[var(--tertiary-color)] hover:text-[var(--primary-color)] cursor-pointer"
					title={TACDialogTitle}
					description={TACDialogDescription}
				>
					<TAC
						onClick={() => {
							View.agreed = true;
							View.openDialog = false;
						}}
					/>
				</Dialog>
				<Button disabled={View.registerDisabled} full={true} type="submit" formAction="?/register"
					>Register</Button
				>
			</form>
			<Button
				disabled={View.isLoading}
				withBg={false}
				textColor="dark"
				full={true}
				onClick={() => View.switchForm()}>Already have an account?</Button
			>
		</Card>
	{:else}
		<Card>
			<h2 class="mb-3 text-2xl font-bold text-[var(--tertiary-color)]">Login</h2>
			<form use:enhance={View.onLoginSubmit} action="?/login" method="post" class="flex flex-col gap-4">
				<Input
					err={View.emailError}
					type="email"
					name="email"
					placeholder="Email"
					id="email"
					bind:value={View.email}
				/>
				<InputSecret
					id="password"
					placeholder="Password"
					name="password"
					bind:value={View.password}
				/>
				<Link theme="dark" href="/reset">Forgot Password?</Link>
				<Button disabled={View.loginDisabled} full={true} type="submit" formAction="?/login"
					>Login</Button
				>
				<Button formAction="?/googlelogin" type="submit" full={true} disabled={View.isLoading}
					>Login With <Google width={24} height={24} /></Button
				>
			</form>
			<Button
				disabled={View.isLoading}
				withBg={false}
				full={true}
				textColor="dark"
				onClick={() => View.switchForm()}>Want to create account?</Button
			>
		</Card>
	{/if}
</div>
