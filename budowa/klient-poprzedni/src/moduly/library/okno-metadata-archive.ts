import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import { poleWyboru } from '../../modele/kontrolki-formularza';
import { utworzAdministracjeRepozytorium } from './administracja-repozytorium';
import { utworzArchiwumRepozytorium } from './archiwum-repozytorium';
import { utworzPanelPakowania } from './archiwum-pakowanie';
import { utworzPanelMaterialu } from './material-panel';
import type { NarzedziaMaterialu } from './material-narzedzia';
import { BEZ_KOMENDY_METADANE } from './etykiety-biblioteki';
import { utworzHigienaRepozytorium } from './higiena-repozytorium';
import { utworzMetadanePliku } from './metadane-pliku';
import { utworzPanelAkcji } from './panel-akcji';
import type { StanBiblioteki } from './stan-biblioteki';
import { utworzStanOkna } from './stan-okna';
import type { ZrodloOtoczenia } from './zrodlo-otoczenia';

/**
 * Okno Metadata & Archive Panel zarządcy skupia trzy obszary pracy z zasobem:
 * Metadane dotyczą pliku wskazanego, a Archiwum i Higiena całego repozytorium,
 * przełączane selektorem obszaru, nie zakładkami.
 */
export interface OknoMetadanych {
  element: HTMLElement;
  odswiez(): void;
  /** Pierwszy odczyt okna: katalog akcji rdzenia. */
  wczytaj(): Promise<void>;
}

/** Trzy obszary panelu okna zarządcy; Metadane są obszarem wyjściowym, widocznym od razu po otwarciu okna. */
type ObszarPanelu = 'metadane' | 'archiwum' | 'higiena';

const OBSZARY: ReadonlyArray<{ kod: ObszarPanelu; etykieta: string }> = [
  { kod: 'metadane', etykieta: 'Obszar: Metadane' },
  { kod: 'archiwum', etykieta: 'Obszar: Archiwum' },
  { kod: 'higiena', etykieta: 'Obszar: Higiena' },
];

export function utworzOknoMetadanych(
  stan: StanBiblioteki,
  otoczenie: ZrodloOtoczenia,
  material: NarzedziaMaterialu,
): OknoMetadanych {
  const rama = utworzRameOkna({
    kod: 'metadata-archive-panel',
    tytul: 'Metadata & Archive Panel',
    rola: 'zarządca',
    modul: 'Library',
    dodatkiNaglowka: [
      utworzDymekObjasnienia(
        'Opis, utrwalenie i higiena zasobu. Metadane dotyczą pliku wskazanego w Library ' +
          'Explorer; Archiwum i Higiena — całego wykazu, więc wskazania nie wymagają.',
        { powloka: 'ml-dymek', znak: 'ml-dymek__znak' },
      ),
    ],
  });
  const okno = utworzStanOkna();
  const panel = utworzPanelAkcji(otoczenie, stan, BEZ_KOMENDY_METADANE);

  const metadane = utworzMetadanePliku(stan);
  const archiwum = utworzArchiwumRepozytorium(stan);
  const higiena = utworzHigienaRepozytorium(stan);
  // Warstwa czwarta stoi przy obszarze Archiwum, bo to sterowanie repozytorium, nie praca nad plikiem.
  const administracja = utworzAdministracjeRepozytorium(stan);
  // Trzy czynności arsenału stoją przy Archiwum, bo mówią o utrwaleniu treści, nie o bibliotece.
  const materialPanel = utworzPanelMaterialu(material);
  const pakowanie = utworzPanelPakowania(material);

  let obszar: ObszarPanelu = 'metadane';

  const selektor = poleWyboru(
    {
      etykieta: 'Obszar panelu',
      opis:
        'Metadane opisują plik wskazany w Explorerze; Archiwum i Higiena obejmują cały ' +
        'wykaz i nie potrzebują wskazania.',
    },
    OBSZARY.map((pozycja) => ({ wartosc: pozycja.kod, etykieta: pozycja.etykieta })),
  );
  selektor.kontrolka.dataset['ster'] = 'obszar-panelu';
  selektor.kontrolka.addEventListener('change', () => {
    obszar = odczytajObszar(selektor.kontrolka.value);
    odswiez();
  });

  okno.tresc.append(
    metadane.element,
    archiwum.element,
    materialPanel.element,
    pakowanie.element,
    administracja.element,
    higiena.element,
  );
  rama.narzedzia.append(selektor.element);
  rama.cialo.append(okno.element, panel.element);

  // Obszar nieczynny znika z widoku, ale zostaje w drzewie: wpisany opis przeżywa zmianę obszaru.
  function ustawWidocznosc(): void {
    metadane.element.hidden = obszar !== 'metadane';
    archiwum.element.hidden = obszar !== 'archiwum';
    // Czynności arsenału pracują na ścieżce z dysku, lecz widoczność wiąże je z obszarem Archiwum.
    materialPanel.element.hidden = obszar !== 'archiwum';
    pakowanie.element.hidden = obszar !== 'archiwum';
    administracja.element.hidden = obszar !== 'archiwum';
    higiena.element.hidden = obszar !== 'higiena';
  }

  // Stany okna czytane z fazy wykazu; zależność wejściowa różni się obszarem, więc różni się stan pusty.
  function ustawStan(): void {
    if (stan.faza() === 'odczyt') {
      okno.ladowanie('Rdzeń odczytuje wykaz plików — panel opisuje to, co z niego przyszło.');
      return;
    }
    if (stan.faza() === 'blad') {
      okno.blad(stan.powod());
      return;
    }
    if (stan.faza() === 'spoczynek') {
      okno.puste(
        'Wykaz plików nieodczytany',
        'Panel opisuje zasoby przyniesione przez Library Explorer, a te nie przyszły ' +
          'jeszcze z rdzenia. Odśwież wykaz w oknie wiodącym.',
      );
      return;
    }
    if (stan.pliki().length === 0) {
      okno.puste(
        'Repozytorium jest puste',
        'Rdzeń nie ma ani jednego pliku, więc nie ma czego opisać, utrwalić ani sprawdzić.',
      );
      return;
    }
    if (obszar === 'metadane' && stan.czynny() === null) {
      okno.puste(
        'Nic nie wybrano',
        'Zakładka Metadane opisuje jeden zasób. Wskaż plik w Library Explorer albo ' +
          'przełącz obszar na Archiwum lub Higienę — te obejmują cały wykaz.',
      );
      return;
    }
    okno.gotowe();
  }

  function odswiez(): void {
    selektor.kontrolka.value = obszar;
    ustawWidocznosc();
    metadane.odswiez();
    archiwum.odswiez();
    administracja.odswiez();
    higiena.odswiez();
    ustawStan();
  }

  // Okno wychodzi z wytwórni w stanie spójnym, inaczej trzy obszary stałyby widoczne do odświeżenia.
  ustawWidocznosc();

  // Przycisku pierwszej akcji okno nie ustawia świadomie, bo cztery stany prowadzą do różnych czynności.

  return {
    element: rama.element,
    odswiez,
    wczytaj: () => panel.wczytaj(),
  };
}

/** Przekład wartości selektora na obszar panelu; wartość spoza wykazu bierze obszar Metadane jako domyślny. */
function odczytajObszar(wartosc: string): ObszarPanelu {
  const pozycja = OBSZARY.find((wpis) => wpis.kod === wartosc);
  return pozycja === undefined ? 'metadane' : pozycja.kod;
}
