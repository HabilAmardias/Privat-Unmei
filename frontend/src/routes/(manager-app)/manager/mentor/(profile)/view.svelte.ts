export class MentorDetailView {
	isDesktop = $state<boolean>(false);
	size = $derived.by<number>(() => {
		if (this.isDesktop) {
			return 150;
		}
		return 100;
	});
	setIsDesktop(b: boolean) {
		this.isDesktop = b;
	}
}
