import { przyciskBezKomendy } from '../../modele/kontrolki-formularza';
import { KATEGORIE_OPERACJI } from './kategorie-operacji';
import { KROJE, POWOD_NIETRWALOSCI, STOPNIE, type Wyrownanie } from './nastawy-wizualne';
import { NOSNIKI } from './nastawy-strony';
import { FORMATY_KONWERSJI } from './konwersja-dokumentu';
import { NARZEDZIA_TEKSTU, type NarzedzieTekstu } from './znaczniki-markdown';
import {
  opiszNastawe,
  SKALA_SUWAKA,
  WIELKOSCI_CIAGLE,
  type NastawySuwakow,
} from './suwaki-koncepcyjne';
import { STYLE_NAZWANE, type RodzajBloku } from './zapis-formatowany';
import type { TrybAdiustacji, TrybWidoku } from './powierzchnia-dokumentu';

/**
 * Czynności wstążki zlecane oknu, pogrupowane wedle zakładki — plik, narzędzia
 * główne, wstawianie, projektowanie, recenzja, widok i asystent — tak jak
 * grupuje je sama wstążka.
 */
export interface CzynnosciWstazki {
  /* Plik */
  naZapis(): void;
  naGalerie(): void;
  naWydanie(format: string): void;
  naPrzekazanie(): void;
  /* Narzędzia główne */
  naNarzedzieTekstu(narzedzie: NarzedzieTekstu): void;
  naStyl(rodzaj: RodzajBloku): void;
  naKrój(krój: string): void;
  naStopien(stopien: number): void;
  naInterlinie(interlinia: number): void;
  naWciecie(milimetry: number): void;
  naOdstep(milimetry: number): void;
  naWyrownanie(wyrownanie: Wyrownanie): void;
  naBarwe(barwa: string): void;
  naSzukanie(): void;
  /* Wstawianie */
  naPodzialStrony(): void;
  /* Projektowanie i układ */
  naProfil(idProfilu: string): void;
  naZapisProfilu(nazwa: string): void;
  naNosnik(oznaczenie: string): void;
  naOrientacje(pozioma: boolean): void;
  naMargines(strona: 'gora' | 'dol' | 'lewy' | 'prawy', milimetry: number): void;
  naNaglowekStrony(tresc: string): void;
  naStopkeStrony(tresc: string): void;
  naNumeracje(czynna: boolean): void;
  naSkale(procent: number): void;
  naKolumny(kolumny: number): void;
  /* Recenzja */
  naSledzenie(czynne: boolean): void;
  naDecyzjeWszystkich(przyjmij: boolean): void;
  naSkokZmiany(wPrzod: boolean): void;
  naAdiustacje(tryb: TrybAdiustacji): void;
  naNowyKomentarz(): void;
  naSkokKomentarza(wPrzod: boolean): void;
  naDymki(): void;
  naPorownanie(): void;
  naRoznicaWygladu(): void;
  naOchrone(): void;
  /* Widok */
  naTryb(tryb: TrybWidoku): void;
  naRedaktora(): void;
  /* Asystent */
  naSuwak(kod: string, wartosc: number): void;
  naOperacje(idAkcji: string): void;
  naDecyzjePropozycji(przyjmij: boolean, fragmenty: readonly number[]): void;
}

/**
 * Gniazda wstążki — kawałki złożone poza nią, takie jak formularz wczytania,
 * panel znajdź/zamień czy pola różnicy wersji, wstawiane gotowe w jej grupy.
 */
export interface GniazdaWstazki {
  /** Formularz wskazania dokumentu do wczytania. */
  wczytanie: HTMLElement;
  /** Panel znajdź/zamień. */
  szukanie: HTMLElement;
  /** Pola porównania wersji. */
  roznica: readonly HTMLElement[];
  /** Wykaz fragmentów różnicy wraz z decyzją wybiórczą. */
  fragmenty: HTMLElement;
  /** Pasek widoku powierzchni jest nieobowiązkowy: gdy podany, wstążka pokazuje go w zakładce widoku. */
  widok?: HTMLElement;
}

/**
 * Wstążka wraz z jej sterowaniem: przełączaniem zakładki czynnej, odświeżaniem
 * suwaków asystenta, znacznika śledzenia zmian, trybu widoku i wykazu profili
 * wydania.
 */
export interface WstazkaPracy {
  element: HTMLElement;
  /** Przestawia zakładkę czynną. */
  ustawZakladke(kod: string): void;
  /** Przestawia nastawy suwaków pokazywane na zakładce asystenta. */
  odswiezSuwaki(nastawy: NastawySuwakow): void;
  /** Przestawia znacznik śledzenia zmian. */
  odswiezSledzenie(czynne: boolean): void;
  /** Przestawia przełącznik trybu źródłowego wedle trybu obowiązującego. */
  odswiezTryb(tryb: TrybWidoku): void;
  /** Przestawia wykaz profili wydania. */
  ustawProfile(profile: readonly { id: string; nazwa: string }[]): void;
}

/**
 * Buduje wstążkę okna pracy z dokumentem złożoną z ośmiu zakładek z nazwanymi
 * grupami czynności; nie woła rdzenia i nie zna stanu modułu, zgłasza jedynie
 * naciśnięcia oknu przez `CzynnosciWstazki`.
 */
export function utworzWstazkePracy(
  czynnosci: CzynnosciWstazki,
  gniazda: GniazdaWstazki,
): WstazkaPracy {
  const zakladki = document.createElement('div');
  zakladki.className = 'ms-wstazka__zakladki';
  zakladki.setAttribute('role', 'tablist');

  const plansze = document.createElement('div');
  plansze.className = 'ms-wstazka__plansze';

  const element = document.createElement('div');
  element.className = 'ms-wstazka';
  element.setAttribute('aria-label', 'Wstążka okna pracy z dokumentem');
  element.append(zakladki, plansze);

  const przyciskiZakladek = new Map<string, HTMLButtonElement>();
  const planszeZakladek = new Map<string, HTMLElement>();

  function dodajZakladke(kod: string, nazwa: string, grupy: readonly HTMLElement[]): void {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'ms-wstazka__zakladka';
    przycisk.textContent = nazwa;
    przycisk.dataset['zakladka'] = kod;
    przycisk.setAttribute('role', 'tab');
    przycisk.addEventListener('click', () => ustawZakladke(kod));
    zakladki.append(przycisk);
    przyciskiZakladek.set(kod, przycisk);

    const plansza = document.createElement('div');
    plansza.className = 'ms-wstazka__plansza';
    plansza.dataset['zakladka'] = kod;
    plansza.setAttribute('role', 'tabpanel');
    plansza.append(...grupy);
    plansza.hidden = true;
    plansze.append(plansza);
    planszeZakladek.set(kod, plansza);
  }

  function ustawZakladke(kod: string): void {
    for (const [nazwa, przycisk] of przyciskiZakladek) {
      przycisk.dataset['czynna'] = nazwa === kod ? 'tak' : 'nie';
      przycisk.setAttribute('aria-selected', nazwa === kod ? 'true' : 'false');
    }
    for (const [nazwa, plansza] of planszeZakladek) plansza.hidden = nazwa !== kod;
  }

  /* ── Plik ────────────────────────────────────────────────────────────────── */

  const zapisz = przyciskWstazki('Zapisz i załóż wersję', 'zapisz', () => czynnosci.naZapis());
  const galeria = przyciskWstazki('Galeria szablonów', 'galeria', () => czynnosci.naGalerie());
  const przekaz = przyciskWstazki('Przekaż do Library', 'do-library', () =>
    czynnosci.naPrzekazanie(),
  );

  const formatWydania = wybor(
    'Format wydania',
    FORMATY_KONWERSJI.map((format) => ({ wartosc: format, etykieta: format })),
    () => undefined,
  );
  const wydaj = przyciskWstazki('Wydaj w formacie', 'wydaj', () =>
    czynnosci.naWydanie(formatWydania.kontrolka.value),
  );

  dodajZakladke('plik', 'Plik', [
    grupa('Dokument', [gniazda.wczytanie, zapisz, galeria]),
    grupa('Wydanie', [formatWydania.element, wydaj, przekaz]),
    grupa('Bez pokrycia w kontrakcie', [
      przyciskBezKomendy(
        'Zmień format dokumentu',
        'Format dokumentu nadaje rdzeń przy wczytaniu i oddaje go w polu StudioDocument.format. ' +
          'Komenda studio.document.save przyjmuje documentId, content, title i createVersion — ' +
          'pola format nie ma. Wydanie w innym formacie jest czym innym i działa: zamiana ' +
          'formatu z obszaru document oddaje NOWY zasób.',
      ),
    ]),
  ]);

  /* ── Narzędzia główne ────────────────────────────────────────────────────── */

  const styl = wybor(
    'Styl nazwany',
    STYLE_NAZWANE.map((pozycja) => ({ wartosc: pozycja.rodzaj, etykieta: pozycja.nazwa })),
    (wartosc) => czynnosci.naStyl(wartosc as RodzajBloku),
  );
  const krój = wybor('Krój', [...KROJE], (wartosc) => czynnosci.naKrój(wartosc));
  const stopien = wybor(
    'Stopień',
    STOPNIE.map((punkty) => ({ wartosc: String(punkty), etykieta: `${punkty} pt` })),
    (wartosc) => czynnosci.naStopien(Number(wartosc)),
  );
  stopien.kontrolka.value = '11';

  const znaczniki = document.createElement('div');
  znaczniki.className = 'ms-wstazka__rzad';
  for (const narzedzie of NARZEDZIA_TEKSTU) {
    znaczniki.append(
      przyciskWstazki(narzedzie.nazwa, narzedzie.kod, () => czynnosci.naNarzedzieTekstu(narzedzie), narzedzie.opis),
    );
  }

  const wyrownania: readonly { kod: Wyrownanie; nazwa: string }[] = [
    { kod: 'lewo', nazwa: 'Do lewej' },
    { kod: 'srodek', nazwa: 'Do środka' },
    { kod: 'prawo', nazwa: 'Do prawej' },
    { kod: 'obustronne', nazwa: 'Obustronnie' },
  ];
  const rzadWyrownan = document.createElement('div');
  rzadWyrownan.className = 'ms-wstazka__rzad';
  for (const pozycja of wyrownania) {
    rzadWyrownan.append(
      przyciskWstazki(pozycja.nazwa, `wyrownanie-${pozycja.kod}`, () =>
        czynnosci.naWyrownanie(pozycja.kod),
      ),
    );
  }

  const interlinia = wybor(
    'Interlinia',
    ['1', '1.15', '1.5', '2'].map((wartosc) => ({ wartosc, etykieta: wartosc })),
    (wartosc) => czynnosci.naInterlinie(Number(wartosc)),
  );
  interlinia.kontrolka.value = '1.5';

  const wciecie = liczba('Wcięcie pierwszego wiersza (mm)', 0, (wartosc) =>
    czynnosci.naWciecie(wartosc),
  );
  const odstep = liczba('Odstęp akapitowy (mm)', 3, (wartosc) => czynnosci.naOdstep(wartosc));

  const barwa = document.createElement('input');
  barwa.type = 'color';
  barwa.className = 'ms-wstazka__barwa';
  barwa.setAttribute('aria-label', 'Barwa pisma akapitu');
  barwa.addEventListener('change', () => czynnosci.naBarwe(barwa.value));

  const oNietrwalosci = document.createElement('p');
  oNietrwalosci.className = 'dn-pole-opis';
  oNietrwalosci.textContent = POWOD_NIETRWALOSCI;

  dodajZakladke('glowne', 'Narzędzia główne', [
    grupa('Styl', [styl.element]),
    grupa('Czcionka', [krój.element, stopien.element, barwa, znaczniki]),
    grupa('Akapit', [rzadWyrownan, interlinia.element, wciecie.element, odstep.element]),
    grupa('Szukanie', [
      przyciskWstazki('Znajdź i zamień', 'szukanie', () => czynnosci.naSzukanie()),
      gniazda.szukanie,
    ]),
    grupa('Trwałość nastaw', [oNietrwalosci]),
  ]);

  /* ── Wstawianie ──────────────────────────────────────────────────────────── */

  const wstawiane = NARZEDZIA_TEKSTU.filter((narzedzie) =>
    ['tabela', 'linia', 'odnosnik', 'obraz', 'kod-blok'].includes(narzedzie.kod),
  );
  const rzadWstawiania = document.createElement('div');
  rzadWstawiania.className = 'ms-wstazka__rzad';
  for (const narzedzie of wstawiane) {
    rzadWstawiania.append(
      przyciskWstazki(narzedzie.nazwa, `wstaw-${narzedzie.kod}`, () =>
        czynnosci.naNarzedzieTekstu(narzedzie), narzedzie.opis),
    );
  }

  dodajZakladke('wstawianie', 'Wstawianie', [
    grupa('Elementy treści', [rzadWstawiania]),
    grupa('Strony', [
      przyciskWstazki('Podział strony', 'podzial-strony', () => czynnosci.naPodzialStrony(),
        'Wstawia jawny podział strony. Zapisuje się w treści komentarzem HTML, bo markdown ' +
        'podziału strony nie ma — przeżywa więc zapis do rdzenia.'),
    ]),
  ]);

  /* ── Projektowanie ───────────────────────────────────────────────────────── */

  const profile = wybor('Profil wydania', [{ wartosc: '', etykieta: 'bez profilu' }], (wartosc) =>
    czynnosci.naProfil(wartosc),
  );
  const nazwaProfilu = document.createElement('input');
  nazwaProfilu.type = 'text';
  nazwaProfilu.className = 'dn-pole-kontrolka';
  nazwaProfilu.placeholder = 'nazwa profilu do zapisania';
  nazwaProfilu.setAttribute('aria-label', 'Nazwa profilu wydania');

  const naglowekStrony = pole('Nagłówek strony', (wartosc) => czynnosci.naNaglowekStrony(wartosc));
  const stopkaStrony = pole('Stopka strony', (wartosc) => czynnosci.naStopkeStrony(wartosc));

  const numeracja = document.createElement('input');
  numeracja.type = 'checkbox';
  numeracja.className = 'dn-przelacznik';
  numeracja.checked = true;
  numeracja.setAttribute('aria-label', 'Numeracja stron');
  numeracja.addEventListener('change', () => czynnosci.naNumeracje(numeracja.checked));

  dodajZakladke('projektowanie', 'Projektowanie', [
    grupa('Profil wydania', [
      profile.element,
      nazwaProfilu,
      przyciskWstazki('Zapisz profil', 'zapisz-profil', () =>
        czynnosci.naZapisProfilu(nazwaProfilu.value.trim()),
        'Zapisuje nastawy strony w rdzeniu komendą studio.export.profile.save — to jedyne ' +
        'miejsce w kontrakcie, w którym kartka, marginesy, nagłówek i stopka trwają.'),
    ]),
    grupa('Nagłówek i stopka', [naglowekStrony.element, stopkaStrony.element, numeracja]),
  ]);

  /* ── Układ ───────────────────────────────────────────────────────────────── */

  const nosnik = wybor(
    'Nośnik',
    NOSNIKI.map((pozycja) => ({
      wartosc: pozycja.oznaczenie,
      etykieta: `${pozycja.oznaczenie} (${pozycja.szerokoscMm}×${pozycja.wysokoscMm} mm)`,
    })),
    (wartosc) => czynnosci.naNosnik(wartosc),
  );
  nosnik.kontrolka.value = 'A4';

  const orientacja = wybor(
    'Orientacja',
    [
      { wartosc: 'pionowa', etykieta: 'Pionowa' },
      { wartosc: 'pozioma', etykieta: 'Pozioma' },
    ],
    (wartosc) => czynnosci.naOrientacje(wartosc === 'pozioma'),
  );

  const marginesy = document.createElement('div');
  marginesy.className = 'ms-wstazka__rzad';
  for (const strona of ['gora', 'dol', 'lewy', 'prawy'] as const) {
    const kontrolka = liczba(`Margines ${strona} (mm)`, 20, (wartosc) =>
      czynnosci.naMargines(strona, wartosc),
    );
    marginesy.append(kontrolka.element);
  }

  const skala = liczba('Skala widoku (%)', 100, (wartosc) => czynnosci.naSkale(wartosc));
  const kolumny = liczba('Kartek obok siebie', 1, (wartosc) => czynnosci.naKolumny(wartosc));

  dodajZakladke('uklad', 'Układ', [
    grupa('Kartka', [nosnik.element, orientacja.element]),
    grupa('Marginesy', [marginesy]),
    grupa('Widok kartek', [skala.element, kolumny.element]),
  ]);

  /* ── Recenzja ────────────────────────────────────────────────────────────── */

  const sledzenie = document.createElement('input');
  sledzenie.type = 'checkbox';
  sledzenie.className = 'dn-przelacznik';
  sledzenie.setAttribute('aria-label', 'Śledzenie zmian dokumentu');
  sledzenie.addEventListener('change', () => czynnosci.naSledzenie(sledzenie.checked));

  const adiustacja = wybor(
    'Pokazywanie adiustacji',
    [
      { wartosc: 'cala', etykieta: 'Cała adiustacja' },
      { wartosc: 'po-zmianach', etykieta: 'Tekst po zmianach' },
      { wartosc: 'okienko', etykieta: 'Okienko recenzowania' },
    ],
    (wartosc) => czynnosci.naAdiustacje(wartosc as TrybAdiustacji),
  );

  dodajZakladke('recenzja', 'Recenzja', [
    grupa('Śledzenie zmian', [
      sledzenie,
      przyciskWstazki('Przyjmij wszystkie', 'przyjmij-wszystkie', () =>
        czynnosci.naDecyzjeWszystkich(true),
      ),
      przyciskWstazki('Odrzuć wszystkie', 'odrzuc-wszystkie', () =>
        czynnosci.naDecyzjeWszystkich(false),
      ),
      przyciskWstazki('Poprzednia zmiana', 'zmiana-wstecz', () => czynnosci.naSkokZmiany(false)),
      przyciskWstazki('Następna zmiana', 'zmiana-dalej', () => czynnosci.naSkokZmiany(true)),
      adiustacja.element,
    ]),
    grupa('Komentarze', [
      przyciskWstazki('Nowy komentarz', 'komentarz-nowy', () => czynnosci.naNowyKomentarz()),
      przyciskWstazki('Poprzedni', 'komentarz-wstecz', () => czynnosci.naSkokKomentarza(false)),
      przyciskWstazki('Następny', 'komentarz-dalej', () => czynnosci.naSkokKomentarza(true)),
      przyciskWstazki('Pokaż komentarze', 'komentarze-widok', () => czynnosci.naDymki()),
      przyciskBezKomendy(
        'Usuń komentarz',
        'Kontrakt niesie studio.comment.add, studio.comment.list i studio.comment.resolve — ' +
          'komendy usuwającej komentarz nie ma. Wątek zamyka się rozwiązaniem, które da się ' +
          'cofnąć; usunięcie byłoby czynnością bez drogi do rdzenia.',
      ),
    ]),
    grupa('Porównanie', [
      ...gniazda.roznica,
      przyciskWstazki('Porównaj i wyszukaj', 'porownaj', () => czynnosci.naPorownanie()),
      przyciskWstazki('Różnica wyglądu', 'roznica-wygladu', () => czynnosci.naRoznicaWygladu(),
        'Porównanie po WYRYSIE stron komendą studio.diff.visual — widzi przesunięcie akapitu ' +
        'i zmianę łamania, których różnica tekstowa nie widzi.'),
      gniazda.fragmenty,
    ]),
    grupa('Ochrona', [
      przyciskWstazki('Czynności ochrony dokumentu', 'ochrona', () => czynnosci.naOchrone(),
        'Szyfrowanie, podpis, weryfikacja podpisu, redakcja, zdjęcie metadanych i wykrycie ' +
        'danych wrażliwych stoją w oknie Warsztatu dokumentu wraz z formularzami swoich ' +
        'parametrów. Ten przycisk przenosi tam ognisko — drugi formularz tych samych ' +
        'komend rozjechałby się z pierwszym.'),
    ]),
    grupa('Bez drogi w module Studio', [
      przyciskBezKomendy(
        'Tezaurus',
        'Wyrazów bliskoznacznych nie ma w kontrakcie ani w rdzeniu: wymagają słownika języka, ' +
          'którego arsenał wkompilowany w serwer nie niesie. Zamiast liczby zmyślonej — brak ' +
          'nazwany. Zbliżone zadanie wykonuje operacja stylu zlecana modelowi.',
      ),
      przyciskBezKomendy(
        'Czytanie na głos',
        'Kanału mowy kontrakt nie niesie w module Studio. Czynność zostaje brakiem nazwanym.',
      ),
      przyciskBezKomendy(
        'Sprawdzenie ułatwień dostępu',
        'Sprawdzianu ułatwień dostępu dokumentu rdzeń nie ma. Ułatwienia samego okna są ' +
          'pilnowane osobno, w warstwie wizualnej produktu.',
      ),
      przyciskBezKomendy(
        'Tłumaczenie w miejscu',
        'Tłumaczenie należy do modułu Translate i tam ma swoje komendy; obszar studio niesie ' +
          'wyłącznie operację kontekstową „tłumaczenie zaznaczenia" zlecaną modelowi. ' +
          'Tłumaczenia w miejscu Studio nie dorabia — zgłoszone jako potrzeba.',
      ),
    ]),
  ]);

  /* ── Widok ───────────────────────────────────────────────────────────────── */

  // Tryb źródłowy stoi osobno, przełącznikiem nad trybami widoku, wydruku
  // i różnicy.
  const tryby: readonly { kod: TrybWidoku; nazwa: string; opis: string }[] = [
    {
      kod: 'formatowany',
      nazwa: 'Widok formatowany',
      opis: 'Kartka o rozmiarze nośnika, z marginesami, nagłówkiem, stopką i paginacją.',
    },
    {
      kod: 'wydanie',
      nazwa: 'Podgląd wydruku',
      opis:
        'Tryb TEGO okna, nie osobne okno: dokument w postaci, w jakiej wyjdzie z drukarki — ' +
        'nośnik i orientacja naprawdę ustawione, paginacja, nagłówek, stopka i numeracja. ' +
        'Pisanie jest w nim wyłączone.',
    },
    {
      kod: 'roznica',
      nazwa: 'Różnica na treści',
      opis: 'Fragmenty różnicy nałożone na treść w miejscu, zamiast obok niej.',
    },
  ];
  const rzadTrybow = document.createElement('div');
  rzadTrybow.className = 'ms-wstazka__rzad';
  for (const tryb of tryby) {
    rzadTrybow.append(
      przyciskWstazki(tryb.nazwa, `tryb-${tryb.kod}`, () => czynnosci.naTryb(tryb.kod), tryb.opis),
    );
  }

  const przelacznikZrodlowy = document.createElement('input');
  przelacznikZrodlowy.type = 'checkbox';
  przelacznikZrodlowy.className = 'dn-przelacznik';
  przelacznikZrodlowy.setAttribute('aria-label', 'Tryb źródłowy ze znacznikami');
  // Znacznik czynności zostaje ten sam co przy przycisku trybu — zmienia się
  // kształt, nie znaczenie.
  przelacznikZrodlowy.dataset['czynnosc'] = 'tryb-zrodlowy';
  przelacznikZrodlowy.title =
    'Pokazuje treść ze znacznikami, tak jak jedzie do rdzenia. Przełącznik, nie tryb domyślny — ' +
    'zdjęcie go wraca do widoku formatowanego na kartce.';
  przelacznikZrodlowy.addEventListener('change', () =>
    czynnosci.naTryb(przelacznikZrodlowy.checked ? 'zrodlowy' : 'formatowany'),
  );
  const etykietaZrodlowego = document.createElement('label');
  etykietaZrodlowego.className = 'ms-wstazka__pole';
  const napisZrodlowego = document.createElement('span');
  napisZrodlowego.className = 'ms-wstazka__etykieta';
  napisZrodlowego.textContent = 'Tryb źródłowy (znaczniki)';
  etykietaZrodlowego.append(napisZrodlowego, przelacznikZrodlowy);

  const grupyWidoku: HTMLElement[] = [
    grupa('Tryb widoku', [rzadTrybow]),
    grupa('Adiustacja i znaczniki', [etykietaZrodlowego]),
  ];
  if (gniazda.widok !== undefined) {
    grupyWidoku.push(grupa('Powierzchnia i kartki', [gniazda.widok]));
  } else {
    const wskazanie = document.createElement('p');
    wskazanie.className = 'dn-pole-opis';
    wskazanie.textContent =
      'Skala widoku wraz z nastawami gotowymi, układ kartek (jedna, obok siebie, rozkładówka), ' +
      'przewijanie ciągłe albo strona po stronie, linijki, jednostka podziałki, granice ' +
      'marginesów i skok o kartkę stoją na pasku widoku PRZY powierzchni — pod ręką, bo sięga ' +
      'się po nie co chwilę. Drugiego ich miejsca wstążka nie zakłada, żeby nastawy nie ' +
      'rozjechały się między dwoma paskami.';
    grupyWidoku.push(grupa('Powierzchnia i kartki', [wskazanie]));
  }
  grupyWidoku.push(
    grupa('Panele na żądanie', [
      przyciskWstazki('Panel Redaktora', 'redaktor', () => czynnosci.naRedaktora(),
        'Wchodzi nakładką nad treścią i schodzi drugim naciśnięciem — powierzchnia należy do ' +
        'dokumentu. Stałe kolumny są trybem do wyboru na pasku widoku.'),
      przyciskWstazki('Dymki komentarzy', 'dymki', () => czynnosci.naDymki(),
        'Tak samo nakładką: komentarze wchodzą, kiedy są potrzebne, i nie zabierają kartce ' +
        'szerokości, kiedy nie są.'),
      przyciskWstazki('Wstążka PDF i cyfryzacja', 'wstazka-pdf', () => czynnosci.naOchrone(),
        'Otwiera zakładkę kontekstową PDF: strony, nakładanie, treść, bezpieczeństwo oraz ' +
        'narzędziownia cyfryzacji. Wstążka jest domyślnie ukryta i przy dokumencie PDF wchodzi ' +
        'sama — nie zajmuje miejsca, gdy Operator nad PDF-em nie pracuje.'),
    ]),
  );

  dodajZakladke('widok', 'Widok', grupyWidoku);

  /* ── Asystent ────────────────────────────────────────────────────────────── */

  const suwaki = new Map<string, HTMLInputElement>();
  const zdaniaSuwakow = new Map<string, HTMLElement>();
  const rzadSuwakow = document.createElement('div');
  rzadSuwakow.className = 'ms-wstazka__suwaki';
  for (const wielkosc of WIELKOSCI_CIAGLE) {
    const etykieta = document.createElement('label');
    etykieta.className = 'ms-suwak__etykieta';
    etykieta.textContent = `${wielkosc.nazwa}: ${wielkosc.koniecDolny} ↔ ${wielkosc.koniecGorny}`;

    const kontrolka = document.createElement('input');
    kontrolka.type = 'range';
    kontrolka.className = 'ms-suwak__kontrolka';
    kontrolka.min = String(SKALA_SUWAKA.dol);
    kontrolka.max = String(SKALA_SUWAKA.gora);
    kontrolka.step = String(SKALA_SUWAKA.krok);
    kontrolka.value = String(wielkosc.neutralna);
    kontrolka.dataset['suwak'] = wielkosc.kod;
    kontrolka.title = wielkosc.opis;

    const zdanie = document.createElement('p');
    zdanie.className = 'dn-pole-opis ms-suwak__zdanie';
    zdanie.textContent = opiszNastawe(wielkosc, wielkosc.neutralna);

    kontrolka.addEventListener('input', () => {
      const wartosc = Number(kontrolka.value);
      zdanie.textContent = opiszNastawe(wielkosc, wartosc);
      czynnosci.naSuwak(wielkosc.kod, wartosc);
    });

    const pozycja = document.createElement('div');
    pozycja.className = 'ms-suwak';
    pozycja.append(etykieta, kontrolka, zdanie);
    rzadSuwakow.append(pozycja);
    suwaki.set(wielkosc.kod, kontrolka);
    zdaniaSuwakow.set(wielkosc.kod, zdanie);
  }

  const operacje = document.createElement('select');
  operacje.className = 'dn-pole-kontrolka';
  operacje.setAttribute('aria-label', 'Operacja kontekstowa do uruchomienia');
  for (const kategoria of KATEGORIE_OPERACJI) {
    const grupaOpcji = document.createElement('optgroup');
    grupaOpcji.label = kategoria.nazwa;
    for (const operacja of kategoria.operacje) {
      const opcja = document.createElement('option');
      opcja.value = operacja.id;
      opcja.textContent = operacja.nazwa;
      grupaOpcji.append(opcja);
    }
    operacje.append(grupaOpcji);
  }

  dodajZakladke('asystent', 'Asystent', [
    grupa('Wielkości ciągłe', [rzadSuwakow]),
    grupa('Operacje kontekstowe', [
      operacje,
      przyciskWstazki('Uruchom operację', 'uruchom', () => czynnosci.naOperacje(operacje.value)),
    ]),
    grupa('Decyzja o wyniku', [
      przyciskWstazki('Przyjmij wynik operacji', 'przyjmij', () =>
        czynnosci.naDecyzjePropozycji(true, []),
      ),
      przyciskWstazki('Odrzuć wynik operacji', 'odrzuc', () =>
        czynnosci.naDecyzjePropozycji(false, []),
      ),
    ]),
  ]);

  ustawZakladke('glowne');

  return {
    element,
    ustawZakladke,

    odswiezSuwaki(nastawy) {
      for (const wielkosc of WIELKOSCI_CIAGLE) {
        const wartosc = nastawy[wielkosc.kod] ?? wielkosc.neutralna;
        const kontrolka = suwaki.get(wielkosc.kod);
        if (kontrolka !== undefined && Number(kontrolka.value) !== wartosc) {
          kontrolka.value = String(wartosc);
        }
        const zdanie = zdaniaSuwakow.get(wielkosc.kod);
        if (zdanie !== null && zdanie !== undefined) {
          zdanie.textContent = opiszNastawe(wielkosc, wartosc);
        }
      }
    },

    odswiezSledzenie(czynne) {
      sledzenie.checked = czynne;
    },

    odswiezTryb(tryb) {
      przelacznikZrodlowy.checked = tryb === 'zrodlowy';
    },

    ustawProfile(wykaz) {
      profile.kontrolka.replaceChildren();
      const bez = document.createElement('option');
      bez.value = '';
      bez.textContent = 'bez profilu — nastawy okna';
      profile.kontrolka.append(bez);
      for (const pozycja of wykaz) {
        const opcja = document.createElement('option');
        opcja.value = pozycja.id;
        opcja.textContent = pozycja.nazwa;
        profile.kontrolka.append(opcja);
      }
    },
  };
}

/* ── Kawałki wspólne wstążki ──────────────────────────────────────────────── */

/**
 * Nazwana grupa czynności — jedno miejsce składania kafla wstążki złożonego
 * z tytułu grupy i rzędu jej kontrolek lub przycisków.
 */
function grupa(nazwa: string, elementy: readonly HTMLElement[]): HTMLElement {
  const tytul = document.createElement('p');
  tytul.className = 'ms-wstazka__grupa-tytul';
  tytul.textContent = nazwa;

  const pas = document.createElement('div');
  pas.className = 'ms-wstazka__grupa-pas';
  pas.append(...elementy);

  const element = document.createElement('section');
  element.className = 'ms-wstazka__grupa';
  element.dataset['grupa'] = nazwa;
  element.append(pas, tytul);
  return element;
}

/**
 * Przycisk wstążki wraz z jego objaśnieniem: tytuł i opis dostępności pokazują
 * to samo objaśnienie, żeby czynność miała jedno, spójne uzasadnienie.
 */
function przyciskWstazki(
  nazwa: string,
  kod: string,
  czynnosc: () => void,
  objasnienie?: string,
): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  przycisk.textContent = nazwa;
  przycisk.dataset['czynnosc'] = kod;
  if (objasnienie !== undefined) {
    przycisk.title = objasnienie;
    przycisk.setAttribute('aria-description', objasnienie);
  }
  przycisk.addEventListener('click', czynnosc);
  return przycisk;
}

/**
 * Lista wyboru wstążki wraz z etykietą, budowana z podanych pozycji i zgłaszająca
 * zmianę wybranej wartości wywołującemu bez własnej pamięci stanu.
 */
function wybor(
  etykieta: string,
  pozycje: readonly { wartosc: string; etykieta: string }[],
  naZmiane: (wartosc: string) => void,
): { element: HTMLElement; kontrolka: HTMLSelectElement } {
  const kontrolka = document.createElement('select');
  kontrolka.className = 'dn-pole-kontrolka ms-wstazka__wybor';
  kontrolka.setAttribute('aria-label', etykieta);
  for (const pozycja of pozycje) {
    const opcja = document.createElement('option');
    opcja.value = pozycja.wartosc;
    opcja.textContent = pozycja.etykieta;
    kontrolka.append(opcja);
  }
  kontrolka.addEventListener('change', () => naZmiane(kontrolka.value));
  return { element: obudowa(etykieta, kontrolka), kontrolka };
}

/**
 * Pole liczbowe wstążki — milimetry, procenty albo liczba kartek — z wartością
 * domyślną i zgłoszeniem zmiany do wywołującego przy każdej edycji.
 */
function liczba(
  etykieta: string,
  domyslna: number,
  naZmiane: (wartosc: number) => void,
): { element: HTMLElement; kontrolka: HTMLInputElement } {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'number';
  kontrolka.className = 'dn-pole-kontrolka ms-wstazka__liczba';
  kontrolka.value = String(domyslna);
  kontrolka.min = '0';
  kontrolka.setAttribute('aria-label', etykieta);
  kontrolka.addEventListener('change', () => naZmiane(Number(kontrolka.value)));
  return { element: obudowa(etykieta, kontrolka), kontrolka };
}

/**
 * Pole tekstowe wstążki z etykietą dostępności, zgłaszające wywołującemu każdą
 * zmianę wpisanej treści bez własnej pamięci stanu.
 */
function pole(
  etykieta: string,
  naZmiane: (wartosc: string) => void,
): { element: HTMLElement; kontrolka: HTMLInputElement } {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'text';
  kontrolka.className = 'dn-pole-kontrolka';
  kontrolka.setAttribute('aria-label', etykieta);
  kontrolka.addEventListener('change', () => naZmiane(kontrolka.value));
  return { element: obudowa(etykieta, kontrolka), kontrolka };
}

/**
 * Obudowa kontrolki wstążki: etykieta widoczna nad kontrolką, wspólna dla
 * wszystkich pól wstążki niezależnie od rodzaju kontrolki.
 */
function obudowa(etykieta: string, kontrolka: HTMLElement): HTMLElement {
  const napis = document.createElement('span');
  napis.className = 'ms-wstazka__etykieta';
  napis.textContent = etykieta;

  const element = document.createElement('label');
  element.className = 'ms-wstazka__pole';
  element.append(napis, kontrolka);
  return element;
}
