/* ============================================================================
   NARZĘDZIA OKIEN — warstwa pod składnikami

   Trzy rzeczy, których potrzebuje każde okno składane ze składników i których
   nie powinno wymyślać po raz drugi: budowanie węzła, sięganie po łańcuch
   z katalogu treści i podstawianie danych w łańcuch.

   Nic tu nie wie o żadnym oknie — to warstwa niżej niż składniki.

       var N = DanacoNarzedzia.zwiaz(katalogTresci);
       N.el('p', { klasa: 'we-lid', tekst: N.tekst('krok1.lid') });
   ============================================================================ */
(function () {
'use strict';

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

/* Węzeł zbudowany ze znacznika. Używany wyłącznie dla ikon z zestawu — nigdy
   dla danych z zewnątrz. */
function zeZnacznika(html) {
  var t = document.createElement('template');
  t.innerHTML = String(html).trim();
  return t.content.firstElementChild;
}

/* Podstawienie danych w łańcuch: `podstaw('Etap {numer} z {ile}', {numer: 2})`.
   Nazwy w nawiasach są częścią kontraktu z tłumaczem — nie tłumaczy się ich. */
function podstaw(wzor, dane) {
  return String(wzor).replace(/\{(\w+)\}/g, function (calosc, nazwa) {
    return dane && dane[nazwa] !== undefined ? dane[nazwa] : calosc;
  });
}

/* Związanie narzędzi z katalogiem jednego okna. Katalog przekazuje się raz,
   przy montażu — składniki dostają już gotowe `tekst`. */
function zwiaz(katalog) {
  /* Sięgnięcie po łańcuch ścieżką kluczy: `tekst('krok3.warianty.x64.opis')`.
     Brak klucza jest usterką katalogu, nie sytuacją do obsłużenia po cichu —
     zwracana jest sama ścieżka, żeby brak rzucał się w oczy w oknie. */
  function tekst(sciezka) {
    var w = katalog;
    var czesci = String(sciezka).split('.');
    for (var i = 0; i < czesci.length && w != null; i++) w = w[czesci[i]];
    return w === undefined || w === null ? '⟨' + sciezka + '⟩' : w;
  }
  return { el: el, zeZnacznika: zeZnacznika, podstaw: podstaw, tekst: tekst, katalog: katalog };
}

/* Odczyt żetonu do postaci, którą przyjmuje płótno. Próbka jest niewidoczna
   i znika po zbudowaniu palety — służy wyłącznie temu, żeby przeglądarka
   rozwinęła wartość żetonu do składowych. Potrzebuje jej każdy składnik
   rysujący na płótnie, więc stoi tu, a nie w trzech kopiach. */
function czytnikBarw(host) {
  var probka = document.createElement('span');
  probka.style.cssText = 'position:absolute;width:0;height:0;visibility:hidden';
  host.appendChild(probka);
  var styl = getComputedStyle(host);
  return {
    barwa: function (zeton, zapas) {
      var v = styl.getPropertyValue(zeton).trim() || zapas;
      probka.style.color = '';
      probka.style.color = v;
      var m = getComputedStyle(probka).color.match(/[\d.]+/g) || ['0', '0', '0'];
      var rgb = [Math.round(+m[0]), Math.round(+m[1]), Math.round(+m[2])];
      var a = m.length > 3 ? +m[3] : 1;
      return { rgb: rgb, a: a, css: 'rgba(' + rgb[0] + ',' + rgb[1] + ',' + rgb[2] + ',' + a + ')' };
    },
    zdejmij: function () { if (probka.parentNode) probka.parentNode.removeChild(probka); }
  };
}

/* Zakładanie składnika na wszystkich polach danego rodzaju. Każdy składnik
   rysujący powtarzał ten sam blok: wykaz pól, pętla, nasłuch na wczytaniu
   dokumentu. */
function polaSkladnika(uchwyt, zaloz) {
  function wszystkie(zakres) {
    var pola = (zakres || document).querySelectorAll('[' + uchwyt + ']');
    for (var i = 0; i < pola.length; i++) zaloz(pola[i]);
  }
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function () { wszystkie(); });
  } else {
    wszystkie();
  }
  return wszystkie;
}

window.DanacoNarzedzia = {
  el: el, zeZnacznika: zeZnacznika, podstaw: podstaw, zwiaz: zwiaz,
  czytnikBarw: czytnikBarw, polaSkladnika: polaSkladnika
};
})();
