import type { ExportFormat, TranslationPanel } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { rozbieznoscOdpowiedzi } from './zgodnosc-odpowiedzi';
import type { ZrodloPaneli } from './zrodlo-paneli';

/**
 * Sześć czynności wykonywanych na jednym panelu języka — sama treść, bez
 * elementów. Pasek narzędzi wie, jak wyglądają przyciski; te funkcje wiedzą, co
 * znaczy odpowiedź rdzenia.
 *
 * Każda czynność kończy się zdaniem. Odmowa rdzenia jest tu wynikiem, nie
 * wyjątkiem: wraca sprawozdaniem o wydźwięku negatywnym, które pasek pokazuje
 * w swoim wierszu odpowiedzi.
 *
 * Status `ok` nie wystarcza za wynik. Trzy odpowiedzi bywają puste mimo
 * powodzenia i rozpoznaje się to po samej odpowiedzi, bez wiedzy o wnętrzu
 * rdzenia: pusta ścieżka w `translate.panel.export` i `translate.speech.synthesize`
 * oraz przekład zwrotny tożsamy z treścią panelu. Sprawozdanie wraca wtedy
 * negatywne i mówi, czego brakuje.
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
 * `translate.backtranslation.run`.
 *
 * `trescPanelu` jest tym, co widać w polu tłumaczenia. Gdy rdzeń oddaje
 * dokładnie tę samą treść, przekładu zwrotnego nie było — sprawozdanie mówi to
 * wprost, zamiast podstawiać kopię pod nazwę czynności.
 *
 * `idKanalu` jest wskazaniem ze steru kanału; pusty znaczy „kanał czynny okna"
 * i wtedy żądanie pola `channelId` nie niesie.
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
  // Wykaz pusty znaczy „rdzeń nic nie zgłosił" i jest wynikiem, nie pustką.
  // Zdanie przypisuje ten wynik rdzeniowi, zamiast orzekać o jakości przekładu:
  // klient nie wie, ilu rodzajów niezgodności rdzeń szuka, więc „bez zastrzeżeń"
  // byłoby zapewnieniem szerszym niż odpowiedź.
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
  // Panel oddany przez rdzeń wchodzi do stanu zawsze, także przy rozbieżności:
  // prawdą o panelu jest to, co rdzeń ma u siebie, a nie to, o co go proszono.
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
 * Sprawozdanie z czynności, której wynikiem ma być plik.
 *
 * Pusta ścieżka nie jest ścieżką — gałąź negatywna wyzwala się wyłącznie pustym
 * polem `path` tej jednej odpowiedzi, więc gdy rdzeń ścieżkę odda, okno przestaje
 * meldować brak bez żadnej zmiany tutaj.
 *
 * Gałąź pozytywna potwierdza wskazanie, nie plik: klient dysku nie czyta, wie
 * tylko, że rdzeń ścieżkę podał.
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
