/* ============================================================================
   PRZEPŁYW WEJŚCIA — montaż okien

   Składa dwa okna ze składników i wstawia je w miejsca wskazane w podglądzie:

       <div data-wejscie-okno="uruchomienie"></div>
       <div data-wejscie-okno="dostep"></div>

   Okno jest siatką o trzech wierszach: belka, korpus, pas działań. W wierszu
   pasa stoją dwie rzeczy — pas przez całą szerokość i nota wydawcy w kolumnie
   tożsamości. Panele i pasy noszą ten sam `data-widok`, więc przełączają się
   razem.

   Po zmontowaniu okna zgłasza się zdarzenie `wejscie-gotowe` — mechanika okna
   rusza dopiero na nie, bo wcześniej nie ma czego obsługiwać.
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;

var OKNA = {
  uruchomienie: {
    tytul: 'okno.uruchamianie',
    odslona: 'uruchamianie',
    panele: function (N) { return W.ekrany.uruchomienie(N); },
    pasy: function (N) { return W.ekrany.uruchomieniePasy(N); }
  },
  dostep: {
    tytul: 'okno.dostep',
    odslona: 'dostep',
    panele: function (N) { return W.ekrany.dostep(N); },
    pasy: function (N) { return W.ekrany.dostepPasy(N); }
  },
  /* Okno przygotowania nie ma belki systemowej — stoi już wewnątrz ramy
     aplikacji, a ta niesie własną. Nie ma też odsłon do przełączania: jest
     jedno, więc panel i pas nie noszą `data-widok`. */
  przygotowanie: {
    belka: false,
    nota: 'przygotowanie.nota',
    odslona: 'przygotowanie',
    panele: function (N) { return [W.ekrany.przygotowanie(N)]; },
    pasy: function (N) { return [W.ekrany.przygotowaniePas(N)]; }
  }
};

function zbuduj(N, nazwa) {
  var o = OKNA[nazwa];
  var S = W.skladniki;
  var dzieci = [];
  if (o.belka !== false) dzieci.push(S.belkaOkna(N, { tytul: o.tytul }));
  dzieci.push(S.kolumnaTozsamosci(N, { odslona: o.odslona }));
  dzieci = dzieci.concat(o.panele(N)).concat(o.pasy(N));
  dzieci.push(o.nota ? S.notaPasa(N, o.nota) : S.notaWydawcy(N));
  return N.el('div', {
    klasa: 'we-okno' + (o.belka === false ? '' : ' we-okno--z-belka')
  }, dzieci.filter(Boolean));
}

W.zamontuj = function (korzen) {
  if (!W.tresci) {
    korzen.textContent = '⟨brak katalogu treści — wepnij okna/wejscie/tresci.js⟩';
    return null;
  }
  var N = window.DanacoNarzedzia.zwiaz(W.tresci);
  var pola = (korzen || document).querySelectorAll('[data-wejscie-okno]');
  for (var i = 0; i < pola.length; i++) {
    var nazwa = pola[i].dataset.wejscieOkno;
    if (!OKNA[nazwa]) continue;
    pola[i].appendChild(zbuduj(N, nazwa));
  }
  document.dispatchEvent(new CustomEvent('wejscie-gotowe', { bubbles: true }));
};

/* Samoczynny montaż. Dzięki niemu plik podglądu nie zawiera ani jednej linii
   skryptu — wyłącznie wpięcia i miejsca, w których okna mają stanąć. */
function start() { W.zamontuj(document); }
if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', start);
else start();
})();
