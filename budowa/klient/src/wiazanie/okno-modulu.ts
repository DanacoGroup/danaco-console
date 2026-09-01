/**
 * Wiązanie okna modułu z rdzeniem. Znacznik niesie biblioteka Właściciela;
 * ten plik wiąże kartę okna roboczego z oknem komunikacji rdzenia, opisuje
 * nagłówek okna wartościami rejestru i wystawia wpisy, z których korzystają
 * wiązania szczegółowe modułów.
 */
import { Command } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { przypiszOknoKomunikacji } from './okna-robocze.ts';

/**
 * Wiąże wnętrze karty modułu, dla której okna komunikacji w rdzeniu jeszcze
 * nie ma — okno zakłada dopiero praca podjęta w karcie. Kod modułu nie wchodzi
 * tu w wiązanie: karta bez okna komunikacji nie ma czego z rdzeniem uzgodnić.
 * Nazwa środowiska wchodzi w nagłówek okna; węzły szuka się od korzenia karty.
 */
export function zwiazOkno(
  kanal: Kanal,
  _kodModulu: string,
  nazwaSrodowiska: string,
  korzen: ParentNode,
): void {
  void opiszNaglowek(kanal, nazwaSrodowiska, korzen);
}

/**
 * Wiąże wnętrze okna, które w rdzeniu już stoi: karta bieżąca dostaje jego
 * identyfikator, więc powrót do karty sesji nie zakłada okna drugi raz,
 * a zamknięcie karty zamyka okno także w rdzeniu.
 */
export function zwiazOknoStojace(
  kanal: Kanal,
  kodModulu: string,
  nazwaSrodowiska: string,
  idOkna: string,
  korzen: ParentNode,
): void {
  if (!przypiszOknoKomunikacji(kodModulu, idOkna)) {
    oglos('Okno modułu', 'Karta bieżąca niesie inny moduł niż okno wskazane przez rdzeń — '
      + 'zamknięcie karty nie zamknie tego okna.', 'ostrzezenie');
  }
  void opiszNaglowek(kanal, nazwaSrodowiska, korzen);
}

/**
 * Opisuje nagłówek okna komunikacji karty środowiskiem wejścia i modelem
 * kanału; nagłówek szuka się od korzenia karty, bo karty stoją w płótnie obok
 * siebie. Pozostałe podpisy prototypu — wysiłek, wykonawca, pamięć, rola,
 * format, tura — schodzą: ich wartości niesie dopiero praca podjęta w oknie.
 */
async function opiszNaglowek(
  kanal: Kanal,
  nazwaSrodowiska: string,
  korzen: ParentNode,
): Promise<void> {
  const naglowek = korzen.querySelector('.sta-kom-naglowek');
  if (naglowek === null) return;
  const wynik = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const kanalModelu = wynik.udany ? wynik.wynik?.channels[0] : undefined;
  const wartosci: Record<string, string> = {
    'Środowisko': nazwaSrodowiska,
    Model: kanalModelu?.model ?? kanalModelu?.name ?? '',
  };
  for (const pole of [...naglowek.querySelectorAll('.sta-kom-pole')]) {
    const podpis = Object.keys(wartosci).find(
      (kandydat) => pole.textContent?.startsWith(kandydat) === true,
    );
    wpiszPole([pole], podpis ?? '', podpis === undefined ? '' : wartosci[podpis]);
  }
}

/** Wpisuje wartość w pole nagłówka rozpoznane po jego podpisie; wartość pusta zdejmuje całe pole, bo pole bez wartości niczego nie mówi. */
export function wpiszPole(pola: Element[], podpis: string, wartosc: string): void {
  const pole = podpis === ''
    ? pola[0]
    : pola.find((kandydat) => kandydat.textContent?.startsWith(podpis) === true);
  if (pole === undefined) return;
  if (wartosc === '') {
    pole.remove();
    return;
  }
  /* Wartość stoi w prototypie raz w `.dane`, raz w samym wyróżnieniu — pole
     środowiska nie niesie klasy danych, a wpisać się musi tak samo. */
  const dane = pole.querySelector('.dane') ?? pole.querySelector('b');
  if (dane !== null) dane.textContent = wartosc;
}

/** Wpisuje wartość w pierwszy niepusty węzeł tekstowy, zostawiając ikonę i przyciski znacznika; fałsz znaczy węzeł bez miejsca na tekst. */
export function wpiszTekst(wezel: Element, tekst: string): boolean {
  for (const dziecko of wezel.childNodes) {
    if (dziecko.nodeType !== Node.TEXT_NODE) continue;
    if ((dziecko.nodeValue ?? '').trim() === '') continue;
    dziecko.nodeValue = tekst;
    return true;
  }
  return false;
}

/**
 * Wpisuje wartość w węzeł tekstowy znacznika, zostawiając ikonę i przyciski.
 * Wartość pusta zostawia węzeł pusty, nie zdejmuje go: pole, które kontrakt
 * zna, czeka na wartość następną — a wpis pusty zachowuje węzeł tekstowy,
 * więc wartość następna ma dokąd wejść.
 */
export function opiszWezel(wezel: Element | null, wartosc: string): void {
  if (wezel === null) return;
  const teksty = [...wezel.childNodes].filter((dziecko) => dziecko.nodeType === Node.TEXT_NODE);
  const cel = teksty.find((dziecko) => (dziecko.nodeValue ?? '').trim() !== '')
    ?? teksty[teksty.length - 1];
  if (cel === undefined) return;
  cel.nodeValue = wartosc;
}

/** Klon węzła wzorcowego, odporny na jego brak w znaczniku. */
export function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}

/** Godzina znacznika czasu rdzenia w zapisie, którego używa prototyp — godziny i minuty. */
export function godzina(znacznik: number): string {
  return new Date(znacznik).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
}

/** Data znacznika czasu rdzenia w zapisie dziennym prototypu. */
export function data(znacznik: number): string {
  return new Date(znacznik).toLocaleDateString('pl-PL');
}

/**
 * Zdejmuje z okna komunikacji podpisy przykładowe wspólne wszystkim oknom
 * modułów. Sterowanie zostaje: wybór modelu, nakładu i urządzenia dźwięku
 * należy do okna Właściciela, a nie do tego, co rdzeń dziś obsługuje.
 */
export function zdejmijSterowanieWspolne(cialo: HTMLElement): void {
  // Podpis piątego trybu opisuje domyślność, której rejestr uprawnień nie zna;
  // to treść przykładowa, nie element sterowania.
  cialo.querySelector('#pop-tryb-upr .sta-popover-wiersz small')?.remove();
}

/**
 * Uzgadnia znaczniki przełączników paneli ze stanem paneli w znaczniku.
 * Panel, którego rdzeń jeszcze nie wypełnia, ZOSTAJE — okno jest kompozycją
 * Właściciela, a nie wyborem tego, co dziś ma pokrycie; panel bez treści stoi
 * pusty i tym mówi prawdę, usunięty kłamałby o układzie okna.
 */
export function uzgodnijPrzelacznikiPaneli(cialo: HTMLElement, _bezPokrycia: string[]): void {
  for (const przelacznik of cialo.querySelectorAll<HTMLElement>('[data-panel-toggle]')) {
    const panel = cialo.querySelector(`#${przelacznik.dataset.panelToggle ?? ''}`);
    przelacznik.setAttribute('aria-checked', String(panel?.hasAttribute('hidden') === false));
  }
}

/** Zdejmuje z okna komunikacji treść przykładową wspólną wszystkim oknom modułów: znacznik pracy, historię rozmowy i znaczniki menu zawężania. */
export function zdejmijTrescWspolna(cialo: HTMLElement): void {
  // Wykazy tracą wiersze przykładowe; ich pojemniki, nagłówki i sterowanie zostają.
  cialo.querySelector('.sta-kom-historia')?.replaceChildren();
  cialo.querySelector('.sta-kom-monitor .sta-kom-monitor-tresc')?.replaceChildren();
  for (const znacznik of cialo.querySelectorAll('#menu-filtr .sta-menu-poz[aria-checked]')) {
    znacznik.removeAttribute('aria-checked');
  }
}

/** Miejsca węzłów zdjętych dla braku wartości; bez nich węzeł nie miałby dokąd wrócić. */
const MIEJSCA = new WeakMap<Element, { rodzic: Element; przed: Node | null }>();

/**
 * Zdejmuje węzeł albo stawia go z powrotem na jego miejscu. Brak wartości nie
 * jest brakiem pola: kontrakt pole niesie, więc wartość może dojść przy
 * kolejnym odczycie i węzeł wraca tam, skąd zszedł. Węzeł, dla którego pola
 * nie ma wcale, schodzi wprost — bez tej drogi powrotnej.
 */
export function zdejmijAlboPostaw(wezel: Element | null, obecny: boolean): void {
  if (wezel === null) return;
  if (!obecny) {
    if (wezel.parentElement === null) return;
    MIEJSCA.set(wezel, { rodzic: wezel.parentElement, przed: wezel.nextSibling });
    wezel.remove();
    return;
  }
  if (wezel.parentElement !== null) return;
  const miejsce = MIEJSCA.get(wezel);
  miejsce?.rodzic.insertBefore(wezel, miejsce.przed);
}

/** Wpisuje wartość w cały węzeł; wartość pusta zdejmuje węzeł do czasu, gdy rdzeń ją poda. */
export function wpiszAlboZdejmij(wezel: Element | null, wartosc: string): void {
  if (wezel === null) return;
  if (wartosc !== '') wezel.textContent = wartosc;
  zdejmijAlboPostaw(wezel, wartosc !== '');
}

/** Wpisuje wartość w pierwszy niepusty węzeł tekstowy, zostawiając ikonę; wartość pusta zdejmuje węzeł do czasu, gdy rdzeń ją poda. */
export function wpiszTekstAlboZdejmij(wezel: Element | null, wartosc: string): void {
  if (wezel === null) return;
  if (wartosc !== '' && !wpiszTekst(wezel, wartosc)) {
    wezel.remove();
    return;
  }
  zdejmijAlboPostaw(wezel, wartosc !== '');
}
