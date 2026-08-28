import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { nazwaWyniku, type NarzedziaMaterialu } from './material-narzedzia';

/**
 * Panel spakowania archiwum komendą `archive.pack` wydaje pracę jednym plikiem:
 * pakuje wskazany katalog albo plik ze stanowiska operatora i oddaje wynik jako
 * zasób magazynu Designu wraz z jego rozmiarem oraz liczbą pozycji.
 */
export interface PanelPakowania {
  element: HTMLElement;
}

/**
 * Postaci archiwum wymienione w opisie komendy kontraktu. Pole postaci jest
 * napisem, a nie wyliczeniem zamkniętym, więc wartość pusta pozostawia
 * rozstrzygnięcie postaci domyślnej rdzeniowi.
 */
const POSTACI: readonly { wartosc: string; etykieta: string }[] = [
  { wartosc: '', etykieta: 'zip — postać domyślna rdzenia' },
  { wartosc: '7z', etykieta: '7z' },
  { wartosc: 'tar.gz', etykieta: 'tar.gz' },
];

export function utworzPanelPakowania(narzedzia: NarzedziaMaterialu): PanelPakowania {
  const odpowiedz = utworzWierszOdpowiedzi();

  const sciezka = poleTekstowe({
    etykieta: 'Katalog albo plik do spakowania',
    opis:
      'Ścieżka na dysku Operatora. Treść jest wciągana do magazynu pod sumą kontrolną. ' +
      'Zasobów repozytorium biblioteki tą drogą spakować nie można — do tego służy paczka ' +
      'migracyjna i migawka repozytorium powyżej.',
  });

  const postac = poleWyboru(
    {
      etykieta: 'Postać archiwum',
      opis: 'Puste zostawia rozstrzygnięcie rdzeniowi, który bierze zip.',
    },
    POSTACI,
  );

  const nazwa = poleTekstowe({
    etykieta: 'Nazwa archiwum',
    opis: 'Puste zostawia nazwę rdzeniowi. Nazwa własna wraca w opisie wyniku.',
  });

  sciezka.kontrolka.dataset['pole'] = 'sciezka';
  postac.kontrolka.dataset['pole'] = 'postac';
  nazwa.kontrolka.dataset['pole'] = 'nazwa';

  const spakuj = przycisk('Spakuj archiwum', 'dn-btn dn-btn--sm dn-btn--atrament');
  spakuj.dataset['czynnosc'] = 'archiwum-pakowanie';
  spakuj.addEventListener('click', () => void wykonaj());

  async function wykonaj(): Promise<void> {
    const droga = sciezka.kontrolka.value.trim();
    if (droga === '') {
      odpowiedz.pokaz('Wpisz katalog albo plik — bez wskazania nie ma czego spakować.', false);
      return;
    }
    odpowiedz.pokaz(
      `Rdzeń pakuje ${droga} — trwa. Wynikiem będzie jeden plik odłożony w magazynie.`,
      true,
    );
    spakuj.disabled = true;
    const wynik = await narzedzia.spakuj({
      sciezka: droga,
      format: postac.kontrolka.value,
      nazwa: nazwa.kontrolka.value.trim(),
    });
    spakuj.disabled = false;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Spakowanie archiwum', wynik.blad), false);
      return;
    }
    const tresc = wynik.wynik;
    // Archiwum o zerowej liczbie pozycji jest odpowiedzią, nie awarią.
    const zawartosc =
      tresc.entries === 0
        ? 'Archiwum nie ma ani jednej pozycji — wskazany katalog był pusty.'
        : `Pozycji w archiwum: ${tresc.entries}.`;
    odpowiedz.pokaz(
      `Archiwum spakowane: ${nazwaWyniku(tresc.asset)}, rozmiar ${tresc.sizeBytes} B. ` +
        `${zawartosc} Bajty leżą w magazynie zasobów pod sumą kontrolną; zasób nie pojawi się ` +
        'w wykazie żadnego okna Assets Panel, bo okna wyniku nie podajemy — kontrakt żąda tu ' +
        'okna modułu Design, a moduł Library zna wyłącznie okno komunikacji sesji.',
      true,
    );
  }

  const pasek = document.createElement('div');
  pasek.className = 'ml-archiwum__pasek';
  pasek.append(spakuj);

  const naglowek = document.createElement('h4');
  naglowek.className = 'ml-metadane__naglowek';
  naglowek.textContent = 'Spakowanie archiwum z dysku Operatora';

  const element = document.createElement('div');
  element.className = 'ml-pakowanie';
  element.append(naglowek, sciezka.element, postac.element, nazwa.element, pasek, odpowiedz.element);

  return { element };
}
