function createUser() {
  const username = prompt('请输入新用户名:');
  if (!username) return;
  const password = prompt('请输入新用户密码:');
  if (!password) return;
  fetch('/api/user/create', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
    credentials: 'include'
  })
    .then(res => res.json())
    .then(data => {
      if (data.success) {
        alert('用户创建成功');
        refreshList && refreshList();
      } else {
        alert(data.error || '用户创建失败');
      }
    });
}
let sid = '';
let currentFolder = 0;
let pathStack = [{id: 0, name: '根目录'}];

function setCookie(name, value) {
  document.cookie = name + '=' + value + ';path=/';
}

function getCookie(name) {
  let arr = document.cookie.match(new RegExp('(^| )' + name + '=([^;]*)(;|$)'));
  return arr ? arr[2] : '';
}

function login() {
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;
  fetch('/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
    credentials: 'include'
  })
    .then(res => res.json())
    .then(data => {
      if (data.sid) {
        sid = data.sid;
        setCookie('sid', sid);
        document.getElementById('login-area').classList.add('hidden');
        document.getElementById('main-area').classList.remove('hidden');
        refreshList();
      } else {
        document.getElementById('login-msg').innerText = data.error || '登录失败';
      }
    });
}

function logout() {
  sid = '';
  setCookie('sid', '', -1); // 立即过期
  document.getElementById('login-area').classList.remove('hidden');
  document.getElementById('main-area').classList.add('hidden');
}

function refreshList() {
  fetch(`/api/file/${currentFolder}/list`, {
    credentials: 'include'
  })
    .then(res => res.json())
    .then(data => {
      renderTable(data.files || data);
      renderPath();
    });
}


function renderTable(files) {
  const tbody = document.querySelector('#file-table tbody');
  tbody.innerHTML = '';
  files.forEach(file => {
    const tr = document.createElement('tr');
    if (file.type === 'folder' || file.type === '') {
      // 文件夹或空类型，点击名称进入
      tr.innerHTML = `
        <td style="cursor:pointer;color:#1976d2;text-decoration:underline" onclick="enterFolder(${file.id},'${file.name}')">${file.name}</td>
        <td>文件夹</td>
        <td>
          <button class='op-btn' onclick='renameFile(${file.id},"${file.name}")'>重命名</button>
          <button class='op-btn' onclick='deleteFile(${file.id})'>删除</button>
        </td>
      `;
    } else {
      // 普通文件
      tr.innerHTML = `
        <td>${file.name}</td>
        <td>文件</td>
        <td>
          <button class='op-btn' onclick='downloadFile(${file.id})'>下载</button>
          <button class='op-btn' onclick='renameFile(${file.id},"${file.name}")'>重命名</button>
          <button class='op-btn' onclick='deleteFile(${file.id})'>删除</button>
        </td>
      `;
    }
    tbody.appendChild(tr);
  });
}

function renderPath() {
  const pathElem = document.getElementById('current-path');
  pathElem.innerHTML = '';
  pathStack.forEach((p, idx) => {
    if (idx > 0) pathElem.innerHTML += ' / ';
    if (idx === pathStack.length - 1) {
      pathElem.innerHTML += `<span>${p.name}</span>`;
    } else {
      pathElem.innerHTML += `<a href=\"javascript:goBackFolder(${idx})\" style=\"color:#1976d2;text-decoration:underline\">${p.name}</a>`;
    }
  });
}


function enterFolder(id, name) {
  currentFolder = id;
  pathStack.push({id, name});
  refreshList();
}

function goBackFolder(index) {
  // 跳转到某一级目录
  if (index < 0 || index >= pathStack.length) return;
  currentFolder = pathStack[index].id;
  pathStack = pathStack.slice(0, index + 1);
  refreshList();
}
window.goBackFolder = goBackFolder;

function uploadFile() {
  const fileInput = document.getElementById('upload-file');
  const file = fileInput.files[0];
  if (!file) return;
  const formData = new FormData();
  formData.append('file', file);
  fetch(`/api/file/${currentFolder}/upload`, {
    method: 'POST',
    body: formData,
    credentials: 'include'
  }).then(() => {
    fileInput.value = '';
    refreshList();
  });
}

function createFolder() {
  const name = prompt('请输入新文件夹名称:');
  if (!name) return;
  fetch(`/api/file/${currentFolder}/new`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name }),
    credentials: 'include'
  }).then(() => refreshList());
}

function renameFile(id, oldName) {
  const newName = prompt('请输入新名称:', oldName);
  if (!newName || newName === oldName) return;
  fetch(`/api/file/${id}/rename`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ new_name: newName }),
    credentials: 'include'
  }).then(() => refreshList());
}

function deleteFile(id) {
  if (!confirm('确定要删除吗？')) return;
  fetch(`/api/file/${id}`, {
    method: 'DELETE',
    credentials: 'include'
  }).then(() => refreshList());
}

function downloadFile(id) {
  window.open(`/api/file/${id}/content?sid=${sid}`);
}

// 添加到 window 方便按钮调用
window.createUser = createUser;
