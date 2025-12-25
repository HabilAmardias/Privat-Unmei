import type { EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';

export class CourseDetailView {
	deleteDialogOpen = $state<boolean>(false);
	detailState = $state<'description' | 'detail'>('description');
	capitalizeFirstLetter(s: string) {
		if (s.length === 0) {
			return s;
		}
		if (s.length === 1) {
			return s.toUpperCase();
		}
		return s.charAt(0).toUpperCase() + s.slice(1);
	}
	onDeleteCourse = () => {
		const loadID = CreateToast('loading', 'deleting course....');
		return async ({ result, update }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'redirect') {
				CreateToast('success', 'successfully delete course');
				await update();
			}
		};
	};
}
