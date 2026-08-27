import type { LibraryFile } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { StanBiblioteki } from './stan-biblioteki';

/**
 * Zakładka Higiena panelu Metadata & Archive Panel — raporty stanu repozytorium
 * i przeliczenie wskaźnika znaczenia, liczone z wykazu, który rdzeń już oddał.
 */
export interface HigienaRepozytorium {
  element: HTMLElement;
  /** Przelicza raporty z wykazu w pamięci modułu. */
  odswiez(): void;
}

export function utworzHigienaRepozytorium(stan: StanBiblioteki): HigienaRepozytorium {
  const odpowiedz = utworzWierszOdpowiedzi();

  const pulpit = document.createElement('dl');
  pulpit.className = 'ml-higiena__pulpit';

  const raporty = document.createElement('div');
  raporty.className = 'ml-higiena__raporty';

  const przeliczOdNowa = przycisk(
    'Przelicz wskaźnik znaczenia od nowa',
    'dn-btn dn-btn--sm dn-btn--zarys',
  );
  przeliczOdNowa.dataset['czynnosc'] = 'wskaznik-przebudowa';
  przeliczOdNowa.addEventListener('click', () => void przelicz(true));

  const przeliczPrzyrostowo = przycisk(
    'Dopisz nowe do wskaźnika',
    'dn-btn dn-btn--sm dn-btn--zarys',
  );
  przeliczPrzyrostowo.dataset['czynnosc'] = 'wskaznik-przyrost';
  przeliczPrzyrostowo.addEventListener('click', () => void przelicz(false));

  async function przelicz(odNowa: boolean): Promise<void> {
    odpowiedz.pokaz(
      odNowa
        ? 'Przebudowa wskaźnika znaczenia dla zakresu biblioteki…'
        : 'Dopisywanie nowych dokumentów do wskaźnika znaczenia…',
      true,
    );
    const wynik = await stan.wskaznik.przelicz(odNowa);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Wskaźnik znaczenia', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    const model = wynik.wynik.model === undefined ? '' : ` Model osadzeń: ${wynik.wynik.model}.`;
    // Zero wprowadzonych pozycji przy niepustym wykazie to odpowiedź, nie awaria.
    const uwaga =
      wynik.wynik.indexed === 0 && stan.pliki().length > 0
        ? ' Wykaz nie jest pusty, a wskaźnik nie przyjął ani jednej pozycji — do wskaźnika ' +
          'wchodzą wyłącznie pliki, których treść rdzeń umie odczytać jako tekst.'
        : '';
    odpowiedz.pokaz(
      `Wskaźnik znaczenia: wprowadzono ${wynik.wynik.indexed}, pozycji po przebiegu ` +
        `${wynik.wynik.total}.${model}${uwaga}`,
      true,
    );
  }

  const pulpitRdzenia = document.createElement('p');
  pulpitRdzenia.className = 'dn-pole-opis ml-higiena__pulpit-rdzenia';
  pulpitRdzenia.textContent =
    'Pulpit z rdzenia jeszcze nieodczytany — poniższe liczby opisują odczytaną stronę wykazu.';

  const wynikiRdzenia = document.createElement('div');
  wynikiRdzenia.className = 'ml-higiena__wyniki';

  /** Wypisuje wiersze wyniku czynności rdzenia w bloku pod paskiem. */
  function pokazWiersze(tytul: string, wiersze: readonly string[]): void {
    const naglowek = document.createElement('strong');
    naglowek.className = 'ml-higiena__tytul';
    naglowek.textContent = tytul;
    const lista = document.createElement('ul');
    lista.className = 'ml-higiena__lista';
    for (const wiersz of wiersze) {
      const pozycja = document.createElement('li');
      pozycja.textContent = wiersz;
      lista.append(pozycja);
    }
    wynikiRdzenia.replaceChildren(naglowek, lista);
  }

  const skanuj = przycisk('Skanuj duplikaty (rdzeń)', 'dn-btn dn-btn--sm dn-btn--zarys');
  skanuj.dataset['czynnosc'] = 'duplikaty-skan';
  skanuj.addEventListener('click', () => void skanujDuplikaty());

  async function skanujDuplikaty(): Promise<void> {
    odpowiedz.pokaz('Rdzeń porównuje zbiór: sumy kontrolne, treść i obraz…', true);
    const wynik = await stan.zrodlo.skanujDuplikaty(['exact', 'nearText', 'nearImage']);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Skan duplikatów', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const grupy = wynik.wynik.groups;
    pokazWiersze(
      `Grupy duplikatów: ${wynik.wynik.total}`,
      grupy.map(
        (grupa) =>
          `${grupa.kind}: ${grupa.fileIds.length} zasobów` +
          (grupa.score === undefined ? '' : ` (trafność ${grupa.score}%)`) +
          ` — ${grupa.fileIds.join(', ')}`,
      ),
    );
    odpowiedz.pokaz(
      `Rozpoznano grup: ${wynik.wynik.total}. Pominięto z braku sumy kontrolnej: ` +
        `${wynik.wynik.skippedWithoutChecksum} — brak sumy znaczy niewiadomą, nie brak duplikatu.`,
      true,
    );
  }

  const integralnosc = przycisk('Weryfikuj integralność', 'dn-btn dn-btn--sm dn-btn--zarys');
  integralnosc.dataset['czynnosc'] = 'fixity';
  integralnosc.addEventListener('click', () => void sprawdzIntegralnosc());

  async function sprawdzIntegralnosc(): Promise<void> {
    odpowiedz.pokaz('Rdzeń przelicza sumy kontrolne z bajtów zasobów…', true);
    const wynik = await stan.zrodlo.sprawdzIntegralnosc();
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Weryfikacja integralności', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    const niezgodne = wynik.wynik.results.filter((pozycja) => !pozycja.matched);
    pokazWiersze(
      `Sprawdzono zasobów: ${wynik.wynik.checkedCount}`,
      niezgodne.map(
        (pozycja) =>
          `${pozycja.fileId}: ${pozycja.missing ? 'treści nie ma pod odwołaniem' : 'suma niezgodna'}`,
      ),
    );
    odpowiedz.pokaz(
      `Integralność: sprawdzono ${wynik.wynik.checkedCount}, niezgodnych ` +
        `${wynik.wynik.mismatchedCount}, bez treści ${wynik.wynik.missingCount}.`,
      wynik.wynik.mismatchedCount === 0 && wynik.wynik.missingCount === 0,
    );
  }

  const normalizacjaProbna = przycisk('Normalizuj nazwy (próbnie)', 'dn-btn dn-btn--sm dn-btn--zarys');
  normalizacjaProbna.dataset['czynnosc'] = 'nazwy-probnie';
  normalizacjaProbna.addEventListener('click', () => void normalizuj(true));

  const normalizacjaZapis = przycisk('Normalizuj nazwy (zapisz)', 'dn-btn dn-btn--sm dn-btn--zarys');
  normalizacjaZapis.dataset['czynnosc'] = 'nazwy-zapis';
  normalizacjaZapis.addEventListener('click', () => void normalizuj(false));

  async function normalizuj(probny: boolean): Promise<void> {
    const zaznaczone = stan.zaznaczone();
    odpowiedz.pokaz(
      probny ? 'Przebieg próbny normalizacji nazw…' : 'Zapisywanie znormalizowanych nazw…',
      true,
    );
    const wynik = await stan.zrodlo.normalizujNazwy(zaznaczone, probny);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Normalizacja nazw', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    pokazWiersze(
      probny ? 'Nazwy do zmiany' : 'Nazwy zmienione',
      wynik.wynik.results
        .filter((pozycja) => pozycja.previousName !== pozycja.newName)
        .map((pozycja) => `${pozycja.previousName} → ${pozycja.newName}`),
    );
    if (!probny) {
      for (const plik of wynik.wynik.results) {
        void plik;
      }
      await stan.odczytaj('', '');
    }
    odpowiedz.pokaz(
      `${probny ? 'Przebieg próbny' : 'Zapis'}: nazw ${probny ? 'do zmiany' : 'zmienionych'} ` +
        `${wynik.wynik.changedCount}. Zakres: ${zaznaczone.length === 0 ? 'całe repozytorium' : `zaznaczenie (${zaznaczone.length})`}.`,
      true,
    );
  }

  const dziennik = przycisk('Dziennik audytu', 'dn-btn dn-btn--sm dn-btn--zarys');
  dziennik.dataset['czynnosc'] = 'audyt';
  dziennik.addEventListener('click', () => void pokazDziennik());

  async function pokazDziennik(): Promise<void> {
    odpowiedz.pokaz('Rdzeń oddaje wpisy dziennika audytu…', true);
    const plik = stan.czynny();
    const wynik = await stan.zrodlo.dziennikAudytu(plik === null ? undefined : plik.id);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Dziennik audytu', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    pokazWiersze(
      `Dziennik audytu: wpisów ${wynik.wynik.total}`,
      wynik.wynik.entries.map(
        (wpis) =>
          `${new Date(wpis.at).toLocaleString('pl')} · ${wpis.action} · ${wpis.actor}` +
          (wpis.detail === undefined ? '' : ` — ${wpis.detail}`),
      ),
    );
    odpowiedz.pokaz(`Dziennik audytu: ${wynik.wynik.total} wpisów.`, true);
  }

  /** Odczytuje pulpit stanu liczony przez rdzeń po całym zbiorze. */
  async function odczytajPulpit(): Promise<void> {
    const wynik = await stan.zrodlo.pulpitStanu();
    if (!wynik.udany || wynik.wynik === undefined) {
      pulpitRdzenia.textContent = opisOdmowy(
        'Pulpit stanu',
        wynik.blad?.code,
        wynik.blad?.message,
      );
      return;
    }
    const dane = wynik.wynik.stats;
    pulpitRdzenia.textContent =
      `Rdzeń po całym zbiorze: zasobów czynnych ${dane.fileCount}, w archiwum ` +
      `${dane.archivedCount}, łącznie ${dane.totalBytes} B, osieroconych ${dane.orphanCount}, ` +
      `w grupach duplikatów ${dane.duplicateCount}, bez sumy kontrolnej ` +
      `${dane.missingChecksumCount}.`;
  }

  const pasek = document.createElement('div');
  pasek.className = 'ml-higiena__pasek';
  pasek.append(
    przeliczPrzyrostowo,
    przeliczOdNowa,
    skanuj,
    integralnosc,
    normalizacjaProbna,
    normalizacjaZapis,
    dziennik,
  );

  const naglowekPulpitu = document.createElement('h4');
  naglowekPulpitu.className = 'ml-metadane__naglowek';
  naglowekPulpitu.textContent = 'Pulpit stanu repozytorium — z bieżącego wykazu';

  const naglowekRaportow = document.createElement('h4');
  naglowekRaportow.className = 'ml-metadane__naglowek';
  naglowekRaportow.textContent = 'Raporty higieny';

  const element = document.createElement('div');
  element.className = 'ml-higiena';
  element.append(
    naglowekPulpitu,
    pulpitRdzenia,
    pulpit,
    naglowekRaportow,
    raporty,
    pasek,
    wynikiRdzenia,
    odpowiedz.element,
  );

  return {
    element,

    odswiez() {
      const pliki = stan.pliki();
      void odczytajPulpit();
      pulpit.replaceChildren(...wierszePulpitu(pliki));
      raporty.replaceChildren(
        kartaRaportu(
          'Duplikaty dokładne',
          duplikatyPoSumie(pliki),
          'Pliki o identycznej sumie kontrolnej wyliczonej przez rdzeń.',
          'Ani jednej pary o wspólnej sumie kontrolnej. Pliki bez sumy do raportu nie ' +
            'wchodzą — o ich treści rdzeń nic nie powiedział.',
          stan,
          odpowiedz.pokaz,
        ),
        kartaRaportu(
          'Pliki osierocone',
          osierocone(pliki),
          'Pliki bez etykiety i bez kolekcji — poza wszelkim porządkiem repozytorium.',
          'Każdy plik wykazu ma etykietę albo kolekcję.',
          stan,
          odpowiedz.pokaz,
        ),
        kartaRaportu(
          'Bez sumy kontrolnej',
          pliki.filter((plik) => (plik.checksum ?? '') === '').map((plik) => plik.id),
          'Pliki, dla których rdzeń nie podał sumy kontrolnej — poza raportem duplikatów ' +
            'i poza kontrolą integralności.',
          'Każda pozycja wykazu ma sumę kontrolną.',
          stan,
          odpowiedz.pokaz,
        ),
      );
    },
  };
}

/** Wiersze pulpitu: rozmiar zbioru, rozkład rodzajów treści i modułów wytwórców, liczba etykiet i kolekcji. */
function wierszePulpitu(pliki: readonly LibraryFile[]): HTMLElement[] {
  const bajty = pliki.reduce((suma, plik) => suma + (plik.sizeBytes ?? 0), 0);
  const wpisy: Array<[string, string]> = [
    ['Plików w wykazie', String(pliki.length)],
    ['Rozmiar wg metryki rdzenia', `${bajty} B`],
    ['Rodzaje treści', rozklad(pliki, (plik) => plik.mimeType ?? 'bez rodzaju')],
    ['Moduły wytwórcy', rozklad(pliki, (plik) => plik.sourceModuleId ?? 'bez modułu')],
    ['Etykiet w użyciu', String(new Set(pliki.flatMap((plik) => plik.tags ?? [])).size)],
    ['Kolekcji w użyciu', String(new Set(pliki.flatMap((plik) => plik.collectionIds ?? [])).size)],
  ];
  return wpisy.flatMap(([nazwa, wartosc]) => {
    const podpis = document.createElement('dt');
    podpis.className = 'ml-metadane__podpis';
    podpis.textContent = nazwa;

    const opis = document.createElement('dd');
    opis.className = 'ml-metadane__wartosc';
    opis.textContent = wartosc;
    return [podpis, opis];
  });
}

/** Rozkład liczebności po kluczu; uporządkowany malejąco, a dalej po polsku, gdy liczby są sobie równe. */
function rozklad(
  pliki: readonly LibraryFile[],
  klucz: (plik: LibraryFile) => string,
): string {
  if (pliki.length === 0) return 'wykaz jest pusty';
  const liczby = new Map<string, number>();
  for (const plik of pliki) {
    const kod = klucz(plik);
    liczby.set(kod, (liczby.get(kod) ?? 0) + 1);
  }
  return [...liczby]
    .sort((pierwszy, drugi) => drugi[1] - pierwszy[1] || pierwszy[0].localeCompare(drugi[0], 'pl'))
    .map(([kod, liczba]) => `${kod} ${liczba}`)
    .join(' · ');
}

/**
 * Duplikaty dokładne — pliki dzielące sumę kontrolną z co najmniej jednym innym; plik
 * bez sumy ma własną kartę, bo brak sumy znaczy nie wiadomo, nie bez duplikatu.
 */
function duplikatyPoSumie(pliki: readonly LibraryFile[]): string[] {
  const poSumie = new Map<string, string[]>();
  for (const plik of pliki) {
    const suma = plik.checksum ?? '';
    if (suma === '') continue;
    const zbior = poSumie.get(suma);
    if (zbior === undefined) poSumie.set(suma, [plik.id]);
    else zbior.push(plik.id);
  }
  return [...poSumie.values()].filter((zbior) => zbior.length > 1).flat();
}

/** Pliki bez etykiety i bez kolekcji — poza porządkiem repozytorium, nieodnalezione żadną istniejącą ścieżką. */
function osierocone(pliki: readonly LibraryFile[]): string[] {
  return pliki
    .filter((plik) => (plik.tags ?? []).length === 0 && (plik.collectionIds ?? []).length === 0)
    .map((plik) => plik.id);
}

/**
 * Karta raportu wraz z odsyłaczem do wykazu: naciśnięcie zawęża Library Explorer do
 * wskazanego zbioru, odwracalnie, aż do nowego odczytu wykazu.
 */
function kartaRaportu(
  tytul: string,
  kody: readonly string[],
  opisZbioru: string,
  zdanieZgodne: string,
  stan: StanBiblioteki,
  pokaz: (tresc: string, powodzenie: boolean) => void,
): HTMLElement {
  const naglowek = document.createElement('strong');
  naglowek.className = 'ml-higiena__tytul';
  naglowek.textContent = tytul;

  const licznik = document.createElement('span');
  licznik.className = 'ml-higiena__licznik';
  licznik.textContent = String(kody.length);

  const opis = document.createElement('span');
  opis.className = 'ml-higiena__opis';
  opis.textContent = kody.length === 0 ? zdanieZgodne : opisZbioru;

  const pokazWWykazie = przycisk('Pokaż w wykazie', 'dn-btn dn-btn--sm dn-btn--duch');
  pokazWWykazie.dataset['czynnosc'] = 'zawezenie';
  pokazWWykazie.addEventListener('click', () => {
    if (kody.length === 0) {
      pokaz(`Raport „${tytul}" nie wskazuje ani jednej pozycji — nie ma czym zawęzić wykazu.`, false);
      return;
    }
    stan.ustawZawezenie({ kody: [...kody], opis: tytul.toLowerCase() });
    pokaz(
      `Library Explorer zawężony do raportu „${tytul}" — pozycji ${kody.length}. ` +
        'Zawężenie zdejmuje przycisk w pasku narzędzi Explorera albo każdy nowy odczyt wykazu.',
      true,
    );
  });

  const element = document.createElement('div');
  element.className = 'dn-karta ml-higiena__karta';
  element.dataset['zgodny'] = kody.length === 0 ? 'tak' : 'nie';
  element.append(naglowek, licznik, opis, pokazWWykazie);
  return element;
}
