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
 * Rama wspólna Frontend Workspace i Backend Workspace, obu okien wiodących
 * modułu Apps: jedna komenda obsługuje oba okna, a różnica sprowadza się do
 * parametru warstwy, zbioru komponentów i wykazu czynności bez drogi
 * w kontrakcie.
 */
export interface OknoWarsztatu {
  element: HTMLElement;
  odswiez(): void;
}

/** Czym różnią się oba warsztaty — wszystko, czego ta rama wspólna nie jest w stanie wywnioskować sama. */
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
  // Rozwijanie z biblioteki kontrolek: wykaz wchodzi dopiero przy odświeżeniu, z kanwy, nie z kontraktu.
  const komponent = utworzWyborZMenu('Komponent architektury');
  const komponentWiersz = wierszWyboru('Komponent architektury', komponent);
  const tresc = poleWielowierszowe({ etykieta: 'Treść pliku po edycji' }, 8);
  const zapisz = przycisk('Zapisz zmianę warsztatu', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odczytaj = przycisk('Odczytaj pliki warsztatu', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  // Wybór pliku z wykazu wstawia go do formularza i przepisuje pole treści, nawet niezapisane.
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
  // Narzędzia zależą od warstwy: frontend i backend mają każdy własny, dobrany do siebie zestaw.
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

  /** Odczyt plików warsztatu z rdzenia, zawężony do własnej warstwy, bo okno jest oknem jednej warstwy. */
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
    // Plik potwierdzony przez rdzeń wchodzi do zbioru warstwy natychmiast, bez czekania na zdarzenie.
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
      // Dwa różne zdania o pustce, bo to dwa różne stany: po odczycie i przed nim.
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
 * Czym zapis oddany przez rdzeń różni się od zamówionego — pusty łańcuch,
 * gdy niczym. Zestawiane są ścieżka i warstwa, bo one rozstrzygają, który
 * plik został nadpisany; treści nie zestawiamy, bo kontrakt nie obiecuje
 * dosłownego echa żądania.
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
 * Ścieżka wypisana tak, żeby dwie różne ścieżki nie wyglądały tak samo:
 * znaki niewidoczne — sterujące, formatujące i spacje inne niż zwykła —
 * wypisujemy kodem, a wszystko pozostałe zostaje bez zmian.
 */
function czytelnaSciezka(sciezka: string): string {
  // Zakresy: znaki sterujące, miękki dywiz, spacje niestandardowe oraz znaczniki kierunku i kolejności.
  const niewidoczne =
    /[\u0000-\u001F\u007F-\u00A0\u00AD\u1680\u2000-\u200F\u2028\u2029\u202F\u205F\u2060\u3000\uFEFF]/g;
  return sciezka.replace(
    niewidoczne,
    (znak) => `\\u${znak.codePointAt(0)?.toString(16).padStart(4, '0') ?? '????'}`,
  );
}

/** Podgląd wyniku pokazuje plik dokładnie tak, jak oddał go rdzeń po zapisie, bez żadnych własnych poprawek. */
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
