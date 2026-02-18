export type monthlyIncomeReport = {
	total_session: number;
	total_cost: number;
};

export type monthlyCostReport = {
	month: number;
	total_cost: number;
};

export type monthlySessionReport = {
	month: number;
	total_session: number;
};

export type historyIncomeReport = {
	session_report: monthlySessionReport[];
	cost_report: monthlyCostReport[];
};

export type monthlyMentorReport = {
	name: string;
	email: string;
	total_cost: number;
	total_session: number;
};
