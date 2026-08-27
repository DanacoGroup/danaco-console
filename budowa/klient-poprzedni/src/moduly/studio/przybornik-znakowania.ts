import './przybornik-znakowania.css';

import {
  StudioAuthor,
  StudioMarkupKind,
  StudioMarkupState,
  type StudioMarkup,
  type StudioMarkupType,
} from '../../../../shared/contract';
import { poleTekstowe, poleWyboru } from '../../modele/kontrolki-formularza';
import {
  NAZWY_RODZAJOW,
  RodzajZnakowania,
  przybornikOpiszWykaz,
  przybornikPrzefiltruj,
  przybornikFiltrPelny,
  type FiltrZnakowan,
  type PozycjaZnakowania,
} from './przybornik-wykaz';
import {
  BARWY_ZNACZNIKOW,
  NAZWY_ZNACZNIKOW_GOTOWE,
  przybornikOpiszZnakowanieRdzenia,
  przybornikZetonBarwy,
} from './przybornik-znaczniki';
import type { ZrodloKontroliStudio } from './zrodlo-kontroli-studio';


/** Czynności przybornika zlecane oknu: znakowanie, propozycja, decyzja i wskazanie miejsca w treści dokumentu. */
export interface CzynnosciPrzybornika {
  /** Zakłada komentarz przypięty do zaznaczenia. */
  naKomentarz(tresc: string): void;
  /** Zleca modelowi propozycję brzmienia dla zaznaczenia — bez wejścia w treść. */
  naPropozycje(polecenie: string): void;
  /** Zakłada adnotację przy wskazanym fragmencie różnicy. */
  naAdnotacje(numerFragmentu: number, tresc: string): void;
  /** Zakłada znacznik własny o podanej nazwie i barwie. */
  naZnacznik(nazwa: string, barwa: string): void;
  /** Przestawia odhaczenie pozycji wykazu; dotyczy znaczników i wątków. */
  naOdhaczenie(kod: string, rodzaj: RodzajZnakowania, zamknij: boolean): void;
  /** Zdejmuje znacznik własny. */
  naZdjecieZnacznika(kod: string): void;
  /** Wyróżnia zaznaczenie barwą wskazaną z palety. */
  naWyroznienie(barwa: string): void;
  /** Zakłada zakładkę w miejscu kursora, do powrotu. */
  naZakladke(nazwa: string): void;
  /** Skacze do wskazanej pozycji wykazu — przewija treść do jej miejsca. */
  naSkok(pozycja: PozycjaZnakowania): void;
}

/** Przybornik znakowania wraz z jego odświeżeniem, elementem osadzanym w oknie i zdaniem o zapleczu rdzenia. */
export interface PrzybornikZnakowania {
  element: HTMLElement;
  /** Przerysowuje wykaz znakowań wraz z zawężeniem. */
  pokaz(pozycje: readonly PozycjaZnakowania[]): void;
  /** Mówi przybornikowi, ile znaków jest zaznaczonych. */
  ustawZaznaczenie(dlugosc: number): void;
  /** Wstawia zakładki sesji do wykazu powrotów. */
  ustawZakladki(zakladki: readonly { kod: string; nazwa: string }[]): void;
  /** Skacze do pozycji następnej albo poprzedniej wykazu widocznego. */
  skocz(wPrzod: boolean): PozycjaZnakowania | null;
  /** Wypisuje odpowiedź rdzenia — powodzenie albo odmowę nazwaną. */
  pokazOdpowiedz(tresc: string, udana: boolean): void;
  /** Bez podanego zaplecza rdzenia funkcja nic nie robi i mówi to wprost Operatorowi. */
  odswiezZnakowaniaRdzenia(): Promise<void>;
  przestawWidocznosc(): void;
  widoczny(): boolean;
}

/** Zaplecze znakowania w rdzeniu, siedem komend rodziny studio.markup — null znaczy, że drogi do rdzenia jeszcze nie wpięto. */
export interface RdzenZnakowania {
  zrodlo: ZrodloKontroliStudio;
  /** Dokument czynny; puste znaczy „nie ma na czym znakować". */
  idDokumentu(): string;
  /** Zaznaczenie w treści; `null` znaczy brak zaznaczenia. */
  zaznaczenie(): { poczatek: number; koniec: number } | null;
  /** Przyjmuje postać dokumentu oddaną po znakowaniu, gdy rdzeń ją zwrócił. */
  naSkutek(tresc: string | undefined, postac: unknown): void;
  /** Przewija treść do miejsca znakowania. */
  naMiejsce(poczatek: number, koniec: number): void;
}

/** Czynności, których rdzeń nie ma, wraz z nazwaniem braku — wykaz jawny i widoczny w oknie, bo ukrycie tych pozycji sprawiłoby, że przybornik wygląda na kompletny. */
export const BRAKI_PRZYBORNIKA: readonly { nazwa: string; czego: string }[] = [
  {
    nazwa: 'Przypis dolny i końcowy',
    czego:
      'Aparat dokumentu ma dziś w rdzeniu własną rodzinę komend, ale droga z tego przybornika do ' +
      'niej należy do odcinka aparatu i wstawiania, nie do odcinka kontroli pracy. Przypis ' +
      'wstawiony jako sam tekst w nawiasie nie przenumerowałby się po wstawieniu przypisu przed ' +
      'nim, więc nie jest przypisem — i dlatego przybornik go nie udaje.',
  },
  {
    nazwa: 'Odsyłacz i odwołanie wzajemne',
    czego:
      'Odwołanie musi wskazywać byt dokumentu (nagłówek, tabelę, ilustrację) i przetrwać jego ' +
      'przeniesienie. Droga do aparatu dokumentu z przybornika jeszcze nie stoi, więc odsyłacz ' +
      'byłby tu napisem, a nie odwołaniem.',
  },
  {
    nazwa: 'Wstawienie tabeli i pola',
    czego:
      'Tabela i pole obliczane są postacią dokumentu i mają w rdzeniu swoje komendy, ale ich ' +
      'wołacz stoi w odcinku tabel i obiektów. Tabela rysowana znakami w treści zniknęłaby przy ' +
      'pierwszym wydaniu do docx.',
  },
  {
    nazwa: 'Zakładka w miejscu trwała',
    czego:
      'Zakładka należy do aparatu dokumentu. Powrót do miejsca działa przez sesję okna; zapis ' +
      'zakładki w dokumencie idzie komendą aparatu, której ten przybornik nie woła.',
  },
];

/** Zakłada przybornik; dyktafon osadza przycisk dyktowania do treści dokumentu, null znaczy, że stanowisko go nie ma. */
export function utworzPrzybornikZnakowania(
  czynnosci: CzynnosciPrzybornika,
  dyktafon: HTMLElement | null,
  rdzen: RdzenZnakowania | null = null,
): PrzybornikZnakowania {
  let otwarty = true;
  let filtr: FiltrZnakowan = przybornikFiltrPelny();
  let wszystkie: readonly PozycjaZnakowania[] = [];
  let widocznePozycje: PozycjaZnakowania[] = [];
  let wskazana = -1;
  let dlugoscZaznaczenia = 0;

  /* ── Czynności na fragmencie ─────────────────────────────────────────────── */

  const trescZnakowania = document.createElement('input');
  trescZnakowania.type = 'text';
  trescZnakowania.className = 'dn-pole-kontrolka';
  trescZnakowania.placeholder = 'treść komentarza, propozycji albo adnotacji';
  trescZnakowania.setAttribute('aria-label', 'Treść znakowania zakładanego na fragmencie');

  const zakresZdanie = document.createElement('p');
  zakresZdanie.className = 'dn-pole-opis ms-przybornik__zakres';

  const komentarz = przyciskCzynnosci(
    'Komentarz do fragmentu',
    'komentarz',
    'studio.comment.add — komentarz NIE niesie brzmienia: mówi o fragmencie, treści nie zmienia.',
  );
  komentarz.addEventListener('click', () => czynnosci.naKomentarz(trescZnakowania.value));

  const propozycja = przyciskCzynnosci(
    'Propozycja brzmienia',
    'propozycja',
    'studio.contextual.op — model oddaje brzmienie fragmentu, które staje NA MARGINESIE. ' +
      'W treści go jeszcze nie ma; Operator je przyjmuje, odrzuca albo poprawia ' +
      '(studio.proposal.decide).',
  );
  propozycja.addEventListener('click', () => czynnosci.naPropozycje(trescZnakowania.value));

  const numerFragmentu = document.createElement('input');
  numerFragmentu.type = 'number';
  numerFragmentu.min = '0';
  numerFragmentu.value = '0';
  numerFragmentu.className = 'dn-pole-kontrolka ms-przybornik__numer';
  numerFragmentu.setAttribute('aria-label', 'Numer fragmentu różnicy, którego adnotacja dotyczy');

  const adnotacja = przyciskCzynnosci(
    'Adnotacja przy fragmencie różnicy',
    'adnotacja',
    'studio.annotation.add — droga własna komendy, nie generyczna window.action. Adnotacja wisi ' +
      'przy NUMERZE fragmentu porównania, nie przy znaku treści bieżącej.',
  );
  adnotacja.addEventListener('click', () =>
    czynnosci.naAdnotacje(Number(numerFragmentu.value), trescZnakowania.value),
  );

  const nazwaZnacznika = document.createElement('input');
  nazwaZnacznika.type = 'text';
  nazwaZnacznika.className = 'dn-pole-kontrolka';
  nazwaZnacznika.placeholder = 'nazwa znacznika, np. „wymaga źródła"';
  nazwaZnacznika.setAttribute('list', 'ms-przybornik-nazwy');
  nazwaZnacznika.setAttribute('aria-label', 'Nazwa znacznika własnego');

  const gotoweNazwy = document.createElement('datalist');
  gotoweNazwy.id = 'ms-przybornik-nazwy';
  for (const nazwa of NAZWY_ZNACZNIKOW_GOTOWE) {
    const wpis = document.createElement('option');
    wpis.value = nazwa;
    gotoweNazwy.append(wpis);
  }

  const paleta = poleWyboru(
    {
      etykieta: 'Barwa znacznika i wyróżnienia',
      opis: 'Barwy są żetonami motywu, nie wartościami zapisanymi wprost — oba motywy obsługują się same.',
    },
    BARWY_ZNACZNIKOW.map((barwa) => ({ wartosc: barwa.kod, etykieta: barwa.nazwa })),
  );

  const znacznik = przyciskCzynnosci(
    'Znacznik własny',
    'znacznik',
    'Nazwa i barwa nadana fragmentowi. Kontrakt nie ma na to pola, więc znacznik żyje przez ' +
      'sesję okna — wykaz, filtr i odhaczanie działają, trwałości nie ma.',
  );
  znacznik.addEventListener('click', () =>
    czynnosci.naZnacznik(nazwaZnacznika.value, paleta.kontrolka.value),
  );

  const wyroznienie = przyciskCzynnosci(
    'Wyróżnij fragment barwą',
    'wyroznienie',
    'Wyróżnienie działa W OKNIE. Bez komendy postaci dokumentu NIE dojeżdża do rdzenia i ginie ' +
      'przy zapisie — to jest brak nazwany, nie funkcja.',
  );
  wyroznienie.addEventListener('click', () => czynnosci.naWyroznienie(paleta.kontrolka.value));

  const nazwaZakladki = document.createElement('input');
  nazwaZakladki.type = 'text';
  nazwaZakladki.className = 'dn-pole-kontrolka';
  nazwaZakladki.placeholder = 'nazwa zakładki do powrotu';
  nazwaZakladki.setAttribute('aria-label', 'Nazwa zakładki w miejscu kursora');

  const zakladka = przyciskCzynnosci(
    'Zakładka w miejscu',
    'zakladka',
    'Powrót do miejsca w treści. Działa przez sesję okna; zapisu zakładki w dokumencie kontrakt ' +
      'nie niesie.',
  );
  zakladka.addEventListener('click', () => czynnosci.naZakladke(nazwaZakladki.value));

  const wykazZakladek = document.createElement('ul');
  wykazZakladek.className = 'ms-przybornik__zakladki';

  /* ── Wykaz znakowań i jego zawężenie ─────────────────────────────────────── */

  const filtrRodzaju = poleWyboru({ etykieta: 'Rodzaj znakowania' }, [
    { wartosc: 'wszystkie', etykieta: 'Wszystkie rodzaje' },
    ...Object.values(RodzajZnakowania).map((rodzaj) => ({
      wartosc: rodzaj,
      etykieta: NAZWY_RODZAJOW[rodzaj].nazwa,
    })),
  ]);

  const filtrAutora = poleWyboru({ etykieta: 'Autor' }, [
    { wartosc: 'wszyscy', etykieta: 'Operator i model' },
    { wartosc: StudioAuthor.Uzytkownik, etykieta: 'Tylko Operator' },
    { wartosc: StudioAuthor.Model, etykieta: 'Tylko model' },
  ]);

  const filtrStanu = poleWyboru({ etykieta: 'Stan' }, [
    { wartosc: 'wszystkie', etykieta: 'Otwarte i zamknięte' },
    { wartosc: 'otwarte', etykieta: 'Tylko otwarte' },
    { wartosc: 'zamkniete', etykieta: 'Tylko rozwiązane i odhaczone' },
  ]);

  const podsumowanie = document.createElement('p');
  podsumowanie.className = 'dn-pole-opis ms-przybornik__podsumowanie';

  const wykaz = document.createElement('ul');
  wykaz.className = 'ms-przybornik__wykaz';

  const poprzednie = document.createElement('button');
  poprzednie.type = 'button';
  poprzednie.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  poprzednie.textContent = 'Poprzednie znakowanie';
  poprzednie.dataset['czynnosc'] = 'skok-wstecz';
  poprzednie.addEventListener('click', () => skocz(false));

  const nastepne = document.createElement('button');
  nastepne.type = 'button';
  nastepne.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  nastepne.textContent = 'Następne znakowanie';
  nastepne.dataset['czynnosc'] = 'skok-wprzod';
  nastepne.addEventListener('click', () => skocz(true));

  const odpowiedzPrzybornika = document.createElement('p');
  odpowiedzPrzybornika.className = 'dn-pole-opis ms-przybornik__odpowiedz';

  function odpowiedzRdzenia(tresc: string, udana: boolean): void {
    odpowiedzPrzybornika.textContent = tresc;
    odpowiedzPrzybornika.dataset['udana'] = udana ? 'tak' : 'nie';
  }

  /* ── Znakowanie TRWAŁE — siedem komend rodziny studio.markup.* ───────────── */

  const rodzajRdzenia = poleWyboru(
    {
      etykieta: 'Rodzaj znakowania trwałego',
      opis:
        'Trzy rodzaje rozłączne: wyróżnienie barwą, znacznik własny Operatora i propozycja zmiany ' +
        'na marginesie wraz z brzmieniem. Propozycja NIE wchodzi w treść, dopóki Operator jej nie ' +
        'przyjmie.',
    },
    [
      { wartosc: StudioMarkupKind.Highlight, etykieta: 'Wyróżnienie barwą' },
      { wartosc: StudioMarkupKind.Mark, etykieta: 'Znacznik własny' },
      { wartosc: StudioMarkupKind.Suggestion, etykieta: 'Propozycja zmiany na marginesie' },
    ],
  );

  const rodzajZnacznikaRdzenia = poleWyboru(
    {
      etykieta: 'Rodzaj znacznika własnego',
      opis: 'Wykaz idzie z rdzenia — rodzaje zakłada Operator, a nie są wpisane w kod okna.',
    },
    [{ wartosc: '', etykieta: 'Rodzaje nie zostały jeszcze odczytane z rdzenia' }],
  );

  const brzmienieRdzenia = poleTekstowe({
    etykieta: 'Brzmienie proponowane',
    podpowiedz: 'brzmienie, które ma stanąć na marginesie',
    opis:
      'Dotyczy wyłącznie propozycji. Brzmienie wchodzi do treści dopiero przy jej PRZYJĘCIU — ' +
      'i wtedy jako zmiana śledzona, którą można cofnąć.',
  });

  const zalozZnakowanie = przyciskCzynnosci(
    'Załóż znakowanie trwałe na zaznaczeniu',
    'znakowanie-rdzenia',
    'studio.markup.add — znakowanie leży w rdzeniu, więc przetrwa zamknięcie karty i widzi je ' +
      'model. Znakowanie modelu jest podpisane jako model.',
  );
  zalozZnakowanie.addEventListener('click', () => {
    void zalozZnakowanieRdzenia();
  });

  const odczytajZnakowania = przyciskCzynnosci(
    'Odczytaj znakowania z rdzenia',
    'znakowania-odczyt',
    'studio.markup.list — wykaz znakowań dokumentu jako spis do przejścia, wraz z komentarzami ' +
      'i zmianami śledzonymi.',
  );
  odczytajZnakowania.addEventListener('click', () => {
    void odswiezZnakowaniaRdzenia();
  });

  const nazwaRodzaju = poleTekstowe({
    etykieta: 'Nazwa rodzaju znacznika do założenia albo zmiany',
    podpowiedz: 'na przykład „wymaga źródła"',
  });

  const zapiszRodzaj = przyciskCzynnosci(
    'Zapisz rodzaj znacznika',
    'rodzaj-zapisz',
    'studio.markup.type.save — rodzaj wraz z nazwą i barwą. Rodzaje są sprawą Operatora, nie ' +
      'wykazu wpisanego w kod.',
  );
  zapiszRodzaj.addEventListener('click', () => {
    void zapiszRodzajZnacznika();
  });

  const usunRodzaj = przyciskCzynnosci(
    'Usuń rodzaj znacznika',
    'rodzaj-usun',
    'studio.markup.type.delete — rodzaju FABRYCZNEGO rdzeń nie usunie i odpowie odmową nazywającą ' +
      'powód. Okno tej odmowy nie uprzedza martwą kontrolką.',
  );
  usunRodzaj.addEventListener('click', () => {
    void usunRodzajZnacznika();
  });

  const podsumowanieRdzenia = document.createElement('p');
  podsumowanieRdzenia.className = 'dn-pole-opis';
  podsumowanieRdzenia.textContent =
    rdzen === null
      ? 'Droga do znakowania trwałego NIE jest w tym oknie wpięta: rdzeń niesie siedem komend ' +
        'rodziny studio.markup.*, a to okno ich jeszcze nie woła. Znacznik założony niżej żyje ' +
        'przez sesję okna. Brak jest po stronie montażu okna, nie po stronie rdzenia.'
      : 'Znakowania trwałe nie zostały jeszcze odczytane z rdzenia.';

  const wykazRdzenia = document.createElement('ul');
  wykazRdzenia.className = 'ms-przybornik__wykaz';

  /* ── Legenda trzech bytów marginesu ──────────────────────────────────────── */

  const legenda = document.createElement('ul');
  legenda.className = 'ms-przybornik__legenda';
  for (const rodzaj of Object.values(RodzajZnakowania)) {
    const nazwa = document.createElement('strong');
    nazwa.textContent = NAZWY_RODZAJOW[rodzaj].nazwa;
    const czym = document.createElement('span');
    czym.className = 'dn-pole-opis';
    czym.textContent = NAZWY_RODZAJOW[rodzaj].czym;
    const pozycja = document.createElement('li');
    pozycja.dataset['rodzaj'] = rodzaj;
    pozycja.append(nazwa, czym);
    legenda.append(pozycja);
  }

  const brakiWykaz = document.createElement('ul');
  brakiWykaz.className = 'ms-przybornik__braki';
  for (const brak of BRAKI_PRZYBORNIKA) {
    const nazwa = document.createElement('strong');
    nazwa.textContent = brak.nazwa;
    const czego = document.createElement('span');
    czego.className = 'dn-pole-opis';
    czego.textContent = brak.czego;
    const pozycja = document.createElement('li');
    pozycja.append(nazwa, czego);
    brakiWykaz.append(pozycja);
  }

  const braki = document.createElement('details');
  braki.className = 'ms-przybornik__czesc';
  const podpisBrakow = document.createElement('summary');
  podpisBrakow.textContent = `Bez zaplecza w rdzeniu (${BRAKI_PRZYBORNIKA.length}) — nazwane wprost`;
  braki.append(podpisBrakow, brakiWykaz);

  const element = document.createElement('aside');
  element.className = 'ms-przybornik';
  element.setAttribute('aria-label', 'Przybornik znakowania dokumentu');
  element.append(
    czescPrzybornika('Znakowanie fragmentu', [
      zakresZdanie,
      trescZnakowania,
      komentarz,
      propozycja,
      numerFragmentu,
      adnotacja,
    ]),
    czescPrzybornika('Znacznik, wyróżnienie, zakładka', [
      nazwaZnacznika,
      gotoweNazwy,
      paleta.element,
      znacznik,
      wyroznienie,
      nazwaZakladki,
      zakladka,
      wykazZakladek,
    ]),
    czescPrzybornika('Wykaz znakowań do przejścia', [
      filtrRodzaju.element,
      filtrAutora.element,
      filtrStanu.element,
      podsumowanie,
      poprzednie,
      nastepne,
      wykaz,
    ]),
    czescPrzybornika('Znakowanie TRWAŁE — leży w rdzeniu, widzi je model', [
      podsumowanieRdzenia,
      rodzajRdzenia.element,
      rodzajZnacznikaRdzenia.element,
      brzmienieRdzenia.element,
      zalozZnakowanie,
      odczytajZnakowania,
      nazwaRodzaju.element,
      zapiszRodzaj,
      usunRodzaj,
      wykazRdzenia,
    ]),
    czescPrzybornika('Trzy byty marginesu — czym się różnią', [legenda]),
    braki,
    odpowiedzPrzybornika,
  );

  if (dyktafon !== null) {
    const zdanieDyktafonu = document.createElement('p');
    zdanieDyktafonu.className = 'dn-pole-opis';
    zdanieDyktafonu.textContent =
      'Podyktowany tekst wchodzi w miejsce kursora jako treść dokumentu — nie jako polecenie dla ' +
      'modelu. Nagranie idzie do magazynu rdzenia i go nie opuszcza.';
    element.insertBefore(
      czescPrzybornika('Dyktowanie do treści', [dyktafon, zdanieDyktafonu]),
      braki,
    );
  }

  for (const kontrolka of [filtrRodzaju.kontrolka, filtrAutora.kontrolka, filtrStanu.kontrolka]) {
    kontrolka.addEventListener('change', () => {
      filtr = {
        rodzaj: filtrRodzaju.kontrolka.value as FiltrZnakowan['rodzaj'],
        autor: filtrAutora.kontrolka.value as FiltrZnakowan['autor'],
        stan: filtrStanu.kontrolka.value as FiltrZnakowan['stan'],
      };
      przerysuj();
    });
  }

  function przerysuj(): void {
    widocznePozycje = przybornikPrzefiltruj(wszystkie, filtr);
    if (wskazana >= widocznePozycje.length) wskazana = widocznePozycje.length - 1;
    podsumowanie.textContent = przybornikOpiszWykaz(wszystkie, widocznePozycje);
    wykaz.replaceChildren(
      ...widocznePozycje.map((pozycja, numer) =>
        wierszZnakowania(pozycja, numer === wskazana, czynnosci),
      ),
    );
  }

  function skocz(wPrzod: boolean): PozycjaZnakowania | null {
    if (widocznePozycje.length === 0) {
      odpowiedzPrzybornika.textContent =
        'Nie ma po czym skakać: wykaz po zawężeniu jest pusty. Zdejmij zawężenie albo zaznacz ' +
        'fragment i założ znakowanie.';
      return null;
    }
    wskazana = wPrzod
      ? (wskazana + 1) % widocznePozycje.length
      : (wskazana - 1 + widocznePozycje.length) % widocznePozycje.length;
    const pozycja = widocznePozycje[wskazana];
    if (pozycja === undefined) return null;
    przerysuj();
    czynnosci.naSkok(pozycja);
    return pozycja;
  }

  /* ── Czynności znakowania trwałego ───────────────────────────────────────── */

  /** Zaplecze wraz z dokumentem albo odpowiedź nazywająca brak. */
  function zapleczeRdzenia(): { rdzen: RdzenZnakowania; idDokumentu: string } | null {
    if (rdzen === null) {
      odpowiedzRdzenia(
        'Znakowania trwałego nie ma czym założyć: droga do rodziny studio.markup.* nie jest w tym ' +
          'oknie wpięta. Rdzeń te komendy niesie — brak jest po stronie montażu okna.',
        false,
      );
      return null;
    }
    const kod = rdzen.idDokumentu();
    if (kod === '') {
      odpowiedzRdzenia(
        'Nie ma dokumentu czynnego, więc nie ma czego znakować. Otwórz dokument albo załóż nowy.',
        false,
      );
      return null;
    }
    return { rdzen, idDokumentu: kod };
  }

  async function zalozZnakowanieRdzenia(): Promise<void> {
    const zaplecze = zapleczeRdzenia();
    if (zaplecze === null) return;
    const zakres = zaplecze.rdzen.zaznaczenie();
    if (zakres === null || zakres.koniec <= zakres.poczatek) {
      odpowiedzRdzenia(
        'Znakowanie trwałe dotyczy FRAGMENTU: rdzeń wymaga zakresu w znakach. Zaznacz fragment ' +
          'w treści — inaczej nie ma czego znakować.',
        false,
      );
      return;
    }
    const rodzaj = rodzajRdzenia.kontrolka.value as StudioMarkupKind;
    const wynik = await zaplecze.rdzen.zrodlo.kontrolaZnakowanieDodaj(
      zaplecze.idDokumentu,
      rodzaj,
      zakres,
      {
        barwa: paleta.kontrolka.value,
        rodzajZnacznika: rodzajZnacznikaRdzenia.kontrolka.value,
        uzasadnienie: trescZnakowania.value.trim(),
        brzmienie: brzmienieRdzenia.kontrolka.value.trim(),
      },
      {},
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedzRdzenia(`Znakowania nie założono: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    // Propozycja na fragmencie zablokowanym przechodzi — droga do fragmentu, którego nie wolno tknąć.
    const oBlokadzie =
      rodzaj === StudioMarkupKind.Suggestion
        ? ' Propozycja przechodzi także na fragmencie zablokowanym — treści nie zmienia, więc ' +
          'blokady nie omija; wnosi ją dopiero Twoje przyjęcie.'
        : '';
    odpowiedzRdzenia(
      `Znakowanie założone w rdzeniu: ${przybornikOpiszZnakowanieRdzenia(wynik.wynik.markup)}.` +
        oBlokadzie,
      true,
    );
    if (wynik.wynik.form !== undefined) zaplecze.rdzen.naSkutek(undefined, wynik.wynik.form);
    await odswiezZnakowaniaRdzenia();
  }

  async function odswiezZnakowaniaRdzenia(): Promise<void> {
    if (rdzen === null) return;
    const kod = rdzen.idDokumentu();
    if (kod === '') return;
    const [znakowania, rodzaje] = await Promise.all([
      rdzen.zrodlo.kontrolaZnakowaniaWykaz(kod, {
        ...(filtr.autor === 'wszyscy' ? {} : { autor: filtr.autor }),
        ...(filtr.stan === 'otwarte' ? { stan: StudioMarkupState.Open } : {}),
        zKomentarzami: false,
        zeZmianami: false,
      }),
      rdzen.zrodlo.kontrolaRodzajeZnacznikow(kod),
    ]);

    if (!znakowania.udany || znakowania.wynik === undefined) {
      podsumowanieRdzenia.textContent =
        `Znakowań trwałych nie udało się odczytać: ${znakowania.blad?.message ?? ''}`;
      wykazRdzenia.replaceChildren();
    } else {
      const tresc = znakowania.wynik;
      podsumowanieRdzenia.textContent =
        tresc.total === 0
          ? 'Dokument nie ma ani jednego znakowania trwałego. Pusty wykaz znaczy tu „nic nie ' +
            'oznaczono", a nie „znakowanie nie działa".'
          : `Znakowań trwałych: ${tresc.total}, z tego czeka na przejście ${tresc.openCount}.`;
      wykazRdzenia.replaceChildren(...tresc.markups.map(wierszZnakowaniaRdzenia));
    }

    if (rodzaje.udany && rodzaje.wynik !== undefined) {
      ustawRodzajeZnacznikow(rodzaje.wynik.markupTypes);
    }
  }

  /** Wstawia rodzaje znaczników oddane przez rdzeń do listy wyboru. */
  function ustawRodzajeZnacznikow(rodzaje: readonly StudioMarkupType[]): void {
    const wybrany = rodzajZnacznikaRdzenia.kontrolka.value;
    rodzajZnacznikaRdzenia.kontrolka.replaceChildren();
    const pozycje: readonly { wartosc: string; etykieta: string }[] =
      rodzaje.length === 0
        ? [{ wartosc: '', etykieta: 'Rdzeń nie ma ani jednego rodzaju — założ go polem niżej' }]
        : rodzaje.map((rodzaj) => ({
            wartosc: rodzaj.name,
            etykieta:
              `${rodzaj.label === '' ? rodzaj.name : rodzaj.label}` +
              `${rodzaj.builtin === true ? ' (fabryczny)' : ''}` +
              `${rodzaj.usageCount === undefined ? '' : ` — użyć: ${rodzaj.usageCount}`}`,
          }));
    for (const pozycja of pozycje) {
      const wpis = document.createElement('option');
      wpis.value = pozycja.wartosc;
      wpis.textContent = pozycja.etykieta;
      rodzajZnacznikaRdzenia.kontrolka.append(wpis);
    }
    if (pozycje.some((pozycja) => pozycja.wartosc === wybrany)) {
      rodzajZnacznikaRdzenia.kontrolka.value = wybrany;
    }
  }

  async function zapiszRodzajZnacznika(): Promise<void> {
    if (rdzen === null) {
      odpowiedzRdzenia(
        'Rodzajów znaczników nie ma czym zapisać: droga do studio.markup.type.save nie jest w tym ' +
          'oknie wpięta.',
        false,
      );
      return;
    }
    const nazwa = nazwaRodzaju.kontrolka.value.trim();
    if (nazwa === '') {
      odpowiedzRdzenia('Rodzaj znacznika musi mieć nazwę — bez niej nie ma czego zapisać.', false);
      return;
    }
    const wynik = await rdzen.zrodlo.kontrolaRodzajZnacznikaZapisz(
      nazwa,
      nazwa,
      paleta.kontrolka.value,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedzRdzenia(`Rodzaju nie zapisano: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    odpowiedzRdzenia(`Rodzaj znacznika „${wynik.wynik.markupType.name}" zapisany.`, true);
    await odswiezZnakowaniaRdzenia();
  }

  async function usunRodzajZnacznika(): Promise<void> {
    if (rdzen === null) {
      odpowiedzRdzenia(
        'Rodzajów znaczników nie ma czym usunąć: droga do studio.markup.type.delete nie jest w tym ' +
          'oknie wpięta.',
        false,
      );
      return;
    }
    const nazwa =
      nazwaRodzaju.kontrolka.value.trim() === ''
        ? rodzajZnacznikaRdzenia.kontrolka.value
        : nazwaRodzaju.kontrolka.value.trim();
    if (nazwa === '') {
      odpowiedzRdzenia('Wskaż rodzaj do usunięcia — nazwą w polu albo wyborem z wykazu.', false);
      return;
    }
    const wynik = await rdzen.zrodlo.kontrolaRodzajZnacznikaUsun(nazwa);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedzRdzenia(`Rodzaju nie usunięto: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    odpowiedzRdzenia(
      wynik.wynik.deleted
        ? `Rodzaj „${nazwa}" usunięty.`
        : `Rodzaju „${nazwa}" rdzeń nie usunął — rodzaju fabrycznego się nie usuwa.`,
      wynik.wynik.deleted,
    );
    await odswiezZnakowaniaRdzenia();
  }

  async function zdejmijZnakowanieRdzenia(znakowanie: StudioMarkup): Promise<void> {
    const zaplecze = zapleczeRdzenia();
    if (zaplecze === null) return;
    const wynik = await zaplecze.rdzen.zrodlo.kontrolaZnakowanieZdejmij(
      zaplecze.idDokumentu,
      znakowanie.id,
      {},
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedzRdzenia(`Znakowania nie zdjęto: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    odpowiedzRdzenia(
      wynik.wynik.removed ? 'Znakowanie zdjęte.' : 'Rdzeń znakowania nie zdjął.',
      wynik.wynik.removed,
    );
    if (wynik.wynik.form !== undefined) zaplecze.rdzen.naSkutek(undefined, wynik.wynik.form);
    await odswiezZnakowaniaRdzenia();
  }

  async function rozstrzygnijPropozycje(
    znakowanie: StudioMarkup,
    przyjmij: boolean,
  ): Promise<void> {
    const zaplecze = zapleczeRdzenia();
    if (zaplecze === null) return;
    const wynik = await zaplecze.rdzen.zrodlo.kontrolaZnakowanieRozstrzygnij(
      zaplecze.idDokumentu,
      [znakowanie.id],
      przyjmij,
      // Brzmienie poprawione jedzie tylko, gdy Operator je wpisał — puste pole przyjmuje propozycję wprost.
      brzmienieRdzenia.kontrolka.value.trim(),
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedzRdzenia(`Propozycji nie rozstrzygnięto: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    const tresc = wynik.wynik;
    odpowiedzRdzenia(
      `Rozstrzygnięto propozycji: ${tresc.decided}. ` +
        (przyjmij
          ? 'Brzmienie weszło do treści jako zmiana śledzona, więc da się je jeszcze cofnąć.'
          : 'Treść została nietknięta.') +
        ` Bilans: weszło w ${tresc.balance.applied} miejscach, stanęło w ${tresc.balance.skippedCount}.`,
      tresc.decided > 0,
    );
    zaplecze.rdzen.naSkutek(tresc.document.content, tresc.form);
    await odswiezZnakowaniaRdzenia();
  }

  /** Jeden wiersz wykazu znakowań trwałych wraz z jego czynnościami. */
  function wierszZnakowaniaRdzenia(znakowanie: StudioMarkup): HTMLElement {
    const glowa = document.createElement('p');
    glowa.className = 'ms-przybornik__glowa';
    glowa.textContent = przybornikOpiszZnakowanieRdzenia(znakowanie);

    const tresc = document.createElement('p');
    tresc.className = 'ms-przybornik__tresc';
    tresc.textContent =
      znakowanie.suggestedText !== undefined && znakowanie.suggestedText !== ''
        ? `Proponowane brzmienie: ${znakowanie.suggestedText}`
        : (znakowanie.body ?? '');

    const pas = document.createElement('div');
    pas.className = 'ms-przybornik__pas';

    const doMiejsca = document.createElement('button');
    doMiejsca.type = 'button';
    doMiejsca.className = 'dn-btn dn-btn--sm dn-btn--duch';
    doMiejsca.textContent = 'Pokaż miejsce';
    doMiejsca.addEventListener('click', () =>
      rdzen?.naMiejsce(znakowanie.rangeStart, znakowanie.rangeEnd),
    );
    pas.append(doMiejsca);

    if (znakowanie.kind === StudioMarkupKind.Suggestion) {
      const przyjmij = document.createElement('button');
      przyjmij.type = 'button';
      przyjmij.className = 'dn-btn dn-btn--sm dn-btn--atrament';
      przyjmij.textContent = 'Przyjmij propozycję';
      przyjmij.dataset['czynnosc'] = 'propozycja-przyjmij';
      przyjmij.addEventListener('click', () => {
        void rozstrzygnijPropozycje(znakowanie, true);
      });

      const odrzuc = document.createElement('button');
      odrzuc.type = 'button';
      odrzuc.className = 'dn-btn dn-btn--sm dn-btn--zarys';
      odrzuc.textContent = 'Odrzuć propozycję';
      odrzuc.dataset['czynnosc'] = 'propozycja-odrzuc';
      odrzuc.addEventListener('click', () => {
        void rozstrzygnijPropozycje(znakowanie, false);
      });
      pas.append(przyjmij, odrzuc);
    }

    const zdejmij = document.createElement('button');
    zdejmij.type = 'button';
    zdejmij.className = 'dn-btn dn-btn--sm dn-btn--duch';
    zdejmij.textContent = 'Zdejmij znakowanie';
    zdejmij.dataset['czynnosc'] = 'znakowanie-zdejmij';
    zdejmij.addEventListener('click', () => {
      void zdejmijZnakowanieRdzenia(znakowanie);
    });
    pas.append(zdejmij);

    const pozycja = document.createElement('li');
    pozycja.dataset['znakowanieRdzenia'] = znakowanie.id;
    pozycja.dataset['rodzaj'] = znakowanie.kind;
    pozycja.dataset['autor'] = znakowanie.author;
    pozycja.dataset['stan'] = znakowanie.state;
    pozycja.append(glowa, tresc, pas);
    return pozycja;
  }

  function ustawZaznaczenie(dlugosc: number): void {
    dlugoscZaznaczenia = dlugosc;
    zakresZdanie.textContent =
      dlugoscZaznaczenia > 0
        ? `Znakowanie obejmie zaznaczenie — ${dlugoscZaznaczenia} znaków.`
        : 'Nic nie jest zaznaczone, więc znakowanie obejmie CAŁY dokument. Zaznacz fragment, ' +
          'jeśli ma dotyczyć tylko jego — przybornik jest narzędziem do pracy na fragmentach.';
  }

  ustawZaznaczenie(0);

  return {
    element,

    pokaz(pozycje) {
      wszystkie = pozycje;
      przerysuj();
    },

    ustawZaznaczenie,

    ustawZakladki(zakladki) {
      if (zakladki.length === 0) {
        wykazZakladek.replaceChildren();
        return;
      }
      wykazZakladek.replaceChildren(
        ...zakladki.map((wpis) => {
          const powrot = document.createElement('button');
          powrot.type = 'button';
          powrot.className = 'dn-btn dn-btn--sm dn-btn--duch';
          powrot.textContent = `Wróć do „${wpis.nazwa}"`;
          powrot.dataset['zakladka'] = wpis.kod;
          const pozycja = document.createElement('li');
          pozycja.append(powrot);
          return pozycja;
        }),
      );
    },

    skocz,

    pokazOdpowiedz: odpowiedzRdzenia,

    odswiezZnakowaniaRdzenia,

    przestawWidocznosc() {
      otwarty = !otwarty;
      element.hidden = !otwarty;
    },

    widoczny: () => otwarty,
  };
}

/** Przycisk czynności przybornika wraz z jego zdaniem o zapleczu rdzenia i stanem klikalności przycisku. */
function przyciskCzynnosci(nazwa: string, kod: string, objasnienie: string): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--sm dn-btn--atrament';
  przycisk.textContent = nazwa;
  przycisk.dataset['czynnosc'] = kod;
  przycisk.title = objasnienie;
  return przycisk;
}

/** Jeden wiersz wykazu znakowań przybornika, z rodzajem znacznika, jego autorem i miejscem w dokumencie. */
function wierszZnakowania(
  pozycja: PozycjaZnakowania,
  wskazana: boolean,
  czynnosci: CzynnosciPrzybornika,
): HTMLElement {
  const glowa = document.createElement('p');
  glowa.className = 'ms-przybornik__glowa';
  const autor = pozycja.autor === StudioAuthor.Model ? 'model' : 'Operator';
  const gdzie =
    pozycja.zakres === null
      ? 'bez zakotwiczenia w treści'
      : `znaki ${pozycja.zakres.poczatek}–${pozycja.zakres.koniec}`;
  glowa.textContent = `${pozycja.tytul} · ${autor} · ${gdzie}`;

  const tresc = document.createElement('p');
  tresc.className = 'ms-przybornik__tresc';
  tresc.textContent = pozycja.tresc;

  const podstawa = document.createElement('p');
  podstawa.className = 'dn-pole-opis';
  podstawa.textContent = pozycja.podstawa;

  const doMiejsca = document.createElement('button');
  doMiejsca.type = 'button';
  doMiejsca.className = 'dn-btn dn-btn--sm dn-btn--duch';
  doMiejsca.textContent = 'Pokaż miejsce';
  doMiejsca.addEventListener('click', () => czynnosci.naSkok(pozycja));

  const odhacz = document.createElement('button');
  odhacz.type = 'button';
  odhacz.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  odhacz.textContent = pozycja.otwarta ? 'Odhacz' : 'Otwórz ponownie';
  odhacz.title =
    'Odhaczenie zamyka pozycję w wykazie: wątek komentarza rozwiązuje studio.comment.resolve, ' +
    'znacznik zamyka przybornik. Zmiany śledzonej odhaczyć nie można — jej stan rozstrzyga ' +
    'decyzja o przyjęciu albo odrzuceniu.';
  odhacz.disabled = pozycja.rodzaj === RodzajZnakowania.Zmiana;
  odhacz.addEventListener('click', () =>
    czynnosci.naOdhaczenie(pozycja.kod, pozycja.rodzaj, pozycja.otwarta),
  );

  const pas = document.createElement('div');
  pas.className = 'ms-przybornik__pas';
  pas.append(doMiejsca, odhacz);

  if (pozycja.rodzaj === RodzajZnakowania.Znacznik) {
    const zdejmij = document.createElement('button');
    zdejmij.type = 'button';
    zdejmij.className = 'dn-btn dn-btn--sm dn-btn--duch';
    zdejmij.textContent = 'Zdejmij znacznik';
    zdejmij.addEventListener('click', () => czynnosci.naZdjecieZnacznika(pozycja.kod));
    pas.append(zdejmij);
  }

  const wiersz = document.createElement('li');
  wiersz.dataset['znakowanie'] = pozycja.kod;
  wiersz.dataset['rodzaj'] = pozycja.rodzaj;
  wiersz.dataset['autor'] = pozycja.autor;
  wiersz.dataset['otwarta'] = pozycja.otwarta ? 'tak' : 'nie';
  wiersz.dataset['wskazana'] = wskazana ? 'tak' : 'nie';
  wiersz.append(glowa, tresc, podstawa, pas);
  return wiersz;
}

/** Część przybornika wraz z jej tytułem i elementami czynności, które ta część przybornika w sobie zawiera. */
function czescPrzybornika(tytul: string, elementy: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('p');
  naglowek.className = 'ms-przybornik__tytul';
  naglowek.textContent = tytul;

  const sekcja = document.createElement('section');
  sekcja.className = 'ms-przybornik__czesc';
  sekcja.append(naglowek, ...elementy);
  return sekcja;
}

/** Żeton barwy do nadania wyróżnieniu w oknie znakowania; jedno miejsce przekładu koloru na nazwę żetonu. */
export { przybornikZetonBarwy };
