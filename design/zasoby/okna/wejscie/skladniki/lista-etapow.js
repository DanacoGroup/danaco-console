/* Wykaz etapów uruchomienia pokazuje cztery etapy nawiązania połączenia, każdy ze znakiem podążającym za jego stanem, a nie pozostawionym z etapu poprzedniego.
   Znak etapu idzie za stanem: ptaszek dla zrobionego, tętno dla trwającego,
   krzyżyk dla nieudanego, numer dla czekającego.

   Właściwości:
     stany   wykaz stanów, po jednym na etap: 'gotowy' | 'pracuje' | 'blad' | 'oczekuje'
     miary   wykaz kluczy katalogu z opisem stanu po prawej
   ============================================================================ */
(function () {
'use strict';

/* Stan kroku maluje odmiana biblioteki; „oczekuje" nie ma odmiany, bo to
   wygląd spoczynkowy kroku. */
function odmianaStanu(stan) {
  if (stan === 'gotowy') { return ' dn-krok--poprawny'; }
  if (stan === 'pracuje') { return ' dn-krok--pracuje'; }
  if (stan === 'blad') { return ' dn-krok--wstrzymany'; }
  return '';
}
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

function znakEtapu(N, stan, numer) {
  if (stan === 'pracuje') return N.el('span', { klasa: 'pt-tetno', 'aria-hidden': 'true' });
  if (stan === 'gotowy' || stan === 'blad') {
    var i = N.zeZnacznika(W.ikony[stan === 'gotowy' ? 'ptaszek' : 'krzyzyk']);
    i.setAttribute('aria-hidden', 'true');
    return i;
  }
  return document.createTextNode(String(numer));
}

W.skladniki.listaEtapow = function (N, w) {
  var nazwy = N.tekst('uruchomienie.etapy');
  return N.el('ol', {
    klasa: 'dn-kolejka dn-kolejka--pola', 'aria-live': 'polite',
    'aria-label': N.tekst('uruchomienie.laczenie.tytul')
  }, nazwy.map(function (nazwa, i) {
    var stan = w.stany[i] || 'oczekuje';
    return N.el('li', { klasa: 'dn-krok dn-krok--pole' + odmianaStanu(stan), dane: { stan: stan } }, [
      N.el('span', { klasa: 'dn-krok-znak', 'aria-hidden': 'true' }, [znakEtapu(N, stan, i + 1)]),
      N.el('span', { tekst: nazwa }),
      N.el('span', { klasa: 'dn-krok-meta', tekst: N.tekst('uruchomienie.stany.' + (w.miary[i] || 'oczekuje')) })
    ]);
  }));
};
})();
