/* Składnik listy etapów postępu przedstawia przebieg zapisu instalacji rozbity na kolejne etapy wraz z ich stanem i bieżącą czynnością.

   Przebieg zapisu rozbity na etapy. Stan etapu niesie znak I słowo naraz:
   ✓ gotowe, wskaźnik pracy w toku, pusty pierścień oczekuje — stan nigdy nie
   stoi na samej barwie. Opis pod nazwą mówi, co dzieje się z komputerem.

   Miary etapu (czasownik licznika, cel, jednostka, pozostały czas) jadą
   w atrybutach `data-*`: należą do etapu, a nie do kodu, który go rysuje.

   Właściwości:
     naglowek   klucz katalogu — nagłówek bloku (opcjonalny)
     etapy      klucz katalogu — tablica { nazwa, opis, licznik, pozostalo }
     stany      klucz katalogu — słowa stanów { gotowe, wToku, oczekuje }
     poczatkowe tablica stanów ('gotowe' | 'w-toku' | 'oczekuje') na start */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

var ZNAK_GOTOWE = null;

K.skladniki.listaEtapow = function (w) {
  ZNAK_GOTOWE = ZNAK_GOTOWE || K.ikony.ptaszek;
  var slowa = tekst(w.stany);
  var SLOWO = { 'gotowe': slowa.gotowe, 'w-toku': slowa.wToku, 'oczekuje': slowa.oczekuje };

  var pozycje = tekst(w.etapy).map(function (e, i) {
    var stan = (w.poczatkowe && w.poczatkowe[i]) || 'oczekuje';
    var znak = el('span', { klasa: 'dn-krok-znak', dane: { 'etap-znak': true }, 'aria-hidden': 'true' });
    if (stan === 'gotowe') znak.innerHTML = ZNAK_GOTOWE;
    else if (stan === 'w-toku') znak.appendChild(el('span', { klasa: 'dn-spinner' }));

    var dane = { etap: stan };
    if (e.licznik) {
      dane['licznik-czasownik'] = e.licznik.czasownik;
      dane['licznik-cel'] = e.licznik.cel;
      if (e.licznik.jednostka) dane['licznik-jednostka'] = e.licznik.jednostka;
      if (e.licznik.rzecz) dane['licznik-rzecz'] = e.licznik.rzecz;
    }
    if (e.pozostalo) dane['pozostalo'] = e.pozostalo;

    return el('li', {
      klasa: 'dn-krok' + (stan === 'gotowe' ? ' dn-krok--poprawny' : (stan === 'w-toku' ? ' dn-krok--pracuje' : '')),
      dane: dane
    }, [
      znak,
      el('span', { tekst: e.nazwa }),
      el('span', { klasa: 'dn-krok-stan', dane: { 'etap-slowo': true }, tekst: SLOWO[stan] }),
      el('span', { klasa: 'dn-krok-opis', tekst: e.opis })
    ]);
  });

  var kolejka = el('ol', { klasa: 'dn-kolejka', dane: { etapy: true } }, pozycje);
  if (!w.naglowek) return kolejka;
  return el('div', { dane: { blok: 'kolejka' } }, [K.skladniki.naglowekBloku({ klucz: w.naglowek }), kolejka]);
};
})();
