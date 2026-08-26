/* ============================================================================
   EKRANY OKNA DOSTĘPU DO KONTA

   Siedem odsłon jednego okna: logowanie i jego odsłona z niepowodzeniem,
   zakładanie konta, potwierdzenie adresu, oraz trzy kroki odzyskiwania dostępu.
   Leżą razem, bo dzielą oprawę — belkę, kolumnę tożsamości i zakładki nad
   treścią.

   Zakładki niosą wyłącznie dwie drogi równorzędne: logowanie i rejestrację.
   Odzyskiwanie dostępu nie jest trzecią drogą — jest wyjściem z logowania,
   więc w zakładkach zostaje zaznaczone logowanie.
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.ekrany = W.ekrany || {};

function panel(N, nazwa, aktywny, pigulka, dzieci, nazwaDostepna) {
  return N.el('div', {
    klasa: 'we-panel', id: 's-' + nazwa, role: 'tabpanel', tabindex: '0',
    /* Obszar o roli `tabpanel` musi mieć nazwę — bez niej czytnik ekranu
       ogłasza „panel" i nic więcej. Zakładki rusztowania mają nieregularne
       identyfikatory, więc nazwa bierze się z tytułu odsłony, nie z nich. */
    'aria-label': nazwaDostepna ? N.tekst(nazwaDostepna) : null,
    dane: {
      widok: nazwa, 'grupa-widoku': 'stan',
      'widok-aktywny': aktywny ? 'tak' : 'nie',
      pigulka: pigulka
    }
  }, [W.skladniki.zakladkiPigulki(N, { wybrana: pigulka })].concat(dzieci));
}

W.ekrany.dostep = function (N) {
  var S = W.skladniki;
  /* Adres w zdaniu stoi krojem maszynowym — łatwiej porównać go ze skrzynką
     znak po znaku. Zdanie składa się więc z węzłów, nie z jednego łańcucha. */
  function lidKodu() {
    var czesci = N.tekst('dostep.kod.lid').split('{adres}');
    return [
      N.podstaw(czesci[0], N.tekst('dostep.kod.lidDane')),
      N.el('span', { klasa: 'au-kod-adres', tekst: N.tekst('dostep.kod.lidDane').adres }),
      N.podstaw(czesci[1] || '', N.tekst('dostep.kod.lidDane'))
    ];
  }

  function polaLogowania(przedrostek, blad) {
    return N.el('div', { klasa: 'au-pola' }, [
      S.poleTekstowe(N, {
        etykieta: 'dostep.logowanie.login', id: przedrostek + '-login', uzupelnij: 'username',
        bledne: blad, opisuje: blad ? 'blad-logowania' : null
      }),
      S.poleHasla(N, {
        etykieta: 'dostep.logowanie.haslo', id: przedrostek + '-haslo',
        opis: blad ? 'dostep.logowanieBlad.capsLock' : null,
        bledne: blad, opisuje: blad ? 'blad-logowania ' + przedrostek + '-haslo-opis' : null
      })
    ]);
  }

  function resetHasla() {
    return S.frazaNawigacyjna(N, {
      pytanie: 'dostep.logowanie.reset.pytanie', czynnosc: 'dostep.logowanie.reset.czynnosc',
      cel: 'odzyskiwanie-adres'
    });
  }

  return [
    /* ── logowanie ─────────────────────────────────────────────────────── */
    panel(N, 'logowanie', true, 'logowanie', [
      S.naglowekEkranu(N, { tytul: 'dostep.logowanie.tytul', lid: 'dostep.logowanie.lid' }),
      N.el('div', { klasa: 'we-komunikaty', 'aria-live': 'polite' }),
      polaLogowania('log', false),
      S.poleSesji(N),
      resetHasla()
    ].concat(S.metodyLogowania(N, { dostepne: ['email'] })).concat([
      S.frazaNawigacyjna(N, {
        pytanie: 'dostep.logowanie.fraza.pytanie', czynnosc: 'dostep.logowanie.fraza.czynnosc',
        cel: 'rejestracja'
      })
    ]), 'dostep.logowanie.tytul'),

    /* ── logowanie, dane nierozpoznane ─────────────────────────────────── */
    panel(N, 'logowanie-blad', false, 'logowanie', [
      S.naglowekEkranu(N, { tytul: 'dostep.logowanie.tytul' }),
      N.el('div', { klasa: 'we-komunikaty', 'aria-live': 'polite' }, [
        S.baner(N, {
          rodzaj: 'blad', ikona: 'ostrzezenie', id: 'blad-logowania',
          glowa: 'dostep.logowanieBlad.baner.glowa', tresc: 'dostep.logowanieBlad.baner.tresc'
        })
      ]),
      polaLogowania('blad', true),
      resetHasla()
    ].concat(S.metodyLogowania(N, { dostepne: ['email'] })).concat([
      S.frazaNawigacyjna(N, { klucz: 'dostep.logowanieBlad.fraza' })
    ]), 'dostep.logowanie.tytul'),

    /* ── zakładanie konta ──────────────────────────────────────────────── */
    panel(N, 'rejestracja', false, 'rejestracja', [
      S.naglowekEkranu(N, { tytul: 'dostep.rejestracja.tytul', lid: 'dostep.rejestracja.lid' }),
      N.el('div', { klasa: 'we-komunikaty', 'aria-live': 'polite' }),
      N.el('div', { klasa: 'au-pola' }, [
        N.el('div', { klasa: 'au-para' }, [
          S.poleTekstowe(N, { etykieta: 'dostep.rejestracja.login', id: 'rej-login', uzupelnij: 'username' }),
          S.poleTekstowe(N, { etykieta: 'dostep.rejestracja.email', id: 'rej-email', typ: 'email', uzupelnij: 'email' })
        ]),
        N.el('div', { klasa: 'au-para' }, [
          N.el('div', {}, [
            S.poleHasla(N, { etykieta: 'dostep.rejestracja.haslo', id: 'rej-haslo', uzupelnij: 'new-password' }),
            S.miernikSily(N, { dla: 'rej-haslo' })
          ]),
          S.poleHasla(N, { etykieta: 'dostep.rejestracja.hasloPowtorz', id: 'rej-haslo-2', uzupelnij: 'new-password' })
        ])
      ]),
      S.poleSesji(N),
      S.frazaNawigacyjna(N, {
        pytanie: 'dostep.rejestracja.fraza.pytanie', czynnosc: 'dostep.rejestracja.fraza.czynnosc',
        cel: 'logowanie'
      })
    ], 'dostep.rejestracja.tytul'),

    /* ── potwierdzenie adresu po założeniu konta ───────────────────────── */
    panel(N, 'kod', false, 'rejestracja', [
      S.naglowekEkranu(N, { tytul: 'dostep.kod.tytul', lidWezly: lidKodu() })
    ].concat(S.poleKodu(N, { odliczanie: '09:12', czynnosc: 'wklej' })).concat([
      N.el('div', { klasa: 'we-komunikaty', 'aria-live': 'polite' }, [
        S.baner(N, {
          rodzaj: 'informacja', ikona: 'informacja',
          glowa: 'dostep.kod.pomoc.glowa', tresc: 'dostep.kod.pomoc.tresc',
          dane: N.tekst('dostep.kod.pomocDane')
        })
      ]),
      S.frazaNawigacyjna(N, { czynnosc: 'dostep.kod.zmienAdres', cel: 'rejestracja' })
    ]), 'dostep.kod.tytul'),

    /* ── logowanie wstrzymane po pięciu nieudanych próbach ─────────────── */
    panel(N, 'logowanie-wstrzymane', false, 'logowanie', [
      S.naglowekEkranu(N, {
        tytul: 'dostep.logowanieWstrzymane.tytul',
        lid: 'dostep.logowanieWstrzymane.lid'
      }),
      S.baner(N, {
        rodzaj: 'ostrzezenie', ikona: 'zegar',
        glowa: 'dostep.logowanieWstrzymane.baner.glowa',
        tresc: 'dostep.logowanieWstrzymane.baner.tresc',
        daneGlowy: N.tekst('dostep.logowanieWstrzymane.banerDane')
      }),
      S.frazaNawigacyjna(N, {
        czynnosc: 'dostep.logowanieWstrzymane.odzyskaj', cel: 'odzyskiwanie-adres'
      })
    ], 'dostep.logowanieWstrzymane.tytul'),

    /* ── wysyłanie kodu wstrzymane po pięciu wysłaniach ─────────────────── */
    panel(N, 'odzyskiwanie-wstrzymane', false, 'logowanie', [
      S.naglowekEkranu(N, {
        tytul: 'dostep.odzyskiwanieWstrzymane.tytul',
        poTytule: [S.krokiOdzyskiwania(N, { biezacy: 1 })],
        lid: 'dostep.odzyskiwanieWstrzymane.lid'
      }),
      S.baner(N, {
        rodzaj: 'ostrzezenie', ikona: 'zegar',
        glowa: 'dostep.odzyskiwanieWstrzymane.baner.glowa',
        tresc: 'dostep.odzyskiwanieWstrzymane.baner.tresc',
        daneGlowy: N.tekst('dostep.odzyskiwanieWstrzymane.banerDane')
      }),
      S.frazaNawigacyjna(N, { czynnosc: 'dostep.odzyskiwanie.powrot', cel: 'logowanie' })
    ], 'dostep.odzyskiwanieWstrzymane.tytul'),

    /* ── odzyskiwanie: adres ───────────────────────────────────────────── */
    panel(N, 'odzyskiwanie-adres', false, 'logowanie', [
      S.naglowekEkranu(N, {
        tytul: 'dostep.odzyskiwanie.adres.tytul',
        poTytule: [S.krokiOdzyskiwania(N, { biezacy: 1 })],
        lid: 'dostep.odzyskiwanie.adres.lid'
      }),
      N.el('div', { klasa: 'au-pola' }, [
        S.poleTekstowe(N, { etykieta: 'dostep.odzyskiwanie.adres.pole', id: 'odz-email', typ: 'email', uzupelnij: 'email' })
      ]),
      N.el('div', { klasa: 'we-komunikaty', 'aria-live': 'polite' }, [
        S.baner(N, {
          rodzaj: 'informacja', ikona: 'tarcza',
          glowa: 'dostep.odzyskiwanie.adres.ostrzezenie.glowa',
          tresc: 'dostep.odzyskiwanie.adres.ostrzezenie.tresc'
        })
      ]),
      S.frazaNawigacyjna(N, { czynnosc: 'dostep.odzyskiwanie.powrot', cel: 'logowanie' })
    ], 'dostep.odzyskiwanie.adres.tytul'),

    /* ── odzyskiwanie: kod ─────────────────────────────────────────────── */
    panel(N, 'odzyskiwanie-kod', false, 'logowanie', [
      S.naglowekEkranu(N, {
        tytul: 'dostep.kod.tytul',
        poTytule: [S.krokiOdzyskiwania(N, { biezacy: 2 })],
        lidWezly: lidKodu()
      }),
      N.el('div', { klasa: 'we-komunikaty', 'aria-live': 'polite' })
    ].concat(S.poleKodu(N, { odliczanie: '08:41', czynnosc: 'ponow' })).concat([
      S.baner(N, {
        rodzaj: 'informacja', ikona: 'tarcza',
        glowa: 'dostep.kod.ostrzezenie.glowa', tresc: 'dostep.kod.ostrzezenie.tresc'
      }),
      S.frazaNawigacyjna(N, { czynnosc: 'dostep.kod.zmienAdres', cel: 'odzyskiwanie-adres' })
    ]), 'dostep.kod.tytul'),

    /* ── odzyskiwanie: nowe hasło ──────────────────────────────────────── */
    panel(N, 'odzyskiwanie-haslo', false, 'logowanie', [
      S.naglowekEkranu(N, {
        tytul: 'dostep.odzyskiwanie.haslo.tytul',
        poTytule: [S.krokiOdzyskiwania(N, { biezacy: 3 })],
        lid: 'dostep.odzyskiwanie.haslo.lid'
      }),
      N.el('div', { klasa: 'we-komunikaty', 'aria-live': 'polite' }),
      N.el('div', { klasa: 'au-pola' }, [
        N.el('div', { klasa: 'au-para' }, [
          S.poleHasla(N, { etykieta: 'dostep.odzyskiwanie.haslo.nowe', id: 'odz-haslo', uzupelnij: 'new-password' }),
          S.poleHasla(N, { etykieta: 'dostep.odzyskiwanie.haslo.powtorz', id: 'odz-haslo-2', uzupelnij: 'new-password' }),
          S.miernikSily(N, { dla: 'odz-haslo' })
        ])
      ]),
      S.poleSesji(N),
      S.frazaNawigacyjna(N, { czynnosc: 'dostep.odzyskiwanie.powrot', cel: 'logowanie' })
    ], 'dostep.odzyskiwanie.haslo.tytul')
  ];
};

W.ekrany.dostepPasy = function (N) {
  var S = W.skladniki;
  var zamknij = { klucz: 'dzialania.zamknijAplikacje', komunikat: 'zamkniecie' };
  var glowna = 'dn-btn dn-btn--sygnal';
  /* Czynność główna prowadzi albo do kolejnego etapu, albo do kolejnej odsłony
     tego samego okna. Bez celu przycisk nie dostaje uchwytu `data-idz` i jest
     martwy — cały tor odzyskiwania dostępu stał tak od pierwszego kroku. */
  function pas(widok, aktywny, klucz, cel, grupa) {
    return S.pasDzialan(N, {
      widok: widok, grupa: 'stan', aktywny: aktywny,
      czynnosci: [zamknij, { klucz: klucz, klasa: glowna, cel: cel, grupa: grupa || 'etap' }]
    });
  }
  return [
    pas('logowanie', true, 'dzialania.zaloguj', 'przygotowanie'),
    pas('logowanie-blad', false, 'dzialania.zalogujPonownie', 'przygotowanie'),
    pas('rejestracja', false, 'dzialania.utworzKonto', 'kod', 'stan'),
    pas('kod', false, 'dzialania.potwierdzKonto', 'przygotowanie'),
    pas('odzyskiwanie-adres', false, 'dzialania.wyslijKod', 'odzyskiwanie-kod', 'stan'),
    pas('odzyskiwanie-kod', false, 'dzialania.potwierdzKod', 'odzyskiwanie-haslo', 'stan'),
    pas('odzyskiwanie-haslo', false, 'dzialania.potwierdzHaslo', 'logowanie', 'stan'),
    /* Odsłony wstrzymania nie mają czynności głównej — nie ma czego wykonać,
       póki godzina nie minie. Zostaje wyjście z programu. */
    S.pasDzialan(N, { widok: 'logowanie-wstrzymane', grupa: 'stan', czynnosci: [zamknij] }),
    S.pasDzialan(N, { widok: 'odzyskiwanie-wstrzymane', grupa: 'stan', czynnosci: [zamknij] })
  ];
};
})();
