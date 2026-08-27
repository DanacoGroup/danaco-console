import './library.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzNarzedziaMaterialu } from './material-narzedzia';
import { utworzOknoExplorer } from './okno-library-explorer';
import { utworzOknoMetadanych } from './okno-metadata-archive';
import { utworzOknoEtykiet } from './okno-tags-collections';
import { utworzOknoPodgladu } from './okno-file-preview';
import { utworzOknoWersji } from './okno-versioning-panel';
import { utworzPasekKontekstu } from './pasek-kontekstu';
import { utworzStanBiblioteki } from './stan-biblioteki';
import { utworzStrazOdmow } from './straz-odmow';
import { utworzZrodloBiblioteki } from './zrodlo-biblioteki';
import { utworzZrodloOtoczenia } from './zrodlo-otoczenia';
import { utworzZrodloZnaczenia } from './zrodlo-znaczenia';

/**
 * Moduł Library łączy pięć okien operacyjnych w jednym układzie: Library
 * Explorer wiodące na całą szerokość, a pozostałe cztery jako zarządcy w pasie
 * drugim, dzielące jeden stan i jedno zaznaczenie.
 */
export interface ModulLibrary {
  /** Element osadzany w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Zleca odczyt kontekstu okna, wykazu plików i katalogu akcji. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza subskrypcje zdarzeń rdzenia. */
  rozlacz(): void;
}

export function utworzModulLibrary(kanal: Kanal): ModulLibrary {
  const straz = utworzStrazOdmow(kanal);
  const stan = utworzStanBiblioteki(
    utworzZrodloBiblioteki(kanal, straz),
    utworzZrodloZnaczenia(straz),
  );
  const otoczenie = utworzZrodloOtoczenia(kanal, straz);

  const kontekst = utworzPasekKontekstu(stan, otoczenie);
  const explorer = utworzOknoExplorer(stan, otoczenie);
  const podglad = utworzOknoPodgladu(stan, otoczenie);
  const wersje = utworzOknoWersji(stan, otoczenie);
  const etykiety = utworzOknoEtykiet(stan, otoczenie);
  const metadane = utworzOknoMetadanych(stan, otoczenie, utworzNarzedziaMaterialu(kanal));

  const pasZarzadcow = document.createElement('div');
  pasZarzadcow.className = 'ml-modul__pas ml-modul__pas--zarzadcy';
  pasZarzadcow.append(podglad.element, wersje.element, etykiety.element, metadane.element);

  const element = document.createElement('div');
  element.className = 'ml-modul';
  element.dataset['modul'] = 'library';
  element.setAttribute('aria-label', 'Moduł Library — okna operacyjne');
  element.append(kontekst.element, explorer.element, pasZarzadcow);

  function odswiezWszystkie(): void {
    explorer.odswiez();
    podglad.odswiez();
    wersje.odswiez();
    etykiety.odswiez();
    metadane.odswiez();
  }

  const odsubskrybuj = stan.obserwuj(odswiezWszystkie);
  odswiezWszystkie();

  return {
    element,

    async wczytaj(idSesji) {
      // Kontekst okna idzie pierwszy: bez niego przenoszenie kontekstu nie ma okna źródłowego.
      await kontekst.wczytaj(idSesji);
      // Wykaz plików i katalog akcji panelu metadanych są niezależne, więc idą razem równolegle.
      await Promise.all([explorer.wczytaj(), metadane.wczytaj()]);
    },

    rozlacz() {
      odsubskrybuj();
      explorer.rozlacz();
      stan.rozlacz();
      straz.rozlacz();
    },
  };
}
