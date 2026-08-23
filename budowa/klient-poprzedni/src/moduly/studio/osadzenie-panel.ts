import type { LibraryFile, LibraryPreview } from '../../../../shared/contract';
import { poleWyboru } from '../../modele/kontrolki-formularza';
import { osadzenieOpiszPochodzenia, type Pochodzenie } from './osadzenie-pochodzenie';
import type { StanStudio } from './stan-studio';
import {
  wstawieniaOpiszBilans,
  type ZrodloWstawienStudio,
} from './zrodlo-wstawien-studio';

/**
 * Przeglądarka i Biblioteka osadzone w oknie pracy — bez wychodzenia z dokumentu.
 *
 * ── Panel osadzony, nie przełączenie modułu ─────────────────────────────────
 * Wymaganie wprost: dokument zostaje widoczny, a Operatora nie wolno wyrzucać
 * z edytora do innego modułu i kazać mu wracać. Panel ma dwa położenia — obok
 * treści i na całej powierzchni — i wybór należy do Operatora, wedle zasady
 * ogólnej zlecenia.
 *
 * ── Wniesienie wprost do dokumentu ──────────────────────────────────────────
 * „Nie do kolejki i nie do zasobów — do dokumentu, w miejsce kursora." Panel
 * pokazuje treść podglądu i wnosi ją stąd wprost w miejsce kursora, wraz
 * z wierszem pochodzenia. Wciągnięcie strony idzie `studio.ingest.url`, która
 * dotąd kończyła się na kolejce; panel bierze z jej odpowiedzi pole `text`
 * i prowadzi je do dokumentu. Gdy rdzeń tekstu jeszcze nie ma — bo pozycja
 * czeka na rozpoznanie pisma — panel mówi to stanem pozycji, zamiast wnosić
 * pustkę.
 *
 * ── Pochodzenie zapisuje RDZEŃ, nie okno ────────────────────────────────────
 * Komendy `studio.insert.from.library` i `studio.insert.from.web` wnoszą fragment
 * wprost do dokumentu I ODDAJĄ ZAPIS POCHODZENIA — adres albo plik, wersję, czas
 * sięgnięcia i autora. Panel woła je, gdy dostał źródło wstawień, i pokazuje
 * pochodzenie oddane przez rdzeń wraz z bilansem czynności. Wiersz pochodzenia
 * wnoszony do treści zostaje jako droga druga, bo przeżywa wydanie dokumentu do
 * formatu, który zapisu pochodzenia nie niesie — i tak jest opisany
 * w `osadzenie-pochodzenie.ts`.
 *
 * ── Czego panel nie robi ────────────────────────────────────────────────────
 * Nie rysuje strony sieciowej. Migawka `browser.snapshot.get` oddaje adres,
 * tytuł, treść renderowaną i źródło — czyli tekst, nie obraz; zrzut ekranu jest
 * odnośnikiem zasobu, a komendy pobierającej jego bajty do przeglądarki kontrakt
 * nie niesie. Panel pokazuje więc to, co rdzeń naprawdę oddaje, i mówi wprost,
 * czego nie oddaje. Ramka z cudzą stroną wewnątrz okna byłaby drugą
 * przeglądarką, a moduł Browser jest jeden.
 */

/** Czynności panelu zlecane oknu. */
export interface CzynnosciOsadzenia {
  /** Odczytuje pliki Biblioteki wedle frazy. */
  naSzukanieBiblioteki(fraza: string): void;
  /** Zleca podgląd pliku Biblioteki. */
  naPodglad(idPliku: string, strona: number): void;
  /** Wnosi treść podglądu do dokumentu w miejsce kursora. */
  naWniesienieZBiblioteki(plik: LibraryFile, podglad: LibraryPreview, tresc: string): void;
  /** Wciąga stronę oknem Studia i wnosi jej treść do dokumentu. */
  naWciagniecieStrony(adres: string, zObrazami: boolean): void;
  /** Bierze migawkę strony otwartej w module Browser. */
  naMigawke(idOknaPrzegladarki: string): void;
  /** Otwiera adres w module Browser. */
  naOtwarcieStrony(idOknaPrzegladarki: string, adres: string): void;
  /** Wnosi treść migawki do dokumentu w miejsce kursora. */
  naWniesienieZeStrony(tytul: string, adres: string, tresc: string): void;
}

/** Panel wraz z jego odświeżeniem. */
export interface OsadzeniePanel {
  element: HTMLElement;
  /** Pokazuje pliki oddane przez Bibliotekę. */
  pokazPliki(pliki: readonly LibraryFile[]): void;
  /** Pokazuje podgląd pliku wraz z jego treścią. */
  pokazPodglad(plik: LibraryFile, podglad: LibraryPreview): void;
  /** Pokazuje treść strony gotową do wniesienia. */
  pokazStrone(tytul: string, adres: string, tresc: string, skad: string): void;
  /** Przerysowuje wykaz pochodzeń wniesionych fragmentów. */
  pokazPochodzenia(zapisy: readonly Pochodzenie[]): void;
  /** Wypisuje odpowiedź rdzenia — powodzenie albo odmowę nazwaną. */
  pokazOdpowiedz(tresc: string, udana: boolean): void;
  przestawWidocznosc(): void;
  widoczny(): boolean;
  /** Położenie panelu: obok treści albo na całej powierzchni. */
  polozenie(): 'obok' | 'calosc';
}

/** Ile znaków podglądu panel prosi od Biblioteki. */
const ZNAKOW_PODGLADU = 4000;

/**
 * Wniesienie przez rdzeń — źródło i stan modułu podane razem.
 *
 * Nieobowiązkowe: bez niego panel prowadzi wniesienie czynnościami okna, tak jak
 * dotąd. Z nim woła `studio.insert.from.*` i zapis pochodzenia robi rdzeń.
 */
export interface WniesienieRdzeniem {
  stan: StanStudio;
  wstawienia: ZrodloWstawienStudio;
}

export function utworzOsadzeniePanel(
  czynnosci: CzynnosciOsadzenia,
  wniesienie?: WniesienieRdzeniem,
): OsadzeniePanel {
  let otwarty = false;
  let gdzie: 'obok' | 'calosc' = 'obok';
  let plikPodgladany: LibraryFile | null = null;
  let podgladBiezacy: LibraryPreview | null = null;
  let stronaBiezaca: { tytul: string; adres: string; tresc: string } | null = null;

  /* ── Biblioteka ──────────────────────────────────────────────────────────── */

  const frazaBiblioteki = document.createElement('search');
  const poleFrazy = document.createElement('input');
  poleFrazy.type = 'search';
  poleFrazy.className = 'dn-pole-kontrolka';
  poleFrazy.placeholder = 'czego szukać w Bibliotece — nazwa, treść, etykieta';
  poleFrazy.setAttribute('aria-label', 'Fraza wyszukiwania w Bibliotece plików');

  const szukaj = document.createElement('button');
  szukaj.type = 'button';
  szukaj.className = 'dn-btn dn-btn--sm dn-btn--atrament';
  szukaj.textContent = 'Szukaj w Bibliotece';
  szukaj.dataset['czynnosc'] = 'szukaj-biblioteka';
  szukaj.title = 'Idzie komendą library.file.list — rodzina library.* jest zbudowana i zamknięta.';
  szukaj.addEventListener('click', () => czynnosci.naSzukanieBiblioteki(poleFrazy.value.trim()));
  poleFrazy.addEventListener('change', () =>
    czynnosci.naSzukanieBiblioteki(poleFrazy.value.trim()),
  );
  frazaBiblioteki.append(poleFrazy, szukaj);

  const wykazPlikow = document.createElement('ul');
  wykazPlikow.className = 'ms-osadzenie__pliki';

  const podglad = document.createElement('div');
  podglad.className = 'ms-osadzenie__podglad';

  /* ── Strona sieciowa ─────────────────────────────────────────────────────── */

  const adres = document.createElement('input');
  adres.type = 'url';
  adres.className = 'dn-pole-kontrolka';
  adres.placeholder = 'adres strony do sprawdzenia albo wciągnięcia';
  adres.setAttribute('aria-label', 'Adres strony sieciowej');

  const zObrazami = document.createElement('input');
  zObrazami.type = 'checkbox';
  zObrazami.className = 'dn-przelacznik';
  zObrazami.dataset['czynnosc'] = 'z-obrazami';
  zObrazami.setAttribute('aria-label', 'Wciągnąć także obrazy strony do magazynu');

  const etykietaObrazow = document.createElement('label');
  etykietaObrazow.className = 'ms-osadzenie__zawezenie';
  etykietaObrazow.append(zObrazami, document.createTextNode('wciągnij także obrazy'));

  const wciagnij = document.createElement('button');
  wciagnij.type = 'button';
  wciagnij.className = 'dn-btn dn-btn--sm dn-btn--atrament';
  wciagnij.textContent = 'Wciągnij stronę oknem Studia';
  wciagnij.dataset['czynnosc'] = 'wciagnij-strone';
  wciagnij.title =
    'studio.ingest.url pobiera stronę i oczyszcza ją z nawigacji i reklam. Nie wymaga okna ' +
    'modułu Browser, więc działa zawsze. Treść wchodzi WPROST do dokumentu, nie do kolejki.';
  wciagnij.addEventListener('click', () =>
    czynnosci.naWciagniecieStrony(adres.value.trim(), zObrazami.checked),
  );

  const oknoPrzegladarki = document.createElement('input');
  oknoPrzegladarki.type = 'text';
  oknoPrzegladarki.className = 'dn-pole-kontrolka';
  oknoPrzegladarki.placeholder = 'okno modułu Browser — identyfikator';
  oknoPrzegladarki.setAttribute('aria-label', 'Okno modułu Browser, z którego brana jest migawka');

  const otworz = document.createElement('button');
  otworz.type = 'button';
  otworz.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  otworz.textContent = 'Otwórz w przeglądarce';
  otworz.dataset['czynnosc'] = 'otworz-strone';
  otworz.title = 'browser.navigate — strona otwiera się w oknie modułu Browser i oddaje migawkę.';
  otworz.addEventListener('click', () =>
    czynnosci.naOtwarcieStrony(oknoPrzegladarki.value.trim(), adres.value.trim()),
  );

  const migawka = document.createElement('button');
  migawka.type = 'button';
  migawka.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  migawka.textContent = 'Weź migawkę otwartej strony';
  migawka.dataset['czynnosc'] = 'migawka';
  migawka.title =
    'browser.snapshot.get oddaje adres, tytuł i treść renderowaną strony otwartej w module ' +
    'Browser. Zrzut ekranu zostaje odnośnikiem zasobu — komendy pobierającej jego bajty ' +
    'kontrakt nie niesie, więc panel obrazu nie pokazuje.';
  migawka.addEventListener('click', () => czynnosci.naMigawke(oknoPrzegladarki.value.trim()));

  const trescStrony = document.createElement('div');
  trescStrony.className = 'ms-osadzenie__strona';

  /* ── Pochodzenie ─────────────────────────────────────────────────────────── */

  const pochodzenia = document.createElement('ul');
  pochodzenia.className = 'ms-osadzenie__pochodzenia';

  const zdaniePochodzen = document.createElement('p');
  zdaniePochodzen.className = 'dn-pole-opis';

  /* ── Położenie panelu ────────────────────────────────────────────────────── */

  const wyborPolozenia = poleWyboru(
    {
      etykieta: 'Położenie panelu',
      opis: 'Obok treści albo na całej powierzchni — oba równorzędne, wybór należy do Operatora.',
    },
    [
      { wartosc: 'obok', etykieta: 'Obok treści' },
      { wartosc: 'calosc', etykieta: 'Na całej powierzchni' },
    ],
  );

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'dn-pole-opis ms-osadzenie__odpowiedz';

  const element = document.createElement('section');
  element.className = 'ms-osadzenie-zrodel';
  element.hidden = true;
  element.dataset['polozenie'] = gdzie;
  element.setAttribute('aria-label', 'Przeglądarka i Biblioteka osadzone w oknie pracy');
  element.append(
    wyborPolozenia.element,
    czescOsadzenia('Biblioteka plików', [frazaBiblioteki, wykazPlikow, podglad]),
    czescOsadzenia('Strona sieciowa', [
      adres,
      etykietaObrazow,
      wciagnij,
      oknoPrzegladarki,
      otworz,
      migawka,
      trescStrony,
    ]),
    czescOsadzenia('Pochodzenie wniesionych fragmentów', [zdaniePochodzen, pochodzenia]),
    odpowiedz,
  );

  /* ── Wniesienie komendą rdzenia ──────────────────────────────────────────── */

  /** Miejsce wniesienia: koniec zaznaczenia albo koniec treści roboczej. */
  function miejsceKursora(stan: StanStudio): number {
    const zakres = stan.zaznaczenie();
    return zakres === null ? stan.trescRobocza().length : zakres.koniec;
  }

  /** Wnosi plik Biblioteki komendą rdzenia wraz z zapisem pochodzenia. */
  async function wniesZBibliotekiRdzeniem(plik: LibraryFile): Promise<void> {
    if (wniesienie === undefined) return;
    const dokument = wniesienie.stan.dokument();
    if (dokument === null) {
      pokazOdpowiedzWewnetrznie(BRAK_DOKUMENTU_OSADZENIA, false);
      return;
    }
    const wynik = await wniesienie.wstawienia.wniesZBiblioteki({
      documentId: dokument.id,
      offset: miejsceKursora(wniesienie.stan),
      libraryFileId: plik.id,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      pokazOdpowiedzWewnetrznie(
        `Rdzeń odmówił wniesienia pliku „${plik.name}" do dokumentu` +
          (wynik.blad?.message === undefined ? '.' : `: ${wynik.blad.message}`),
        false,
      );
      return;
    }
    const zapis = wynik.wynik.provenance;
    pokazOdpowiedzWewnetrznie(
      `Plik „${plik.name}" wniesiony WPROST do dokumentu na znaku ` +
        `${miejsceKursora(wniesienie.stan)}. ${wstawieniaOpiszBilans(wynik.wynik.balance)} ` +
        `Pochodzenie zapisał rdzeń: ${zapis.id} · znaki ${zapis.rangeStart}–${zapis.rangeEnd}` +
        (zapis.sourceVersion === undefined ? '' : ` · wersja ${zapis.sourceVersion}`) +
        '. Zapis przeżywa zamknięcie okna, bo stoi w rdzeniu, nie w karcie.',
      wynik.wynik.balance.applied > 0,
    );
  }

  /** Wnosi fragment strony komendą rdzenia wraz z zapisem pochodzenia. */
  async function wniesZeStronyRdzeniem(
    adresStrony: string,
    fragment: string,
  ): Promise<void> {
    if (wniesienie === undefined) return;
    const dokument = wniesienie.stan.dokument();
    if (dokument === null) {
      pokazOdpowiedzWewnetrznie(BRAK_DOKUMENTU_OSADZENIA, false);
      return;
    }
    const wynik = await wniesienie.wstawienia.wniesZeStrony({
      documentId: dokument.id,
      offset: miejsceKursora(wniesienie.stan),
      url: adresStrony,
      ...(fragment === '' ? {} : { text: fragment }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      pokazOdpowiedzWewnetrznie(
        `Rdzeń odmówił wniesienia fragmentu strony ${adresStrony} do dokumentu` +
          (wynik.blad?.message === undefined ? '.' : `: ${wynik.blad.message}`),
        false,
      );
      return;
    }
    const zapis = wynik.wynik.provenance;
    pokazOdpowiedzWewnetrznie(
      `Fragment strony ${adresStrony} wniesiony WPROST do dokumentu. ` +
        `${wstawieniaOpiszBilans(wynik.wynik.balance)} Pochodzenie zapisał rdzeń: ${zapis.id}` +
        (zapis.retrievedAt === undefined
          ? ''
          : ` · sięgnięto ${new Date(zapis.retrievedAt).toLocaleString('pl-PL')}`) +
        (zapis.sourceTitle === undefined ? '' : ` · „${zapis.sourceTitle}"`) +
        '. Zapis przeżywa zamknięcie okna, bo stoi w rdzeniu, nie w karcie.',
      wynik.wynik.balance.applied > 0,
    );
  }

  /** Wypisuje odpowiedź w wierszu panelu — używane też z wnętrza panelu. */
  function pokazOdpowiedzWewnetrznie(tresc: string, udana: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.dataset['udana'] = udana ? 'tak' : 'nie';
  }

  wyborPolozenia.kontrolka.addEventListener('change', () => {
    gdzie = wyborPolozenia.kontrolka.value === 'calosc' ? 'calosc' : 'obok';
    element.dataset['polozenie'] = gdzie;
  });

  return {
    element,

    pokazPliki(pliki) {
      if (pliki.length === 0) {
        const puste = document.createElement('li');
        puste.className = 'dn-pole-opis';
        puste.textContent =
          'Biblioteka nie oddała ani jednego pliku przy tej frazie. To odpowiedź rdzenia, nie brak ' +
          'odpowiedzi — zmień frazę albo wnieś plik w module Library.';
        wykazPlikow.replaceChildren(puste);
        return;
      }
      wykazPlikow.replaceChildren(
        ...pliki.map((plik) => {
          const glowa = document.createElement('p');
          glowa.className = 'ms-osadzenie__glowa';
          glowa.textContent =
            `${plik.name} · ${plik.mimeType ?? 'rodzaj nieznany'} · ` +
            `${plik.sizeBytes ?? 0} bajtów · wersja ${plik.versionId ?? 'bez wskazania'}`;

          const wez = document.createElement('button');
          wez.type = 'button';
          wez.className = 'dn-btn dn-btn--sm dn-btn--zarys';
          wez.textContent = 'Pokaż podgląd';
          wez.dataset['plik'] = plik.id;
          wez.addEventListener('click', () => czynnosci.naPodglad(plik.id, 1));

          const pozycja = document.createElement('li');
          pozycja.dataset['plik'] = plik.id;
          pozycja.append(glowa, wez);
          return pozycja;
        }),
      );
    },

    pokazPodglad(plik, wynikPodgladu) {
      plikPodgladany = plik;
      podgladBiezacy = wynikPodgladu;
      const tresc = wynikPodgladu.text ?? '';

      const glowa = document.createElement('p');
      glowa.className = 'ms-osadzenie__glowa';
      glowa.textContent =
        `${plik.name} · podgląd rodzaju ${wynikPodgladu.kind} · strona ` +
        `${wynikPodgladu.page ?? 1} z ${wynikPodgladu.pageCount ?? 1}` +
        (wynikPodgladu.truncated === true ? ' · podgląd SKRÓCONY przez rdzeń' : '');

      if (tresc === '') {
        const bezTekstu = document.createElement('p');
        bezTekstu.className = 'dn-pole-opis';
        bezTekstu.textContent =
          'Podgląd nie niesie treści tekstowej' +
          (wynikPodgladu.imageRef === undefined
            ? '. Do dokumentu nie ma więc czego wnieść — plik jest obrazem albo postacią, ' +
              'z której rdzeń tekstu nie wydobył.'
            : `, tylko odnośnik podglądu graficznego (${wynikPodgladu.imageRef}). Komendy ` +
              'pobierającej bajty zasobu do przeglądarki kontrakt nie niesie, więc obrazu panel ' +
              'nie pokaże ani nie wniesie — to brak nazwany.');
        podglad.replaceChildren(glowa, bezTekstu);
        return;
      }

      const czytanie = document.createElement('pre');
      czytanie.className = 'ms-osadzenie__tresc';
      czytanie.textContent = tresc;

      const wnies = document.createElement('button');
      wnies.type = 'button';
      wnies.className = 'dn-btn dn-btn--sm dn-btn--sygnal';
      wnies.textContent = 'Wnieś do dokumentu w miejsce kursora';
      wnies.dataset['czynnosc'] = 'wnies-biblioteka';
      wnies.title =
        'Treść wchodzi WPROST do dokumentu, wraz z wierszem pochodzenia: nazwa pliku, ' +
        'identyfikator, wersja i suma kontrolna.';
      wnies.addEventListener('click', () => {
        if (plikPodgladany === null || podgladBiezacy === null) return;
        // Droga rdzenia ma pierwszeństwo: pochodzenie zapisane po stronie rdzenia
        // przeżywa zamknięcie okna, a wykaz sesji ginie z kartą.
        if (wniesienie !== undefined) {
          void wniesZBibliotekiRdzeniem(plikPodgladany);
          return;
        }
        czynnosci.naWniesienieZBiblioteki(plikPodgladany, podgladBiezacy, tresc);
      });

      const dalejStrona = document.createElement('button');
      dalejStrona.type = 'button';
      dalejStrona.className = 'dn-btn dn-btn--sm dn-btn--duch';
      dalejStrona.textContent = 'Następna strona podglądu';
      dalejStrona.addEventListener('click', () =>
        czynnosci.naPodglad(plik.id, (wynikPodgladu.page ?? 1) + 1),
      );

      podglad.replaceChildren(glowa, czytanie, wnies, dalejStrona);
    },

    pokazStrone(tytul, adresStrony, tresc, skad) {
      stronaBiezaca = { tytul, adres: adresStrony, tresc };

      const glowa = document.createElement('p');
      glowa.className = 'ms-osadzenie__glowa';
      glowa.textContent = `${tytul} · ${adresStrony} · ${skad} · ${tresc.length} znaków`;

      if (tresc === '') {
        const bezTresci = document.createElement('p');
        bezTresci.className = 'dn-pole-opis';
        bezTresci.textContent =
          'Rdzeń nie oddał treści tej strony. Przy wciągnięciu oknem Studia znaczy to, że pozycja ' +
          'kolejki czeka na rozpoznanie pisma albo że strona nie miała warstwy tekstowej — panel ' +
          'nie wnosi pustki do dokumentu i mówi to wprost.';
        trescStrony.replaceChildren(glowa, bezTresci);
        return;
      }

      const czytanie = document.createElement('pre');
      czytanie.className = 'ms-osadzenie__tresc';
      czytanie.textContent = tresc;

      const wnies = document.createElement('button');
      wnies.type = 'button';
      wnies.className = 'dn-btn dn-btn--sm dn-btn--sygnal';
      wnies.textContent = 'Wnieś do dokumentu w miejsce kursora';
      wnies.dataset['czynnosc'] = 'wnies-strone';
      wnies.addEventListener('click', () => {
        if (stronaBiezaca === null) return;
        if (wniesienie !== undefined) {
          void wniesZeStronyRdzeniem(stronaBiezaca.adres, stronaBiezaca.tresc);
          return;
        }
        czynnosci.naWniesienieZeStrony(
          stronaBiezaca.tytul,
          stronaBiezaca.adres,
          stronaBiezaca.tresc,
        );
      });

      trescStrony.replaceChildren(glowa, czytanie, wnies);
    },

    pokazPochodzenia(zapisy) {
      zdaniePochodzen.textContent = osadzenieOpiszPochodzenia(zapisy);
      pochodzenia.replaceChildren(
        ...zapisy.map((zapis) => {
          const pozycja = document.createElement('li');
          pozycja.dataset['pochodzenie'] = zapis.kod;
          pozycja.dataset['zrodlo'] = zapis.zrodlo;
          pozycja.textContent =
            `${zapis.nazwa} · ${zapis.wskazanie} · ${zapis.znakow} znaków · wstawiono na znaku ` +
            `${zapis.wstawionoNaZnaku} · ${new Date(zapis.czas).toLocaleString('pl-PL')}`;
          return pozycja;
        }),
      );
    },

    pokazOdpowiedz(tresc, udana) {
      odpowiedz.textContent = tresc;
      odpowiedz.dataset['udana'] = udana ? 'tak' : 'nie';
    },

    przestawWidocznosc() {
      otwarty = !otwarty;
      element.hidden = !otwarty;
    },

    widoczny: () => otwarty,
    polozenie: () => gdzie,
  };
}

/** Ile znaków podglądu panel prosi od Biblioteki — jedno miejsce nastawy. */
export const OSADZENIE_ZNAKOW_PODGLADU = ZNAKOW_PODGLADU;

/** Część panelu wraz z jej tytułem. */
function czescOsadzenia(tytul: string, elementy: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('p');
  naglowek.className = 'ms-osadzenie__tytul';
  naglowek.textContent = tytul;

  const sekcja = document.createElement('section');
  sekcja.className = 'ms-osadzenie__czesc';
  sekcja.append(naglowek, ...elementy);
  return sekcja;
}

const BRAK_DOKUMENTU_OSADZENIA =
  'Fragment wchodzi WPROST do dokumentu — wczytaj go albo załóż nowy. Komendy ' +
  'studio.insert.from.library i studio.insert.from.web przyjmują identyfikator dokumentu jako pole ' +
  'obowiązkowe i bez niego nie ma gdzie wnieść treści.';
