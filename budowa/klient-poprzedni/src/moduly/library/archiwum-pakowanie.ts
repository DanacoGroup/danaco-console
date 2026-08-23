import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { nazwaWyniku, type NarzedziaMaterialu } from './material-narzedzia';

/**
 * Spakowanie archiwum (`archive.pack`) — wydanie pracy Operatorowi jednym
 * plikiem.
 *
 * Czynność stoi w obszarze Archiwum panelu Metadata & Archive Panel, obok
 * utrwalenia i paczki migracyjnej, i jest od nich rozdzielona zdaniem, bo
 * pracuje na innym zbiorze:
 *
 *   — `library.preservation.run` i `library.package.export` obejmują zasoby
 *     REPOZYTORIUM biblioteki i oddają wynik jako zasób biblioteki,
 *   — `archive.pack` pakuje katalog albo plik z dysku Operatora i oddaje wynik
 *     jako zasób magazynu Designu.
 *
 * Nazwanie jednego drugim byłoby obietnicą, że spakowana została biblioteka.
 * Wykaz zasobów (`assetIds`) tą drogą nie jedzie z tego samego powodu, z jakiego
 * nie jedzie w `media.*`: identyfikator pliku biblioteki wraca z magazynu
 * Designu odmową `not_found`, a przycisk pewnej odmowy nie jest funkcją.
 *
 * Formatu archiwum kontrakt nie zamyka wyliczeniem — pole jest napisem, a brak
 * bierze zip. Wykaz trzech postaci pochodzi z opisu komendy w kontrakcie
 * („zip, 7z, tar.gz"), nie z domysłu.
 */
export interface PanelPakowania {
  element: HTMLElement;
}

/** Postaci archiwum wymienione w opisie komendy; brak wartości bierze zip. */
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
    // Archiwum o zerowej liczbie pozycji jest odpowiedzią, nie awarią: tak
    // wraca spakowany katalog pusty. Zdanie nazywa to wprost, bo plik istnieje
    // i Operator ma prawo wiedzieć, że nic w nim nie ma.
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
