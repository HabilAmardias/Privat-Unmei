import type { EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';

export class CourseDetailView {
	deleteDialogOpen = $state<boolean>(false);
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
