// i18n.js — phone UI strings; language follows tablet Lobby (GET /api/phone-ui).
var strings = {};
var language = 'en';
var ready = false;

export function t(key) {
  if (strings[key]) return strings[key];
  return key;
}

export function tf(key) {
  var s = t(key);
  for (var i = 1; i < arguments.length; i++) {
    s = s.split('%' + i).join(String(arguments[i]));
  }
  return s;
}

export function currentLanguage() {
  return language;
}

export function i18nReady() {
  return ready;
}

export function isPhoneUA() {
  return /iPhone|iPod/.test(navigator.userAgent) ||
    (/Android/.test(navigator.userAgent) && /Mobile/.test(navigator.userAgent));
}

export function applyPasteButtonLabel() {
  var btn = document.getElementById('typing-paste');
  if (!btn) return;
  btn.textContent = isPhoneUA() ? t('paste.fromPhone') : t('paste.fromHere');
}

function applyStatic() {
  document.documentElement.lang = language;
  var nodes = document.querySelectorAll('[data-i18n]');
  for (var i = 0; i < nodes.length; i++) {
    var el = nodes[i];
    var key = el.getAttribute('data-i18n');
    if (key) el.textContent = t(key);
  }
  var titles = document.querySelectorAll('[data-i18n-title]');
  for (var j = 0; j < titles.length; j++) {
    var te = titles[j];
    var tk = te.getAttribute('data-i18n-title');
    if (tk) te.title = t(tk);
  }
  var ph = document.querySelectorAll('[data-i18n-placeholder]');
  for (var k = 0; k < ph.length; k++) {
    var pe = ph[k];
    var pk = pe.getAttribute('data-i18n-placeholder');
    if (pk) pe.placeholder = t(pk);
  }
  applyPasteButtonLabel();
}

export function formatSyncAgo(lastSyncAt) {
  if (!lastSyncAt) return '';
  var sec = Math.max(0, Math.floor(Date.now() / 1000) - lastSyncAt);
  if (sec < 60) return t('sync.agoJustNow');
  var mins = Math.floor(sec / 60);
  if (mins === 1) return t('sync.agoMinute');
  if (mins < 60) return tf('sync.agoMinutes', mins);
  var hours = Math.floor(mins / 60);
  if (hours === 1) return t('sync.agoHour');
  if (hours < 48) return tf('sync.agoHours', hours);
  var days = Math.floor(hours / 24);
  if (days === 1) return t('sync.agoDay');
  return tf('sync.agoDays', days);
}

export function loadI18n() {
  return fetch('/api/phone-ui', { credentials: 'same-origin' })
    .then(function(r) { return r.ok ? r.json() : null; })
    .then(function(data) {
      if (data && data.strings) {
        language = data.language || 'en';
        strings = data.strings;
      }
      ready = true;
      applyStatic();
      return language;
    })
    .catch(function() {
      ready = true;
      applyStatic();
      return language;
    });
}
