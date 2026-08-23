import { Command, ConfigScope, type WorkspaceInstructions } from '../../../../shared/contract';
import type { WykazBrakow } from './braki-kontraktu';
import { KLUCZ_INSTRUKCJI } from './wynik-czastkowy';
import { POZIOMY_ZASIEGU } from './poziomy-zasiegu';
import { opisWarstwy, SZABLONY } from './warstwy-instrukcji';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  poleTresci,
  przyciskAkcji as przycisk,
  wiersz,
  wybor,
} from '../../modele/kontrolki-formularza';
import type { StanProjektu } from './stan-projektu';
import { utworzStanTresci } from './stany-okna';
import type { ZrodloWorkspace } from './zrodlo-workspace';

/**
 * Instructions Panel — okno pomocnicze modułu Workspace: edycja instrukcji
 * systemowych i przypisanie ich do projektu.
 *
 * Instrukcje obowiązują domyślnie w zakresie projektu, a współdzieleniem jest
 * poziom zasięgu: zapis na poziomie szerszym niż projekt obowiązuje w każdym
 * projekcie tego zasięgu.
 *
 * Dziedziczenie warstwowe panel pokazuje, ale go nie przelicza. Warstwę
 * obowiązującą rozstrzyga rdzeń i oddaje ją w wyniku zapisu wraz ze wskazaniem
 * poziomu, z którego pochodzi; drugi rozstrzygacz po stronie okna rozjechałby
 * się z rdzeniem przy pierwszej zmianie reguł.
 *
 * „Podgląd warstwy” pyta `config.get` o zapis jednego poziomu. To odczyt, nie
 * rozstrzyganie: odpowiada na pytanie, czy ta warstwa ma własną treść.
 */
export interface OknoInstrukcji {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoInstrukcji(
  zrodlo: ZrodloWorkspace,
  stan: StanProjektu,
  braki: WykazBrakow,
): OknoInstrukcji {
  const rama = utworzRameOkna({
    tytul: 'Instructions Panel',
    kod: 'instructions-panel',
    rola: 'pomocnicze',
    przeznaczenie: 'Instrukcje systemowe projektu wraz z warstwą, która po zapisie naprawdę obowiązuje.',
    przedrostek: 'dw',
  });
  const tresc = utworzStanTresci();

  // Wykaz poziomów bierze się z jednej prawdy modułu (`poziomy-zasiegu.ts`),
  // więc ta sama nastawa nazywa się tu tak samo jak w obu panelach pamięci.
  const poziom = wybor(
    'Poziom zasięgu instrukcji',
    POZIOMY_ZASIEGU.map(([wartosc, opis]) => [wartosc, opis]),
  );
  poziom.value = ConfigScope.Project;
  const bytPoziomu = pole('Byt poziomu zasięgu', 'puste = projekt bieżący');
  const edytor = poleTresci('Treść instrukcji systemowych', 12);
  const szablony = wybor('Szablony instrukcji', SZABLONY.map(([nazwa]) => [nazwa, nazwa]));

  const zapisz = przycisk('Zapisz instrukcje', 'dn-btn dn-btn--atrament');
  const podglad = przycisk('Podgląd warstwy');
  const eksport = przycisk('Eksportuj');
  const importuj = przycisk('Importuj');
  const plik = document.createElement('input');
  plik.type = 'file';
  plik.accept = '.md,.txt';
  plik.hidden = true;

  rama.akcje.append(
    zapisz,
    podglad,
    eksport,
    importuj,
    braki.przyciskBraku(
      'Wersje ▾ (Diff, przywrócenie)',
      'Wykaz wersji instrukcji projektu wraz z przywróceniem',
      Command.WorkspaceInstructionsVersionList,
      Command.WorkspaceInstructionsVersionRestore,
    ),
    // Bez nazwy komendy, bo żadna jej nie nosi. Nazwa wpisana „na zapas”
    // przedstawiałaby Operatorowi wymyślony identyfikator jako kandydata.
    braki.przyciskBraku('Test w Chat', 'Próbne przekazanie instrukcji do okna rozmowy'),
    plik,
  );
  rama.cialo.append(
    wiersz('Warstwa zapisu', poziom, { klasa: 'dw-wiersz', objasnienie: 'Domyślnie projekt; poziom najwęższy wygrywa, a brak zapisu na poziomie = obowiązuje wartość poziomu szerszego.' }),
    wiersz('Byt poziomu', bytPoziomu, { klasa: 'dw-wiersz', objasnienie: 'Identyfikator środowiska, modułu, karty sesji albo okna; poziom globalny bytu nie ma.' }),
    wiersz('Szablon', szablony, { klasa: 'dw-wiersz', objasnienie: 'Szablon wstawia treść do edytora; zapis pozostaje osobną decyzją.' }),
    wiersz('Instrukcje', edytor, { klasa: 'dw-wiersz' }),
    tresc.element,
  );

  /**
   * Licznik długości instrukcji. Opracowanie wymienia obok liczby znaków także
   * przybliżoną liczbę tokenów — tej okno nie pokazuje, bo nie ma tokenizatora
   * modelu, a liczba oszacowana regułą własną byłaby miarą wymyśloną, nie zmierzoną.
   */
  function pokazDlugosc(): void {
    rama.ustawZnacznik(`znaki: ${edytor.value.length}`);
  }
  edytor.addEventListener('input', pokazDlugosc);
  pokazDlugosc();

  function bytZadania(): string {
    const wskazany = bytPoziomu.value.trim();
    if (wskazany !== '') return wskazany;
    return poziom.value === ConfigScope.Project ? stan.projekt() : '';
  }

  function pokazWarstwe(instrukcje: WorkspaceInstructions): void {
    tresc.tresc().append(opisWarstwy(instrukcje, poziom.value as ConfigScope));
  }

  zapisz.addEventListener('click', () => {
    const projekt = stan.projekt();
    if (projekt === '') {
      tresc.blad('Zapis bez projektu nie ma zakresu — wskaż projekt na pulpicie.');
      return;
    }
    tresc.ladowanie('Zapis instrukcji…');
    void zrodlo
      .zapiszInstrukcje({
        projectId: projekt,
        content: edytor.value,
        scope: poziom.value as ConfigScope,
        scopeId: bytZadania(),
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Rdzeń nie przyjął instrukcji.', wynik.blad);
          return;
        }
        const zapisane = wynik.wynik.instructions;
        pokazWarstwe(zapisane);
        // Potwierdzenie mówi to, co oddał rdzeń: samo przyjęcie wywołania nie
        // dowodzi, że zapisana treść jest tą, którą wysłał edytor.
        const rozjazd =
          zapisane.scope === (poziom.value as ConfigScope) && zapisane.content !== edytor.value;
        tresc.potwierdzenie(
          rozjazd
            ? 'Rdzeń przyjął zapis, ale oddał na tej warstwie treść INNĄ niż wysłana — ' +
                'poniżej stoi treść z jego odpowiedzi, nie z edytora.'
            : `Rdzeń zapisał instrukcje; warstwą obowiązującą jest wedle jego odpowiedzi ` +
                `${zapisane.scope}${zapisane.scopeId === undefined ? '' : ` (${zapisane.scopeId})`}.`,
          !rozjazd,
        );
      });
  });

  podglad.addEventListener('click', () => {
    tresc.ladowanie('Odczyt zapisu tej warstwy…');
    void zrodlo.instrukcjeWarstwy(poziom.value as ConfigScope, bytZadania()).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń nie oddał zapisu tej warstwy.', wynik.blad);
        return;
      }
      // Klucz brany ze stałej, którą wysłano w zapytaniu — napis przepisany
      // ręcznie rozjechałby się z nią przy pierwszej zmianie i okno meldowałoby
      // brak zapisu na warstwie, która zapis ma.
      const wpis = wynik.wynik.find((pozycja) => pozycja.key === KLUCZ_INSTRUKCJI);
      if (wpis === undefined) {
        // Rdzeń oddał wykaz bez tego klucza — i tylko to wynika z odpowiedzi.
        // Reguły rozstrzygania warstw ta odpowiedź nie niesie, więc okno o niej
        // nie orzeka.
        tresc.pusto(
          `Rdzeń nie oddał na tej warstwie zapisu klucza ${KLUCZ_INSTRUKCJI} ` +
            `(config.get, poziom ${poziom.value}). Warstwę obowiązującą rozstrzyga rdzeń ` +
            'i pokazuje ją dopiero odpowiedź na zapis.',
        );
        return;
      }
      edytor.value = typeof wpis.value === 'string' ? wpis.value : JSON.stringify(wpis.value);
      // Zapis z kodu nie wywołuje zdarzenia `input`, więc licznik trzeba
      // przeliczyć wprost — inaczej pokazywałby długość treści poprzedniej.
      pokazDlugosc();
      tresc.potwierdzenie('Treść tej warstwy wczytana do edytora.', true);
    });
  });

  szablony.addEventListener('change', () => {
    const szablon = SZABLONY.find(([nazwa]) => nazwa === szablony.value);
    if (szablon === undefined || szablon[1] === '') return;
    edytor.value = edytor.value === '' ? szablon[1] : `${edytor.value}\n${szablon[1]}`;
    pokazDlugosc();
  });

  eksport.addEventListener('click', () => {
    // Pusty edytor daje pusty plik, więc eksport odmawia zamiast potwierdzać
    // czynność, która nie ma czego wykonać.
    if (edytor.value === '') {
      tresc.potwierdzenie('Edytor jest pusty — nie ma czego wyeksportować.', false);
      return;
    }
    pobierzPlik(`${stan.projekt() || 'projekt'}-instrukcje.md`, edytor.value, 'text/markdown');
    tresc.potwierdzenie(
      `Treść edytora (${edytor.value.length} znaków) pobrana jako plik Markdown.`,
      true,
    );
  });

  importuj.addEventListener('click', () => plik.click());
  plik.addEventListener('change', () => {
    const wybrany = plik.files?.[0];
    if (wybrany === undefined) return;
    void wybrany.text().then((zawartosc) => {
      edytor.value = zawartosc;
      pokazDlugosc();
      // Wskazanie kontrolki plikowej trzeba wyczyścić, inaczej ponowny wybór
      // tego samego pliku nie wywoła zdarzenia `change` i import przejdzie bez śladu.
      plik.value = '';
      tresc.potwierdzenie(`Wczytano ${wybrany.name}; zapis pozostaje osobną decyzją.`, true);
    });
  });

  tresc.pusto('Zapisz instrukcje albo obejrzyj warstwę, aby zobaczyć, co obowiązuje w tym projekcie.');
  stan.naZmiane(() => {
    if (poziom.value === ConfigScope.Project) bytPoziomu.placeholder = `puste = ${stan.projekt()}`;
  });

  return {
    element: rama.element,
    odswiez() {
      tresc.pusto('Zapisz instrukcje albo obejrzyj warstwę, aby zobaczyć, co obowiązuje w tym projekcie.');
    },
  };
}
