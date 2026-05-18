document.addEventListener('DOMContentLoaded', function () {
  var pwdInputs = document.querySelectorAll('input[type="password"]');
  pwdInputs.forEach(function (input) {
    var wrapper = document.createElement('div');
    wrapper.style.position = 'relative';
    input.parentNode.insertBefore(wrapper, input);
    wrapper.appendChild(input);

    var toggle = document.createElement('button');
    toggle.type = 'button';
    toggle.setAttribute('aria-label', 'Mostrar/ocultar senha');
    toggle.style.cssText = [
      'position:absolute', 'right:14px', 'top:50%', 'transform:translateY(-50%)',
      'background:none', 'border:none', 'cursor:pointer', 'padding:0',
      'color:#90A4AE', 'font-size:18px', 'line-height:1', 'display:flex',
      'align-items:center'
    ].join(';');
    toggle.innerHTML = '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>';
    wrapper.appendChild(toggle);

    toggle.addEventListener('click', function () {
      if (input.type === 'password') {
        input.type = 'text';
        toggle.style.color = '#1976D2';
      } else {
        input.type = 'password';
        toggle.style.color = '#90A4AE';
      }
    });
  });

  var submitBtn = document.querySelector('.btn-primary, #kc-login, #kc-register');
  if (submitBtn) {
    var form = submitBtn.closest('form');
    if (form) {
      form.addEventListener('submit', function () {
        submitBtn.disabled = true;
        submitBtn.style.opacity = '0.7';
        submitBtn.value = 'Aguarde...';
      });
    }
  }
});
