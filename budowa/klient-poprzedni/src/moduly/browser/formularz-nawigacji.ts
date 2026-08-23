import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleLogiczne, poleTekstowe, przycisk } from '../../modele/kontrolki-formularza';
import { KLASY_DYMKA, OBJASNIENIA, POZYCJE_BEZ_OBSLUGI } from './etykiety-browser';
import { skutekPobrania } from './skutek-zapisu';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Nawigacja Browser Window — formularz przejścia do strony wraz z odczytem
 * migawki i przewijaniem podglądu.
 *
 * Jedna odpowiedzialność: formularz przejścia i dwie komendy, które z niego
 * wychodzą. Podgląd, pasek zaznaczenia i pasek dolny są osobno.
 *
 * Przewijanie przewija migawkę, nie stronę w rdzeniu: przewinięcie po stronie
 * rdzenia ma własną komendę (`browser.scroll`), której formularz jeszcze nie
 * wywołuje. Zdanie pod przyciskami bierze powód z odczytu wykazu komend rdzenia,
 * więc zmieni się samo w dniu dobudowy obsługi.
 *
 * Adres w zdaniu końcowym pochodzi z migawki, nie z pola. `trim()` przeglądarki
 * nie jest tą samą funkcją co `strings.TrimSpace` rdzenia (JavaScript zdejmuje
 * U+FEFF, Go nie; Go zdejmuje U+0085, JavaScript nie), więc zdanie zbudowane
 * z pola byłoby prawdziwe przypadkiem. `browser.navigate` oddaje migawkę
 * z adresem, pod którym strona została pobrana (`skutek-zapisu.ts`).
 */
export interface FormularzNawigacji {
  element: HTMLElement;
}

/** Czego formularz potrzebuje od ramy okna. */
export interface UjsciaNawigacji {
  przewin(kierunek: -1 | 1): void;
  /** Przewija podgląd do wiersza z pierwszym trafieniem wyszukiwania. */
  pokazWiersz(numer: number): void;
  powiedz(tresc: string, powodzenie: boolean): void;
}

export function utworzFormularzNawigacji(
  stan: StanPrzegladania,
  ujscia: UjsciaNawigacji,
): FormularzNawigacji {
  const adres = poleTekstowe({
    etykieta: 'Adres strony',
    podpowiedz: 'https://',
    opis: 'Przejście idzie komendą browser.navigate w oknie przeglądarki tej sesji.',
  });
  const nowaKarta = poleLogiczne({ etykieta: 'Otwórz w nowej karcie' });
  const zeZrodlem = poleLogiczne({ etykieta: 'Migawka ze źródłem strony' });

  const szukane = poleTekstowe({
    etykieta: 'Szukaj w treści strony',
    podpowiedz: 'fragment tekstu',
  });

  const idz = przycisk('Przejdź', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odswiez = przycisk('Odśwież migawkę', 'dn-btn dn-btn--sm dn-btn--zarys');
  const wGore = przycisk('Przewiń w górę', 'dn-btn dn-btn--sm dn-btn--duch');
  const wDol = przycisk('Przewiń w dół', 'dn-btn dn-btn--sm dn-btn--duch');
  const znajdz = przycisk('Znajdź na stronie', 'dn-btn dn-btn--sm dn-btn--zarys');

  const uwaga = document.createElement('p');
  uwaga.className = 'dn-pole-opis mb-uwaga';
  stan.pokrycie.naOdczyt(() => {
    uwaga.textContent = stan.pokrycie.zdanie(
      POZYCJE_BEZ_OBSLUGI.przewijanie.komenda,
      POZYCJE_BEZ_OBSLUGI.przewijanie.czynnosc,
    );
  });

  const przyciski = document.createElement('div');
  przyciski.className = 'mb-nawigacja__przyciski';
  przyciski.append(idz, odswiez, wGore, wDol, znajdz);

  const element = document.createElement('div');
  element.className = 'mb-nawigacja';
  element.append(
    adres.element,
    utworzDymekObjasnienia(OBJASNIENIA.adres, KLASY_DYMKA),
    nowaKarta.element,
    utworzDymekObjasnienia(OBJASNIENIA.nowaKarta, KLASY_DYMKA),
    zeZrodlem.element,
    utworzDymekObjasnienia(OBJASNIENIA.zrodloStrony, KLASY_DYMKA),
    szukane.element,
    utworzDymekObjasnienia(OBJASNIENIA.szukanieNaStronie, KLASY_DYMKA),
    przyciski,
    uwaga,
  );

  /** Okno przeglądarki albo powód jego braku — jedno sprawdzenie na dwie akcje. */
  function idOknaAlboPowod(): string {
    const idOkna = stan.idOkna();
    if (idOkna === '') ujscia.powiedz(stan.powod(), false);
    return idOkna;
  }

  async function przejdz(): Promise<void> {
    const idOkna = idOknaAlboPowod();
    if (idOkna === '') return;
    const cel = adres.kontrolka.value.trim();
    if (cel === '') {
      ujscia.powiedz('Wskaż adres — rdzeń odmówi przejścia bez niego.', false);
      return;
    }
    ujscia.powiedz(`Przejście do ${cel}…`, true);
    const wynik = await stan.zrodlo.przejdz({
      windowId: idOkna,
      url: cel,
      newTab: nowaKarta.kontrolka.checked,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      ujscia.powiedz(opisOdmowy('Przejście do strony', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    stan.wchlonMigawke(wynik.wynik.snapshot);
    const skutek = skutekPobrania(wynik.wynik.snapshot, cel, 'stronę');
    ujscia.powiedz(skutek.zdanie, skutek.udany);
  }

  /**
   * Odświeżenie idzie tą samą drogą co odczyt przy wejściu do modułu — przez
   * stan modułu, nie własnym wywołaniem `browser.snapshot.get`. Dzięki temu
   * o odmowie dowiaduje się także okno wiodące, a nie sam wiersz odpowiedzi.
   */
  async function pobierzMigawke(): Promise<void> {
    const idOkna = idOknaAlboPowod();
    if (idOkna === '') return;
    ujscia.powiedz('Odczyt migawki strony w toku…', true);
    await stan.zaciagnijMigawke(zeZrodlem.kontrolka.checked);
    const uwagi = stan.powodMigawki();
    ujscia.powiedz(uwagi === '' ? 'Migawka strony odczytana z rdzenia.' : uwagi, uwagi === '');
  }

  /**
   * Wyszukiwanie w treści bieżącej strony. Dzieje się w kliencie na migawce,
   * bo migawka jest całą treścią, którą moduł od rdzenia dostał — komendy
   * wyszukiwania w stronie kontrakt nie niesie.
   */
  function znajdzWTresci(): void {
    const migawka = stan.migawka();
    if (migawka === null) {
      ujscia.powiedz('Nie ma czego przeszukiwać — najpierw pobierz migawkę strony.', false);
      return;
    }
    const fraza = szukane.kontrolka.value.trim();
    if (fraza === '') {
      ujscia.powiedz('Wpisz szukany fragment tekstu.', false);
      return;
    }
    const wiersze = (migawka.text ?? '').split('\n');
    const szukana = fraza.toLowerCase();
    const trafienia = wiersze.filter((linia) => linia.toLowerCase().includes(szukana));
    if (trafienia.length === 0) {
      ujscia.powiedz(`Frazy „${fraza}" nie ma w treści migawki.`, false);
      return;
    }
    const pierwszy = wiersze.findIndex((linia) => linia.toLowerCase().includes(szukana)) + 1;
    ujscia.pokazWiersz(pierwszy);
    ujscia.powiedz(
      `Fraza „${fraza}" w ${trafienia.length} wierszach treści; pierwszy to wiersz ${pierwszy}.`,
      true,
    );
  }

  idz.addEventListener('click', () => void przejdz());
  odswiez.addEventListener('click', () => void pobierzMigawke());
  wGore.addEventListener('click', () => ujscia.przewin(-1));
  wDol.addEventListener('click', () => ujscia.przewin(1));
  znajdz.addEventListener('click', znajdzWTresci);

  return { element };
}
