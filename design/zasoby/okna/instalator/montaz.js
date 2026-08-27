/* Moduł montażu składa okno kreatora instalacji z belki, korpusu z szyną kroków i płótnem ekranów, pasa działań oraz okien dialogowych.


   Składa okno z bryły rodziny `.dn-kreator-*` i ze składników atomowych:
   belka, korpus (szyna kroków + płótno z sześcioma ekranami), pas działań
   i okna dialogowe. Nie zna treści — bierze ją z katalogu; nie zna wyglądu —
   bierze go z biblioteki; nie zna przepływu — ten stoi w `instalator.js`.

   Kolejność jest wiążąca: katalog treści musi być wczytany, zanim cokolwiek
   powstanie, bo składnik pyta o łańcuch w chwili budowania.

       DanacoKreator.zamontuj(korzen) → Promise */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;

function belka() {
  var znak = N.zeZnacznika(K.ikony.godlo);
  znak.setAttribute('class', 'dn-kreator-belka-znak');
  znak.setAttribute('aria-hidden', 'true');
  var przycisk = function (ikona, etykieta, klasa) {
    var s = N.zeZnacznika(K.ikony[ikona]);
    s.setAttribute('aria-hidden', 'true');
    return el('button', { klasa: klasa, type: 'button', 'aria-label': tekst(etykieta) }, [s]);
  };
  return el('div', { klasa: 'dn-kreator-belka' }, [
    znak,
    el('p', { klasa: 'dn-kreator-belka-tytul', tekst: tekst('okno.tytul') }),
    el('div', { klasa: 'dn-kreator-belka-sterowanie' }, [
      przycisk('zwin', 'okno.zwin', 'dn-kreator-belka-btn'),
      przycisk('rozwin', 'okno.rozwin', 'dn-kreator-belka-btn'),
      przycisk('zamknij', 'okno.zamknij', 'dn-kreator-belka-btn dn-kreator-belka-btn--zamknij')
    ])
  ]);
}

function szyna() {
  return el('div', { klasa: 'dn-kreator-szyna' }, [
    K.skladniki.nawigacjaKrokow({ klucz: 'szyna.kroki', biezacy: 1 }),
    el('div', {
      klasa: 'dn-bryla dn-kreator-bryla',
      role: 'img',
      'aria-label': tekst('bryla.etykieta'),
      dane: { bryla: true, 'kadr-gora': '0.44' }
    })
  ]);
}

function plotno() {
  var ekrany = [];
  for (var n = 1; n <= 6; n++) ekrany.push(K.ekrany[n]());
  return el('div', { klasa: 'dn-kreator-plotno' }, ekrany);
}

function dzialania() {
  return K.skladniki.pasDzialan({
    dodatkowa: { klucz: 'dzialania.kopiujSzczegoly', dane: { 'blad-kopiuj': true }, ukryty: true },
    poboczna: { klucz: 'dzialania.wstecz', dane: { 'krok-wstecz': true } },
    glowna: { klucz: 'dzialania.dalej', dane: { 'krok-dalej': true } }
  });
}

function dialogi() {
  var S = K.skladniki;
  var pomoc = el('div', { klasa: 'dn-tekst-ciagly' }, [
    el('ol', {}, tekst('dialogi.pomocProcesor.kroki').map(function (k) { return el('li', { tekst: k }); })),
    el('p', { tekst: tekst('dialogi.pomocProcesor.x64') }),
    el('p', { tekst: tekst('dialogi.pomocProcesor.arm') })
  ]);
  return [
    /* Jedno okno na dwa pytania o zamknięcie: treść zależy od tego, co
       instalator zdążył zrobić, a nie od tego, który przycisk ją wywołał. */
    S.oknoDialogowe({
      nazwa: 'potwierdzenie', rola: 'alertdialog', dane: { potwierdzenie: true },
      tytul: 'dialogi.zamkniecie.tytul', tresc: 'dialogi.zamkniecie.tresc',
      daneTytulu: { 'pot-tytul': true }, daneTresci: { 'pot-tresc': true },
      czynnosci: [
        { klucz: 'dialogi.zamkniecie.zostan', klasa: 'dn-btn dn-btn--zarys', dane: { 'potwierdzenie-nie': true, 'pot-zostan': true } },
        { klucz: 'dialogi.zamkniecie.wyjdz', klasa: 'dn-btn dn-btn--niebezpieczny', dane: { 'potwierdzenie-tak': true, 'pot-wyjdz': true } }
      ]
    }),
    S.oknoDialogowe({
      nazwa: 'niezgodnosc', rola: 'alertdialog',
      tytul: 'dialogi.niezgodnosc.tytul', tresc: 'dialogi.niezgodnosc.tresc',
      daneTresci: { 'nz-tresc': true },
      czynnosci: [
        { klucz: 'dialogi.niezgodnosc.popraw', klasa: 'dn-btn dn-btn--atrament', dane: { 'nz-popraw': true } },
        { klucz: 'dialogi.niezgodnosc.mimoTo', klasa: 'dn-btn dn-btn--zarys', dane: { 'nz-mimo-to': true } }
      ]
    }),
    S.oknoDialogowe({
      nazwa: 'pomoc-procesor', rola: 'dialog',
      tytul: 'dialogi.pomocProcesor.tytul', cialo: pomoc,
      czynnosci: [{ klucz: 'dialogi.pomocProcesor.zamknij', klasa: 'dn-btn dn-btn--atrament', dane: { 'pp-zamknij': true } }]
    })
  ];
}

K.zamontuj = function (korzen) {
  /* Katalog treści wpina się skryptem przed montażem. Jego brak to pomylona
     kolejność wpięć w oknie, nie sytuacja do obsłużenia po cichu. */
  if (!K.tresci) {
    korzen.textContent = '⟨brak katalogu treści — wepnij okna/instalator/tresci.js⟩';
    return null;
  }
  var okno = el('div', { klasa: 'dn-kreator' }, [
    belka(),
    el('div', { klasa: 'dn-kreator-korpus' }, [szyna(), plotno()]),
    dzialania()
  ]);
  korzen.appendChild(el('div', { klasa: 'dn-scena' }, [okno]));
  dialogi().forEach(function (d) { korzen.appendChild(d); });
  /* Bryła zakłada się samoczynnie na `DOMContentLoaded`, a wtedy tego pola
     jeszcze nie ma — po zmontowaniu trzeba ją założyć wprost. */
  if (window.DanacoBryla) window.DanacoBryla.zalozWszystkie(korzen);
  korzen.dispatchEvent(new CustomEvent('kreator-gotowy', { bubbles: true }));
  return okno;
};
})();

/* Samoczynny montaż w punkcie zaczepienia `[data-kreator]`. Dzięki temu plik
   podglądu nie zawiera ani jednej linii skryptu — wyłącznie wpięcia i miejsce,
   w którym okno ma stanąć. */
(function () {
  'use strict';
  function start() {
    var korzen = document.querySelector('[data-kreator]');
    if (korzen) window.DanacoKreator.zamontuj(korzen);
  }
  if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', start);
  else start();
})();
