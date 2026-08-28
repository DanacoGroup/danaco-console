/* Ekrany okna uruchomienia obejmują trzy odsłony tego samego okna: łączenie, powrót z ważnym tokenem oraz błąd połączenia, różniące się stanem czterech etapów.
   Leżą razem, bo dzielą całą oprawę — belkę, kolumnę tożsamości i wykaz
   czterech etapów. Różni je stan etapów i to, co stoi pod wykazem.

   Łączenie i powrót z tokenem nie mają czynności głównej: przechodzą dalej
   same, po ostatnim etapie. Błąd ją ma, bo tam jest co rozstrzygnąć.
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.ekrany = W.ekrany || {};

function panel(N, nazwa, aktywny, dzieci, dane) {
  var d = { widok: nazwa, 'grupa-widoku': 'wariant', 'widok-aktywny': aktywny ? 'tak' : 'nie' };
  if (dane) Object.keys(dane).forEach(function (k) { d[k] = dane[k]; });
  return N.el('div', {
    klasa: 'we-panel', id: nazwa, role: 'tabpanel', tabindex: '0',
    'aria-labelledby': 'zak-' + nazwa, dane: d
  }, dzieci);
}

W.ekrany.uruchomienie = function (N) {
  var S = W.skladniki;

  return [
    panel(N, 'w-laczenie', true, [
      S.naglowekEkranu(N, { nadtytul: 'uruchomienie.nadtytul', tytul: 'uruchomienie.laczenie.tytul', lid: 'uruchomienie.laczenie.lid' }),
      S.listaEtapow(N, {
        stany: ['gotowy', 'pracuje', 'oczekuje', 'oczekuje'],
        miary: ['nawiazane', 'wToku', 'oczekuje', 'oczekuje']
      })
    ], { 'po-polaczeniu': 'uwierzytelnienie' }),

    panel(N, 'w-token', false, [
      S.naglowekEkranu(N, { nadtytul: 'uruchomienie.nadtytul', tytul: 'uruchomienie.token.tytul', lid: 'uruchomienie.token.lid' }),
      S.listaEtapow(N, {
        stany: ['gotowy', 'gotowy', 'gotowy', 'pracuje'],
        miary: ['nawiazane', 'zgodna', 'zaufane', 'wToku']
      }),
      S.baner(N, { rodzaj: 'sukces', ikona: 'tarcza', glowa: 'uruchomienie.token.baner.glowa', tresc: 'uruchomienie.token.baner.tresc' }),
      S.frazaNawigacyjna(N, { klucz: 'uruchomienie.token.fraza' })
    ], { 'po-polaczeniu': 'przygotowanie' }),

    panel(N, 'w-blad', false, [
      S.naglowekEkranu(N, { nadtytul: 'uruchomienie.nadtytul', tytul: 'uruchomienie.blad.tytul', lid: 'uruchomienie.blad.lid' }),
      S.listaEtapow(N, {
        stany: ['blad', 'oczekuje', 'oczekuje', 'oczekuje'],
        miary: ['nieudane', 'oczekuje', 'oczekuje', 'oczekuje']
      }),
      S.baner(N, {
        rodzaj: 'ostrzezenie', ikona: 'ostrzezenie',
        glowa: 'uruchomienie.blad.baner.glowa', tresc: 'uruchomienie.blad.baner.tresc',
        odliczanieTresci: 15, postacOdliczania: 'sekundy'
      }),
      S.frazaNawigacyjna(N, { klucz: 'uruchomienie.blad.fraza' })
    ])
  ];
};

W.ekrany.uruchomieniePasy = function (N) {
  var S = W.skladniki;
  var zamknij = { klucz: 'dzialania.zamknijAplikacje', komunikat: 'zamkniecie' };
  return [
    S.pasDzialan(N, { widok: 'w-laczenie', grupa: 'wariant', aktywny: true, czynnosci: [zamknij] }),
    S.pasDzialan(N, { widok: 'w-token', grupa: 'wariant', czynnosci: [zamknij] }),
    S.pasDzialan(N, { widok: 'w-blad', grupa: 'wariant', czynnosci: [
      zamknij,
      { klucz: 'dzialania.ustawieniaPolaczenia', komunikat: 'ustawienia' },
      { klucz: 'dzialania.sprobujPonownie', klasa: 'dn-btn dn-btn--sygnal', id: 'btn-ponow', ikona: 'ponow' }
    ] })
  ];
};
})();
