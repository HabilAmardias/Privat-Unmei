<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Card from '$lib/components/card/Card.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import InputSecret from '$lib/components/form/InputSecret.svelte';
	import { ManagerAuthView } from './view.svelte';
	import { enhance } from '$app/forms';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import toast from 'svelte-french-toast';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { PrivatUnmeiLogo } from '$lib/utils/constants';

	const View = new ManagerAuthView();

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
		};
	});

	function onAdminLoginSubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = toast.loading('logging in.....', { position: 'top-right' });
		return async ({ result, update }: EnhancementReturn) => {
			if (result.type === 'success') {
				await goto('/admin/profile', { replaceState: true });
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
	function onMentorLoginSubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = toast.loading('logging in.....', { position: 'top-right' });
		return async ({ result, update }: EnhancementReturn) => {
			if (result.type === 'success') {
				await goto('/mentor/profile', { replaceState: true });
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
</script>

<svelte:head>
	<title>Manger Login - Privat Unmei</title>
	<meta name="description" content="Login - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-screen w-full flex-col items-center justify-center gap-4 md:flex-row md:gap-0">
	<div class="hidden md:flex md:flex-1">
		<CldImage src={PrivatUnmeiLogo} width={400} height={125} />
	</div>
	<div class="block md:hidden">
		<CldImage src={PrivatUnmeiLogo} width={200} height={60} />
	</div>
	{#if !View.loginAdmin}
		<Card>
			<h2 class="mb-3 text-2xl font-bold text-[var(--tertiary-color)]">Login</h2>
			<form
				use:enhance={onMentorLoginSubmit}
				action="?/loginMentor"
				method="post"
				class="flex flex-col gap-4"
			>
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
				<Button disabled={View.loginDisabled} full={true} type="submit" formAction="?/loginMentor"
					>Login</Button
				>
			</form>
			<Button
				disabled={View.isLoading}
				withBg={false}
				full={true}
				textColor="dark"
				onClick={() => View.switchForm()}>Admin Login</Button
			>
		</Card>
	{:else}
		<Card>
			<h2 class="mb-3 text-2xl font-bold text-[var(--tertiary-color)]">Login</h2>
			<form
				use:enhance={onAdminLoginSubmit}
				action="?/loginAdmin"
				method="post"
				class="flex flex-col gap-4"
			>
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
				<Button disabled={View.loginDisabled} full={true} type="submit" formAction="?/loginAdmin"
					>Login</Button
				>
			</form>
			<Button
				disabled={View.isLoading}
				withBg={false}
				full={true}
				textColor="dark"
				onClick={() => View.switchForm()}>Mentor Login</Button
			>
		</Card>
	{/if}
</div>
