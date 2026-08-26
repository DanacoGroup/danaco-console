/* ============================================================================
   KREATOR INSTALACJI — narzędzia wspólne składników

   Trzy rzeczy, których potrzebuje każdy składnik i których nie powinien
   wymyślać po raz dziesiąty: budowanie węzła, sięganie po łańcuch z katalogu
   treści i podstawianie danych w łańcuch.

   Nic tu nie wie o instalatorze — to warstwa niżej niż składniki.
   ============================================================================ */
(function () {
'use strict';

var K = (window.DanacoKreator = window.DanacoKreator || {});

/* Budowa węzła. Atrybuty rozpoznawane po nazwie: `klasa` → class, `tekst` →
   textContent, `html` → innerHTML (wyłącznie dla treści z katalogu, nigdy dla
   danych), `dane` → komplet atrybutów `data-*`, reszta wprost jako atrybut.
   Wartość `null` albo `false` pomija atrybut — dzięki temu warianty składnika
   piszą się warunkiem, a nie rozgałęzieniem. */
function el(znacznik, atrybuty, dzieci) {
  var n = document.createElement(znacznik);
  atrybuty = atrybuty || {};
  Object.keys(atrybuty).forEach(function (k) {
    var v = atrybuty[k];
    if (v === null || v === undefined || v === false) return;
    if (k === 'klasa') n.className = v;
    else if (k === 'tekst') n.textContent = v;
    else if (k === 'html') n.innerHTML = v;
    else if (k === 'dane') Object.keys(v).forEach(function (d) {
      if (v[d] === null || v[d] === undefined || v[d] === false) return;
      n.setAttribute('data-' + d, v[d] === true ? '' : v[d]);
    });
    else if (v === true) n.setAttribute(k, '');
    else n.setAttribute(k, v);
  });
  (dzieci || []).forEach(function (d) {
    if (d === null || d === undefined || d === false) return;
    n.appendChild(typeof d === 'string' ? document.createTextNode(d) : d);
  });
  return n;
}

/* Węzeł zbudowany ze znacznika. Używany wyłącznie dla ikon z zestawu i dla
   treści dokumentu wstawianej przy budowaniu — nigdy dla danych z zewnątrz. */
function zeZnacznika(html) {
  var t = document.createElement('template');
  t.innerHTML = html.trim();
  return t.content.firstElementChild;
}

/* Sięgnięcie po łańcuch ścieżką kluczy: `tekst('krok1.wymagania.naglowek')`.
   Brak klucza jest usterką katalogu, nie sytuacją do obsłużenia po cichu —
   zwracana jest sama ścieżka, żeby brak rzucał się w oczy w oknie. */
function tekst(sciezka) {
  var w = K.tresci;
  var czesci = String(sciezka).split('.');
  for (var i = 0; i < czesci.length && w != null; i++) w = w[czesci[i]];
  return w === undefined || w === null ? '⟨' + sciezka + '⟩' : w;
}

/* Podstawienie danych w łańcuch: `podstaw('Etap {numer} z {ile}', {numer: 2, ile: 4})`.
   Nazwy w nawiasach są częścią kontraktu z tłumaczem — nie tłumaczy się ich. */
function podstaw(wzor, dane) {
  return String(wzor).replace(/\{(\w+)\}/g, function (calosc, nazwa) {
    return dane && dane[nazwa] !== undefined ? dane[nazwa] : calosc;
  });
}

K.narzedzia = { el: el, zeZnacznika: zeZnacznika, tekst: tekst, podstaw: podstaw };
})();
