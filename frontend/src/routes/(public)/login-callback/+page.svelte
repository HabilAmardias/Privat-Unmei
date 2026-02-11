<script lang="ts">
    import { enhance } from "$app/forms";
    import Button from "$lib/components/button/Button.svelte";
    import Card from "$lib/components/card/Card.svelte";
	import PinInput from "$lib/components/form/PinInput.svelte";
	import { LoginCallbackView } from "./view.svelte";

    const View = new LoginCallbackView()
</script>

<svelte:head>
	<title>Login - Privat Unmei</title>
	<meta name="description" content="Login - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-dvh w-full items-center justify-center">
	<Card>
		<h2 class="mb-3 text-2xl font-bold text-[var(--tertiary-color)]">Enter OTP</h2>
		<form
			use:enhance={View.onLoginSubmit}
			action="?/login"
			method="POST"
			class="flex flex-col gap-4"
		>
            <PinInput bind:value={View.otp}/>
			<Button disabled={View.loginDisabled} type="submit" full={true}>Submit</Button>
		</form>
	</Card>
</div>