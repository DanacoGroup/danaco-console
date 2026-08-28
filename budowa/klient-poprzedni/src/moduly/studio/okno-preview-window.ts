import { Command, type DesignAsset } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { przyciskBrakuDrogi } from '../design/brak-drogi';
import { BRAKI, OKNO_PREVIEW_WINDOW } from '../design/etykiety-designu';
import { nazwaZasobu } from '../design/karta-zasobu';
import { utworzPlytePodgladu, type PlytaPodgladu } from '../design/plyta-podgladu';
import {
  utworzPorownanieWariantow,
  type PorownanieWariantow,
} from '../design/porownanie-wariantow';
import { skutekPrzekazania } from '../design/skutek-designu';
import type { StanDesignu } from '../design/stan-designu';
import { naglowekOkna, utworzStanOkna, type StanOkna } from '../design/stan-okna';

/**
 * Kod modułu Library — jedyny cel przekazania z wariantu Designu.
 *
 * Jedna wartość, nie kopia katalogu modułów: okno przekazuje wyłącznie tam.
 * Wykaz modułów do wyboru swobodnego bierze Assets Panel z komendy `module.list`.
 */
const KOD_LIBRARY = 'library';

/** Preview Window, okno pomocnicze modułu Design pokazujące jedyny dziś wariant: podgląd zasobu wizualnego bez obrazu rdzenia. */
export interface OknoPreviewWindow {
  element: HTMLElement;
  odswiez(): void;
}

/** Wariant okna podglądu, jedyne miejsce, w którym rozchodzą się moduł Studio i moduł Design, każdy z własnym zapleczem. */
export type WariantPodgladu = { rodzaj: 'design'; stan: StanDesignu };

export function utworzOknoPreviewWindow(wariant: WariantPodgladu): OknoPreviewWindow {
  return podgladZasobu(wariant.stan);
}

/** Wariant Designu, podgląd zasobu wizualnego wskazanego w Assets Panel, bez formatu eksportu ani decyzji, bo nie ma dokumentu. */
function podgladZasobu(stan: StanDesignu): OknoPreviewWindow {
  const okno: StanOkna = utworzStanOkna();
  const plyta: PlytaPodgladu = utworzPlytePodgladu();
  const porownanie: PorownanieWariantow = utworzPorownanieWariantow();

  const doLibrary = przycisk('Przekaż do Library', 'dn-btn dn-btn--sm dn-btn--atrament');
  doLibrary.dataset['czynnosc'] = 'do-library';
  const odpowiedz = utworzWierszOdpowiedzi();

  // Przełącznika tła w oknie nie ma: tło podglądu ma sens za treścią obrazu, nie za zdaniem tekstu.
  const oTle = document.createElement('p');
  oTle.className = 'dn-pole-opis md-podglad__o-tle';
  oTle.textContent =
    'Przełącznika tła podglądu tu nie ma: tło służy ocenie przezroczystości OBRAZU, ' +
    'a rdzeń treści obrazu nie oddaje. Wróci razem z polem uri wypełnionym przez rdzeń.';

  const rzadBrakow = document.createElement('div');
  rzadBrakow.className = 'md-braki__rzad';
  rzadBrakow.append(
    przyciskBrakuDrogi(BRAKI.eksportZasobu),
    przyciskBrakuDrogi(BRAKI.akceptacjaWyniku),
  );

  okno.tresc.append(
    plyta.element,
    oTle,
    porownanie.element,
    doLibrary,
    odpowiedz.element,
    rzadBrakow,
  );

  const element = document.createElement('section');
  element.className = 'md-okno md-okno--pomocnicze';
  element.dataset['okno'] = OKNO_PREVIEW_WINDOW.kod;
  element.append(naglowekOkna(OKNO_PREVIEW_WINDOW.nazwa, OKNO_PREVIEW_WINDOW.rola), okno.element);

  async function przekaz(): Promise<void> {
    const zasob = stan.wybrany();
    if (zasob === null) {
      odpowiedz.pokaz('Przekazanie dotyczy zasobu wskazanego — wskaż go w Assets Panel.', false);
      return;
    }
    if (stan.idOkna() === '') {
      odpowiedz.pokaz(
        `Komenda ${Command.ContextTransfer} wymaga okna źródłowego. ${stan.opisOkna()}`,
        false,
      );
      return;
    }
    odpowiedz.pokaz(`Przekazywanie zasobu „${nazwaZasobu(zasob)}" do Library…`, true);
    // Pod czuwaniem: przekazanie zrywa się jak każde inne wywołanie, więc komunikat nie wisi bez końca.
    const wynik = await stan.czuwanie.prowadz(
      'przekazanie zasobu do Library',
      stan.zaplecze.przekaz({
        idOknaZrodlowego: stan.idOkna(),
        kodModuluDocelowego: KOD_LIBRARY,
        komplet: kompletZasobu(zasob),
      }),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Przekazanie do Library', wynik.blad), false);
      return;
    }
    const skutek = skutekPrzekazania(KOD_LIBRARY, wynik.wynik, nazwaZasobu(zasob));
    odpowiedz.pokaz(skutek.zdanie, skutek.udany);
  }

  doLibrary.addEventListener('click', () => void przekaz());

  function odswiez(): void {
    const zasob = stan.wybrany();
    plyta.pokaz(zasob);
    porownanie.odswiez(stan.zasoby(), zasob);
    if (zasob === null) {
      okno.puste('Podgląd bez wskazanego zasobu — wskaż zasób w Assets Panel.');
      return;
    }
    okno.gotowe();
  }

  odswiez();

  return { element, odswiez };
}

/**
 * Komplet kontekstu niosący zasób do modułu docelowego (wariant Designu).
 *
 * `ContextBundle` nie ma pola zasobu wizualnego, więc zasób jedzie polem
 * `documentIds`, a polecenie nazywa go po imieniu. Ta sama rozbieżność występuje
 * w panelu metadanych.
 */
function kompletZasobu(zasob: DesignAsset): { prompt: string; documentIds: string[] } {
  return {
    prompt: `Zasób wizualny modułu Design z podglądu: ${nazwaZasobu(zasob)} (${zasob.kind}).`,
    documentIds: [zasob.id],
  };
}
