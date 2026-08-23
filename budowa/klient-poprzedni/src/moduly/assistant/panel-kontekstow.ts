import type { ConfigScope } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  przelacznik,
  przyciskAkcji,
  utworzWierszOdpowiedzi,
  wiersz,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { BRAKI } from './braki-kontraktu';
import { NAZWY_ZASIEGOW, ODCZYTY, POZIOMY_PAMIECI } from './etykiety-assistant';
import type { StanAssistant } from './stan-assistant';
import type { ZrodloPamieci } from './zrodlo-pamieci';

/**
 * Zakładka kontekstów Memory & Context Manager — poziomy pamięci karty sesji.
 *
 * Kontrakt niesie jeden mechanizm przełączania pamięci w trakcie pracy:
 * `memory.toggle` ustala, które poziomy zasięgu karta sesji widzi i czy zapis
 * pamięci jest czynny. To jest kontekst, który rzeczywiście da się przełączyć —
 * i tak też okno go nazywa.
 *
 * Nazwanych zestawów pamięci („praca", „dom", „projekt X") okno nie udaje.
 * Kontrakt nie ma bytu, w którym taki zestaw miałby zamieszkać: `memory.*`
 * zna wpis i poziom zasięgu, a nie nazwany profil pamięci. Brak nazywa przycisk
 * obok, zamiast pola, którego rdzeń nie zapisze.
 *
 * Przestawienie idzie jednym wywołaniem obejmującym oba pola naraz. Rdzeń
 * odczytuje `levels` jako stan po przestawieniu (pusta lista zostawia poziomy
 * bez zmian), więc wysyłanie samego jednego przełącznika zamieniałoby resztę
 * wykazu w „bez zmian" przy każdym kliknięciu — a Operator odznaczający
 * ostatni poziom nie miałby jak wyłączyć wszystkich.
 */
export interface PanelKontekstow {
  element: HTMLElement;
}

export function utworzPanelKontekstow(
  stan: StanAssistant,
  zrodlo: ZrodloPamieci,
): PanelKontekstow {
  const odpowiedz: WierszOdpowiedzi = utworzWierszOdpowiedzi();

  const poziomy = new Map<ConfigScope, HTMLInputElement>();
  const wykazPoziomow = document.createElement('div');
  wykazPoziomow.className = 'ma-poziomy';
  for (const poziom of POZIOMY_PAMIECI) {
    const kontrolka = przelacznik(`Poziom pamięci: ${NAZWY_ZASIEGOW[poziom]}`);
    kontrolka.dataset['poziom'] = poziom;
    poziomy.set(poziom, kontrolka);
    wykazPoziomow.append(
      wiersz(NAZWY_ZASIEGOW[poziom], kontrolka, { klasa: 'ma-wiersz ma-wiersz--waski' }),
    );
  }

  const zapis = przelacznik('Zapis pamięci czynny');
  zapis.checked = true;

  const przestaw = przyciskAkcji('Przestaw poziomy pamięci', 'dn-btn dn-btn--sm dn-btn--atrament');
  przestaw.addEventListener('click', () => void przestawPoziomy());

  const przyciski = document.createElement('div');
  przyciski.className = 'ma-formularz__przyciski';
  przyciski.append(przestaw);

  const element = document.createElement('div');
  element.className = 'ma-obszar';
  element.dataset['obszar'] = 'konteksty';
  element.append(
    opisPoziomow(),
    wykazPoziomow,
    wiersz('Zapis pamięci czynny', zapis, {
      klasa: 'ma-wiersz',
      objasnienie:
        'Wyłączenie zostawia pamięć do odczytu: asystent widzi ustalenia, ale nowych ' +
        'nie dopisuje.',
    }),
    przyciski,
    odpowiedz.element,
  );

  async function przestawPoziomy(): Promise<void> {
    const sesja = stan.idSesji();
    if (sesja === '') {
      odpowiedz.pokaz(BRAKI.brakSesji, false);
      return;
    }
    const wybrane = [...poziomy.entries()]
      .filter(([, kontrolka]) => kontrolka.checked)
      .map(([poziom]) => poziom);
    odpowiedz.pokaz(ODCZYTY.poziomy, true);
    const wynik = await zrodlo.przestaw({
      sessionId: sesja,
      levels: wybrane,
      writeEnabled: zapis.checked,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Przestawienie poziomów pamięci', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    // Stan pokazany po przestawieniu bierze się z odpowiedzi rdzenia: pusta
    // lista poziomów zostawia je bez zmian, więc kontrolki odbijające samo
    // zamówienie kłamałyby przy każdym takim wywołaniu.
    const po = wynik.wynik;
    for (const [poziom, kontrolka] of poziomy) {
      kontrolka.checked = po.levels.includes(poziom);
    }
    zapis.checked = po.writeEnabled;
    odpowiedz.pokaz(
      `Rdzeń trzyma włączone poziomy: ${nazwijPoziomy(po.levels)}; ` +
        `zapis pamięci ${po.writeEnabled ? 'czynny' : 'wstrzymany'}.`,
      true,
    );
  }

  return { element };
}

/** Wykaz poziomów po nazwach; pusty wykaz nazywamy słowem, nie pustką. */
function nazwijPoziomy(poziomy: readonly ConfigScope[]): string {
  if (poziomy.length === 0) return 'żaden';
  return poziomy.map((poziom) => NAZWY_ZASIEGOW[poziom]).join(', ');
}

/** Zdanie o tym, czym poziomy pamięci są dla bieżącej rozmowy. */
function opisPoziomow(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent =
    'Poziomy włączone dla tej karty sesji rozstrzygają, którą pamięć asystent widzi ' +
    'w rozmowie. Przestawienie obejmuje wszystkie poziomy naraz — wykaz poniżej jest ' +
    'stanem po zmianie, nie różnicą wobec niego.';
  return element;
}
