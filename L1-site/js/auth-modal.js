(async function () {
	'use strict';

	try {
		const [loginHtml, registerHtml] = await Promise.all([
			fetch('/partials/modal-login.html').then(r => r.text()),
			fetch('/partials/modal-register.html').then(r => r.text()),
		]);

		document.body.insertAdjacentHTML('beforeend', loginHtml + registerHtml);
	} catch (err) {
		console.error('auth-modal: не удалось загрузить шаблоны', err);
		return;
	}

	function openModal(id) {
		const m = document.getElementById(id);
		if (m) m.hidden = false;
	}

	function closeAll() {
		document.querySelectorAll('.auth-modal').forEach(m => m.hidden = true);
	}

	document.addEventListener('click', function (e) {
		const opener = e.target.closest('[data-open]');
		if (opener) {
			e.preventDefault();
			closeAll();
			openModal(opener.dataset.open);
			return;
		}
		if (e.target.closest('[data-close]')) {
			e.preventDefault();
			closeAll();
		}
	});

	document.addEventListener('keydown', function (e) {
		if (e.key === 'Escape') closeAll();
	});
})();