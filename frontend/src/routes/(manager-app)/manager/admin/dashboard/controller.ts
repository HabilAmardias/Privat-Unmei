import type { ServerResponse, Fetch } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { monthlyIncomeReport, monthlyMentorReport, historyIncomeReport } from './model';

class DashboardController {
	async getThisMonthIncomeReport(fetch: Fetch) {
		const url = `/api/v1/admins/reports/income-report`;
		const { success, message, status, res } = await FetchData(fetch, url);
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<monthlyIncomeReport> = await res?.json();
		return { success, message, status, resBody };
	}
	async getIncomeHistoryReport(fetch: Fetch) {
		const url = `/api/v1/admins/reports/history-report`;
		const { success, message, status, res } = await FetchData(fetch, url);
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<historyIncomeReport> = await res?.json();
		return { success, message, status, resBody };
	}
	async getThisMonthMentorReport(fetch: Fetch) {
		const url = `/api/v1/admins/reports/mentor-report`;
		const { success, message, status, res } = await FetchData(fetch, url);
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<monthlyMentorReport[]> = await res?.json();
		return { success, message, status, resBody };
	}
}

export const controller = new DashboardController();
