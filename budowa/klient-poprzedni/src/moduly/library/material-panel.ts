import { MediaOperationKind } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import {
  liczbaZPola,
  nazwaWyniku,
  RODZAJE_PRZETWORZENIA,
  type NarzedziaMaterialu,
} from './material-narzedzia';

/**
 * Materiał dźwiękowy i filmowy: rozpoznanie stoi nad przetworzeniem, bo wycięcie
 * bez znajomości czasu trwania i zmiana rozdzielczości bez znajomości proporcji
 * dają wynik zepsuty.
 */
export interface PanelMaterialu {
  element: HTMLElement;
}

export function utworzPanelMaterialu(narzedzia: NarzedziaMaterialu): PanelMaterialu {
  const odpowiedz = utworzWierszOdpowiedzi();

  const sciezka = poleTekstowe({
    etykieta: 'Ścieżka materiału na dysku Operatora',
    opis:
      'Rdzeń wciąga treść wskazanego pliku do magazynu pod sumą kontrolną — nie dowiązuje ' +
      'jej. Zasobu biblioteki tą drogą wskazać nie można: rodzina media.* rozwiązuje ' +
      'identyfikator zasobu przez magazyn modułu Design, a plik biblioteki leży w innym ' +
      'rejestrze i wróciłby stamtąd odmową.',
  });

  const operacja = poleWyboru(
    {
      etykieta: 'Rodzaj przetworzenia',
      opis: 'Pięć rodzajów kontraktu. Pola poniżej dotyczą tylko tych rodzajów, które ich żądają.',
    },
    RODZAJE_PRZETWORZENIA.map((pozycja) => ({ wartosc: pozycja.kod, etykieta: pozycja.nazwa })),
  );

  const format = poleTekstowe({
    etykieta: 'Format docelowy',
    opis: 'Puste zostawia rozstrzygnięcie rdzeniowi. Format rozstrzyga też rodzaj powstałego zasobu.',
  });

  const odMs = poleTekstowe({
    etykieta: 'Początek fragmentu w milisekundach',
    opis: 'Wyłącznie dla wycięcia fragmentu. Puste znaczy brak wartości, a nie zero.',
  });

  const doMs = poleTekstowe({
    etykieta: 'Koniec fragmentu w milisekundach',
    opis: 'Wyłącznie dla wycięcia fragmentu.',
  });

  const szerokosc = poleTekstowe({
    etykieta: 'Szerokość docelowa',
    opis: 'Wyłącznie dla zmiany rozdzielczości.',
  });

  const wysokosc = poleTekstowe({
    etykieta: 'Wysokość docelowa',
    opis: 'Wyłącznie dla zmiany rozdzielczości.',
  });

  sciezka.kontrolka.dataset['pole'] = 'sciezka';
  operacja.kontrolka.dataset['pole'] = 'operacja';
  format.kontrolka.dataset['pole'] = 'format';
  odMs.kontrolka.dataset['pole'] = 'od-ms';
  doMs.kontrolka.dataset['pole'] = 'do-ms';
  szerokosc.kontrolka.dataset['pole'] = 'szerokosc';
  wysokosc.kontrolka.dataset['pole'] = 'wysokosc';

  const zbadaj = przycisk('Rozpoznaj materiał', 'dn-btn dn-btn--sm dn-btn--zarys');
  zbadaj.dataset['czynnosc'] = 'material-rozpoznanie';
  zbadaj.addEventListener('click', () => void rozpoznaj());

  const przetworz = przycisk('Przetwórz materiał', 'dn-btn dn-btn--sm dn-btn--atrament');
  przetworz.dataset['czynnosc'] = 'material-przetworzenie';
  przetworz.addEventListener('click', () => void przetworzMaterial());

  /** Ścieżka wpisana; brak nie jest odmową rdzenia, tylko niedokończonym formularzem. */
  function wskazanaSciezka(): string {
    return sciezka.kontrolka.value.trim();
  }

  async function rozpoznaj(): Promise<void> {
    const droga = wskazanaSciezka();
    if (droga === '') {
      odpowiedz.pokaz('Wpisz ścieżkę materiału — bez niej nie ma czego zmierzyć.', false);
      return;
    }
    odpowiedz.pokaz(`Rdzeń mierzy materiał ${droga} — trwa odczyt strumieni…`, true);
    const wynik = await narzedzia.zbadaj(droga);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Rozpoznanie materiału', wynik.blad), false);
      return;
    }
    const tresc = wynik.wynik;
    // Czas trwania zerowy jest odpowiedzią, nie brakiem odpowiedzi: tak wraca strumień żywy.
    const czas =
      tresc.durationMs === 0
        ? 'czasu trwania materiał nie ma zapisanego (strumień żywy albo kontener bez nagłówka czasu)'
        : `czas trwania ${tresc.durationMs} ms`;
    odpowiedz.pokaz(
      `Materiał rozpoznany: kontener ${tresc.format}, ${czas}, rozmiar ${tresc.sizeBytes} B. ` +
        `Strumienie: ${tresc.streams}`,
      true,
    );
  }

  async function przetworzMaterial(): Promise<void> {
    const droga = wskazanaSciezka();
    if (droga === '') {
      odpowiedz.pokaz('Wpisz ścieżkę materiału — bez niej nie ma czego przetworzyć.', false);
      return;
    }
    const rodzaj = odczytajRodzaj(operacja.kontrolka.value);
    odpowiedz.pokaz(
      `Rdzeń przetwarza materiał ${droga} — trwa. Wynik będzie nowym zasobem w magazynie.`,
      true,
    );
    przetworz.disabled = true;
    const wynik = await narzedzia.przetworz({
      sciezka: droga,
      operacja: rodzaj,
      format: format.kontrolka.value.trim(),
      ...przedzial(),
      ...wymiary(),
    });
    przetworz.disabled = false;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Przetworzenie materiału', wynik.blad), false);
      return;
    }
    const tresc = wynik.wynik;
    const czas = tresc.durationMs === undefined ? '' : ` Czas trwania wyniku: ${tresc.durationMs} ms.`;
    odpowiedz.pokaz(
      `Materiał przetworzony. Wynik: ${nazwaWyniku(tresc.asset)}.${czas} Bajty leżą ` +
        'w magazynie zasobów pod sumą kontrolną. Zasób nie pojawi się w wykazie żadnego okna ' +
        'Assets Panel, bo okna wyniku nie podajemy: kontrakt żąda tu okna modułu Design, ' +
        'a moduł Library zna wyłącznie okno komunikacji sesji.',
      true,
    );
  }

  /** Przedział czasu wyłącznie dla wycięcia fragmentu — pola pozostałych rodzajów nie jadą. */
  function przedzial(): { odMs?: number; doMs?: number } {
    const od = liczbaZPola(odMs.kontrolka.value);
    const doo = liczbaZPola(doMs.kontrolka.value);
    return {
      ...(od === undefined ? {} : { odMs: od }),
      ...(doo === undefined ? {} : { doMs: doo }),
    };
  }

  function wymiary(): { szerokosc?: number; wysokosc?: number } {
    const sz = liczbaZPola(szerokosc.kontrolka.value);
    const wy = liczbaZPola(wysokosc.kontrolka.value);
    return {
      ...(sz === undefined ? {} : { szerokosc: sz }),
      ...(wy === undefined ? {} : { wysokosc: wy }),
    };
  }

  const pasek = document.createElement('div');
  pasek.className = 'ml-archiwum__pasek';
  pasek.append(zbadaj, przetworz);

  const naglowek = document.createElement('h4');
  naglowek.className = 'ml-metadane__naglowek';
  naglowek.textContent = 'Materiał dźwiękowy i filmowy';

  const granica = document.createElement('p');
  granica.className = 'dn-pole-opis';
  granica.textContent =
    'Miejsca dla tej pary czynności nie nazywa żadne opracowanie modułu — stoi ona tutaj do ' +
    'rozstrzygnięcia Właściciela, przy pozostałej pracy nad treścią zasobu.';

  const element = document.createElement('div');
  element.className = 'ml-material';
  element.append(
    naglowek,
    granica,
    sciezka.element,
    operacja.element,
    format.element,
    odMs.element,
    doMs.element,
    szerokosc.element,
    wysokosc.element,
    pasek,
    odpowiedz.element,
  );

  return { element };
}

/** Przekład wartości selektora na rodzaj przetworzenia; wartość spoza wykazu bierze rodzaj zmiany formatu jako domyślny. */
function odczytajRodzaj(wartosc: string): MediaOperationKind {
  const pozycja = RODZAJE_PRZETWORZENIA.find((wpis) => wpis.kod === wartosc);
  return pozycja === undefined ? MediaOperationKind.Convert : pozycja.kod;
}
