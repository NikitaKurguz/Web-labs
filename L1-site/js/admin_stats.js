(function () {
	'use strict';

	const aggTbody = document.getElementById('agg-tbody');
	const aggCounter = document.getElementById('agg-counter');
	const selGroup = document.getElementById('filter-group');
	const selAction = document.getElementById('filter-action');
	const selSort = document.getElementById('filter-sort');
	const btnAggApply = document.getElementById('btn-agg-apply');

	async function loadAgg() {
		const params = new URLSearchParams();
		params.set('group', selGroup.value);
		params.set('sort', selSort.value);
		params.set('order', 'desc');
		if (selAction.value) params.set('action', selAction.value);

		aggTbody.innerHTML = `<tr><td colspan="2" class="text-center text-muted py-4">Загрузка...</td></tr>`;

		try {
			const resp = await fetch('/api/admin/stats?' + params.toString());
			if (!resp.ok) {
				aggTbody.innerHTML = `<tr><td colspan="2" class="text-center text-danger py-4">Ошибка</td></tr>`;
				return;
			}
			const rows = await resp.json();
			renderAgg(rows);
		} catch (err) {
			console.error(err);
		}
	}

	function renderAgg(rows) {
		if (!rows.length) {
			aggTbody.innerHTML = `<tr><td colspan="2" class="text-center text-muted py-4">Нет данных</td></tr>`;
			aggCounter.textContent = '';
			return;
		}
		aggTbody.innerHTML = rows.map(r => `
			<tr>
				<td>${formatPeriod(r.period)}</td>
				<td>${r.count}</td>
			</tr>
		`).join('');
		aggCounter.textContent = 'Строк: ' + rows.length;
	}

	function formatPeriod(iso) {
		const d = new Date(iso);
		const pad = n => String(n).padStart(2, '0');
		const g = selGroup.value;
		if (g === 'hour') {
			return `${pad(d.getDate())}.${pad(d.getMonth()+1)}.${d.getFullYear()} ${pad(d.getHours())}:00`;
		}
		if (g === 'month') {
			const m = ['янв','фев','мар','апр','май','июн','июл','авг','сен','окт','ноя','дек'];
			return `${m[d.getMonth()]} ${d.getFullYear()}`;
		}
		return `${pad(d.getDate())}.${pad(d.getMonth()+1)}.${d.getFullYear()}`;
	}

	btnAggApply.addEventListener('click', loadAgg);

	const logTbody = document.getElementById('log-tbody');
	const logCounter = document.getElementById('log-counter');
	const logAction = document.getElementById('log-action');
	const logSort = document.getElementById('log-sort');
	const logOrder = document.getElementById('log-order');
	const btnLogApply = document.getElementById('btn-log-apply');

	async function loadLog() {
		const params = new URLSearchParams();
		params.set('sort', logSort.value);
		params.set('order', logOrder.value);
		params.set('limit', '200');
		if (logAction.value) params.set('action', logAction.value);

		logTbody.innerHTML = `<tr><td colspan="6" class="text-center text-muted py-4">Загрузка...</td></tr>`;

		try {
			const resp = await fetch('/api/admin/actions?' + params.toString());
			if (!resp.ok) {
				logTbody.innerHTML = `<tr><td colspan="6" class="text-center text-danger py-4">Ошибка</td></tr>`;
				return;
			}
			const rows = await resp.json();
			renderLog(rows);
		} catch (err) {
			console.error(err);
		}
	}

	function renderLog(rows) {
		if (!rows.length) {
			logTbody.innerHTML = `<tr><td colspan="6" class="text-center text-muted py-4">Нет данных</td></tr>`;
			logCounter.textContent = '';
			return;
		}
		logTbody.innerHTML = rows.map(r => `
			<tr>
				<td>${r.id}</td>
				<td>${formatDate(r.createdAt)}</td>
				<td>${r.userId || '—'}</td>
				<td>${escapeHtml(r.action)}</td>
				<td>${escapeHtml(r.target) || '—'}</td>
				<td>${r.ip || '—'}</td>
			</tr>
		`).join('');
		logCounter.textContent = 'Записей: ' + rows.length;
	}

	function formatDate(iso) {
		const d = new Date(iso);
		const pad = n => String(n).padStart(2, '0');
		return `${pad(d.getDate())}.${pad(d.getMonth()+1)}.${d.getFullYear()} ` +
			`${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
	}

	function escapeHtml(s) {
		if (s == null) return '';
		return String(s)
			.replaceAll('&', '&amp;')
			.replaceAll('<', '&lt;')
			.replaceAll('>', '&gt;')
			.replaceAll('"', '&quot;');
	}

	btnLogApply.addEventListener('click', loadLog);

	loadAgg();
	loadLog();
})();