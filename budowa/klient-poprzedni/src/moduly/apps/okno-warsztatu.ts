import type { AppWorkspaceLayer, DeveloperFile } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWielowierszowe,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { utworzWykazBrakow, type BrakFunkcji } from './braki-kontraktu';
import { opiszPole } from './dymek-objasnienia';
import { BEZ_OKNA_MODULU } from './etykiety-apps';
import { utworzRameApps } from './rama-okna';
import { narzedziaWarsztatu } from './narzedzia-apps';
import { utworzPrzybornikApps } from './przybornik-apps';
import type { StanProduktu } from './stan-produktu';
import { utworzWykazKomponentowWarstwy, komponentyWarstwy } from './wykaz-komponentow-warstwy';
import { utworzWykazPlikowWarsztatu } from './wykaz-plikow-warsztatu';
import { utworzWyborZMenu, wierszWyboru } from './wybor-z-menu';

/**
 * Rama wspólna Frontend Workspace i Backend Workspace — obu okien wiodących
 * modułu Apps.
 *
 * Jedna komenda obsługuje oba okna: `apps.workspace.update` niesie pole `layer`
 * o dwóch wartościach, a różnica między oknami sprowadza się do wartości tego
 * pola, do zbioru komponentów warstwy i do wykazu czynności bez drogi
 * w kontrakcie. Różnice przychodzą parametrem, nie drugim plikiem widoku.
 *
 * Podgląd wyniku pokazuje to, co potwierdził rdzeń: odpowiedź komendy niesie
 * `DeveloperFile` — ścieżkę, treść po zapisie, rozmiar i wersję. Podgląd na żywo
 * (hot reload) nie ma w kontrakcie ani komendy, ani zdarzenia, więc stoi
 * w wykazie braków zamiast w udawanym oknie podglądu.
 *
 * Ścieżka jest tożsamością pliku, więc okno nie przycina jej po swojemu i nie
 * odmawia w imieniu kontraktu. Klucz `(okno, warstwa, ścieżka)` rozstrzyga,
 * który plik zostanie nadpisany (`migracja_051_aplikacje.sql`), a `trim()`
 * przeglądarki nie jest tym samym przycięciem co `strings.TrimSpace` rdzenia:
 * ścieżki złożonej z samych znaków niewidocznych rdzeń za pustą nie uznaje.
 * Okno zatrzymuje więc wyłącznie pole dosłownie puste i mówi wtedy o sobie,
 * a nie o kontrakcie; o każdej innej wartości rozstrzyga rdzeń.
 *
 * Zdanie o zapisie zestawia wpisane z oddanym. Sama ścieżka z odpowiedzi to za
 * mało: gdy rdzeń zapisze plik pod ścieżką inną niż wpisana, potwierdzenie
 * wymieniające ścieżkę oddaną jest prawdziwe, a mimo to zostawia czytającego
 * w przekonaniu, że zapisał to, co wpisał. Różnicę okno nazywa.
 */
export interface OknoWarsztatu {
  element: HTMLElement;
  odswiez(): void;
}

/** Czym różnią się oba warsztaty — wszystko, czego rama nie wywnioskuje sama. */
export interface OpisWarsztatu {
  kodOkna: string;
  tytul: string;
  warstwa: AppWorkspaceLayer;
  /** Nagłówek wykazu komponentów: drzewo komponentów albo mapa zależności usług. */
  tytulWykazu: string;
  /** Objaśnienie pola ścieżki właściwe warstwie. */
  objasnienieSciezki: string;
  braki: readonly BrakFunkcji[];
}

export function utworzOknoWarsztatu(stan: StanProduktu, opis: OpisWarsztatu): OknoWarsztatu {
  const rama = utworzRameApps(opis.kodOkna, opis.tytul, 'wiodące');
  const wykaz = utworzWykazKomponentowWarstwy(opis.warstwa, opis.tytulWykazu);

  const sciezka = poleTekstowe({ etykieta: 'Ścieżka pliku', podpowiedz: 'src/…' });
  // Rozwijanie z biblioteki kontrolek, nie natywny `<select>`. Wykaz wchodzi
  // dopiero w `odswiez`, bo pochodzi z kanwy, nie z wyliczenia kontraktu.
  const komponent = utworzWyborZMenu('Komponent architektury');
  const komponentWiersz = wierszWyboru('Komponent architektury', komponent);
  const tresc = poleWielowierszowe({ etykieta: 'Treść pliku po edycji' }, 8);
  const zapisz = przycisk('Zapisz zmianę warsztatu', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odczytaj = przycisk('Odczytaj pliki warsztatu', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  // Wybór pliku z wykazu wstawia go do formularza — i mówi o tym wprost, bo
  // przepisuje pole treści, w którym mogła już stać niezapisana zmiana.
  const pliki = utworzWykazPlikowWarsztatu('Pliki warsztatu w rdzeniu', (plik) => {
    sciezka.kontrolka.value = plik.path;
    tresc.kontrolka.value = plik.content ?? '';
    podglad.textContent = opiszPlik(plik);
    odpowiedz.pokaz(
      `Do formularza wstawiono plik ${czytelnaSciezka(plik.path)} z wykazu warsztatu. ` +
        (plik.content === undefined
          ? 'Rdzeń nie oddał jego treści, więc pole treści jest puste — zapis pod tą ścieżką ' +
            'NADPISAŁBY plik pustką.'
          : 'Zapis pod tą ścieżką nadpisze plik istniejący.'),
      plik.content !== undefined,
    );
  });

  const podglad = document.createElement('pre');
  podglad.className = 'mp-podglad';

  rama.akcje.append(utworzWykazBrakow('Bez drogi w kontrakcie', opis.braki));
  // Narzędzia zależą od warstwy: podgląd i mapa routingu należą do frontendu,
  // eksplorator punktów końcowych i schemat bazy — do backendu.
  rama.akcje.append(
    utworzPrzybornikApps('Narzędzia warstwy', narzedziaWarsztatu(stan, opis.warstwa)).element,
  );
  rama.tresc.append(
    wykaz.element,
    odczytaj,
    pliki.element,
    opiszPole(sciezka.element, opis.objasnienieSciezki),
    opiszPole(
      komponentWiersz,
      'Komponent, którego dotyczy zmiana. Wykaz pochodzi z kanwy Architecture ' +
        'Designera; pusty wybór wysyła żądanie bez pola componentId.',
    ),
    tresc.element,
    zapisz,
    odpowiedz.element,
    podglad,
  );

  zapisz.addEventListener('click', () => void wyslij());
  odczytaj.addEventListener('click', () => void odczytajPliki());

  /**
   * Odczyt plików warsztatu z rdzenia.
   *
   * Zawężamy do własnej warstwy, bo okno jest oknem jednej warstwy. Odczyt bez
   * pola `layer` przyniósłby oba warsztaty naraz i Frontend Workspace pokazałby
   * pliki backendu jako swoje, a klucz `(okno, warstwa, ścieżka)` czyni z nich
   * byty osobne.
   */
  async function odczytajPliki(): Promise<void> {
    rama.ladowanie('Odczyt plików warsztatu w toku…');
    odpowiedz.pokaz('Odczyt plików warsztatu: żądanie wysłane do rdzenia…', true);
    await stan.odczytajWarsztat(opis.warstwa);
    const powod = stan.powodOdczytuWarsztatu(opis.warstwa);
    if (powod !== '') {
      rama.blad(powod);
      odpowiedz.pokaz(powod, false);
      return;
    }
    const ile = stan.plikiWarsztatu(opis.warstwa).length;
    odpowiedz.pokaz(
      ile === 0
        ? 'Odczyt plików warsztatu: rdzeń nie zna ani jednego pliku tej warstwy. To jest ' +
          'pustka POTWIERDZONA, nie brak odczytu.'
        : `Odczyt plików warsztatu: rdzeń oddał ${ile} plików tej warstwy.`,
      true,
    );
    rama.gotowe();
    odswiez();
  }

  async function wyslij(): Promise<void> {
    const idOkna = stan.idOkna();
    if (idOkna === '') {
      odpowiedz.pokaz(BEZ_OKNA_MODULU, false);
      return;
    }
    const zamowionaSciezka = sciezka.kontrolka.value;
    if (zamowionaSciezka === '') {
      odpowiedz.pokaz(
        'Okno nie wysyła pustego pola ścieżki — wpisz ją. O tym, czy wpisana ' +
          'wartość jest ścieżką, rozstrzyga rdzeń, nie to okno.',
        false,
      );
      return;
    }
    rama.ladowanie('Zapis warsztatu w toku…');
    odpowiedz.pokaz('Zapis warsztatu: żądanie wysłane do rdzenia…', true);
    const wynik = await stan.zrodlo.zapiszPlik({
      idOkna,
      warstwa: opis.warstwa,
      sciezka: zamowionaSciezka,
      tresc: tresc.kontrolka.value,
      idKomponentu: komponent.wartosc(),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      const zdanie = opisOdmowy('Zapis warsztatu', wynik.blad?.code, wynik.blad?.message);
      rama.blad(zdanie);
      odpowiedz.pokaz(zdanie, false);
      return;
    }
    podglad.textContent = opiszPlik(wynik.wynik.file);
    // Plik potwierdzony przez rdzeń wchodzi do zbioru warstwy natychmiast —
    // bez czekania na zdarzenie `apps.workspace.changed`, którego rdzeń nie
    // musi odesłać nadawcy zmiany. Wykaz pokazuje wtedy stan po zapisie zamiast
    // stanu sprzed niego.
    stan.wchlonZapisWarsztatu(wynik.wynik.layer, wynik.wynik.file);
    const rozbieznosc = rozbieznoscZapisu(
      zamowionaSciezka,
      opis.warstwa,
      wynik.wynik.layer,
      wynik.wynik.file,
    );
    if (rozbieznosc !== '') {
      const zdanie = `Zapis warsztatu: rdzeń zapisał CO INNEGO, NIŻ ZAMÓWIONO — ${rozbieznosc}.`;
      rama.blad(zdanie);
      odpowiedz.pokaz(zdanie, false);
      return;
    }
    odpowiedz.pokaz(`Rdzeń potwierdził zapis pliku ${czytelnaSciezka(wynik.wynik.file.path)}.`, true);
    rama.gotowe();
  }

  function odswiez(): void {
    const komponenty = stan.komponenty();
    const ile = wykaz.nanies(komponenty);
    const ilePlikow = pliki.nanies(stan.plikiWarsztatu(opis.warstwa));
    komponent.ustawPozycje([
      { wartosc: '', etykieta: 'bez wskazania komponentu' },
      ...komponentyWarstwy(opis.warstwa, komponenty).map((pozycja) => ({
        wartosc: pozycja.id,
        etykieta: pozycja.name,
      })),
    ]);
    if (rama.faza() === 'blad' || rama.faza() === 'ladowanie') return;
    if (ile === 0 && ilePlikow === 0 && podglad.textContent === '') {
      // Dwa różne zdania o pustce, bo to dwa różne stany: po udanym odczycie
      // wiadomo, że rdzeń warsztatu nie ma, a przed nim wyłącznie tyle, że nikt
      // o warsztat nie pytał.
      rama.puste(
        stan.czyWarsztatCzytany(opis.warstwa)
          ? 'Rdzeń nie zna ani jednego pliku tej warstwy ani komponentu do niej należącego. ' +
            'Zacznij w Architecture Designerze albo zapisz pierwszy plik tutaj.'
          : 'Bez zestawionych jeszcze zasobów warstwy — zacznij w Architecture Designerze ' +
            'albo naciśnij „Odczytaj pliki warsztatu", żeby zobaczyć, co rdzeń już ma.',
      );
      return;
    }
    rama.gotowe();
  }

  odswiez();
  return { element: rama.element, odswiez };
}

/**
 * Czym zapis oddany przez rdzeń różni się od zamówionego — pusty łańcuch, gdy
 * niczym.
 *
 * Zestawiane są ścieżka i warstwa, bo to one rozstrzygają, który plik został
 * nadpisany: klucz naturalny warsztatu to `(okno, warstwa, ścieżka)`
 * (`migracja_051_aplikacje.sql`), więc rozejście się choćby jednego z tych pól
 * znaczy nadpisanie innego pliku niż zamierzony. Treści nie zestawiamy —
 * kontrakt nie obiecuje, że `content` w odpowiedzi jest dosłownym echem
 * żądania. Podgląd pokazuje treść oddaną w całości.
 */
function rozbieznoscZapisu(
  zamowionaSciezka: string,
  zamowionaWarstwa: AppWorkspaceLayer,
  oddanaWarstwa: AppWorkspaceLayer,
  plik: DeveloperFile,
): string {
  const rozejscia: string[] = [];
  if (plik.path !== zamowionaSciezka) {
    rozejscia.push(
      `wpisano ścieżkę „${czytelnaSciezka(zamowionaSciezka)}", rdzeń zapisał plik pod ` +
        `„${czytelnaSciezka(plik.path)}"`,
    );
  }
  if (oddanaWarstwa !== zamowionaWarstwa) {
    rozejscia.push(`zamówiono warstwę ${zamowionaWarstwa}, rdzeń oddał ${oddanaWarstwa}`);
  }
  return rozejscia.join('; ');
}

/**
 * Ścieżka wypisana tak, żeby dwie różne ścieżki nie wyglądały tak samo.
 *
 * Rdzeń przyjmuje ścieżkę poprzedzoną znakiem niewidocznym i zapisuje ją
 * dosłownie, a na ekranie wygląda ona identycznie jak ścieżka bez tego znaku —
 * choć klucz `(okno, warstwa, ścieżka)` czyni z nich dwa różne pliki.
 * Potwierdzenie zapisu byłoby wtedy prawdziwe co do znaku i mylące co do rzeczy.
 * Znaki niewidoczne — sterujące, formatujące i spacje inne niż zwykła —
 * wypisujemy więc kodem, a wszystko pozostałe zostaje bez zmian.
 */
function czytelnaSciezka(sciezka: string): string {
  // Zakresy: sterujące C0/C1, miękki dywiz, spacje inne niż U+0020, znaczniki
  // kierunku i złączenia oraz znacznik kolejności bajtów.
  const niewidoczne =
    /[\u0000-\u001F\u007F-\u00A0\u00AD\u1680\u2000-\u200F\u2028\u2029\u202F\u205F\u2060\u3000\uFEFF]/g;
  return sciezka.replace(
    niewidoczne,
    (znak) => `\\u${znak.codePointAt(0)?.toString(16).padStart(4, '0') ?? '????'}`,
  );
}

/** Podgląd wyniku: plik tak, jak oddał go rdzeń po zapisie. */
function opiszPlik(plik: DeveloperFile): string {
  const wiersze = [
    `ścieżka: ${plik.path}`,
    `język: ${plik.language ?? 'nieokreślony przez rdzeń'}`,
    `rozmiar: ${plik.sizeBytes ?? 0} B`,
    `wersja: ${plik.versionId ?? 'bez wersji w odpowiedzi'}`,
    '',
    plik.content ?? '(rdzeń nie oddał treści pliku)',
  ];
  return wiersze.join('\n');
}
