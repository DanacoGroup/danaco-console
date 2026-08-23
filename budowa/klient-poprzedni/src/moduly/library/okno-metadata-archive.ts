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
 * Metadata & Archive Panel — okno zarządcy modułu
 * (`library.metadata-archive-panel`).
 *
 * Okno skupia trzy obszary pracy z zasobem, które w innym układzie rozpraszają
 * się między pozostałe okna: opis, utrwalenie i higienę. Przełącza je selektor
 * obszaru, a nie pasek zakładek — tak opisuje ten element dokumentacja modułu
 * i taki jest sens rozróżnienia: zakładki obiecują trzy równorzędne widoki
 * jednego bytu, a tu Metadane mówią o pliku wskazanym, Archiwum i Higiena — o
 * całym repozytorium.
 *
 * Zależność wejściowa różni się więc obszarem. Metadane bez wskazania pliku nie
 * mają o czym mówić i okno stoi wtedy w stanie pustym. Archiwum i Higiena
 * pracują na wykazie i wskazania nie potrzebują — potrzebują odczytanego
 * wykazu, którego dostarcza Library Explorer.
 *
 * Kod okna nie jest dziś w katalogu okien operacyjnych rdzenia: migracja
 * `migracja_031_okna_modulow.sql` przypisuje modułowi Library cztery okna.
 * Okno buduje się mimo to, bo dokumentacja modułu wymienia je jako siódme okno
 * operacyjne; rozjazd z katalogiem rdzenia należy do rdzenia, nie do klienta.
 */
export interface OknoMetadanych {
  element: HTMLElement;
  odswiez(): void;
  /** Pierwszy odczyt okna: katalog akcji rdzenia. */
  wczytaj(): Promise<void>;
}

/** Trzy obszary panelu; Metadane są obszarem wyjściowym. */
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
  // Warstwa czwarta modułu — reguły, schemat, retencja, nasłuchy, sugestie
  // i porównanie. Stoi przy obszarze Archiwum, bo to sterowanie repozytorium,
  // nie praca nad wskazanym plikiem.
  const administracja = utworzAdministracjeRepozytorium(stan);
  // Trzy czynności arsenału nad treścią z dysku Operatora: rozpoznanie
  // materiału, jego przetworzenie i spakowanie archiwum. Stoją przy obszarze
  // Archiwum, bo obie mówią o utrwaleniu treści, i są od czynności
  // repozytorium rozdzielone własnym nagłówkiem oraz zdaniem o zbiorze,
  // którego dotyczą — wynik idzie do magazynu zasobów, nie do biblioteki.
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

  /**
   * Obszar nieczynny znika z widoku, ale zostaje w drzewie: opis wpisany
   * w formularz Dublin Core przeżywa zajrzenie do Higieny, tak samo jak
   * odpowiedź pokazana przy wywozie manifestu.
   */
  function ustawWidocznosc(): void {
    metadane.element.hidden = obszar !== 'metadane';
    archiwum.element.hidden = obszar !== 'archiwum';
    // Czynności arsenału nie potrzebują ani wskazanego pliku, ani odczytanego
    // wykazu — pracują na ścieżce z dysku Operatora. Widoczność wiąże je
    // z obszarem Archiwum, bo o utrwaleniu treści mówią, ale zależności
    // wejściowej obszaru nie dziedziczą.
    materialPanel.element.hidden = obszar !== 'archiwum';
    pakowanie.element.hidden = obszar !== 'archiwum';
    administracja.element.hidden = obszar !== 'archiwum';
    higiena.element.hidden = obszar !== 'higiena';
  }

  /**
   * Stany obowiązkowe okna czytane z fazy wykazu plików.
   *
   * Okno nie ma własnego odczytu poza przeliczeniem wskaźnika znaczenia:
   * kontrakt nie ma komendy metadanych zasobu, więc wszystko, co panel
   * pokazuje, pochodzi z wykazu przyniesionego przez Library Explorer i
   * z odpowiedzi o treści odłożonej przez File Preview.
   *
   * Zależność wejściowa różni się obszarem, więc stan pusty też się różni —
   * jedno zdanie dla wszystkich trzech obszarów orzekałoby o wskazaniu pliku
   * także tam, gdzie wskazanie nie jest do niczego potrzebne.
   */
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

  // Okno wychodzi z wytwórni w stanie spójnym: bez tego wywołania wszystkie
  // trzy obszary stałyby w drzewie widoczne naraz aż do pierwszego odświeżenia.
  ustawWidocznosc();

  // Przycisku pierwszej akcji okno nie ustawia świadomie: jego cztery stany
  // puste prowadzą do czterech różnych czynności (odśwież wykaz, dodaj plik,
  // wskaż plik, przełącz obszar), a jeden przycisk o stałym napisie kierowałby
  // w trzech z nich w złą stronę. Każdy stan mówi to zdaniem.

  return {
    element: rama.element,
    odswiez,
    wczytaj: () => panel.wczytaj(),
  };
}

/** Przekład wartości selektora na obszar; wartość spoza wykazu bierze Metadane. */
function odczytajObszar(wartosc: string): ObszarPanelu {
  const pozycja = OBSZARY.find((wpis) => wpis.kod === wartosc);
  return pozycja === undefined ? 'metadane' : pozycja.kod;
}
