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

function panel(N, nazwa, aktywny, pigulka, dzieci) {
  return N.el('div', {
    klasa: 'we-panel', id: 's-' + nazwa, role: 'tabpanel', tabindex: '0',
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
    ])),

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
    ])),

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
    ]),

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
    ])),

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
    ]),

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
    ])),

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
    ])
  ];
};

W.ekrany.dostepPasy = function (N) {
  var S = W.skladniki;
  var zamknij = { klucz: 'dzialania.zamknijAplikacje', komunikat: 'zamkniecie' };
  var glowna = 'dn-btn dn-btn--sygnal';
  function pas(widok, aktywny, klucz, cel) {
    return S.pasDzialan(N, {
      widok: widok, grupa: 'stan', aktywny: aktywny,
      czynnosci: [zamknij, { klucz: klucz, klasa: glowna, cel: cel, grupa: 'etap' }]
    });
  }
  return [
    pas('logowanie', true, 'dzialania.zaloguj', 'przygotowanie'),
    pas('logowanie-blad', false, 'dzialania.zalogujPonownie', 'przygotowanie'),
    pas('rejestracja', false, 'dzialania.utworzKonto', 'przygotowanie'),
    pas('kod', false, 'dzialania.potwierdzKonto', 'przygotowanie'),
    pas('odzyskiwanie-adres', false, 'dzialania.wyslijKod', null),
    pas('odzyskiwanie-kod', false, 'dzialania.potwierdzKod', null),
    pas('odzyskiwanie-haslo', false, 'dzialania.potwierdzHaslo', null)
  ];
};
})();
