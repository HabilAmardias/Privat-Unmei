<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Card from '$lib/components/card/Card.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import InputSecret from '$lib/components/form/InputSecret.svelte';
	import { View } from './view.svelte';
	import { PUBLIC_RECAPTCHA_SITE_KEY } from '$env/static/public';
	import { enhance } from '$app/forms';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import toast from 'svelte-french-toast';
	import Link from '$lib/components/button/Link.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';
	import Google from '$lib/components/icons/Google.svelte';
	import Image from '$lib/components/image/Image.svelte';
	import LandingIcons from '$lib/images/website-maintenance.png';

	onMount(() => {
		View.setIsDesktop(window.innerWidth >= 768);
		function isDesktop() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		if (loadingStore.logOutLoadID) {
			toast.dismiss(loadingStore.logOutLoadID);
			loadingStore.removeLogOutLoadID();
			toast.success('log out success', { position: 'top-right' });
		}
		window.addEventListener('resize', isDesktop);
		return () => {
			window.removeEventListener('resize', isDesktop);
			View.removeGoogleScript();
		};
	});

	async function onRegisterSubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = toast.loading('Creating account....', {
			position: 'top-right'
		});
		await new Promise((resolve) => {
			grecaptcha.ready(resolve);
		});
		const token = await grecaptcha.execute(PUBLIC_RECAPTCHA_SITE_KEY, { action: 'submit' });
		args.formData.append('token', token);

		return async ({ result, update }: EnhancementReturn) => {
			View.setIsLoading(false);
			toast.dismiss(loadID);
			if (result.type === 'success') {
				toast.success('Successfully registered', {
					position: 'top-right'
				});
				View.switchForm();
			}
			if (result.type === 'failure') {
				toast.error(result.data?.message, {
					position: 'top-right'
				});
			}
			update();
		};
	}

	function onLoginSubmit(args: EnhancementArgs) {
		if (args.action.search === '?/login') {
			View.setIsLoading(true);
			const loadID = toast.loading('logging in.....', { position: 'top-right' });
			return async ({ result, update }: EnhancementReturn) => {
				if (result.type === 'success') {
					localStorage.setItem('status', result.data?.userStatus);
					await goto('/home', { replaceState: true });
					View.setIsLoading(false);
					toast.dismiss(loadID);
					toast.success('login success', {
						position: 'top-right'
					});
				}
				if (result.type === 'failure') {
					View.setIsLoading(false);
					toast.dismiss(loadID);
					toast.error(result.data?.message, {
						position: 'top-right'
					});
				}
				update();
			};
		}
	}
</script>

<svelte:head>
	<title>Login</title>
	<meta name="description" content="Login - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<script
		bind:this={View.googleScript}
		src="https://www.google.com/recaptcha/api.js?render={PUBLIC_RECAPTCHA_SITE_KEY}"
	></script>
</svelte:head>

<div class="flex h-full w-full items-center justify-center md:justify-between">
	{#if View.isDesktop}
		<Image src={LandingIcons} width={500} height={500} />
	{/if}
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
					onBlur={() => View.nameOnBlur()}
					type="text"
					name="name"
					placeholder="Username"
					id="name"
					bind:value={View.name}
				/>
				<Input
					err={View.emailError}
					onBlur={() => View.emailOnBlur()}
					type="email"
					name="email"
					placeholder="Email"
					id="email"
					bind:value={View.email}
				/>
				<InputSecret
					err={View.passwordError}
					onBlur={() => View.passwordOnBlur()}
					id="password"
					placeholder="Password"
					name="password"
					bind:value={View.password}
				/>
				<InputSecret
					err={View.repeatPasswordError}
					onBlur={() => View.repeatPasswordOnBlur()}
					id="repeat-password"
					placeholder="Repeat Password"
					name="repeat-password"
					bind:value={View.repeatPassword}
				/>
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
			<form use:enhance={onLoginSubmit} action="?/login" method="post" class="flex flex-col gap-4">
				<Input
					err={View.emailError}
					onBlur={() => View.emailOnBlur()}
					type="email"
					name="email"
					placeholder="Email"
					id="email"
					bind:value={View.email}
				/>
				<InputSecret
					err={View.passwordError}
					onBlur={() => View.passwordOnBlur()}
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
