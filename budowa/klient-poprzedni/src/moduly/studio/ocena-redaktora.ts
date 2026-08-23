import { policzTresc, type LicznikiDokumentu } from './liczniki-dokumentu';

/**
 * Pomiary panelu Redaktora — wyłącznie to, co da się policzyć z treści.
 *
 * ── Zasada tego pliku ───────────────────────────────────────────────────────
 * Każda liczba, którą panel pokazuje, jest tu policzona z napisu, i nic poza
 * tym. Panel Redaktora w pakiecie biurowym pokazuje „ocenę 87 %" liczoną
 * słownikiem języka, korpusem i regułami gramatyki — czego ten produkt nie ma
 * i czego nie zdobędzie programem spoza instalki. Zamiast oceny wymyślonej
 * panel pokazuje ocenę POLICZONĄ z miar czytelności i nazywa, z czego ona jest;
 * pozycje bez pomiaru mówią wprost, że pomiaru nie ma, i wskazują, co byłoby
 * potrzebne, żeby był.
 *
 * ── Skąd wskaźnik czytelności ───────────────────────────────────────────────
 * Ze średniej długości zdania i średniej długości słowa — dwóch wielkości, które
 * liczy się z samego tekstu, bez słownika. Jest to rodzina wskaźników
 * mglistości (Gunning fog, FOG-PL): im dłuższe zdania i im dłuższe wyrazy, tym
 * wyżej wykształcenia potrzeba, żeby tekst przeczytać bez potykania się.
 * Wartość nie jest oceną jakości i panel tego nie udaje.
 *
 * Plik nie zna DOM ani rdzenia.
 */

/** Pomiar nazwany wraz z tym, czym został policzony. */
export interface PomiarRedaktora {
  kod: string;
  nazwa: string;
  /** Wartość policzona; `null` znaczy „tej wielkości nie mierzymy". */
  wartosc: number | null;
  /** Jednostka albo miano wartości. */
  miano: string;
  /** Z czego wartość policzono — albo czego brakuje, żeby ją policzyć. */
  podstawa: string;
}

/** Wynik pomiaru czytelności wraz z jego składnikami. */
export interface Czytelnosc {
  /** Średnia liczba słów w zdaniu. */
  slowNaZdanie: number;
  /** Średnia liczba znaków w słowie. */
  znakowNaSlowo: number;
  /** Udział słów długich (co najmniej 10 znaków) w procentach. */
  udzialSlowDlugich: number;
  /** Wskaźnik mglistości — lata nauki potrzebne do swobodnego czytania. */
  mglistosc: number;
}

/** Granica, od której słowo liczy się jako długie. Wartość nazwana, nie magiczna. */
const DLUGIE_SLOWO = 10;

/** Liczy czytelność treści; treść pusta oddaje zera, nie dzielenie przez zero. */
export function policzCzytelnosc(tresc: string): Czytelnosc {
  const liczniki = policzTresc(tresc);
  const slowa = tresc.split(/\s+/u).filter((slowo) => slowo !== '');
  if (liczniki.slowa === 0 || liczniki.zdania === 0) {
    return { slowNaZdanie: 0, znakowNaSlowo: 0, udzialSlowDlugich: 0, mglistosc: 0 };
  }
  const dlugie = slowa.filter((slowo) => slowo.replace(/\W/gu, '').length >= DLUGIE_SLOWO).length;
  const slowNaZdanie = liczniki.slowa / liczniki.zdania;
  const udzial = (dlugie / liczniki.slowa) * 100;
  return {
    slowNaZdanie: zaokraglij(slowNaZdanie),
    znakowNaSlowo: zaokraglij(liczniki.znakiBezOdstepow / liczniki.slowa),
    udzialSlowDlugich: zaokraglij(udzial),
    // Wzór mglistości: 0,4 × (średnia długość zdania + udział słów długich).
    mglistosc: zaokraglij(0.4 * (slowNaZdanie + udzial)),
  };
}

/** Zaokrąglenie do jednego miejsca po przecinku — panel nie pokazuje szumu. */
function zaokraglij(wartosc: number): number {
  return Math.round(wartosc * 10) / 10;
}

/**
 * Ocena dokumentu wyrażona liczbą — punkty czytelności od 0 do 100.
 *
 * Nie jest to „jakość pisma": to odwrotność mglistości przełożona na skalę
 * setną. Mglistość 8 (tekst prasowy) daje wynik wysoki, mglistość 20 (zdania
 * wielokrotnie złożone z terminami) — niski. Panel podaje przy liczbie jej
 * podstawę, bo liczba bez podstawy jest kopertą.
 */
export function ocenaCzytelnosci(czytelnosc: Czytelnosc): number {
  if (czytelnosc.mglistosc === 0) return 0;
  // Granice: mglistość 6 i niżej to 100 punktów, 26 i wyżej to 0.
  const punkty = ((26 - czytelnosc.mglistosc) / 20) * 100;
  return Math.max(0, Math.min(100, Math.round(punkty)));
}

/** Zdanie o ocenie wraz z jej podstawą. */
export function opiszOcene(czytelnosc: Czytelnosc): string {
  if (czytelnosc.mglistosc === 0) {
    return 'Oceny nie ma, bo nie ma czego mierzyć — dokument jest pusty.';
  }
  return (
    `Ocena ${ocenaCzytelnosci(czytelnosc)} / 100 policzona z mglistości ${czytelnosc.mglistosc} ` +
    `(0,4 × [${czytelnosc.slowNaZdanie} słów na zdanie + ${czytelnosc.udzialSlowDlugich} % słów ` +
    `co najmniej dziesięcioznakowych]). Nie jest to ocena jakości pisma, tylko miara trudności ` +
    'czytania — jedyna, którą da się policzyć bez słownika języka.'
  );
}

/* ── Korekty: co da się policzyć bez słownika ──────────────────────────────── */

/**
 * Pomiary korekty — wykaz wzorców, które rozpoznaje się w samym zapisie.
 *
 * Pisowni i gramatyki nie mierzymy: jedno wymaga słownika języka, drugie
 * analizy składniowej, a rdzeń nie ma ani jednego, ani drugiego. Wiersze
 * pisowni i gramatyki zostają więc w panelu z wartością niepodaną i z powodem —
 * usunięcie ich kazałoby panelowi wyglądać na kompletny.
 *
 * Interpunkcję mierzymy częściowo i tylko tam, gdzie wzorzec jest pewny:
 * odstęp przed znakiem przestankowym, brak odstępu po nim, podwójny odstęp,
 * podwójny znak przestankowy, spójnik na końcu wiersza.
 */
export function pomiaryKorekty(tresc: string): PomiarRedaktora[] {
  const odstepPrzed = (tresc.match(/\s+[,.;:!?]/gu) ?? []).length;
  const bezOdstepu = (tresc.match(/[,;:](?=\p{L})/gu) ?? []).length;
  const podwojnyOdstep = (tresc.match(/\p{L} {2,}\p{L}/gu) ?? []).length;
  const podwojnyZnak = (tresc.match(/[,.;:]{2,}/gu) ?? []).length;
  const wiszacySpojnik = (tresc.match(/\s(?:i|w|z|a|o|u|na|do|od|za|po|ze)\n/giu) ?? []).length;

  return [
    {
      kod: 'interpunkcja',
      nazwa: 'Interpunkcja — wzorce pewne',
      wartosc: odstepPrzed + bezOdstepu + podwojnyZnak,
      miano: 'trafień',
      podstawa:
        `Odstęp przed znakiem przestankowym: ${odstepPrzed}; brak odstępu po znaku: ` +
        `${bezOdstepu}; znak przestankowy podwojony: ${podwojnyZnak}. Liczone wzorcem ` +
        'w treści bufora, bez słownika i bez rdzenia.',
    },
    {
      kod: 'sklad',
      nazwa: 'Skład — odstępy i wiszące spójniki',
      wartosc: podwojnyOdstep + wiszacySpojnik,
      miano: 'trafień',
      podstawa:
        `Odstęp podwójny między słowami: ${podwojnyOdstep}; spójnik jednoliterowy na końcu ` +
        `wiersza: ${wiszacySpojnik}. Reguła zapisu, nie reguła języka.`,
    },
    {
      kod: 'pisownia',
      nazwa: 'Pisownia',
      wartosc: null,
      miano: '',
      podstawa:
        'BEZ POMIARU. Rozpoznanie błędu pisowni wymaga słownika języka polskiego, a rdzeń go nie ' +
        'niesie; program słownikowy spoza instalki jest zakazany. Pomiar będzie, gdy słownik ' +
        'wejdzie do arsenału wkompilowanego w serwer — dopóki go nie ma, liczba w tym wierszu ' +
        'byłaby zmyślona. Sprawdzenie pisowni zleca się modelowi operacją korekty.',
    },
    {
      kod: 'gramatyka',
      nazwa: 'Gramatyka',
      wartosc: null,
      miano: '',
      podstawa:
        'BEZ POMIARU. Zgodność przypadków i czasów wymaga analizy składniowej, której rdzeń nie ' +
        'ma. Zleca się ją modelowi operacją korekty ortograficzno-gramatycznej — wynik wchodzi ' +
        'do treści jako zmiana śledzona, a nie jako liczba w tym panelu.',
    },
  ];
}

/* ── Uściślenia: wielkości ciągłe pracy nad stylem ─────────────────────────── */

/** Jedno uściślenie panelu — nazwa, pomiar wskazujący potrzebę i akcja. */
export interface Uscislenie {
  kod: string;
  nazwa: string;
  /** Zdanie o tym, co pomiar mówi o tekście. */
  wskazanie: string;
  /** Identyfikator akcji z `kategorie-operacji.ts`, którą uściślenie się wykonuje. */
  idAkcji: string;
}

/**
 * Uściślenia liczone z tekstu.
 *
 * Każde ma pomiar, który mówi, po co je uruchamiać: zwięzłość opiera się na
 * długości zdań, słownictwo na powtórzeniach, rejestr na udziale słów długich,
 * język literacki na stronie biernej rozpoznawanej po formach „został/została/
 * zostały" wraz z imiesłowem. Uściślenie bez pomiaru byłoby przyciskiem
 * z ładną nazwą.
 */
export function uscisleniaTresci(tresc: string, liczniki: LicznikiDokumentu): Uscislenie[] {
  const czytelnosc = policzCzytelnosc(tresc);
  const bierne = (tresc.match(/\b(?:został|została|zostały|zostało|jest|są)\s+\p{L}+[nyt]\w*/giu) ?? [])
    .length;
  const powtorzenia = policzPowtorzenia(tresc);

  return [
    {
      kod: 'zwiezlosc',
      nazwa: 'Zwięzłość',
      wskazanie:
        `Średnio ${czytelnosc.slowNaZdanie} słów na zdanie przy ${liczniki.zdania} zdaniach. ` +
        'Zdania dłuższe niż 20 słów czyta się z potknięciem.',
      idAkcji: 'studio.styl.skrocenie',
    },
    {
      kod: 'slownictwo',
      nazwa: 'Słownictwo',
      wskazanie:
        `Słów powtórzonych co najmniej trzy razy: ${powtorzenia}. Liczone na słowach ` +
        'dłuższych niż cztery znaki, bez form odmienionych — rdzeń nie ma lematyzacji.',
      idAkcji: 'studio.korekta.powtorzenia',
    },
    {
      kod: 'jezyk-literacki',
      nazwa: 'Język literacki',
      wskazanie:
        `Zwrotów w stronie biernej rozpoznanych wzorcem: ${bierne}. Strona bierna nie jest ` +
        'błędem — w piśmie urzędowym bywa właściwa, więc to wskazanie, nie zarzut.',
      idAkcji: 'studio.styl.uproszczenie',
    },
    {
      kod: 'rejestr',
      nazwa: 'Rejestr',
      wskazanie:
        `Słów co najmniej dziesięcioznakowych: ${czytelnosc.udzialSlowDlugich} % treści. ` +
        'Wysoki udział znaczy rejestr urzędowy albo techniczny.',
      idAkcji: 'studio.styl.rejestr',
    },
  ];
}

/** Liczy słowa powtórzone co najmniej trzy razy; słowa krótkie się nie liczą. */
function policzPowtorzenia(tresc: string): number {
  const wystapienia = new Map<string, number>();
  for (const slowo of tresc.toLowerCase().match(/\p{L}{5,}/gu) ?? []) {
    wystapienia.set(slowo, (wystapienia.get(slowo) ?? 0) + 1);
  }
  let powtorzone = 0;
  for (const liczba of wystapienia.values()) {
    if (liczba >= 3) powtorzone += 1;
  }
  return powtorzone;
}
