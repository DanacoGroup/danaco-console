import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';
import { WARTOSC_KANALU_OKNA } from './etykiety-translate';
import type { ZrodloKanalowTranslate } from './zrodlo-kanalow-translate';

/**
 * Ster kanału modelu — wskazanie, który model wykonuje przekład.
 *
 * To jest ster, nie wyświetlacz: etykieta uchwytu niesie wartość bieżącą
 * nastawy (nazwę kanału, a nie słowo „Kanał"), kliknięcie rozwija wybór, wybór
 * zmienia nastawę. Nazwa rodzajowa idzie do `aria-label` mechanizmu, bo czytnik
 * ekranu musi wiedzieć, czego dotyczy wartość, której nazwa sama tego nie mówi.
 *
 * Mechanizm pochodzi z biblioteki, własnego menu tu nie ma. Całe rozwijanie,
 * szukanie, znacznik wyboru i obsługa klawiatury idą z
 * `komponenty/menu-drzewo.ts`; ten plik wyłącznie obsadza go danymi rejestru
 * i tłumaczy klucz pozycji na identyfikator kanału.
 *
 * Podpis jest blokiem, nie `<label>`. Etykieta bez `for` wiąże się z pierwszym
 * potomkiem, który da się etykietować, a tym jest uchwyt menu. Skutek byłby
 * dwojaki i oba razy zły: kliknięcie w podpis otwierałoby menu (podpis nie jest
 * sterem), a nazwa dostępna uchwytu konkurowałaby z `aria-label`, które
 * mechanizm ustawia sam i które niesie bieżącą wartość. To samo rozstrzygnięcie
 * co w `poczta/wiersz-steru.ts`; przepisane, bo tamten plik należy do innego
 * modułu i niesie jego klasy rozkładu.
 *
 * Pusty wybór jest wyborem, nie brakiem. Pozycja „kanał czynny okna" stoi
 * w drzewie zawsze i jest domyślna, bo dokładnie to opisuje kontrakt: brak
 * wskazania bierze kanał czynny okna. Dzięki niej jest droga powrotna do
 * zachowania sprzed wskazania, a żądanie bez `channelId` nie jest stanem
 * nieosiągalnym z ekranu.
 *
 * Rejestr, który nie dotarł, nie odbiera czynności: drzewo ma wtedy samą
 * pozycję domyślną, a powód odmowy stoi przy niej jako opis. Ster nie wygasza
 * się, nie znika i nie zatrzymuje przycisku obok.
 */
export interface SterKanalu {
  /** Wiersz z podpisem i sterem, osadzany w formularzu okna. */
  element: HTMLElement;
  /** Kanał wskazany przez Operatora; pusty napis znaczy „kanał czynny okna". */
  wybrany(): string;
  /** Przepisuje drzewo i etykietę z rejestru. */
  odswiez(): void;
  /**
   * Zwija menu i zdejmuje jego nasłuchy dokumentu — obowiązkowe przy usunięciu.
   *
   * Menu rozwinięte wiesza `pointerdown` na dokumencie (tak zamyka się po
   * kliknięciu poza sobą). Instancja panelu, która znika z wykazu rdzenia, jest
   * z dokumentu usuwana — gdyby jej ster był wtedy rozwinięty, nasłuch zostałby
   * na dokumencie na zawsze i wołałby do elementu, którego już nie ma. Panele
   * przychodzą i znikają z każdą zmianą wykazu, więc to nie jest przypadek
   * teoretyczny.
   */
  zwin(): void;
}

/** Przedrostek klucza pozycji — mechanizm oddaje klucz, nie identyfikator. */
const KLUCZ_KANALU = 'kanal:';

export interface OpisSteruKanalu {
  /** Podpis nad sterem — mówi, czego dotyczy wskazanie. */
  podpis: string;
  /** Nazwa czynności w zdaniu opisu, np. „przekład panelu". */
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

  /**
   * Przepisanie drzewa wraz z etykietą uchwytu.
   *
   * Wskazanie, którego nie ma już w wykazie, wraca do domyślnego. Kanał
   * wyłączony w innym oknie znika z `channel.list(enabledOnly)`, a rdzeń
   * odmówiłby przekładu wskazaniem na niego. Cichy powrót byłby jednak zmianą
   * nastawy bez wiedzy Operatora, więc towarzyszy mu zdanie przy pozycji
   * domyślnej.
   */
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

/** Zdanie przy pozycji domyślnej — mówi, co się dzieje przy braku wskazania. */
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

/** Zdanie po zniknięciu wskazanego kanału z wykazu kanałów czynnych. */
function zdanieZnikniecia(): string {
  return (
    'Kanału wskazanego wcześniej nie ma już wśród czynnych — wskazanie wróciło do kanału czynnego ' +
    'okna, żeby rdzeń nie odmówił przekładu kanałem, którego ktoś w międzyczasie wyłączył.'
  );
}

/** Opis pozycji kanału: rodzaj i model, bo sama nazwa nie mówi, czym przełoży. */
function opisKanalu(rodzaj: string, model: string | undefined): string {
  const nazwaModelu = (model ?? '').trim();
  if (nazwaModelu === '') return `${rodzaj} — rdzeń nie podał modelu tego kanału.`;
  return `${rodzaj} — model ${nazwaModelu}.`;
}
