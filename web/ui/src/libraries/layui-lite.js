const layer = {
  msg(text, opts = {}) {
    const duration = opts.time || 2000;
    const el = document.createElement('div');
    el.className = 'layui-layer layui-layer-msg';
    el.style.top = '20px';
    el.style.left = '50%';
    el.style.transform = 'translateX(-50%)';
    el.style.padding = '12px 24px';
    el.style.zIndex = 2000;
    el.innerText = text;
    document.body.appendChild(el);
    setTimeout(() => {
      if (el.parentNode) {
        el.parentNode.removeChild(el);
      }
    }, duration);
  },
  confirm(text, options = {}, yes) {
    const shade = document.createElement('div');
    shade.className = 'layui-layer-shade';
    const box = document.createElement('div');
    box.className = 'layui-layer';
    const title = document.createElement('div');
    title.className = 'layui-layer-title';
    title.innerText = options.title || 'Confirm';
    const content = document.createElement('div');
    content.className = 'layui-layer-content';
    content.innerText = text;
    const btnWrap = document.createElement('div');
    btnWrap.className = 'layui-layer-btn';
    const okBtn = document.createElement('button');
    okBtn.className = 'layui-btn';
    okBtn.innerText = 'OK';
    const cancelBtn = document.createElement('button');
    cancelBtn.className = 'layui-btn layui-btn-primary';
    cancelBtn.innerText = 'Cancel';
    btnWrap.appendChild(cancelBtn);
    btnWrap.appendChild(okBtn);
    box.appendChild(title);
    box.appendChild(content);
    box.appendChild(btnWrap);
    document.body.appendChild(shade);
    document.body.appendChild(box);
    const closeAll = () => {
      if (shade.parentNode) shade.parentNode.removeChild(shade);
      if (box.parentNode) box.parentNode.removeChild(box);
    };
    okBtn.addEventListener('click', () => {
      closeAll();
      yes && yes();
    });
    cancelBtn.addEventListener('click', closeAll);
  },
  close() {}
};

window.layui = window.layui || {};
window.layui.layer = window.layui.layer || layer;

export default window.layui;
