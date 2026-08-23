import { Command, type DesignAsset } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { dopnijDymek } from './dymek-objasnienia';
import { nazwaZasobu } from './karta-zasobu';
import { rozbijEtykiety } from './przyciecie-pol';
import { skutekEtykiet } from './skutek-designu';
import type { StanDesignu } from './stan-designu';

/**
 * Nadanie etykiet zasobowi wskazanemu w Assets Panel.
 *
 * Zbiera pełny zestaw etykiet jednego zasobu i oddaje go rdzeniowi komendą
 * `design.asset.tag.set`. Kontrolka stoi przy filtrze wykazu, bo pole zasila
 * to samo zawężenie co `tags` żądania `design.asset.list`.
 *
 * Zestaw zastępuje poprzedni, nie dokłada się do niego — tak stanowi kontrakt.
 * Pole pokazuje więc stan bieżący zasobu przy każdym przewybraniu, a
 * wyczyszczenie pola zdejmuje wszystkie etykiety; pusty zestaw jest drogą
 * udaną, nie odmową.
 *
 * Etykietę przycina `przyciecie-pol.ts`, nie `trim()`: rdzeń zapisuje etykietę
 * co do znaku, a `trim()` przeglądarki zostawia NEL (U+0085), więc zestaw
 * złożony z samego NEL dałby etykietę niewidoczną o długości jednego znaku.
 * Przycięcie jest wspólne z filtrem wykazu, żeby etykieta nadana i szukana
 * były jednym ciągiem znaków.
 *
 * Rdzeń oddaje etykiety odczytane z bazy po zapisie, nie echo żądania, więc
 * zdanie skutku stoi na odpowiedzi, a różnica wobec zestawu zamówionego jest
 * odmową (`skutek-designu.ts`).
 */
export interface NadanieEtykiet {
  element: HTMLElement;
  /** Wstawia do pola etykiety zasobu wskazanego; woła je odświeżenie panelu. */
  odswiez(): void;
}

export function utworzNadanieEtykiet(stan: StanDesignu): NadanieEtykiet {
  const pole = poleTekstowe({
    etykieta: 'Etykiety zasobu wskazanego',
    podpowiedz: 'kampania, zima',
    opis:
      'Rozdzielone przecinkiem. Zestaw ZASTĘPUJE poprzedni — puste pole zdejmuje ' +
      'wszystkie etykiety. Te same etykiety zawężają odczyt w filtrze powyżej. ' +
      'Brzegowe znaki białe są zdejmowane; rdzeń zapisuje etykietę co do znaku.',
  });
  dopnijDymek(pole.element, `Pole tags żądania ${Command.DesignAssetTagSet}. Pełny zestaw, nie dokładka.`);

  const nadaj = przycisk('Ustaw etykiety zasobu', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odpowiedz = utworzWierszOdpowiedzi();

  const element = document.createElement('div');
  element.className = 'md-etykietowanie';
  element.append(pole.element, nadaj, odpowiedz.element);

  /** Zasób wskazany w chwili odświeżenia — pamiętany, by rozpoznać przewybranie. */
  let wskazany: string | null = null;

  async function nadajEtykiety(): Promise<void> {
    const zasob = stan.wybrany();
    if (zasob === null) {
      odpowiedz.pokaz('Wskaż zasób w wykazie — etykiety nadaje się jednemu zasobowi.', false);
      return;
    }
    const etykiety = rozbijEtykiety(pole.kontrolka.value);
    odpowiedz.pokaz(zdanieZlecenia(zasob, etykiety), true);

    // Czuwanie pilnuje, by zerwane gniazdo nie zostawiło wiersza na zdaniu
    // „Ustawianie … etykiet…" bez końca. Cisza kanału nie jest odmową nadania
    // — rdzeń mógł etykiety zapisać (`czuwanie-rdzenia.ts`).
    const wynik = await stan.czuwanie.prowadz(
      'nadanie etykiet',
      stan.zrodlo.ustawEtykiety({ idZasobu: zasob.id, etykiety }),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Nadanie etykiet', wynik.blad), false);
      return;
    }
    // Zasób wchodzi do jednego zbioru modułu tą samą drogą co wynik
    // generowania — panel metadanych i wykaz zobaczą nowe etykiety bez
    // ponownego odczytu całej listy.
    //
    // Wciągnięcie gubi pole `promptId`: odpowiedź `design.asset.tag.set` niesie
    // te same klucze co `design.asset.list`, bez `promptId`, bo repozytorium nie
    // ma przekładu klucza wiersza promptu na kod kontraktu
    // (`adapter_modul_design_etykiety.go`). Zasób świeżo wygenerowany traci więc
    // po otagowaniu wskazanie promptu w Preview Window; okno mówi o tym wprost
    // w wierszu „Prompt źródłowy" (`plyta-podgladu.ts`) zamiast dorabiać wartość
    // z poprzedniej odpowiedzi.
    stan.wchlon(wynik.wynik.asset);
    const skutek = skutekEtykiet(etykiety, wynik.wynik.asset, nazwaZasobu(wynik.wynik.asset));
    odpowiedz.pokaz(skutek.zdanie, skutek.udany);
  }

  nadaj.addEventListener('click', () => void nadajEtykiety());

  return {
    element,

    odswiez() {
      const zasob = stan.wybrany();
      const id = zasob?.id ?? null;
      // Pole przepisuje się wyłącznie przy zmianie zasobu; inaczej odświeżenie
      // wywołane czymkolwiek innym kasowałoby etykiety właśnie wpisywane.
      if (id === wskazany) return;
      wskazany = id;
      pole.kontrolka.value = (zasob?.tags ?? []).join(', ');
      odpowiedz.wyczysc();
    },
  };
}

/** Zdanie o zleceniu — mówi, ile etykiet idzie i że zestaw zastępuje poprzedni. */
function zdanieZlecenia(zasob: DesignAsset, etykiety: readonly string[]): string {
  if (etykiety.length === 0) {
    return `Zdejmowanie wszystkich etykiet zasobu „${nazwaZasobu(zasob)}"…`;
  }
  return `Ustawianie ${etykiety.length} etykiet zasobu „${nazwaZasobu(zasob)}" (zestaw pełny)…`;
}
