import type { LibraryVersion } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';
import type { StanTresci } from './dostepnosc-tresci';

/**
 * Jeden wiersz historii wersji dokumentu.
 *
 * Wersje narastają przy każdej zmianie dokumentu w dowolnym module, w tym przy
 * zmianie wykonanej przez model — dlatego wiersz pokazuje sprawcę osobno od etykiety.
 *
 * Osiągalność treści bierze się z odpowiedzi rdzenia (`dostepnosc-tresci.ts`), nie
 * z obecności sumy kontrolnej w opisie wersji: `checksum` bywa obecna przy wersji,
 * której treści rdzeń nie oddaje, i nieobecna przy wersji, którą oddaje w całości.
 * Odpowiedź dotyczy treści, którą dokument niesie jako bieżącą, więc werdykt siada
 * wyłącznie na wierszu bieżącym — kontrakt nie ma komendy pytającej o treść wersji
 * niebieżącej. Werdykt `odmowa` znaczy, że rdzeń treści nie oddał, ale też o niej
 * nie orzekł, i nie odbiera przycisku „Przywróć".
 */
export interface OpisWierszaWersji {
  /** Czy wersja jest tą, którą plik niesie jako bieżącą. */
  biezaca: boolean;
  /** Odpowiedź rdzenia o treści dokumentu, którego ta wersja dotyczy. */
  tresc: StanTresci;
  naPrzywrocenie(): void;
}

export function utworzWierszWersji(
  wersja: LibraryVersion,
  opis: OpisWierszaWersji,
): HTMLElement {
  const etykieta = document.createElement('span');
  etykieta.className = 'ml-wersja__etykieta';
  etykieta.textContent = wersja.label ?? wersja.id;

  const metryka = document.createElement('span');
  metryka.className = 'ml-wersja__metryka';
  metryka.textContent = opisMetryki(wersja);

  const przywroc = przycisk('Przywróć', 'dn-btn dn-btn--sm dn-btn--zarys');
  przywroc.dataset['wersja'] = wersja.id;
  przywroc.addEventListener('click', () => opis.naPrzywrocenie());

  const element = document.createElement('li');
  element.className = 'ml-wersja';
  element.dataset['biezaca'] = opis.biezaca ? 'tak' : 'nie';
  // Wiersz niebieżący nie jest „nieznany" z braku odpowiedzi — nie był pytany.
  element.dataset['tresc'] = opis.biezaca ? opis.tresc.werdykt : 'niepytana';
  element.append(etykieta, metryka, przywroc);

  const zdanie = opis.biezaca ? zdanieOTresci(opis.tresc) : '';
  if (zdanie !== '') {
    const powod = document.createElement('span');
    powod.className = 'ml-wersja__powod';
    powod.textContent = zdanie;
    element.append(powod);
  }

  return element;
}

/**
 * Zdanie wiersza o treści dokumentu — jedno na werdykt, żadne bez werdyktu.
 *
 * Wołane wyłącznie dla wiersza bieżącego: odpowiedź rdzenia o treści dokumentu
 * jest odpowiedzią o treści tej właśnie wersji. „Nie ma do czego wrócić" pada
 * tylko po werdykcie `brak`; `odmowa` oznacza nieudany odczyt, a `odwolanie` —
 * wskazanie miejsca zamiast bajtów.
 */
function zdanieOTresci(tresc: StanTresci): string {
  if (tresc.werdykt === 'osiagalna') return 'wersja bieżąca — rdzeń oddaje jej treść';
  if (tresc.werdykt === 'brak') return `nie ma do czego wrócić — ${tresc.powod}`;
  if (tresc.werdykt === 'odmowa') return `nie wiadomo, czy jest do czego wracać — ${tresc.powod}`;
  if (tresc.werdykt === 'odwolanie') return `treści nie widziałem — ${tresc.powod}`;
  return '';
}

/** Metryka wiersza: sprawca, rozmiar i chwila powstania wersji. */
function opisMetryki(wersja: LibraryVersion): string {
  const czesci: string[] = [];
  if (wersja.author !== undefined) czesci.push(`sprawca: ${wersja.author}`);
  if (wersja.sizeBytes !== undefined) czesci.push(`${wersja.sizeBytes} B`);
  czesci.push(new Date(wersja.createdAt).toLocaleString('pl'));
  return czesci.join(' — ');
}
