import { zrodloAkcjiKanalu } from '../okno-komunikacji/katalog-akcji';
import type { KontekstOkna } from '../okno-komunikacji/panel-kontekstu';
import type { Kanal } from '../protokol/kanal';
import { sledzModulOkna } from './modul-okna';
import { utworzRozmowe, type OpcjeRozmowy, type Rozmowa } from './rozmowa';
import type { PolitykaUlotnosci } from './ulotnosc';
import type { WidokZapisu } from './widok-zapisu';
import { utworzWidokRozmowy } from './widok-rozmowy';
import { dolozenieKanalu, nasluchDolozen, zrodloWykazuKanalu } from './zrodlo-wykazu-ukosnika';

/** Ustawienia montażu: rozmowa oraz to, czym okno różni się w module. */
export interface OpcjeMontazu extends OpcjeRozmowy {
  /** Moduł okna znany powłoce; brak = okno zapyta rdzeń `window.state.get`. */
  modul?: string;
  /** Parametry wykonania okna pokazywane w panelu kontekstu. */
  kontekst?: KontekstOkna;
  /**
   * Polityka pamięci rozmowy dla kodu modułu (`ulotnosc.ts`).
   *
   * Podaje ją warstwa składająca, bo tylko ona zna moduły i ich konteksty
   * robocze. Pominięta znaczy „każdy moduł z pamięcią".
   */
  politykaModulu?: (kod: string) => PolitykaUlotnosci;
  /**
   * Gotowy pasek zlecenia — rząd sterów koperty, stawiany w rzędzie akcji pod
   * polem wypowiedzi.
   *
   * Montaż go nie buduje i nie zagląda do środka — przepuszcza element do widoku
   * nietknięty. Buduje go ten, kto ma komplet sterowania okna
   * (`aplikacja/wiazanie-gniazda.ts`); podgląd rozmowy go nie ma i pomija.
   */
  pasekZlecenia?: HTMLElement;
}

/** Rozmowa wraz z jej widokiem, gotowa do osadzenia w oknie. */
export interface ZamontowanaRozmowa {
  /** Warstwa rozmowy — wysyłka, przerwanie, wpisy. */
  rozmowa: Rozmowa;
  /** Element widoku wstawiony do wskazanego kontenera. */
  element: HTMLElement;
  /** Przestawia okno na wskazany moduł; historia wątku zostaje. */
  ustawModul(kod: string): void;
  /** Moduł, w którym okno pracuje w tej chwili. */
  modul(): string;
  /**
   * Przestawia widok transkryptu — cztery tryby z `WIDOKI_ZAPISU`.
   *
   * To jest wejście dla menu sesji powłoki. Powłoka nie musi znać ani wpisów,
   * ani warstw: podaje kod trybu i dostaje przerysowany zapis. Wywołanie nigdy
   * nie jest odrzucane i nigdy nie kosztuje rundy do rdzenia.
   */
  ustawWidokZapisu(widok: WidokZapisu): void;
  /** Tryb widoku transkryptu, w którym okno pracuje w tej chwili. */
  widokZapisu(): WidokZapisu;
  /** Odłącza subskrypcje kanału i usuwa widok z dokumentu. */
  rozlacz(): void;
}

/**
 * Osadza rozmowę okna w dokumencie.
 *
 * Jedna odpowiedzialność: powiązanie warstwy rozmowy z jej widokiem i wstawienie
 * całości w kontener. To jest punkt styku warstwy rozmowy z powłoką — powłoka
 * nie musi znać ani widoku, ani kontraktu.
 *
 * Moduł okna jest śledzony, nie zakładany: okno pyta rdzeń o swój moduł
 * i słucha zmian (`window.changed`), więc `workspace.enter` wykonane gdziekolwiek
 * indziej przestawia to okno samo. Powłoka może też przestawić okno wprost przez
 * `ustawModul`, gdy zna wynik komendy wcześniej.
 *
 * Ognisko ląduje na polu wypowiedzi od razu po złożeniu.
 */
export function zamontujRozmowe(
  kontener: HTMLElement,
  kanal: Kanal,
  idOkna: string,
  opcje: OpcjeMontazu = {},
): ZamontowanaRozmowa {
  // Śledzenie modułu idzie przed rozmową: polityka pamięci zależy od modułu,
  // a rozmowa musi ją znać już w chwili złożenia, bo to ona rozstrzyga, czy okno
  // w ogóle woła `message.list`. Odwrotna kolejność dałaby okno modułu bez
  // pamięci sesyjnej, które odtwarza wątek z rdzenia, zanim się dowie, że nie
  // miało prawa go odtworzyć.
  const sledzenie = sledzModulOkna(kanal, idOkna, opcje.modul ?? '');
  const polityka = opcje.politykaModulu;
  const rozmowa = utworzRozmowe(kanal, idOkna, {
    ...opcje,
    ...(polityka === undefined ? {} : { ulotnosc: polityka(sledzenie.biezacy()) }),
  });
  const widok = utworzWidokRozmowy(rozmowa, {
    modul: sledzenie.biezacy(),
    katalogAkcji: zrodloAkcjiKanalu(kanal),
    // Wykaz po ukośniku dostaje każde okno z kanałem — bez nastawy i bez
    // przełącznika. Jest jedyną drogą doraźnego dostępu do narzędzia, więc nie
    // może zależeć od tego, którym oknem Operator akurat pracuje.
    zrodloWykazuUkosnika: zrodloWykazuKanalu(kanal),
    dolozenieNarzedzia: dolozenieKanalu(kanal),
    ...(opcje.kontekst === undefined ? {} : { kontekst: opcje.kontekst }),
    ...(opcje.pasekZlecenia === undefined
      ? {}
      : { pasekZlecenia: opcje.pasekZlecenia }),
  });

  // Dołożenie narzędzia cudzą ręką. Wybór z tego okna melduje się sam
  // z odpowiedzi komendy; poszerzenie zestawu spoza okna przychodzi wyłącznie
  // zdarzeniem `session.tool.attached`.
  //
  // Powód niepodpięcia nasłuchu nie idzie do wątku przy montażu. Gdy zdarzenia
  // nie ma w wygenerowanym kontrakcie, ten sam brak melduje się zdaniem w chwili
  // sięgnięcia po wykaz (`zrodlo-wykazu-ukosnika.ts`), więc powtarzanie go przy
  // otwarciu każdego okna byłoby hałasem.
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
