import { zrodloAkcjiKanalu } from '../okno-komunikacji/katalog-akcji';
import type { KontekstOkna } from '../okno-komunikacji/panel-kontekstu';
import type { Kanal } from '../protokol/kanal';
import { sledzModulOkna } from './modul-okna';
import { utworzRozmowe, type OpcjeRozmowy, type Rozmowa } from './rozmowa';
import type { PolitykaUlotnosci } from './ulotnosc';
import type { WidokZapisu } from './widok-zapisu';
import { utworzWidokRozmowy } from './widok-rozmowy';
import { dolozenieKanalu, nasluchDolozen, zrodloWykazuKanalu } from './zrodlo-wykazu-ukosnika';

/** Interfejs ustawień montażu rozszerza opcje rozmowy o moduł okna, kontekst wykonania oraz politykę pamięci właściwą temu modułowi. */
export interface OpcjeMontazu extends OpcjeRozmowy {
  /** Moduł okna znany powłoce; brak = okno zapyta rdzeń `window.state.get`. */
  modul?: string;
  /** Parametry wykonania okna pokazywane w panelu kontekstu. */
  kontekst?: KontekstOkna;
  /** Polityka pamięci rozmowy właściwa dla kodu modułu, ustalana przez warstwę składającą. */
  politykaModulu?: (kod: string) => PolitykaUlotnosci;
  /** Gotowy pasek zlecenia stawiany w rzędzie akcji pod polem wypowiedzi, przekazany bez ingerencji. */
  pasekZlecenia?: HTMLElement;
}

/** Interfejs zamontowanej rozmowy udostępnia warstwę rozmowy wraz z widokiem, sterowanie modułem oraz trybem widoku transkryptu. */
export interface ZamontowanaRozmowa {
  /** Warstwa rozmowy — wysyłka, przerwanie, wpisy. */
  rozmowa: Rozmowa;
  /** Element widoku wstawiony do wskazanego kontenera. */
  element: HTMLElement;
  /** Przestawia okno na wskazany moduł; historia wątku zostaje. */
  ustawModul(kod: string): void;
  /** Moduł, w którym okno pracuje w tej chwili. */
  modul(): string;
  /** Przestawia widok transkryptu na jeden z czterech trybów, będąc wejściem dla menu sesji powłoki. */
  ustawWidokZapisu(widok: WidokZapisu): void;
  /** Tryb widoku transkryptu, w którym okno pracuje w tej chwili. */
  widokZapisu(): WidokZapisu;
  /** Odłącza subskrypcje kanału i usuwa widok z dokumentu. */
  rozlacz(): void;
}

/** Funkcja osadza rozmowę okna w dokumencie, wiążąc warstwę rozmowy z jej widokiem, śledząc moduł okna wprost z rdzenia zamiast go zakładać. */
export function zamontujRozmowe(
  kontener: HTMLElement,
  kanal: Kanal,
  idOkna: string,
  opcje: OpcjeMontazu = {},
): ZamontowanaRozmowa {
  // Śledzenie modułu poprzedza rozmowę: polityka pamięci zależy od modułu już przy złożeniu.
  const sledzenie = sledzModulOkna(kanal, idOkna, opcje.modul ?? '');
  const polityka = opcje.politykaModulu;
  const rozmowa = utworzRozmowe(kanal, idOkna, {
    ...opcje,
    ...(polityka === undefined ? {} : { ulotnosc: polityka(sledzenie.biezacy()) }),
  });
  const widok = utworzWidokRozmowy(rozmowa, {
    modul: sledzenie.biezacy(),
    katalogAkcji: zrodloAkcjiKanalu(kanal),
    // Wykaz po ukośniku dostaje każde okno z kanałem, bez osobnej nastawy i przełącznika.
    zrodloWykazuUkosnika: zrodloWykazuKanalu(kanal),
    dolozenieNarzedzia: dolozenieKanalu(kanal),
    ...(opcje.kontekst === undefined ? {} : { kontekst: opcje.kontekst }),
    ...(opcje.pasekZlecenia === undefined
      ? {}
      : { pasekZlecenia: opcje.pasekZlecenia }),
  });

  // Dołożenie narzędzia spoza tego okna zgłasza się wyłącznie zdarzeniem dołożenia narzędzia sesji.
  const odsubskrybujDolozenia = nasluchDolozen(kanal, (zdanie) =>
    rozmowa.zglosKomunikat(zdanie),
  );

  const odsubskrybujModul = sledzenie.naZmiane((kod) => {
    widok.ustawModul(kod);
    if (polityka !== undefined) rozmowa.ustawUlotnosc(polityka(kod));
  });

  kontener.append(widok.element);
  widok.ustawOgnisko();

  return {
    rozmowa,
    element: widok.element,
    ustawModul: (kod) => sledzenie.ustaw(kod),
    modul: widok.modul,
    ustawWidokZapisu: widok.ustawWidokZapisu,
    widokZapisu: widok.widokZapisu,
    rozlacz() {
      odsubskrybujModul();
      odsubskrybujDolozenia();
      sledzenie.rozlacz();
      rozmowa.rozlacz();
      widok.element.remove();
    },
  };
}
