import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const [historyReport, incomeReport, mentorReport] = await Promise.all([
		controller.getIncomeHistoryReport(fetch),
		controller.getThisMonthIncomeReport(fetch),
		controller.getThisMonthMentorReport(fetch)
	]);
	if (!historyReport.success) {
		throw error(historyReport.status, { message: historyReport.message });
	}
	if (!incomeReport.success) {
		throw error(incomeReport.status, { message: incomeReport.message });
	}
	if (!mentorReport.success) {
		throw error(mentorReport.status, { message: mentorReport.message });
	}
	return {
		mentorReport: mentorReport.resBody.data,
		incomeReport: incomeReport.resBody.data,
		historyReport: historyReport.resBody.data
	};
};
