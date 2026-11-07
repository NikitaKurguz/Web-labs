document.addEventListener('DOMContentLoaded', function(){
    var link = document.getElementById('PrintStyles');
    var button = document.getElementById('printToggle');
    if (!link || !button) return;
    link.disabled = true;

    if (localStorage.getItem('primarte') === '1'){
        link.disabled = false;
        button.setAttribute('aria-pressed', 'true');
        button.textContent = 'Обычная версия';
    }

    button.addEventListener('click', function(){
        var on = link.disabled;
        link.disabled = !on;
        button.setAttribute('aria-pressed', String(on));
        button.textContent = on ? 'Обычная версия' : 'Версия для печати';

        let un;
        if (link.disabled) un = '0';
        else un = '1';
        localStorage.setItem('primarte', un);
    });

});