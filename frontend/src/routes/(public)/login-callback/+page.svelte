<script lang="ts">
    import { enhance } from "$app/forms";
	import { LoginCallbackView } from "./view.svelte";
	import { PrivatUnmeiLogo } from "$lib/utils/constants";
    import Button from "$lib/components/button/Button.svelte";
    import Card from "$lib/components/card/Card.svelte";
	import PinInput from "$lib/components/form/PinInput.svelte";
	import CldImage from "$lib/components/image/CldImage.svelte";

    const View = new LoginCallbackView()
</script>

<!-- TODO: Refine OTP Page Design and Add timer to resend OTP -->

<svelte:head>
	<title>Login - Privat Unmei</title>
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
	<Card>
		<h2 class="mb-3 text-2xl font-bold text-[var(--tertiary-color)]">Enter Code</h2>
		<form
			use:enhance={View.onLoginSubmit}
			action="?/submitOTP"
			method="POST"
			class="flex flex-col gap-4"
		>
            <PinInput bind:value={View.otp}/>
			<Button formAction="?/resendOTP" withBg={false} textColor="dark">Resend Code</Button>
			<Button disabled={View.loginDisabled} type="submit" full={true}>Submit</Button>
		</form>
	</Card>
</div>