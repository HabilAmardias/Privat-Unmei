import { monthMap } from './constants';
import type { historyIncomeReport, monthlyIncomeReport, monthlyMentorReport } from './model';

export class AdminDashboardView {
	totalCost = $state<number>();
	totalSession = $state<number>();
	mentorReports = $state<monthlyMentorReport[]>([]);
	costValueHistoryReport = $state<number[]>([]);
	costLabelHistoryReport = $state<string[]>([]);
	sessionValueHistoryReport = $state<number[]>([]);
	sessionLabelHistoryReport = $state<string[]>([]);

	constructor(
		historyReport: historyIncomeReport,
		incomeReport: monthlyIncomeReport,
		mentorReport: monthlyMentorReport[]
	) {
		this.totalCost = incomeReport.total_cost;
		this.totalSession = incomeReport.total_session;
		this.mentorReports = mentorReport;

		historyReport.cost_report.forEach((e) => {
			this.costValueHistoryReport.push(e.total_cost);
			this.costLabelHistoryReport.push(monthMap[e.month - 1]);
		});

		historyReport.session_report.forEach((e) => {
			this.sessionLabelHistoryReport.push(monthMap[e.month - 1]);
			this.sessionValueHistoryReport.push(e.total_session);
		});
	}
}
