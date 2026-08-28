import {
  przycisk,
  przyciskBezKomendy,
  utworzWierszOdpowiedzi,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { PUSTE, zdanieJezykaRdzenia } from './etykiety-translate';
import { naglowekOkna } from './kontrolki-translate';
import { utworzOdczytyZrodla, type OdczytyZrodla } from './odczyty-zrodla';
import { zbudujPolaZrodla } from './pola-zrodla';
import { utworzStanOkna, type StanOkna } from './stan-okna-translate';
import type { StanTranslate } from './stan-translate';
import { utworzRozwiniecie } from './warstwy-translate';
import { utworzZapisyZrodla } from './zapisy-zrodla';

/**
 * Source Panel to okno wiodące modułu Translate; wprowadzenie albo wklejenie tekstu źródłowego
 * uruchamia jednoczesną aktualizację wszystkich Translation Panels.
 */
export interface OknoSourcePanel {
  element: HTMLElement;
  odswiez(): void;
  /** Wpisuje treść do pola tekstu źródłowego bez jej zapisywania, drogą Format Studio. */
  wstawZrodlo(tekst: string): void;
}

const BRAK_IMPORTU =
  'Importu pliku nie ma czym wykonać: kontrakt nie niesie komendy przyjmującej plik ' +
  'źródłowy w obszarze translate. Wklej treść do pola powyżej albo wydobądź tekst ' +
  'dokumentu w Format Studio.';

const BRAK_SCALANIA =
  'Łączenia i dzielenia pojedynczego segmentu kontrakt nie niesie: segment nie jest w nim bytem ' +
  'o własnym identyfikatorze, a podział przyjmuje wyłącznie cały tekst. Wykonalny jest ponowny ' +
  'podział całości przyciskiem „Segmentuj ponownie".';

const BRAK_REGUL_SEGMENTACJI =
  'Reguł segmentacji SRX nie ma czym przekazać: żądanie podziału niesie sam tekst, a żądanie ' +
  'zapisu źródła — tekst, język i znacznik ponownego podziału. Granice segmentów ustala rdzeń ' +
  'i okno ich nie dostroi.';

const BRAK_HISTORII_ZRODLA =
  'Historii zmian tekstu źródłowego kontrakt nie prowadzi: odpowiedź zapisu niesie język, liczbę ' +
  'pozycji podziału i panele, a wersji poprzedniej nie oddaje żadna komenda obszaru.';

export function utworzOknoSourcePanel(stan: StanTranslate): OknoSourcePanel {
  const okno: StanOkna = utworzStanOkna(PUSTE.zrodlo);
  const pola = zbudujPolaZrodla();
  const odpowiedz = utworzWierszOdpowiedzi();

  const zapisy = utworzZapisyZrodla(
    stan,
    {
      tekst: pola.tekst.kontrolka,
      jezyk: pola.jezyk.kontrolka,
      ponownaSegmentacja: pola.ponowna.kontrolka,
    },
    okno,
    odpowiedz,
  );

  // Język źródłowy widziany po stronie rdzenia; to zdanie mówi, co rdzeń trzyma u siebie.
  const jezykRdzenia = document.createElement('p');
  jezykRdzenia.className = 'mt-zrodlo__jezyk-rdzenia';

  const segmenty = document.createElement('ol');
  segmenty.className = 'mt-segmenty';

  const pasek = document.createElement('div');
  pasek.className = 'mt-pasek';
  pasek.append(
    guzik('Wklej ze schowka', () => void wklejZeSchowka(pola.tekst.kontrolka, odpowiedz)),
    guzik('Zapisz źródło', () => void zapisy.zapisz(), true),
    guzik('Segmentuj ponownie', () => void zapisy.segmentuj()),
    guzik('Rozpoznaj język', () => {
      void zapisy.rozpoznaj().then(odswiezJezykRdzenia);
    }),
    guzik('Importuj plik ▾', () => odpowiedz.pokaz(BRAK_IMPORTU, false)),
  );

  const odczyty: OdczytyZrodla = utworzOdczytyZrodla();

  okno.tresc.append(
    pola.tekst.element,
    pola.zaznaczenie,
    odczyty.element,
    pola.jezyk.element,
    pola.podpowiedzi,
    jezykRdzenia,
    pola.ponowna.element,
    pasek,
    odpowiedz.element,
    segmenty,
    menuTekstu(),
    regulySegmentacji(),
  );

  const element = document.createElement('section');
  element.className = 'mt-okno mt-okno--wiodace';
  element.dataset['okno'] = 'source-panel';
  element.append(naglowekOkna('Source Panel', 'wiodące'), okno.element);

  for (const zdarzenie of ['select', 'keyup', 'click'] as const) {
    pola.tekst.kontrolka.addEventListener(zdarzenie, odswiezZaznaczenie);
  }
  // Odczyty liczą się z treści pola, więc idą za pisaniem, a nie za ogłoszeniem stanu modułu.
  pola.tekst.kontrolka.addEventListener('input', odswiezOdczyty);
  pola.jezyk.kontrolka.addEventListener('input', odswiezJezykRdzenia);

  function odswiezZaznaczenie(): void {
    pokazZaznaczenie(pola.tekst.kontrolka, element, pola.zaznaczenie);
  }

  function odswiezOdczyty(): void {
    odczyty.odswiez(pola.tekst.kontrolka.value, stan.segmenty().length);
  }

  /** Zdanie o języku rdzenia ma trzy wyzwalacze: ogłoszenie stanu, pisanie w polu i powrót rozpoznania. */
  function odswiezJezykRdzenia(): void {
    jezykRdzenia.textContent = zdanieJezykaRdzenia(
      stan.tekstZrodlowy(),
      stan.jezykZrodlowy(),
      pola.jezyk.kontrolka.value,
    );
  }

  /** Przerysowanie nie zdejmuje fazy trwającej ani fazy błędu — te zdejmuje czynność, która je postawiła. */
  function odswiez(): void {
    segmenty.replaceChildren(...stan.segmenty().map(wierszSegmentu));
    odswiezJezykRdzenia();
    odswiezZaznaczenie();
    odswiezOdczyty();
    if (okno.faza() === 'ladowanie' || okno.faza() === 'blad') return;
    if (stan.tekstZrodlowy() === '') {
      okno.puste(PUSTE.zrodlo);
      return;
    }
    okno.gotowe();
  }

  odswiezZaznaczenie();
  odswiezJezykRdzenia();
  odswiezOdczyty();

  return {
    element,
    odswiez,

    wstawZrodlo(tekst) {
      pola.tekst.kontrolka.value = tekst;
      odswiezOdczyty();
      odswiezJezykRdzenia();
      pola.tekst.kontrolka.focus();
      odpowiedz.pokaz(
        `Do pola tekstu źródłowego wpisano ${String(tekst.length)} znaków. Zapisz źródło, ` +
          'aby ruszyły panele — okno nie zapisuje go za Operatora.',
        true,
      );
    },
  };
}

/** Warstwa trzecia mieści operacje na tekście, których kontrakt nie niesie: łączenie i dzielenie segmentów oraz historię zmian tekstu. */
function menuTekstu(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 3,
    nazwa: 'Menu operacji na tekście',
    wyjasnienie: 'Łączenie i dzielenie segmentów oraz historia zmian tekstu źródłowego.',
    znacznik: '⋮',
  });
  rozwiniecie.tresc.append(
    przyciskBezKomendy('Połącz segmenty', BRAK_SCALANIA),
    przyciskBezKomendy('Podziel segment', BRAK_SCALANIA),
    przyciskBezKomendy('Historia zmian tekstu źródłowego', BRAK_HISTORII_ZRODLA),
  );
  return rozwiniecie.element;
}

/** Warstwa czwarta niesie reguły segmentacji SRX, wskazujące własne zasady podziału tekstu źródłowego na segmenty tłumaczenia. */
function regulySegmentacji(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Reguły segmentacji SRX',
    wyjasnienie: BRAK_REGUL_SEGMENTACJI,
    znacznik: '☰',
  });
  rozwiniecie.tresc.append(
    przyciskBezKomendy('Własne reguły segmentacji', BRAK_REGUL_SEGMENTACJI),
  );
  return rozwiniecie.element;
}

/**
 * Przycisk paska czynności. Klasa wariantu jest wpisana jawnie, a nie składana
 * w czasie działania: nazwa złożona z `dn-btn--${…}` nie daje się odnaleźć
 * w kontroli pokrycia klas CSS.
 */
function guzik(napis: string, czynnosc: () => void, glowny = false): HTMLButtonElement {
  const element = przycisk(
    napis,
    glowny ? 'dn-btn dn-btn--sm dn-btn--atrament' : 'dn-btn dn-btn--sm dn-btn--zarys',
  );
  element.addEventListener('click', czynnosc);
  return element;
}

/**
 * Wklejenie ze schowka przeglądarki.
 *
 * Schowek bywa niedostępny — brak zgody albo kontekst bez uprawnienia. Przycisk
 * podaje wtedy skrót klawiszowy jako drogę zastępczą.
 */
async function wklejZeSchowka(
  pole: HTMLTextAreaElement,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  const schowek = navigator.clipboard;
  if (schowek === undefined) {
    odpowiedz.pokaz('Przeglądarka nie udostępnia schowka — wklej treść skrótem Ctrl+V.', false);
    return;
  }
  try {
    const tresc = await schowek.readText();
    pole.value = tresc;
    odpowiedz.pokaz(`Wklejono ${tresc.length} znaków. Zapisz źródło, aby ruszyły panele.`, true);
  } catch (blad) {
    odpowiedz.pokaz(`Schowek odmówił odczytu: ${String(blad)}. Wklej treść skrótem Ctrl+V.`, false);
  }
}

/** Stan zaznaczenia oznacza, że zaznaczony fragment tekstu źródłowego jest przedmiotem operacji wykonywanej z paska czynności. */
function pokazZaznaczenie(
  pole: HTMLTextAreaElement,
  okno: HTMLElement,
  wskaznik: HTMLElement,
): void {
  const dlugosc = pole.selectionEnd - pole.selectionStart;
  okno.dataset['zaznaczenie'] = dlugosc > 0 ? 'tak' : 'nie';
  wskaznik.textContent =
    dlugosc > 0
      ? `Zaznaczono ${dlugosc} znaków — fragment jest przedmiotem operacji.`
      : 'Bez zaznaczenia operacja obejmuje całość tekstu.';
}

function wierszSegmentu(tresc: string): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mt-segmenty__wiersz';
  element.textContent = tresc;
  return element;
}
