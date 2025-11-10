import { PrivatUnmeiLogoLight, PrivatUnmeiLogoMini } from '$lib/utils/constants';

export class AppLayoutView {
	isDesktop = $state<boolean>(false);
	logoWidth = $derived.by<number>(() => {
		return this.isDesktop ? 90 : 40;
	});
	logoHeight = $derived.by<number>(() => {
		return 30;
	});
	logoSrc = $derived.by<string>(() => {
		return this.isDesktop ? PrivatUnmeiLogoLight : PrivatUnmeiLogoMini;
	});
	setIsDesktop(b: boolean) {
		this.isDesktop = b;
	}
}
