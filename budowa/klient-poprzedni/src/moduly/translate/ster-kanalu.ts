import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';
import { WARTOSC_KANALU_OKNA } from './etykiety-translate';
import type { ZrodloKanalowTranslate } from './zrodlo-kanalow-translate';

/**
 * Ster kanału modelu wskazuje, który model wykonuje przekład; mechanizm pochodzi z biblioteki menu
 * drzewa, a plik obsadza go danymi rejestru kanałów.
 */
export interface SterKanalu {
  /** Wiersz z podpisem i sterem, osadzany w formularzu okna. */
  element: HTMLElement;
  /** Kanał wskazany przez Operatora; pusty napis znaczy „kanał czynny okna". */
  wybrany(): string;
  /** Przepisuje drzewo i etykietę z rejestru. */
  odswiez(): void;
  /** Zwija menu i zdejmuje jego nasłuchy dokumentu, obowiązkowe przy usunięciu instancji panelu. */
  zwin(): void;
}

/** Przedrostek klucza pozycji menu odróżnia klucz pozycji drzewa od identyfikatora kanału, który ta pozycja opisuje. */
const KLUCZ_KANALU = 'kanal:';

export interface OpisSteruKanalu {
  /** Podpis nad sterem — mówi, czego dotyczy wskazanie. */
  podpis: string;
  /** Nazwa czynności w zdaniu opisu podaje przykładowo przekład panelu jako czynność kanału. */
  czynnosc: string;
}

export function utworzSterKanalu(
  rejestr: ZrodloKanalowTranslate,
  opis: OpisSteruKanalu,
): SterKanalu {
  let wskazany = '';

  const menu = utworzMenuDrzewo({
    nastawa: 'Kanał modelu',
    naWybor: (klucz) => {
      wskazany = klucz.slice(KLUCZ_KANALU.length);
      odswiez();
    },
  });

  const podpis = document.createElement('span');
  podpis.className = 'dn-pole-etykieta';
  podpis.textContent = opis.podpis;

  const element = document.createElement('div');
  element.className = 'mt-ster';
  element.append(podpis, menu.element);

  /** Przepisanie drzewa z etykietą uchwytu wraca do pozycji domyślnej, gdy wskazania nie ma w wykazie. */
  function odswiez(): void {
    const kanaly = rejestr.kanaly();
    const znaleziony = kanaly.find((kanalModelu) => kanalModelu.id === wskazany);
    const zniknal = wskazany !== '' && znaleziony === undefined;
    if (zniknal) wskazany = '';

    const drzewo: PozycjaMenu[] = [
      {
        rodzaj: 'wybor',
        klucz: KLUCZ_KANALU,
        nazwa: WARTOSC_KANALU_OKNA,
        opis: zniknal ? zdanieZnikniecia() : zdanieDomyslne(rejestr, kanaly.length, opis.czynnosc),
        wybrany: wskazany === '',
      },
      ...kanaly.map<PozycjaMenu>((kanalModelu) => ({
        rodzaj: 'wybor',
        klucz: KLUCZ_KANALU + kanalModelu.id,
        nazwa: kanalModelu.name,
        opis: opisKanalu(kanalModelu.kind, kanalModelu.model),
        wybrany: kanalModelu.id === wskazany,
      })),
    ];

    menu.ustaw(znaleziony === undefined ? WARTOSC_KANALU_OKNA : znaleziony.name, drzewo);
  }

  odswiez();

  return { element, wybrany: () => wskazany, odswiez, zwin: () => menu.zwin() };
}

/** Zdanie przy pozycji domyślnej mówi, co się dzieje przy braku wskazania kanału, zależnie od fazy odczytu rejestru. */
function zdanieDomyslne(
  rejestr: ZrodloKanalowTranslate,
  ile: number,
  czynnosc: string,
): string {
  const faza = rejestr.faza();
  if (faza === 'spoczynek') {
    return `Rejestru kanałów jeszcze nie czytano, więc nie ma czego wskazać. ${czynnosc} wykona kanał czynny okna.`;
  }
  if (faza === 'odczyt') {
    return `Rejestr kanałów jest właśnie czytany. Do jego powrotu ${czynnosc} wykona kanał czynny okna.`;
  }
  if (faza === 'blad') {
    return (
      `Rejestru kanałów nie udało się odczytać: ${rejestr.powod()}. Wyboru nie ma czym obsadzić, ` +
      `ale ${czynnosc} wykona się dalej — żądanie idzie bez wskazania, a kanał bierze rdzeń.`
    );
  }
  if (ile === 0) {
    return (
      'Rdzeń nie zna ani jednego kanału czynnego, więc nie ma czego wskazać. Żądanie idzie bez ' +
      'wskazania; jeśli rdzeń nie ma też kanału domyślnego, odmówi i nazwie ten brak.'
    );
  }
  return `Nie wskazujesz kanału — ${czynnosc} wykona kanał, który rdzeń uznaje za czynny dla tego okna.`;
}

/** Zdanie po zniknięciu wskazanego kanału z wykazu kanałów czynnych informuje o powrocie wskazania do kanału czynnego okna. */
function zdanieZnikniecia(): string {
  return (
    'Kanału wskazanego wcześniej nie ma już wśród czynnych — wskazanie wróciło do kanału czynnego ' +
    'okna, żeby rdzeń nie odmówił przekładu kanałem, którego ktoś w międzyczasie wyłączył.'
  );
}

/** Opis pozycji kanału niesie rodzaj i model kanału, bo sama nazwa rodzaju nie mówi, czym przełoży tekst źródłowy. */
function opisKanalu(rodzaj: string, model: string | undefined): string {
  const nazwaModelu = (model ?? '').trim();
  if (nazwaModelu === '') return `${rodzaj} — rdzeń nie podał modelu tego kanału.`;
  return `${rodzaj} — model ${nazwaModelu}.`;
}
