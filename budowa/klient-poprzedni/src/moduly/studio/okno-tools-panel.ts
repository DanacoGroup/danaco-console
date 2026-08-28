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

/** Tools Panel, stały panel operacji jako tryb do wyboru operatora, nie postać domyślna: katalog operacji nie zajmuje stałej kolumny na stałe. */
export interface OknoToolsPanel {
  element: HTMLElement;
  /** Wczytuje katalog akcji zasięgu modułu z rdzenia. */
  wczytaj(): Promise<void>;
  /** Uruchamia operację wskazaną z zewnątrz — z paska zaznaczenia edytora. */
  uruchom(idAkcji: string): Promise<void>;
  /** Przenosi ognisko do panelu operacji, rozwijając go, gdy panel jest w danej chwili zwinięty. */
  przenieOgnisko(): void;
  /** Przestawia tryb wykazu operacji: narzędzia ukryte albo stały panel. */
  ustawTryb(tryb: TrybOperacji): void;
  odswiez(): void;
  /** Zwija wykaz operacji panelu i zdejmuje jego nasłuchy dokumentu, żeby nie zostawić żywego nasłuchu. */
  zamknij(): void;
}

/** Zdanie wskaźnika mówiące, jaki zakres pojedzie do rdzenia, rozróżniające zakres z zaznaczenia, wybór ręczny i brak zaznaczenia. */
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

/** Czym panel rozporządza poza rdzeniem: nastawą trybu wykazu operacji, pamiętaną między sesjami tego okna. */
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
    // Znacznik siedzi na elemencie okna, bo to jego kolumnę rezerwuje pas wiodący modułu.
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
  // Wybór ręczny idzie do stanu, bo ustawienie value z kodu nie wywołuje zdarzenia change.
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
      // Nazwa stanu mówi, czego brakuje: wykaz jest czynny, brakuje wyłącznie przedmiotu operacji.
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
