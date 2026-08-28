import { Command, type DesignAsset } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { dopnijDymek } from './dymek-objasnienia';
import { nazwaZasobu } from './karta-zasobu';
import { rozbijEtykiety } from './przyciecie-pol';
import { skutekEtykiet } from './skutek-designu';
import type { StanDesignu } from './stan-designu';

/**
 * Nadanie etykiet zasobowi wskazanemu w Assets Panel: zbiera pełny zestaw etykiet, który zastępuje
 * poprzedni, nie dokłada się do niego.
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

    // Czuwanie pilnuje, by zerwane gniazdo nie zostawiło wiersza bez końca; cisza nie jest odmową.
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
    // Zasób wchodzi do zbioru modułu tą samą drogą co wynik generowania, bez ponownego odczytu listy.
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
      // Pole przepisuje się wyłącznie przy zmianie zasobu; inaczej odświeżenie kasuje wpisywane etykiety.
      if (id === wskazany) return;
      wskazany = id;
      pole.kontrolka.value = (zasob?.tags ?? []).join(', ');
      odpowiedz.wyczysc();
    },
  };
}

/** Zdanie o zleceniu nadania etykiet — mówi, ile etykiet idzie i że zestaw zastępuje poprzedni w całości. */
function zdanieZlecenia(zasob: DesignAsset, etykiety: readonly string[]): string {
  if (etykiety.length === 0) {
    return `Zdejmowanie wszystkich etykiet zasobu „${nazwaZasobu(zasob)}"…`;
  }
  return `Ustawianie ${etykiety.length} etykiet zasobu „${nazwaZasobu(zasob)}" (zestaw pełny)…`;
}
