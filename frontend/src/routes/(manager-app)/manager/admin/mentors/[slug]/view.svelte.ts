import type { EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';

export class MentorDetailView {
	isDesktop = $state<boolean>(false);
	alertOpen = $state<boolean>(false);
	size = $derived.by<number>(() => {
		if (this.isDesktop) {
			return 150;
		}
		return 100;
	});
	onDeleteMentor = () => {
		const loadID = CreateToast('loading', 'deleting....');
		return async ({ result, update }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'redirect') {
				CreateToast('success', 'Successfully delete mentor');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			update();
		};
	};
	setIsDesktop(b: boolean) {
		this.isDesktop = b;
	}
}
