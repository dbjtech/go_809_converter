let API_BASE = (function(){
  const u = window.location.href;
  const i = u.indexOf('/static');
  return i === -1 ? u.substring(0, u.lastIndexOf('/')) : u.substring(0, i);
})();

const FIELDS = [
  'enable','encryptKey','govServerIP','govServerPort','localServerIP','localServerPort','name',
  'openCrypto','platformId','platformPassword','platformUserId','protocolVersion'
];

const DISPLAY_LABELS = {
  enable: '开启推送',
  encryptKey: '加密秘钥',
  govServerIP: '对方平台IP/域名',
  govServerPort: '对方平台端口',
  localServerIP: '通过本IP/域名连接本平台',
  localServerPort: '通过本端口连接本平台',
  name: '对方平台名称',
  openCrypto: '开启加密功能',
  platformId: '连接对方平台的id',
  platformPassword: '连接对方平台的密码',
  platformUserId: '连接对方平台的用户ID',
  protocolVersion: '协议版本'
};

function getEnv(config){
  try{
    const env = typeof config.env === 'string' ? config.env.trim() : '';
    if (env) return env;
    const keys = Object.keys(config||{});
    const cands = ['develop','online','staging','production'];
    for (const cand of cands){ if (keys.some(k => k === cand)) return cand; }
    return 'develop';
  }catch(e){ return 'develop'; }
}

function showStatus(msg, ok){
  const el = document.getElementById('status');
  el.className = ok ? 'success' : 'error';
  el.textContent = msg;
  el.style.display = 'block';
  if (ok) setTimeout(() => { el.style.display = 'none'; }, 2000);
}

function renderNodes(config){
  const env = getEnv(config);
  const nodesWrap = document.getElementById('nodes');
  nodesWrap.innerHTML = '';
  const conv = ((config[env]||{}).converter)||{};
  Object.keys(conv).forEach(name => {
    const node = conv[name]||{};
    const card = document.createElement('div');
    card.className = 'node';
    const header = document.createElement('div');
    header.className = 'node-header';
    header.innerHTML = `<strong>${env}.converter.${name}</strong>`;
    const act = document.createElement('div');
    const delBtn = document.createElement('button');
    delBtn.className = 'btn btn-danger';
    delBtn.textContent = '删除节点';
    delBtn.addEventListener('click', (e)=>{ e.preventDefault(); deleteNode(env, name); });
    const saveBtn = document.createElement('button');
    saveBtn.className = 'btn btn-primary';
    saveBtn.textContent = '保存该节点';
    act.appendChild(saveBtn);
    act.appendChild(delBtn);
    header.appendChild(act);
    card.appendChild(header);
    const grid = document.createElement('div');
    grid.className = 'grid';
    FIELDS.forEach(f => {
      const item = document.createElement('div');
      item.className = 'item';
      const label = document.createElement('label');
      label.textContent = DISPLAY_LABELS[f] || f;
      let input;
      const v = node[f];
      if (f === 'name') {
        const col = document.createElement('div');
        col.className = 'input-col';
        input = document.createElement('input');
        input.type = 'text';
        input.value = (v==null? '' : v);
        input.pattern = '^[A-Za-z0-9_]+$';
        input.title = '仅限字母、数字、下划线';
        input.setAttribute('data-field', f);
        const err = document.createElement('small');
        err.className = 'field-error';
        col.appendChild(input);
        col.appendChild(err);
        item.appendChild(label);
        item.appendChild(col);
        grid.appendChild(item);
        return;
      }
      if (f === 'govServerIP' || f === 'localServerIP'){
        const col = document.createElement('div');
        col.className = 'input-col';
        input = document.createElement('input');
        input.type = 'text';
        input.value = (v==null? '' : v);
        input.setAttribute('data-field', f);
        const err = document.createElement('small');
        err.className = 'field-error';
        col.appendChild(input);
        col.appendChild(err);
        item.appendChild(label);
        item.appendChild(col);
        grid.appendChild(item);
        return;
      }
      if (f === 'enable' || f === 'openCrypto'){
        input = document.createElement('input');
        input.type = 'checkbox';
        input.checked = !!v;
      }else if (['encryptKey','govServerPort','localServerPort','platformId','platformUserId'].includes(f)){
        if (f === 'govServerPort' || f === 'localServerPort'){
          const col = document.createElement('div');
          col.className = 'input-col';
          input = document.createElement('input');
          input.type = 'number';
          input.value = (v==null? '' : v);
          input.min = '1';
          input.max = '65535';
          input.step = '1';
          input.setAttribute('data-field', f);
          const err = document.createElement('small');
          err.className = 'field-error';
          col.appendChild(input);
          col.appendChild(err);
          item.appendChild(label);
          item.appendChild(col);
          grid.appendChild(item);
          return;
        } else {
          input = document.createElement('input');
          input.type = 'number';
          input.value = (v==null? '' : v);
        }
      }else if (f === 'platformPassword'){
        const wrap = document.createElement('span');
        wrap.className = 'password-wrapper';
        input = document.createElement('input');
        input.type = 'password';
        input.autocomplete = 'new-password';
        input.setAttribute('form','config-form');
        input.value = (v==null? '' : v);
        const btn = document.createElement('button');
        btn.type = 'button';
        btn.className = 'toggle-password';
        btn.setAttribute('data-visible','false');
        btn.textContent = '👁';
        btn.addEventListener('click', (e)=>{
          e.preventDefault();
          const isPwd = input.type === 'password';
          input.type = isPwd ? 'text' : 'password';
          btn.textContent = isPwd ? '🙈' : '👁';
          btn.title = isPwd ? '隐藏' : '显示';
        });
        wrap.appendChild(input);
        wrap.appendChild(btn);
        item.appendChild(label);
        item.appendChild(wrap);
        grid.appendChild(item);
        return;
      }else{
        input = document.createElement('input');
        input.type = 'text';
        input.value = (v==null? '' : v);
        if (f === 'name') {
          input.pattern = '^[A-Za-z0-9_]+$';
          input.title = '仅限字母、数字、下划线';
        }
      }
      input.setAttribute('data-field', f);
      item.appendChild(label);
      item.appendChild(input);
      grid.appendChild(item);
    });
    card.appendChild(grid);
    saveBtn.addEventListener('click', (e)=>{
      e.preventDefault();
      const payload = {};
      FIELDS.forEach(f => {
        const el = grid.querySelector(`[data-field="${f}"]`) || grid.querySelector(`#${f}`);
        if (!el) return;
        if (el.type === 'checkbox') payload[f] = el.checked;
        else if (el.type === 'number') payload[f] = el.value === '' ? null : Number(el.value);
        else payload[f] = el.value;
      });
  const nameEl = grid.querySelector('[data-field="name"]');
      if (!validateNameInput(nameEl)) { nameEl.focus(); shake(nameEl); return; }
      const ipEls = [grid.querySelector('[data-field="govServerIP"]'), grid.querySelector('[data-field="localServerIP"]')];
      for (const el of ipEls){ if (!validateHostInput(el)) { el && el.focus(); shake(el); return; } }
      const portEls = [grid.querySelector('[data-field="govServerPort"]'), grid.querySelector('[data-field="localServerPort"]')];
      for (const el of portEls){ if (!validatePortInput(el)) { el && el.focus(); shake(el); return; } }
      saveNode(env, name, payload);
    });
    nodesWrap.appendChild(card);
  });
}

function saveNode(env, name, node){
  const nested = {};
  nested[env] = { converter: {} };
  nested[env].converter[name] = node;
  $.ajax({
    url: `${API_BASE}/setting/save`, method: 'POST', contentType: 'application/json',
    data: JSON.stringify({ config: nested, operation: 'update' }),
    success: (d)=>{ if (d.success){ showStatus('保存成功', true); load(); } else { showStatus('保存失败: '+d.message, false);} },
    error: ()=> showStatus('网络错误，保存失败', false)
  });
}

function deleteNode(env, name){
  const key = `${env}.converter.${name}`;
  $.ajax({
    url: `${API_BASE}/setting/delete`, method: 'DELETE', contentType: 'application/json',
    data: JSON.stringify({ key, operation: 'delete' }),
    success: (d)=>{ if (d.success){ showStatus('删除成功', true); load(); } else { showStatus('删除失败: '+d.message, false);} },
    error: ()=> showStatus('网络错误，删除失败', false)
  });
}

function addNode(){
  const name = $('#new-name').val().trim();
  if (!name){ showStatus('请输入节点名称', false); return; }
  const newNameEl = document.getElementById('new-name');
  if (!validateNameInput(newNameEl)) { newNameEl.focus(); shake(newNameEl); return; }
  const newGovIp = document.getElementById('new-govServerIP');
  const newLocalIp = document.getElementById('new-localServerIP');
  const newGovPort = document.getElementById('new-govServerPort');
  const newLocalPort = document.getElementById('new-localServerPort');
  if (!validateHostInput(newGovIp)) { newGovIp.focus(); shake(newGovIp); return; }
  if (!validateHostInput(newLocalIp)) { newLocalIp.focus(); shake(newLocalIp); return; }
  if (!validatePortInput(newGovPort)) { newGovPort.focus(); shake(newGovPort); return; }
  if (!validatePortInput(newLocalPort)) { newLocalPort.focus(); shake(newLocalPort); return; }
  const env = window.__ENV__ || 'develop';
  const node = {
    name,
    enable: $('#new-enable').is(':checked'),
    encryptKey: Number($('#new-encryptKey').val()||0),
    govServerIP: $('#new-govServerIP').val(),
    govServerPort: Number($('#new-govServerPort').val()||0),
    localServerIP: $('#new-localServerIP').val(),
    localServerPort: Number($('#new-localServerPort').val()||0),
    openCrypto: $('#new-openCrypto').is(':checked'),
    platformId: Number($('#new-platformId').val()||0),
    platformPassword: $('#new-platformPassword').val(),
    platformUserId: Number($('#new-platformUserId').val()||0),
    protocolVersion: $('#new-protocolVersion').val()
  };
  const nested = {}; nested[env] = { converter: {} }; nested[env].converter[name] = node;
  $.ajax({
    url: `${API_BASE}/setting/save`, method: 'POST', contentType: 'application/json',
    data: JSON.stringify({ config: nested, operation: 'add_subproject' }),
    success: (d)=>{ if (d.success){ showStatus('添加成功', true); clearAddForm(); load(); } else { showStatus('添加失败: '+d.message, false);} },
    error: ()=> showStatus('网络错误，添加失败', false)
  });
}

function clearAddForm(){
  $('#new-name').val('');
  $('#new-enable').prop('checked', true);
  $('#new-encryptKey').val('223344');
  $('#new-govServerIP').val('127.0.0.1');
  $('#new-govServerPort').val('19001');
  $('#new-localServerIP').val('127.0.0.1');
  $('#new-localServerPort').val('1301');
  $('#new-openCrypto').prop('checked', false);
  $('#new-platformId').val('1001');
  $('#new-platformUserId').val('100101');
  $('#new-platformPassword').val('');
  $('#new-protocolVersion').val('1.0.0');
}

function load(){
  $.get(`${API_BASE}/setting/current`, (d)=>{
    if (!d.success){ showStatus('加载配置失败', false); return; }
    const cfg = d.config||{}; const env = getEnv(cfg); window.__ENV__ = env;
    renderNodes(cfg);
  }).fail(()=> showStatus('网络错误，加载失败', false));
}

$(document).ready(function(){
  load();
  $('#btn-add').on('click', function(e){ e.preventDefault(); addNode(); });
  $('#btn-clear').on('click', function(e){ e.preventDefault(); clearAddForm(); });
  $(document).on('click', '.toggle-password[data-target]', function(e){
    e.preventDefault();
    const id = this.getAttribute('data-target');
    const input = document.getElementById(id);
    if (!input) return;
    const isPwd = input.type === 'password';
    input.type = isPwd ? 'text' : 'password';
    this.textContent = isPwd ? '🙈' : '👁';
    this.title = isPwd ? '隐藏' : '显示';
  });
  $(document).on('input', 'input[data-field="name"]', function(){ validateNameInput(this); });
  $('#new-name').on('input', function(){ validateNameInput(this); });
  $(document).on('input', 'input[data-field="govServerIP"], input[data-field="localServerIP"]', function(){ validateHostInput(this); });
  $(document).on('input', 'input[data-field="govServerPort"], input[data-field="localServerPort"]', function(){ validatePortInput(this); });
  $('#new-govServerIP, #new-localServerIP').on('input', function(){ validateHostInput(this); });
  $('#new-govServerPort, #new-localServerPort').on('input', function(){ validatePortInput(this); });
});

function validateNameInput(el){
  if (!el) return false;
  const val = (el.value||'').trim();
  const ok = /^[A-Za-z0-9_]+$/.test(val);
  const errEl = el.parentElement && el.parentElement.querySelector('.field-error');
  if (ok){
    el.classList.remove('input-invalid');
    if (errEl){ errEl.style.display = 'none'; errEl.textContent = ''; }
    return true;
  } else {
    el.classList.add('input-invalid');
    if (errEl){ errEl.style.display = 'block'; errEl.textContent = '仅限字母、数字、下划线'; }
    return false;
  }
}

function validateHostInput(el){
  if (!el) return false;
  const val = (el.value||'').trim();
  const ipv4 = /^((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)\.){3}(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)$/;
  const domain = /^(localhost|([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+)$/;
  const ok = ipv4.test(val) || domain.test(val);
  const errEl = el.parentElement && el.parentElement.querySelector('.field-error');
  if (ok){
    el.classList.remove('input-invalid');
    if (errEl){ errEl.style.display = 'none'; errEl.textContent = ''; }
    return true;
  } else {
    el.classList.add('input-invalid');
    if (errEl){ errEl.style.display = 'block'; errEl.textContent = '请输入合法的IPv4或域名'; }
    return false;
  }
}

function validatePortInput(el){
  if (!el) return false;
  const val = (el.value||'').trim();
  const num = Number(val);
  const ok = Number.isInteger(num) && num >= 1 && num <= 65535;
  const errEl = el.parentElement && el.parentElement.querySelector('.field-error');
  if (ok){
    el.classList.remove('input-invalid');
    if (errEl){ errEl.style.display = 'none'; errEl.textContent = ''; }
    return true;
  } else {
    el.classList.add('input-invalid');
    if (errEl){ errEl.style.display = 'block'; errEl.textContent = '端口需为1-65535的整数'; }
    return false;
  }
}

function shake(el){
  if (!el) return;
  el.classList.remove('shake');
  void el.offsetWidth;
  el.classList.add('shake');
  el.addEventListener('animationend', function(){ el.classList.remove('shake'); }, { once: true });
}