import { Command } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { nazwaZasobu } from './karta-zasobu';
import type { StanDesignu } from './stan-designu';

/**
 * Wydania zasobu wskazanego w Assets Panel — trzy komendy, które wynoszą treść
 * poza rdzeń: `design.asset.content.get`, `design.asset.export`
 * i `design.asset.export.batch`.
 *
 * ── Po co osobna droga do treści ────────────────────────────────────────────
 * Pole `uri` zasobu jest ścieżką w systemie plików RDZENIA, a rdzeń biegnie na
 * innej maszynie niż przeglądarka Operatora — przeglądarka nie wczyta spod
 * niego niczego, także wtedy, gdy zasób powstał bez zarzutu. Kafelek
 * z odwołaniem, którego nie da się otworzyć, wygląda dokładnie tak samo jak
 * zasób bez bajtów; to jest szkoda, którą ten moduł ma w historii. Przycisk
 * „Sprawdź treść w magazynie" jest jedyną drogą, która na to pytanie odpowiada.
 *
 * ── Odpowiedź mówi o BAJTACH, nie o powodzeniu ──────────────────────────────
 * Sprawdzenie treści melduje zmierzoną wielkość i sumę kontrolną oddaną przez
 * rdzeń, a nie „udało się". Wydanie melduje wielkość pliku i jego typ treści.
 * Zdanie „gotowe" bez liczby nie odróżniłoby pliku od pustki.
 *
 * ── Formaty wymienia rdzeń, nie to okno ─────────────────────────────────────
 * Lista formatów niesie te, które rdzeń NAPRAWDĘ zapisuje biblioteką
 * wkompilowaną. Formatu, którego rdzeń odmawia (webp na wyjściu, avif),
 * w liście nie ma — a wpisany ręcznie wraca odmową wymieniającą formaty
 * obsługiwane, i to zdanie okno pokazuje bez skracania.
 */
export interface WydaniaZasobu {
  element: HTMLElement;
  /** Przepisuje podpis pod zasób wskazany; woła je odświeżenie panelu. */
  odswiez(): void;
}

/** Formaty wydania, które rdzeń składa biblioteką wkompilowaną. */
const FORMATY_WYDANIA = [
  { wartosc: 'png', etykieta: 'PNG — bezstratny, z przezroczystością' },
  { wartosc: 'jpeg', etykieta: 'JPEG — stratny, z regulacją jakości' },
  { wartosc: 'pdf', etykieta: 'PDF — obraz osadzony w dokumencie' },
  { wartosc: 'ico', etykieta: 'ICO — ikona, bok do 256' },
  { wartosc: 'svg', etykieta: 'SVG — przechodzi bez rasteryzacji' },
] as const;

export function utworzWydaniaZasobu(stan: StanDesignu): WydaniaZasobu {
  const format = poleWyboru(
    {
      etykieta: 'Format wydania',
      opis:
        'Pole format żądania design.asset.export. Rdzeń zapisuje png, jpeg, pdf, ico i svg — ' +
        'formatu bez biblioteki wkompilowanej odmawia, zamiast oddać png pod cudzą nazwą.',
    },
    [...FORMATY_WYDANIA],
  );

  const skale = poleTekstowe({
    etykieta: 'Krotności skal',
    podpowiedz: 'np. 1, 2, 3',
    opis:
      'Pole scales żądania design.asset.export.batch. Puste bierze skalę naturalną. ' +
      'Wydanie pojedyncze bierze pierwszą krotność z tego pola.',
  });

  const jakosc = poleTekstowe({
    etykieta: 'Jakość kompresji',
    podpowiedz: '1–100',
    opis: 'Pole quality. Dotyczy wyłącznie JPEG — PNG jest bezstratny. Puste bierze domyślną.',
  });

  const sprawdz = przycisk('Sprawdź treść w magazynie', 'dn-btn dn-btn--sm dn-btn--zarys');
  const wydaj = przycisk('Wydaj zasób', 'dn-btn dn-btn--sm dn-btn--atrament');
  const wydajPartie = przycisk('Wydaj cały wykaz', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  const pasek = document.createElement('div');
  pasek.className = 'md-czynnosci__pasek';
  pasek.append(sprawdz, wydaj, wydajPartie);

  const element = document.createElement('div');
  element.className = 'md-wydania';
  element.append(format.element, skale.element, jakosc.element, pasek, odpowiedz.element);

  sprawdz.addEventListener('click', () => void sprawdzTresc());
  wydaj.addEventListener('click', () => void wydajWskazany());
  wydajPartie.addEventListener('click', () => void wydajWykaz());

  /**
   * Rozkłada pole krotności na liczby. Wartości niedodatnie i nieliczby
   * wypadają: skala zero nie jest krotnością, a wysłana do rdzenia wróciłaby
   * odmową o żądaniu, którego Operator nie złożył świadomie.
   */
  function krotnosci(): readonly number[] {
    return skale.kontrolka.value
      .split(',')
      .map((czesc) => Number.parseFloat(czesc.trim()))
      .filter((liczba) => Number.isFinite(liczba) && liczba > 0);
  }

  /** Jakość z pola; zero znaczy „nie wskazano" i pole do rdzenia nie idzie. */
  function stopienJakosci(): number {
    const liczba = Number.parseInt(jakosc.kontrolka.value.trim(), 10);
    return Number.isFinite(liczba) && liczba > 0 ? liczba : 0;
  }

  async function sprawdzTresc(): Promise<void> {
    const zasob = stan.wybrany();
    if (zasob === null) {
      odpowiedz.pokaz('Wskaż zasób w wykazie — treść czyta się dla jednego zasobu.', false);
      return;
    }
    odpowiedz.pokaz(`Odczyt treści zasobu „${nazwaZasobu(zasob)}"…`, true);
    const wynik = await stan.czuwanie.prowadz(
      'odczyt treści zasobu',
      // Granica zero znaczy „bez granicy": okno nie narzuca własnego limitu,
      // bo nie wie, jak wielka jest treść, a odmowa za przekroczenie granicy,
      // której Operator nie ustawił, byłaby odmową bez powodu.
      stan.zrodlo.trescZasobu(zasob.id, 0),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Odczyt treści zasobu', wynik.blad), false);
      return;
    }
    const tresc = wynik.wynik;
    odpowiedz.pokaz(
      `Za zasobem „${nazwaZasobu(zasob)}" leży ${tresc.sizeBytes} bajtów typu ` +
        `${tresc.mediaType}; suma kontrolna ${tresc.checksum}.`,
      true,
    );
  }

  async function wydajWskazany(): Promise<void> {
    const zasob = stan.wybrany();
    if (zasob === null) {
      odpowiedz.pokaz('Wskaż zasób w wykazie — wydanie dotyczy jednego zasobu.', false);
      return;
    }
    const wybrane = krotnosci();
    odpowiedz.pokaz(`Wydanie zasobu „${nazwaZasobu(zasob)}"…`, true);
    const wynik = await stan.czuwanie.prowadz(
      'wydanie zasobu',
      stan.zrodlo.wydajZasob({
        idZasobu: zasob.id,
        format: format.kontrolka.value,
        skala: wybrane[0] ?? 0,
        jakosc: stopienJakosci(),
      }),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Wydanie zasobu', wynik.blad), false);
      return;
    }
    const wydanie = wynik.wynik;
    odpowiedz.pokaz(
      `Rdzeń wydał plik „${wydanie.fileName}" — ${wydanie.sizeBytes} bajtów typu ` +
        `${wydanie.mediaType}.`,
      true,
    );
  }

  async function wydajWykaz(): Promise<void> {
    const zasoby = stan.zasoby();
    if (zasoby.length === 0) {
      odpowiedz.pokaz(
        'Wykaz zasobów jest pusty — najpierw odczytaj zasoby albo wnieś pierwszy.',
        false,
      );
      return;
    }
    odpowiedz.pokaz(`Wydanie partii ${zasoby.length} zasobów…`, true);
    const wynik = await stan.czuwanie.prowadz(
      'wydanie partii zasobów',
      stan.zrodlo.wydajPartie({
        idZasobow: zasoby.map((zasob) => zasob.id),
        format: format.kontrolka.value,
        skale: krotnosci(),
        jakosc: stopienJakosci(),
      }),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Wydanie partii zasobów', wynik.blad), false);
      return;
    }
    odpowiedz.pokaz(opisBilansuPartii(wynik.wynik.exported, wynik.wynik.files.length,
      wynik.wynik.rejected ?? []), (wynik.wynik.rejected ?? []).length === 0);
  }

  return {
    element,

    odswiez() {
      const zasob = stan.wybrany();
      wydaj.textContent =
        zasob === null ? 'Wydaj zasób' : `Wydaj zasób „${nazwaZasobu(zasob)}"`;
      wydajPartie.textContent = `Wydaj cały wykaz (${stan.zasoby().length})`;
    },
  };
}

/**
 * Składa zdanie o bilansie partii.
 *
 * Odrzucenia wchodzą do zdania wraz z powodem, a nie samą liczbą: „dwa
 * odrzucono" nie mówi Operatorowi, czy poprawić żądanie, czy zasoby. Wykaz
 * krótszy bez słowa byłby ciszą, której kontrakt tej komendy wprost zabrania.
 */
function opisBilansuPartii(
  wydanych: number,
  plikow: number,
  odrzucone: readonly { assetId: string; reason: string }[],
): string {
  const rdzen = `Rdzeń wydał ${plikow} plików z ${wydanych} zasobów.`;
  if (odrzucone.length === 0) return rdzen;
  const powody = odrzucone.map((pozycja) => `${pozycja.assetId}: ${pozycja.reason}`).join('; ');
  return `${rdzen} Odrzuconych: ${odrzucone.length} — ${powody}`;
}

/** Komendy wykonywane przez ten panel — czytelne dla katalogu funkcji. */
export const KOMENDY_WYDAN = [
  Command.DesignAssetContentGet,
  Command.DesignAssetExport,
  Command.DesignAssetExportBatch,
] as const;
