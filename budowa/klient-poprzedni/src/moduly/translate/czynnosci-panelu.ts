import type { ExportFormat, TranslationPanel } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { rozbieznoscOdpowiedzi } from './zgodnosc-odpowiedzi';
import type { ZrodloPaneli } from './zrodlo-paneli';

/**
 * Sześć czynności wykonywanych na jednym panelu języka: każda kończy się zdaniem, a odmowa rdzenia jest wynikiem, nie wyjątkiem.
 */
export interface Sprawozdanie {
  tresc: string;
  powodzenie: boolean;
}

export function udane(tresc: string): Sprawozdanie {
  return { tresc, powodzenie: true };
}

export function odmowa(
  czynnosc: string,
  blad?: { code?: string; message?: string },
): Sprawozdanie {
  return { tresc: opisOdmowy(czynnosc, blad?.code, blad?.message), powodzenie: false };
}

/**
 * Przekład zwrotny wykonania translate.backtranslation.run; treść tożsama z panelem znaczy, że przekładu zwrotnego nie było.
 */
export async function tlumaczZwrotnie(
  zrodlo: ZrodloPaneli,
  idPanelu: string,
  trescPanelu: string,
  idKanalu: string,
): Promise<Sprawozdanie> {
  const wynik = await zrodlo.tlumaczZwrotnie(idPanelu, idKanalu);
  if (!wynik.udany || wynik.wynik === undefined) return odmowa('Tłumaczenie zwrotne', wynik.blad);
  const zwrotne = wynik.wynik.text;
  if (zwrotne.trim() === '') {
    return {
      tresc:
        'Tłumaczenie zwrotne: rdzeń odpowiedział, ale kolumna zwrotna wróciła pusta — ' +
        'nie ma czego porównać z przekładem.',
      powodzenie: false,
    };
  }
  if (zwrotne.trim() === trescPanelu.trim()) {
    return {
      tresc:
        'Tłumaczenie zwrotne: rdzeń oddał DOSŁOWNIE treść tego panelu, znak w znak. ' +
        'To jest przepisanie, nie przekład odwrotny — porównanie z nim niczego nie sprawdzi.',
      powodzenie: false,
    };
  }
  return udane(`Tłumaczenie zwrotne: ${zwrotne}`);
}

export async function kontrolaJakosci(
  zrodlo: ZrodloPaneli,
  idPanelu: string,
): Promise<Sprawozdanie> {
  const wynik = await zrodlo.kontrolaJakosci(idPanelu);
  if (!wynik.udany || wynik.wynik === undefined) return odmowa('Kontrola jakości', wynik.blad);
  // Wykaz pusty znaczy, że rdzeń nic nie zgłosił, i jest wynikiem kontroli, nie pustką odpowiedzi.
  const zastrzezenia = wynik.wynik.issues;
  return udane(
    zastrzezenia.length === 0
      ? 'Kontrola jakości: rdzeń nie zgłosił żadnego zastrzeżenia.'
      : `Kontrola jakości: ${zastrzezenia.join(' · ')}`,
  );
}

export async function eksportuj(
  zrodlo: ZrodloPaneli,
  idPanelu: string,
  format: ExportFormat,
): Promise<Sprawozdanie> {
  const wynik = await zrodlo.eksportuj(idPanelu, format);
  if (!wynik.udany || wynik.wynik === undefined) return odmowa('Eksport panelu', wynik.blad);
  return zdanieOSciezce('Eksport panelu', wynik.wynik.path);
}

/**
 * `translate.panel.tone.set`.
 *
 * Zdanie bierze ton z panelu, który wrócił, a nie z wypełnionego pola:
 * `TranslatePanelToneSetResponse` niesie cały panel po zmianie, więc gdy rdzeń
 * ton znormalizuje albo odrzuci, sprawozdanie nadal mówi prawdę o zapisie.
 */
export async function zmienTon(
  zrodlo: ZrodloPaneli,
  idPanelu: string,
  ton: string,
  wchlon: (panel: TranslationPanel) => void,
): Promise<Sprawozdanie> {
  if (ton === '') {
    return { tresc: 'Wskaż ton — rdzeń odmówi zmiany bez jego nazwy.', powodzenie: false };
  }
  const wynik = await zrodlo.ustawTon(idPanelu, ton);
  if (!wynik.udany || wynik.wynik === undefined) return odmowa('Zmiana tonu', wynik.blad);
  const panel = wynik.wynik.panel;
  // Panel oddany przez rdzeń wchodzi do stanu zawsze, także przy rozbieżności z tym, o co proszono.
  wchlon(panel);
  const rozbiezne = rozbieznoscOdpowiedzi('Zmiana tonu', [
    { nazwa: 'ton panelu', zamowione: ton, oddane: panel.tone },
  ]);
  if (rozbiezne !== null) return { tresc: rozbiezne, powodzenie: false };
  return udane(`Ton panelu ustawiony na „${panel.tone ?? ''}" — tak oddał go rdzeń.`);
}

export async function odsluchaj(zrodlo: ZrodloPaneli, idPanelu: string): Promise<Sprawozdanie> {
  const wynik = await zrodlo.odsluchaj(idPanelu);
  if (!wynik.udany || wynik.wynik === undefined) return odmowa('Synteza mowy', wynik.blad);
  return zdanieOSciezce('Synteza mowy', wynik.wynik.path);
}

/**
 * Sprawozdanie z czynności, której wynikiem ma być zapisany plik; pusta ścieżka wyzwala gałąź negatywną.
 */
function zdanieOSciezce(czynnosc: string, sciezka: string): Sprawozdanie {
  if (sciezka.trim() === '') {
    return {
      tresc:
        `${czynnosc}: rdzeń przyjął zlecenie, ale nie oddał ścieżki pliku — pole „path" ` +
        'odpowiedzi jest puste, więc pliku nie ma czym wskazać ani gdzie szukać.',
      powodzenie: false,
    };
  }
  return udane(
    `${czynnosc}: rdzeń wskazał ścieżkę wyniku ${sciezka}. Okno nie czyta dysku — ` +
      'potwierdza wskazanie rdzenia, nie istnienie pliku.',
  );
}

export async function podpowiedzPamieci(
  zrodlo: ZrodloPaneli,
  idPanelu: string,
  segment: string,
): Promise<Sprawozdanie> {
  if (segment === '') {
    return {
      tresc: 'Pamięć tłumaczeń potrzebuje segmentu źródłowego — zapisz najpierw źródło.',
      powodzenie: false,
    };
  }
  const wynik = await zrodlo.podpowiedzPamieci(idPanelu, segment);
  if (!wynik.udany || wynik.wynik === undefined) return odmowa('Podpowiedź pamięci', wynik.blad);
  const podpowiedzi = wynik.wynik.suggestions;
  return udane(
    podpowiedzi.length === 0
      ? 'Pamięć tłumaczeń nie ma dopasowania dla tego segmentu.'
      : `Pamięć tłumaczeń: ${podpowiedzi.join(' · ')}`,
  );
}
