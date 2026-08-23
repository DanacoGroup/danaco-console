import { StudioOperationScope } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleWyboru, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import {
  uruchomOperacje,
  zlozZadanieOperacji,
  type ZapleczeNarzedzi,
} from './czynnosci-narzedzi';
import { utworzOknoStudio } from './okno-studio';
import { TrybOperacji } from './przybornik-uzycie';
import type { StanStudio } from './stan-studio';
import { utworzWykazOperacji } from './wykaz-operacji';
import type { ZrodloAkcjiStudio } from './zrodlo-akcji-studio';

/**
 * Tools Panel — stały panel operacji jako TRYB DO WYBORU, nie postać domyślna.
 *
 * ── Rozstrzygnięcie Właściciela ─────────────────────────────────────────────
 * Katalog operacji nie może zjadać stałej kolumny powierzchni. Drogą domyślną
 * są narzędzia ukryte: pływak przy zaznaczeniu i uchwyt katalogu
 * (`przybornik-plywak.ts`, `przybornik-katalog.ts`). Ten panel zostaje, bo dla
 * części pracy jest wygodniejszy — ale otwiera się wyborem Operatora, a nastawa
 * jest jawna, odwracalna i pamiętana w rdzeniu.
 *
 * Zwinięty panel nie zajmuje kolumny: niesie jeden wiersz z przełącznikiem
 * i znacznikiem `data-tryb-operacji`, po którym arkusz
 * `przybornik-znakowania.css` schodzi pas wiodący do jednej kolumny. Zwinięcie
 * do zera byłoby wygodniejsze wizualnie i błędne: Operator nie miałby czym
 * panelu wrócić.
 *
 * Poza tym: wywołanie operacji kontekstowej AI na zaznaczonym fragmencie albo
 * całym dokumencie. Sterowanie to przełącznik zakresu, wskaźnik zakresu
 * i przycisk uruchomienia.
 *
 * Zaznaczenie w edytorze przestawia przełącznik na „zaznaczenie", a wybór ręczny
 * ma pierwszeństwo do następnej zmiany zaznaczenia. Wybór „zaznaczenie" bez
 * zaznaczenia nie jest blokowany — wskaźnik mówi, czego brakuje.
 *
 * Nastawa zakresu mieszka w stanie modułu, a nie w polu `value` listy: lista
 * pokazuje `stan.zakresZadany()` i zgłasza wybór przez `stan.ustawZakresReczny()`,
 * dzięki czemu wskaźnik obok, pasek zaznaczenia w edytorze i żądanie
 * `studio.contextual.op` czytają jedną nastawę.
 *
 * Wszystkie operacje idą jedną komendą `studio.contextual.op`, rozróżnianą polem
 * `actionId` z rejestru akcji — stąd jedna ścieżka wywołania i wykaz
 * identyfikatorów zamiast osobnej ścieżki na operację.
 */
export interface OknoToolsPanel {
  element: HTMLElement;
  /** Wczytuje katalog akcji zasięgu modułu z rdzenia. */
  wczytaj(): Promise<void>;
  /** Uruchamia operację wskazaną z zewnątrz — z paska zaznaczenia edytora. */
  uruchom(idAkcji: string): Promise<void>;
  /**
   * Przenosi ognisko do panelu; panel zwinięty przy tym się rozwija.
   *
   * Bez rozwinięcia „Stały panel operacji" z pływaka przenosiłoby ognisko do
   * wiersza, w którym nie ma czego wybrać — czyli w nic.
   */
  przenieOgnisko(): void;
  /** Przestawia tryb wykazu operacji: narzędzia ukryte albo stały panel. */
  ustawTryb(tryb: TrybOperacji): void;
  odswiez(): void;
  /**
   * Zwija wykaz operacji i zdejmuje jego nasłuchy dokumentu.
   *
   * Menu wykazu zakłada nasłuch na `document` (zamknięcie kliknięciem obok,
   * obsługa klawiatury), więc moduł zdjęty z ekranu bez tego wywołania
   * zostawiłby po sobie żywy nasłuch.
   */
  zamknij(): void;
}

/**
 * Zdanie wskaźnika: jaki zakres pojedzie do rdzenia, a gdy nie jest to zakres
 * wybrany na liście — także z jakiego powodu.
 *
 * Rozróżnione są trzy stany: zakres wzięty z zaznaczenia, zaznaczenie pominięte
 * wyborem ręcznym oraz wybrane „zaznaczenie", którego w edytorze nie ma.
 */
function opiszZakres(stan: StanStudio): string {
  const wybrany = stan.zaznaczenie();
  if (stan.zakresSkuteczny() === StudioOperationScope.Selection && wybrany !== null) {
    return `Zakres bieżący: zaznaczenie od znaku ${wybrany.poczatek} do ${wybrany.koniec}.`;
  }
  if (wybrany !== null) {
    return (
      `Zakres bieżący: cały dokument — wybór ręczny ma pierwszeństwo, więc zaznaczenie ` +
      `(${wybrany.koniec - wybrany.poczatek} znaków) zostaje pominięte do następnej zmiany zaznaczenia.`
    );
  }
  if (stan.zakresZadany() === StudioOperationScope.Selection) {
    return (
      'Zakres bieżący: cały dokument — wybrano „zaznaczenie", ale w edytorze nic nie jest ' +
      'zaznaczone. Wybór nie jest blokowany; operacja pójdzie na całość.'
    );
  }
  return 'Zakres bieżący: cały dokument — w edytorze nic nie jest zaznaczone.';
}

/** Czym panel rozporządza poza rdzeniem: nastawą trybu wykazu operacji. */
export interface CzynnosciToolsPanelu {
  /** Zapisuje wybrany tryb wykazu operacji — nastawa pamiętana w rdzeniu. */
  naTryb(tryb: TrybOperacji): void;
}

export function utworzOknoToolsPanel(
  stan: StanStudio,
  akcje: ZrodloAkcjiStudio,
  czynnosciZewnetrzne: CzynnosciToolsPanelu,
): OknoToolsPanel {
  const rama = utworzOknoStudio({
    kod: 'studio.tools-panel',
    tytul: 'Tools Panel',
    rola: 'pomocnicze',
    objasnienie:
      'Zestaw operacji kontekstowych AI. Wykaz pochodzi z rejestru akcji rdzenia (action.list); ' +
      'operacja idzie do rdzenia komendą studio.contextual.op z identyfikatorem akcji i zakresem.',
  });

  const zakres = poleWyboru(
    {
      etykieta: 'Zakres operacji',
      opis: 'Zaznaczenie w Studio Editorze przestawia zakres samo; wybór ręczny ma pierwszeństwo do następnego zaznaczenia.',
    },
    [
      { wartosc: StudioOperationScope.Selection, etykieta: 'Zaznaczenie' },
      { wartosc: StudioOperationScope.Document, etykieta: 'Cały dokument' },
    ],
  );

  const wskaznik = document.createElement('p');
  wskaznik.className = 'dn-pole-opis ms-wskaznik';

  const wykaz = utworzWykazOperacji(() => odswiez());
  const uruchomienie = przycisk('Uruchom operację', 'dn-btn dn-btn--sygnal');
  uruchomienie.dataset['czynnosc'] = 'uruchom';
  const odpowiedz = utworzWierszOdpowiedzi();

  /* ── Tryb wykazu: narzędzia ukryte albo ten panel ─────────────────────────── */

  let tryb: TrybOperacji = TrybOperacji.Ukryte;

  const przelacznikTrybu = przycisk('Otwórz stały panel operacji', 'dn-btn dn-btn--sm dn-btn--zarys');
  przelacznikTrybu.dataset['czynnosc'] = 'tryb-panelu';

  const zdanieTrybu = document.createElement('p');
  zdanieTrybu.className = 'dn-pole-opis';

  const zwijane = document.createElement('div');
  zwijane.className = 'ms-panel-operacji__zwijane';
  zwijane.append(zakres.element, wskaznik, wykaz.element, odpowiedz.element);

  rama.pasek.append(uruchomienie);
  rama.stan.tresc.append(przelacznikTrybu, zdanieTrybu, zwijane);

  przelacznikTrybu.addEventListener('click', () => {
    const nowy = tryb === TrybOperacji.Panel ? TrybOperacji.Ukryte : TrybOperacji.Panel;
    ustawTryb(nowy);
    czynnosciZewnetrzne.naTryb(nowy);
  });

  function ustawTryb(nowy: TrybOperacji): void {
    tryb = nowy;
    const otwarty = nowy === TrybOperacji.Panel;
    // Znacznik siedzi na elemencie okna, bo to jego kolumnę rezerwuje pas
    // wiodący modułu — arkusz odcinka schodzi pas do jednej kolumny po nim.
    rama.element.dataset['trybOperacji'] = nowy;
    zwijane.hidden = !otwarty;
    uruchomienie.hidden = !otwarty;
    przelacznikTrybu.textContent = otwarty
      ? 'Zwiń panel — wróć do narzędzi ukrytych'
      : 'Otwórz stały panel operacji';
    zdanieTrybu.textContent = otwarty
      ? 'Stały panel operacji stoi otwarty i zajmuje kolumnę powierzchni. Nastawa jest pamiętana ' +
        'w rdzeniu (config.set) i cofa się tym samym przełącznikiem.'
      : 'Katalog operacji stoi jako narzędzia ukryte: pływak przy zaznaczeniu, uchwyt pełnego ' +
        'katalogu i wiersz polecenia. Ten panel jest trybem do wyboru — otwórz go, jeśli wolisz ' +
        'wykaz pod ręką na stałe.';
  }

  ustawTryb(TrybOperacji.Ukryte);

  const zaplecze: ZapleczeNarzedzi = { stan, pas: rama.stan, odpowiedz };

  async function uruchom(idAkcji: string): Promise<void> {
    await uruchomOperacje(zaplecze, zlozZadanieOperacji(stan, idAkcji), idAkcji);
  }

  uruchomienie.addEventListener('click', () => void uruchom(wykaz.wybrana()));
  // Wybór ręczny idzie do stanu, a nie zostaje w kontrolce. Stan ogłasza zmianę,
  // odświeżenie wraca tu z powrotem i przepisuje wartość listy — pętli nie ma,
  // bo ustawienie `value` z kodu nie wywołuje zdarzenia `change`.
  zakres.kontrolka.addEventListener('change', () => {
    stan.ustawZakresReczny(zakres.kontrolka.value as StudioOperationScope);
  });

  async function wczytaj(): Promise<void> {
    rama.stan.ladowanie('Odczyt rejestru akcji modułu…');
    const wynik = await akcje.katalog();
    if (!wynik.udany || wynik.wynik === undefined) {
      rama.stan.blad(opisOdmowy('Odczyt rejestru akcji', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    wykaz.ustawRejestr(wynik.wynik.actions);
    rama.stan.gotowe();
  }

  function odswiez(): void {
    zakres.kontrolka.value = stan.zakresZadany();
    wskaznik.textContent = opiszZakres(stan);
    if (rama.stan.faza() === 'ladowanie' || rama.stan.faza() === 'blad') return;
    if (stan.dokument() === null) {
      // Nazwa stanu mówi, czego brakuje operacji: wykaz jest czynny, brakuje
      // wyłącznie przedmiotu, na którym operacja ma pracować.
      rama.stan.puste(
        'Operacje bez przedmiotu',
        'Tools Panel zbiera operacje kontekstowe AI z rejestru akcji rdzenia i puszcza ' +
          'wybraną komendą studio.contextual.op. Wykaz poniżej jest czynny; przedmiot — ' +
          'cały dokument albo zaznaczony fragment — daje dokument wczytany w Studio Editorze.',
      );
      return;
    }
    rama.stan.gotowe();
  }

  return {
    element: rama.element,
    wczytaj,
    uruchom,
    przenieOgnisko() {
      if (tryb !== TrybOperacji.Panel) {
        ustawTryb(TrybOperacji.Panel);
        czynnosciZewnetrzne.naTryb(TrybOperacji.Panel);
      }
      rama.element.dataset['ognisko'] = 'tak';
      rama.element.scrollIntoView({ block: 'nearest' });
    },
    ustawTryb,
    odswiez,
    zamknij: () => wykaz.zamknij(),
  };
}
