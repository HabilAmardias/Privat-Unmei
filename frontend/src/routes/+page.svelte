<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import Card from '$lib/components/card/Card.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import InputSecret from '$lib/components/form/InputSecret.svelte';
	import { View } from './view.svelte';
</script>

<svelte:head>
	<title>Login</title>
	<meta name="description" content="Login - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-full w-full items-center justify-center">
	{#if View.cardState === 'Register'}
		<Card title={View.cardState}>
			<form action="?/register" method="post" class="flex flex-col gap-4">
				<Input type="text" name="name" placeholder="Username" id="name" bind:value={View.name} />
				<Input type="email" name="email" placeholder="email" id="email" bind:value={View.email} />
				<InputSecret placeholder="Password" name="password" bind:value={View.password} />
				<InputSecret
					placeholder="Repeat Password"
					name="repeat-password"
					bind:value={View.repeatPassword}
				/>
				<div class="flex justify-end px-2">
					<Button full={true} type="submit" formAction="?/register">Register</Button>
				</div>
			</form>
		</Card>
	{:else}
		<Card title={View.cardState}>
			<form action="?/login" method="post" class="flex flex-col gap-4">
				<Input type="email" name="email" placeholder="Email" id="email" bind:value={View.email} />
				<InputSecret placeholder="Password" name="password" bind:value={View.password} />
				<div class="flex justify-end">
					<Button full={true} type="submit" formAction="?/login">Login</Button>
				</div>
			</form>
		</Card>
	{/if}
</div>
