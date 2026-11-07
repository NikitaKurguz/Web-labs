let items = [];

const properties = {
	type:     'Тип файла',
	color:    'Цвет',
	category: 'Категория',
	price:    'Цена'
};

function loadItems(callback) {
	fetch('/api/products')
		.then(response => response.json())
		.then(data => {
			items = data;
			callback();
		});
}

function allValues(prop) {
	let vals = [];
	for (let item of items) {
		let v = item[prop];
		if (Array.isArray(v)) {
			for (let x of v) {
				if (!vals.includes(x)) vals.push(x);
			}
		} else {
			if (!vals.includes(v)) vals.push(v);
		}
	}
	vals.sort();
	return vals;
}

function printItems(arItems, selector) {

	const startTemplate = '<div class="catalog__grid">';
	const itemTemplate =
		'<article class="product-card">' +
			'<div class="product-card__image">' +
				'<img src="{{img}}" alt="{{name}}" loading="lazy">' +
			'</div>' +
			'<div class="product-card__body">' +
				'<div>' +
					'<h3 class="product-card__title">{{name}}</h3>' +
					'<ul class="product-card__meta">' +
						'<li>{{type}}</li>' +
						'<li>{{color}}</li>' +
						'<li>{{category}}</li>' +
					'</ul>' +
				'</div>' +
				'<span class="product-card__price">{{priceLabel}}</span>' +
			'</div>' +
		'</article>';
	const endTemplate = '</div>';

	let output = startTemplate;
	for (let item of arItems) {
		let cats = Array.isArray(item.category)
			? item.category.join(', ')
			: item.category;

		let tmpLine = itemTemplate
			.replaceAll('{{img}}',        item.img)
			.replaceAll('{{name}}',       item.name)
			.replaceAll('{{type}}',       item.type)
			.replaceAll('{{color}}',      item.color)
			.replaceAll('{{category}}',   cats)
			.replaceAll('{{priceLabel}}', item.priceLabel);
		output += tmpLine;
	}
	output += endTemplate;

	document.querySelector(selector).innerHTML = output;
	document.querySelector('#counter').textContent = 'Найдено: ' + arItems.length + ' из ' + items.length;
	document.querySelector('#empty').hidden = arItems.length !== 0;
}

function printFilters(arItems, arProperties, selector) {

	const startTemplate =
		'<fieldset class="filter-group" data-prop="{{prop}}">' +
			'<legend class="filter-group__title">{{label}}</legend>';
	const itemTemplate =
		'<label class="filter-item">' +
			'<input type="checkbox" name="{{prop}}" value="{{value}}">' +
			'<span class="filter-item__text">{{value}}</span>' +
			'<span class="filter-item__count" data-count="{{prop}}:{{value}}"></span>' +
		'</label>';
	const endTemplate = '</fieldset>';

	let output = '';
	for (let prop in arProperties) {
		let tmpLine = startTemplate
			.replace('{{prop}}',  prop)
			.replace('{{label}}', arProperties[prop]);

		for (let value of allValues(prop)) {
			tmpLine += itemTemplate
				.replaceAll('{{prop}}',  prop)
				.replaceAll('{{value}}', value);
		}
		output += tmpLine;
	}

	output += '<button type="button" class="filter-reset" id="filterReset">Сбросить фильтр</button>';
	document.querySelector(selector).innerHTML = output;
}

function readCurFilters(arProperties) {

	let result = {};
	for (let prop in arProperties) {
		let checked = [];
		let boxes = document.querySelectorAll('#filters input[name="' + prop + '"]:checked');
		for (let cb of boxes) checked.push(cb.value);
		result[prop] = checked;
	}
	return result;
}

function applyFilters(data, filter, arProperties) {

	let result = [];
	for (let item of data) {
		let ok = true;
		for (let prop in arProperties) {
			if (!filter[prop].length) continue;

			let values = Array.isArray(item[prop]) ? item[prop] : [item[prop]];
			let match  = false;
			for (let v of values) {
				if (filter[prop].indexOf(v) !== -1) { match = true; break; }
			}
			if (!match) ok = false;
		}
		if (ok) result.push(item);
	}
	return result;
}

function computeAvailability(data, filter, arProperties) {

	let avail = {};
	for (let prop in arProperties) {
		avail[prop] = {};
		let othersOk = applyFilters(data, filter, arProperties);

		for (let value of allValues(prop)) {
			let found = false;
			for (let item of othersOk) {
				let values = Array.isArray(item[prop]) ? item[prop] : [item[prop]];
				if (values.indexOf(value) !== -1) { found = true; break; }
			}
			avail[prop][value] = found;
		}
	}
	return avail;
}

function update() {

	let curFilter = readCurFilters(properties);
	let avail     = computeAvailability(items, curFilter, properties);

	for (let prop in properties) {
		let boxes = document.querySelectorAll('#filters input[name="' + prop + '"]');
		for (let cb of boxes) {
			let value    = cb.value;
			let possible = avail[prop][value];

			cb.disabled = !possible && !cb.checked;
			cb.closest('.filter-item').classList.toggle('is-disabled', cb.disabled);

			let tempFilter = {};
			for (let k in curFilter) tempFilter[k] = curFilter[k];
			tempFilter[prop] = [value];

			let cnt = applyFilters(items, tempFilter, properties).length;
			let counter = document.querySelector('[data-count="' + prop + ':' + value + '"]');
			if (counter) counter.textContent = cnt ? '(' + cnt + ')' : '(0)';
		}
	}

	let filtered = applyFilters(items, curFilter, properties);
	printItems(filtered, '#elements');
}

document.addEventListener('DOMContentLoaded', function () {

	loadItems(function () {

		printItems(items, '#elements');
		printFilters(items, properties, '#filters');

		document.querySelector('#filters').addEventListener('change', update);

		document.querySelector('#filters').addEventListener('click', function (e) {
			if (e.target.id === 'filterReset') {
				let boxes = document.querySelectorAll('#filters input:checked');
				for (let cb of boxes) cb.checked = false;
				update();
			}
		});

		update();
	});
});