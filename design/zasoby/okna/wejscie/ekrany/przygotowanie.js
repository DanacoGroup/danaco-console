/* Ekran okna przygotowania środowiska pracy stanowi trzeci etap wejścia, w którym program po rozpoznaniu konta odtwarza stan pracy sprzed zamknięcia.
   Okno nie pyta o nic — wykaz etapów mówi, co się dzieje,
   a pas działań daje dwa wyjścia: pominąć przywracanie albo przerwać i się
   wylogować.

   Kolumna tożsamości niesie tu animację powłok zamiast wykazu zalet: na tym
   etapie użytkownik nie wybiera już programu, tylko czeka, aż się złoży.
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
W.ekrany = W.ekrany || {};

/* Stan każdego z pięciu etapów wraz z miarą po prawej. Wykaz jest daną, nie
   rozgałęzieniem w kodzie — przy zmianie przebiegu zmienia się tabela. */
var ETAPY = [
  { stan: 'gotowy',   miara: 'rozpoznane' },
  { stan: 'gotowy',   miara: 'uprawnienia' },
  { stan: 'pracuje',  miara: 'karty' },
  { stan: 'oczekuje', miara: 'oczekuje' },
  { stan: 'oczekuje', miara: 'oczekuje' }
];

W.ekrany.przygotowanie = function (N) {
  var S = W.skladniki;
  var dane = N.tekst('przygotowanie.miaryDane');

  var kroki = N.tekst('przygotowanie.etapy').map(function (nazwa, i) {
    var e = ETAPY[i];
    var znak;
    if (e.stan === 'gotowy') {
      znak = N.zeZnacznika(W.ikony.ptaszek);
      znak.setAttribute('aria-hidden', 'true');
    } else if (e.stan === 'pracuje') {
      znak = N.el('span', { klasa: 'dn-kropka dn-kropka--tetno', 'aria-hidden': 'true' });
    } else {
      znak = document.createTextNode(String(i + 1));
    }
    return N.el('li', { klasa: 'dn-krok dn-krok--pole' + odmianaStanu(e.stan), dane: { stan: e.stan } }, [
      N.el('span', { klasa: 'dn-krok-znak', 'aria-hidden': 'true' }, [znak]),
      N.el('span', { tekst: nazwa }),
      N.el('span', {
        klasa: 'dn-krok-meta',
        tekst: N.podstaw(N.tekst('przygotowanie.miary.' + e.miara), dane)
      })
    ]);
  });

  return N.el('div', { klasa: 'we-panel' }, [
    S.naglowekEkranu(N, {
      nadtytul: 'przygotowanie.nadtytul',
      tytul: 'przygotowanie.tytul',
      lid: 'przygotowanie.lid'
    }),
    N.el('ol', {
      klasa: 'dn-kolejka dn-kolejka--pola', 'aria-live': 'polite',
      'aria-label': N.tekst('przygotowanie.obszarEtapow')
    }, kroki),
    S.pasekPostepu(N, {
      etykieta: 'przygotowanie.postep.etykieta',
      opisPaska: 'przygotowanie.postep.opisPaska',
      wartosc: N.tekst('przygotowanie.postep.wartosc')
    })
  ]);
};

W.ekrany.przygotowaniePas = function (N) {
  return W.skladniki.pasDzialan(N, {
    czynnosci: [
      { klucz: 'dzialania.pominPrzywracanie', komunikat: 'pominiecie' },
      { klucz: 'dzialania.przerwijIWyloguj', cel: 'uwierzytelnienie', grupa: 'etap' }
    ]
  });
};
})();
