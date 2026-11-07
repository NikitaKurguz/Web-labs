(function () {
	'use strict';

	const modal = new bootstrap.Modal(document.getElementById('product-modal'));
	const tbody = document.getElementById('products-tbody');
	const counter = document.getElementById('counter');
	const form = document.getElementById('product-form');
	const modalTitle = document.getElementById('modal-title');

	async function loadProducts() {
		const resp = await fetch('/api/admin/products');
		const products = await resp.json();

		tbody.innerHTML = products.map(p => `
			<tr data-id="${p.id}">
				<td>${p.id}</td>
				<td><img src="/${p.img}" style="width:50px;height:50px;object-fit:cover"></td>
				<td>${p.name}</td>
				<td>${p.type}</td>
				<td>${p.color}</td>
				<td>${p.category.join(', ')}</td>
				<td>${p.price === 0 ? 'Бесплатно' : p.price + ' ₽'}</td>
				<td>
					<button class="btn btn-sm btn-outline-dark" data-action="edit">Изм.</button>
					<button class="btn btn-sm btn-outline-danger" data-action="delete">Удал.</button>
				</td>
			</tr>
		`).join('');

		counter.textContent = 'Всего товаров: ' + products.length;
	}

	document.getElementById('btn-new-product').addEventListener('click', () => {
		form.reset();
		form.querySelector('#form-id').value = '';
		modalTitle.textContent = 'Добавить товар';
		modal.show();
	});

	tbody.addEventListener('click', async (e) => {
		const btn = e.target.closest('button[data-action]');
		if (!btn) return;

		const id = btn.closest('tr').dataset.id;

		if (btn.dataset.action === 'delete') {
			if (!confirm('Удалить?')) return;
			await fetch('/api/admin/products/' + id, { method: 'DELETE' });
			await loadProducts();
			return;
		}

		if (btn.dataset.action === 'edit') {
			const resp = await fetch('/api/admin/products');
			const products = await resp.json();
			const p = products.find(x => x.id === Number(id));
			if (!p) return;

			form.querySelector('#form-id').value = p.id;
			form.querySelector('[name="name"]').value = p.name;
			form.querySelector('[name="img"]').value = p.img;
			form.querySelector('[name="type"]').value = p.type;
			form.querySelector('[name="color"]').value = p.color;
			form.querySelector('[name="category"]').value = p.category.join(',');
			form.querySelector('[name="price"]').value = p.price;

			modalTitle.textContent = 'Редактировать товар';
			modal.show();
		}
	});

	form.addEventListener('submit', async (e) => {
		e.preventDefault();

		const id = form.querySelector('#form-id').value;

		const data = {
			name: form.querySelector('[name="name"]').value.trim(),
			img: form.querySelector('[name="img"]').value.trim(),
			type: form.querySelector('[name="type"]').value.trim(),
			color: form.querySelector('[name="color"]').value.trim(),
			category: form.querySelector('[name="category"]').value.split(',').map(s => s.trim()).filter(Boolean),
			price: Number(form.querySelector('[name="price"]').value) || 0,
		};

		const url = id ? '/api/admin/products/' + id : '/api/admin/products';
		const method = id ? 'PUT' : 'POST';

		const resp = await fetch(url, {
			method: method,
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(data),
		});

		if (!resp.ok) {
			const err = await resp.json().catch(() => ({}));
			alert('Ошибка: ' + (err.error || 'неизвестная'));
			return;
		}

		modal.hide();
		await loadProducts();
	});

	loadProducts();
})();