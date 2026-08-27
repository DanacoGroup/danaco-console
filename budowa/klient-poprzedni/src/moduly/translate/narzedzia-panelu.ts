import type { ExportFormat, TranslationPanel } from '../../../../shared/contract';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  type PoleFormularza,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import {
  eksportuj,
  kontrolaJakosci,
  odsluchaj,
  podpowiedzPamieci,
  tlumaczZwrotnie,
  zmienTon,
  type Sprawozdanie,
} from './czynnosci-panelu';
import { FORMATY_EKSPORTU, OBJASNIENIA, PODPOWIEDZ_TONU } from './etykiety-translate';
import { dopnijDymek, podepnijPodpowiedz } from './kontrolki-translate';
import { utworzSterKanalu } from './ster-kanalu';
import type { ZrodloKanalowTranslate } from './zrodlo-kanalow-translate';
import type { ZrodloPaneli } from './zrodlo-paneli';

/**
 * Sześć narzędzi kontekstowych jednej instancji Translation Panels tworzy pasek czynności nad
 * treścią tłumaczenia, wydzielony z panelu i sterowany osobnymi komendami kontraktu.
 */
export interface KontekstNarzedzi {
  /** Panel, którego dotyczą czynności. */
  idPanelu(): string;
  /** Segment źródłowy dla podpowiedzi pamięci tłumaczeń. */
  segment(): string;
  /** Treść tłumaczenia widoczna w panelu, potrzebna tłumaczeniu zwrotnemu do rozpoznania braku przekładu. */
  tresc(): string;
  /** Wciąga panel oddany przez rdzeń po zmianie tonu. */
  wchlon(panel: TranslationPanel): void;
  /** Rejestr kanałów modelu — obsada steru `channelId` dla tłumaczenia zwrotnego. */
  kanaly: ZrodloKanalowTranslate;
}

export interface NarzedziaPanelu {
  element: HTMLElement;
  /** Przerysowuje pasek danymi rdzenia, obecnie tonem panelu, zgodnym z wartością w nagłówku panelu. */
  odswiez(panel: TranslationPanel): void;
  /** Zwija ster kanału i zdejmuje jego nasłuchy — wołane przy usunięciu panelu. */
  rozlacz(): void;
}

export function utworzNarzedziaPanelu(
  zrodlo: ZrodloPaneli,
  kontekst: KontekstNarzedzi,
  odpowiedz: WierszOdpowiedzi,
  przyrostek: string,
): NarzedziaPanelu {
  const ton = poleTonu(przyrostek);
  const format = poleFormatu();
  const sterKanalu = utworzSterKanalu(kontekst.kanaly, {
    podpis: 'Kanał tłumaczenia zwrotnego',
    czynnosc: 'przekład zwrotny',
  });

  const czynnosci: readonly [string, () => Promise<Sprawozdanie>][] = [
    [
      'Tłumaczenie zwrotne',
      () =>
        tlumaczZwrotnie(zrodlo, kontekst.idPanelu(), kontekst.tresc(), sterKanalu.wybrany()),
    ],
    ['Kontrola jakości', () => kontrolaJakosci(zrodlo, kontekst.idPanelu())],
    [
      'Eksportuj panel',
      () => eksportuj(zrodlo, kontekst.idPanelu(), format.kontrolka.value as ExportFormat),
    ],
    [
      'Ustaw ton',
      () => zmienTon(zrodlo, kontekst.idPanelu(), ton.pole.kontrolka.value.trim(), kontekst.wchlon),
    ],
    ['Odsłuchaj', () => odsluchaj(zrodlo, kontekst.idPanelu())],
    [
      'Podpowiedź pamięci',
      () => podpowiedzPamieci(zrodlo, kontekst.idPanelu(), kontekst.segment()),
    ],
  ];

  const pasek = document.createElement('div');
  pasek.className = 'mt-pasek';
  for (const [nazwa, czynnosc] of czynnosci) {
    const guzik = przycisk(nazwa, 'dn-btn dn-btn--sm dn-btn--zarys');
    guzik.addEventListener('click', () => void wykonaj(nazwa, czynnosc));
    pasek.append(guzik);
  }

  const element = document.createElement('div');
  element.className = 'mt-narzedzia';
  element.append(ton.pole.element, ton.podpowiedzi, format.element, sterKanalu.element, pasek);

  /** Ton z rdzenia wchodzi do pola, chyba że pole jest pod ogniskiem operatora piszącego własną wartość. */
  function odswiez(panel: TranslationPanel): void {
    sterKanalu.odswiez();
    if (document.activeElement === ton.pole.kontrolka) return;
    ton.pole.kontrolka.value = panel.tone ?? '';
  }

  /** Zapowiedź, wykonanie i zdanie o wyniku — naciśnięcie zawsze odpowiada. */
  async function wykonaj(nazwa: string, czynnosc: () => Promise<Sprawozdanie>): Promise<void> {
    if (kontekst.idPanelu() === '') {
      odpowiedz.pokaz(`${nazwa}: panel nie ma jeszcze identyfikatora z rdzenia.`, false);
      return;
    }
    odpowiedz.pokaz(`${nazwa}…`, true);
    const sprawozdanie = await czynnosc();
    odpowiedz.pokaz(sprawozdanie.tresc, sprawozdanie.powodzenie);
  }

  return { element, odswiez, rozlacz: () => sterKanalu.zwin() };
}

/** Pole tonu wraz z podpowiedzią wartości: ton jest w kontrakcie dowolnym napisem, więc pole nie ogranicza go do zamkniętego katalogu wartości. */
function poleTonu(przyrostek: string): {
  pole: PoleFormularza<HTMLInputElement>;
  podpowiedzi: HTMLDataListElement;
} {
  const pole = poleTekstowe({ etykieta: 'Ton panelu', podpowiedz: 'np. formalny' });
  dopnijDymek(pole.element, OBJASNIENIA.tonPanelu);
  return { pole, podpowiedzi: podepnijPodpowiedz(pole.kontrolka, `mt-tony-${przyrostek}`, PODPOWIEDZ_TONU) };
}

/** Lista formatów eksportu odpowiada katalogowi zamkniętemu kontraktu ExportFormat i nie przyjmuje wartości spoza tego katalogu przy wyborze formatu. */
function poleFormatu(): PoleFormularza<HTMLSelectElement> {
  const pole = poleWyboru(
    { etykieta: 'Format eksportu' },
    FORMATY_EKSPORTU.map((pozycja) => ({ wartosc: pozycja.wartosc, etykieta: pozycja.etykieta })),
  );
  dopnijDymek(pole.element, OBJASNIENIA.formatEksportu);
  return pole;
}
