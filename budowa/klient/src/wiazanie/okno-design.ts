/**
 * Wiązanie wnętrza okna modułu Design z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * design.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla których
 * kontrakt nie ma pola, zdejmuje.
 */
import {
  Command,
  EventType,
  type DesignAsset,
  type DesignBoard,
  type DesignPrompt,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Panele prototypu, którym rodzina design.* nie oddaje ani jednego pola; schodzą razem z przełącznikami menu sesji. */
const PANELE_BEZ_POKRYCIA = ['panel-preview', 'panel-plan', 'panel-pliki', 'panel-zadania'];

/** Wzory zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone; klon zachowuje ikonę i klasy nadane przez bibliotekę. */
interface Wzory {
  pozycja: HTMLElement | null;
  miniatura: HTMLElement | null;
  gwiazdka: SVGElement | null;
  plakietka: HTMLElement | null;
  warstwa: SVGGElement | null;
  /* Odstęp podpisu od narożnika prostokąta warstwy, zdjęty ze wzoru kanwy:
     kompozycja nie niesie położenia napisu, niesie je tylko warstwa. */
  odsunieciePodpisu: { x: number; y: number };
}

/**
 * Znak pasa czynności zdjęty z okna wraz z miejscem, w które wraca. Znak bez
 * wartości schodzi, bo znak pusty niczego nie mówi; wartość podana później —
 * kompozycja założona już po otwarciu okna — stawia ten sam znak z powrotem.
 */
interface Znak {
  wezel: HTMLElement;
  rodzic: Element;
  miejsce: number;
  wzorTekstu: string;
}

/** Znaki pasa czynności okna Design: nazwa kompozycji wiodącej i liczba jej wersji. */
interface Znaki {
  kompozycja: Znak | null;
  wersja: Znak | null;
}

/**
 * Wiąże wnętrze okna modułu Design. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazDesign(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  const znaki = zdejmijZnakiCzynnosci(cialo);
  zdejmijSterowanieBezPokrycia(cialo);

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Powtarza odczyt kompozycji i zasobów; wołane przy wejściu i po każdej zmianie ogłoszonej przez rdzeń. */
  const odswiez = (): void => {
    void wypelnijStanowisko(kanal, idOkna, cialo, wzory, znaki);
  };

  zwiazPromptBuilder(kanal, idOkna, cialo);
  zwiazZdarzenia(kanal, idOkna, odswiez);

  odswiez();
  void wypelnijPrompt(kanal, idOkna, cialo);
}

/** Zdejmuje wzory pozycji i warstw, zanim wykazy zostaną wyczyszczone; wzór traci oznaczenia bez pokrycia w kontrakcie. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const pozycja = sklonuj(cialo.querySelector<HTMLElement>('.dn-szyna-modulu-lista .pt-pozycja'));
  // Znak pracy nad kompozycją nie ma pola: kompozycja niesie nazwę, warstwy
  // i czas zmiany, nic o toczącej się robocie.
  pozycja?.querySelector('.pt-tetno')?.remove();
  pozycja?.removeAttribute('aria-current');

  const miniatura = sklonuj(cialo.querySelector<HTMLElement>('#panel-assets .dg-miniatura'));
  // Znak ulubionego schodzi ze wzoru i staje osobno: nadaje go dopiero pole
  // `favorite` zasobu, a wzór ma być miniaturą zasobu dowolnego.
  const znakUlubionego = miniatura?.querySelector<SVGElement>('.podpis svg') ?? null;
  const gwiazdka = sklonuj(znakUlubionego);
  znakUlubionego?.remove();
  // Rysunek miniatury jest treścią przykładową: zasób niesie odnośnik do
  // treści, a nie kształty, którymi dałoby się go tu narysować.
  miniatura?.querySelector(':scope > svg')?.remove();
  miniatura?.removeAttribute('aria-selected');

  return {
    pozycja,
    miniatura,
    gwiazdka,
    plakietka: sklonuj(cialo.querySelector<HTMLElement>('#panel-assets .podpis .dn-plakietka')),
    ...wzorWarstwy(cialo),
  };
}

/** Zdejmuje z kanwy wzór warstwy wraz z odstępem podpisu; ozdoby wzoru schodzą, bo kompozycja niesie prostokąt i nazwę, nie kształty. */
function wzorWarstwy(cialo: HTMLElement): Pick<Wzory, 'warstwa' | 'odsunieciePodpisu'> {
  const warstwa = sklonuj(cialo.querySelector<SVGGElement>('#panel-board .dg-kanwa svg g'));
  if (warstwa === null) return { warstwa: null, odsunieciePodpisu: { x: 0, y: 0 } };
  for (const ozdoba of [...warstwa.querySelectorAll('path, circle')]) ozdoba.remove();
  const prostokat = warstwa.querySelector('rect');
  const podpis = warstwa.querySelector('text');
  const odsuniecie = {
    x: liczba(podpis, 'x') - liczba(prostokat, 'x'),
    y: liczba(podpis, 'y') - liczba(prostokat, 'y'),
  };
  return { warstwa, odsunieciePodpisu: odsuniecie };
}

/** Wartość liczbowa atrybutu węzła kanwy; brak węzła albo atrybutu daje zero. */
function liczba(wezel: Element | null, atrybut: string): number {
  return Number.parseFloat(wezel?.getAttribute(atrybut) ?? '0') || 0;
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz rozmów i historia wiszą na
 * rodzinach spoza design.*; przypięty zasób, znak pracy, monitor generowania
 * i miara różnicy nie mają pola w całym kontrakcie.
 */
function zdejmijTrescPrzykladowa(cialo: HTMLElement): void {
  cialo.querySelector('.sta-kom-historia')?.replaceChildren();
  for (const znacznik of cialo.querySelectorAll('#menu-filtr .sta-menu-poz[aria-checked]')) {
    znacznik.removeAttribute('aria-checked');
  }
  cialo.querySelector('.sta-kom-naglowek .sta-kom-stan')?.remove();
  cialo.querySelector('.sta-kom-monitor')?.remove();
  cialo.querySelector('.sta-kom-kontekst')?.replaceChildren();
  // Znacznik przyciągania kanwy: siatkę ustala design.grid.set, lecz kontrakt
  // nie oddaje jej z powrotem żadną komendą odczytu.
  cialo.querySelector('#panel-board .sta-okno-znacznik')?.remove();
  zdejmijPasyKanwy(cialo);
  zdejmijPasyZasobow(cialo);
  oznaczChipyCzynnosci(cialo);
}

/**
 * Zdejmuje pasy narzędzi kanwy. Warstwy, wyrównanie, siatka, szablon układu,
 * wersje i adnotacja nie mają w oknie węzła, który przyjąłby ich wynik; zoom,
 * zaznaczenie i liczba kursorów nie mają pola w kontrakcie.
 */
function zdejmijPasyKanwy(cialo: HTMLElement): void {
  for (const pas of cialo.querySelectorAll('#panel-board .dg-narzedzia')) pas.remove();
}

/**
 * Zdejmuje pasy Assets Panel: wyszukiwanie po nazwie i etykiecie nie ma pola
 * w design.asset.list, a przejścia do innych modułów nie mają komendy.
 */
function zdejmijPasyZasobow(cialo: HTMLElement): void {
  for (const pas of cialo.querySelectorAll('#panel-assets .dg-narzedzia')) pas.remove();
  const chipWidoku = [...cialo.querySelectorAll('#panel-assets .sta-okno-akcje .sta-chip')];
  for (const chip of chipWidoku) chip.remove();
}

/**
 * Porządkuje pas czynności: znak zasięgu wykonania i miara różnicy schodzą,
 * a przycisk przejścia do Prompt Buildera schodzi razem z nimi, bo rodzina
 * design.* nie niesie ani zasięgu wykonania, ani miary zmian kompozycji.
 */
function oznaczChipyCzynnosci(cialo: HTMLElement): void {
  const chipy = [...cialo.querySelectorAll<HTMLElement>('.sta-kontekst-akcji > .sta-chip')];
  chipy[0]?.remove();
  cialo.querySelector('.sta-kontekst-akcji .sta-chip--diff')?.remove();
  cialo.querySelector('.sta-kontekst-akcji button.dn-btn')?.remove();
}

/** Zdejmuje z pasa czynności znaki nazwy i wersji kompozycji wraz z miejscem, w które wracają, gdy rdzeń poda dla nich wartość. */
function zdejmijZnakiCzynnosci(cialo: HTMLElement): Znaki {
  const chipy = [...cialo.querySelectorAll<HTMLElement>('.sta-kontekst-akcji > .sta-chip')];
  return { kompozycja: zdejmijZnak(chipy[1] ?? null), wersja: zdejmijZnak(chipy[2] ?? null) };
}

/**
 * Zdejmuje sterowanie, którego rodzina design.* nie obsługuje: panele bez
 * pokrycia wraz z ich przełącznikami, wybór modelu i nakładu, urządzenia
 * wejścia dźwięku, wykaz szablonów promptu oraz konektory i wtyczki menu
 * dodawania.
 */
function zdejmijSterowanieBezPokrycia(cialo: HTMLElement): void {
  for (const nazwa of PANELE_BEZ_POKRYCIA) {
    cialo.querySelector(`#${nazwa}`)?.remove();
    cialo.querySelector(`[data-panel-toggle="${nazwa}"]`)?.remove();
  }
  for (const przelacznik of cialo.querySelectorAll<HTMLElement>('[data-panel-toggle]')) {
    const panel = cialo.querySelector(`#${przelacznik.dataset.panelToggle ?? ''}`);
    przelacznik.setAttribute('aria-checked', String(panel?.hasAttribute('hidden') === false));
  }
  cialo.querySelector('#pop-mik')?.closest('.sta-nrz')?.remove();
  cialo.querySelector('#pop-model')?.closest('.sta-nrz')?.remove();
  cialo.querySelector('#pop-wysilek')?.closest('.sta-nrz')?.remove();
  // Podpis piątego trybu opisuje domyślność, której rejestr uprawnień nie zna.
  cialo.querySelector('#pop-tryb-upr .sta-popover-wiersz small')?.remove();
  // Wykaz szablonów promptu ma komendę, lecz nie ma w oknie węzła, który by go
  // przyjął; przycisk bez miejsca na wynik prowadzi donikąd.
  cialo.querySelector('#panel-prompt .sta-okno-akcje .sta-chip')?.remove();
  zdejmijKonektory(cialo);
}

/** Zostawia w menu dodawania sam tytuł i trzy pozycje zasobu; konektory i wtyczki należą do rodzin spoza design.*. */
function zdejmijKonektory(cialo: HTMLElement): void {
  const menu = cialo.querySelector('#pop-plus');
  if (menu === null) return;
  for (const [numer, pozycja] of [...menu.children].entries()) {
    if (numer > 3) pozycja.remove();
  }
}

/** Pyta rdzeń o kompozycje i zasoby okna, po czym wypełnia nimi szynę, kanwę, znaki pasa czynności i Assets Panel. */
async function wypelnijStanowisko(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  znaki: Znaki,
): Promise<void> {
  const kompozycje = await wywolaj(kanal, Command.DesignBoardList, { windowId: idOkna });
  const zasoby = await wywolaj(kanal, Command.DesignAssetList, { windowId: idOkna });
  const wykazKompozycji = kompozycje.udany ? (kompozycje.wynik?.boards ?? []) : [];
  const wykazZasobow = zasoby.udany ? (zasoby.wynik?.assets ?? []) : [];

  const wiodaca = wykazKompozycji[0];
  wypelnijSzyne(cialo, wzory, wykazKompozycji);
  opiszKompozycje(cialo, znaki, wiodaca);
  wypelnijKanwe(cialo, wzory, wiodaca, wykazZasobow);
  wypelnijZasoby(cialo, wzory, wykazZasobow, zasoby.wynik?.total ?? wykazZasobow.length);
  await opiszWersje(kanal, znaki, wiodaca?.id ?? '');
}

/** Stawia w szynie po jednej pozycji na kompozycję okna; nagłówki grup schodzą, bo kompozycji kontrakt nie grupuje. */
function wypelnijSzyne(cialo: HTMLElement, wzory: Wzory, kompozycje: DesignBoard[]): void {
  const lista = cialo.querySelector<HTMLElement>('.dn-szyna-modulu-lista');
  if (lista === null) return;
  lista.replaceChildren();
  if (wzory.pozycja === null) return;
  for (const kompozycja of kompozycje) {
    const nazwa = kompozycja.name ?? '';
    if (nazwa === '') continue;
    const pozycja = wzory.pozycja.cloneNode(true) as HTMLElement;
    pozycja.dataset.kompozycja = kompozycja.id;
    const tytul = pozycja.querySelector('.pt-pozycja-tytul');
    if (tytul !== null) tytul.textContent = nazwa;
    lista.appendChild(pozycja);
  }
}

/** Nanosi nazwę kompozycji na tytuł panelu, tytuł okna komunikacji i znak pasa czynności; brak kompozycji zostawia sam człon stały. */
function opiszKompozycje(
  cialo: HTMLElement,
  znaki: Znaki,
  kompozycja: DesignBoard | undefined,
): void {
  const nazwa = kompozycja?.name ?? '';
  opiszTytul(cialo.querySelector('#panel-board .sta-okno-tytul b'), nazwa);
  opiszTytul(cialo.querySelector('#okno-czat-1 .sta-okno-tytul b'), nazwa);
  opiszZnak(znaki.kompozycja, nazwa, false);
}

/** Stawia na kanwie po jednym prostokącie na warstwę kompozycji; podpis bierze nazwę zasobu warstwy albo jej adnotację. */
function wypelnijKanwe(
  cialo: HTMLElement,
  wzory: Wzory,
  kompozycja: DesignBoard | undefined,
  zasoby: DesignAsset[],
): void {
  const kanwa = cialo.querySelector('#panel-board .dg-kanwa svg');
  if (kanwa === null) return;
  for (const stojaca of [...kanwa.querySelectorAll('g')]) stojaca.remove();
  if (wzory.warstwa === null) return;
  const nazwyZasobow = new Map(zasoby.map((zasob) => [zasob.id, zasob.name ?? '']));
  for (const warstwa of kompozycja?.layers ?? []) {
    const wezel = wzory.warstwa.cloneNode(true) as SVGGElement;
    const prostokat = wezel.querySelector('rect');
    const podpis = wezel.querySelector('text');
    if (prostokat === null) continue;
    prostokat.setAttribute('x', String(warstwa.x ?? 0));
    prostokat.setAttribute('y', String(warstwa.y ?? 0));
    prostokat.setAttribute('width', String(warstwa.width ?? 0));
    prostokat.setAttribute('height', String(warstwa.height ?? 0));
    const opis = warstwa.note ?? nazwyZasobow.get(warstwa.assetId ?? '') ?? '';
    if (podpis !== null && opis === '') podpis.remove();
    else if (podpis !== null) {
      podpis.setAttribute('x', String((warstwa.x ?? 0) + wzory.odsunieciePodpisu.x));
      podpis.setAttribute('y', String((warstwa.y ?? 0) + wzory.odsunieciePodpisu.y));
      podpis.textContent = opis;
    }
    kanwa.appendChild(wezel);
  }
}

/** Wypełnia Assets Panel miniaturami zasobów okna i nanosi ich liczbę na znacznik panelu. */
function wypelnijZasoby(
  cialo: HTMLElement,
  wzory: Wzory,
  zasoby: DesignAsset[],
  razem: number,
): void {
  /* Znacznik idzie samą liczbą: prototyp niesie ją z rzeczownikiem w dopełniaczu
     mnogim, a odmiana zależy od liczby, której znacznik z góry nie zna. */
  const znacznik = cialo.querySelector('#panel-assets .sta-okno-znacznik');
  if (znacznik !== null) znacznik.textContent = String(razem);
  const siatka = cialo.querySelector<HTMLElement>('#panel-assets .dg-siatka');
  if (siatka === null) return;
  siatka.replaceChildren();
  if (wzory.miniatura === null) return;
  for (const zasob of zasoby) {
    const miniatura = wzory.miniatura.cloneNode(true) as HTMLElement;
    miniatura.dataset.zasob = zasob.id;
    const podpis = miniatura.querySelector<HTMLElement>('.podpis');
    if (podpis === null) continue;
    if (!wpiszTekst(podpis, zasob.name ?? zasob.id)) continue;
    if (zasob.favorite === true && wzory.gwiazdka !== null) {
      podpis.prepend(wzory.gwiazdka.cloneNode(true));
    }
    for (const etykieta of zasob.tags ?? []) {
      if (wzory.plakietka === null) break;
      const plakietka = wzory.plakietka.cloneNode(true) as HTMLElement;
      plakietka.textContent = etykieta;
      podpis.appendChild(plakietka);
    }
    siatka.appendChild(miniatura);
  }
}

/** Nanosi na znak pasa czynności liczbę wersji kompozycji; kompozycja bez wersji zdejmuje znak. */
async function opiszWersje(kanal: Kanal, znaki: Znaki, idKompozycji: string): Promise<void> {
  if (idKompozycji === '') {
    opiszZnak(znaki.wersja, '', true);
    return;
  }
  const wynik = await wywolaj(kanal, Command.DesignBoardVersionList, { boardId: idKompozycji });
  const ile = wynik.udany ? (wynik.wynik?.total ?? 0) : 0;
  opiszZnak(znaki.wersja, ile === 0 ? '' : String(ile), true);
}

/**
 * Wypełnia Prompt Builder ostatnim promptem wydanym w oknie. Styl, paleta,
 * proporcje i liczba wariantów schodzą razem ze swoimi listami wyboru: kontrakt
 * niesie te pola jako łańcuchy swobodne, więc wykaz trzech wartości w każdej
 * z list jest treścią przykładową, nie zbiorem dopuszczalnym.
 */
async function wypelnijPrompt(kanal: Kanal, idOkna: string, cialo: HTMLElement): Promise<void> {
  for (const para of cialo.querySelectorAll('#panel-prompt .dg-dwa')) para.remove();
  const wynik = await wywolaj(kanal, Command.DesignPromptHistoryList, {
    windowId: idOkna,
    limit: 1,
  });
  const prompt = wynik.udany ? wynik.wynik?.prompts[0]?.prompt : undefined;
  opiszPoleTekstowe(cialo, '#dg-temat', prompt?.subject ?? '');
  opiszPoleTekstowe(cialo, '#dg-wykluczenia-negative-prompt', prompt?.exclusions ?? '');
  opiszKreatywnosc(cialo, prompt);
  opiszZiarno(cialo, prompt);
}

/** Wpisuje wartość w pole tekstowe Prompt Buildera; wartość pusta zostawia pole puste, bo prompt bez tej części jest promptem bez niej. */
function opiszPoleTekstowe(cialo: HTMLElement, wybor: string, wartosc: string): void {
  const pole = cialo.querySelector<HTMLInputElement>(`#panel-prompt ${wybor}`);
  if (pole !== null) pole.value = wartosc;
}

/** Nanosi kreatywność promptu na suwak i jego podpis; prompt bez kreatywności zdejmuje całe pole. */
function opiszKreatywnosc(cialo: HTMLElement, prompt: DesignPrompt | undefined): void {
  const suwak = cialo.querySelector<HTMLInputElement>('#dg-kreatywnosc-wiernosc-0-7');
  const pole = suwak?.closest('.dg-pole');
  if (suwak === undefined || suwak === null || pole === null || pole === undefined) return;
  const wartosc = prompt?.creativity;
  if (wartosc === undefined) {
    pole.remove();
    return;
  }
  const gorna = Number.parseFloat(suwak.max) || 1;
  suwak.value = String(Math.round(wartosc * gorna));
  opiszTytul(pole.querySelector('label'), String(wartosc));
}

/**
 * Nanosi ziarno generowania na miarę Prompt Buildera. Długość promptu w znakach
 * nie ma pola w kontrakcie, więc miara idzie samym ziarnem, a prompt bez ziarna
 * zdejmuje ją całą.
 */
function opiszZiarno(cialo: HTMLElement, prompt: DesignPrompt | undefined): void {
  const miara = cialo.querySelector('#panel-prompt .dn-meta');
  if (miara === null) return;
  if (prompt?.seed === undefined) {
    miara.remove();
    return;
  }
  miara.textContent = (miara.textContent ?? '').replace(/^.*seed/u, 'seed')
    .replace(/\d+/u, String(prompt.seed));
}

/** Wiąże przycisk generowania Prompt Buildera z komendą wytworzenia zasobu; temat pusty wstrzymuje czynność, bo prompt bez tematu nie mówi, co ma powstać. */
function zwiazPromptBuilder(kanal: Kanal, idOkna: string, cialo: HTMLElement): void {
  const przycisk = cialo.querySelector('#panel-prompt .dn-wykaz-modulu-poz button');
  przycisk?.addEventListener('click', () => {
    const temat = cialo.querySelector<HTMLInputElement>('#dg-temat')?.value.trim() ?? '';
    if (temat === '') return;
    const wykluczenia =
      cialo.querySelector<HTMLInputElement>('#dg-wykluczenia-negative-prompt')?.value ?? '';
    void wywolaj(kanal, Command.DesignAssetGenerate, {
      windowId: idOkna,
      prompt: { subject: temat, exclusions: wykluczenia === '' ? undefined : wykluczenia },
    });
  });
}

/** Nasłuchuje zmian rdzenia: nowy zasób i zmieniona kompozycja odświeżają stanowisko tego okna. */
function zwiazZdarzenia(kanal: Kanal, idOkna: string, odswiez: () => void): void {
  kanal.naZdarzenie(EventType.DesignAssetChanged, (tresc) => {
    if (tresc.asset.windowId !== idOkna) return;
    odswiez();
  });
  kanal.naZdarzenie(EventType.DesignBoardChanged, (tresc) => {
    if (tresc.board.windowId !== idOkna) return;
    odswiez();
  });
}

/** Zapamiętuje znak wraz z miejscem w pasie czynności; znak, którego prototyp nie niesie, daje pustkę. */
function zdejmijZnak(wezel: HTMLElement | null): Znak | null {
  const rodzic = wezel?.parentElement ?? null;
  if (wezel === null || rodzic === null) return null;
  return {
    wezel,
    rodzic,
    miejsce: [...rodzic.children].indexOf(wezel),
    wzorTekstu: wezel.textContent ?? '',
  };
}

/**
 * Wpisuje wartość w znak pasa czynności i stawia go w jego miejscu, gdy stał
 * zdjęty. Znak liczbowy zachowuje podpis wzoru i podmienia w nim samą liczbę;
 * znak nazwy bierze wartość w całości.
 */
function opiszZnak(znak: Znak | null, wartosc: string, liczbowy: boolean): void {
  if (znak === null) return;
  if (wartosc === '') {
    znak.wezel.remove();
    return;
  }
  znak.wezel.textContent = liczbowy ? znak.wzorTekstu.replace(/\d+/u, wartosc) : wartosc;
  if (znak.wezel.isConnected) return;
  znak.rodzic.insertBefore(znak.wezel, znak.rodzic.children[znak.miejsce] ?? null);
}

/** Zostawia w tytule człon stały prototypu i dokłada do niego wartość rdzenia; wartość pusta zostawia sam człon stały. */
function opiszTytul(wezel: Element | null, wartosc: string): void {
  if (wezel === null) return;
  const staly = (wezel.textContent ?? '').split(' — ')[0] ?? '';
  wezel.textContent = wartosc === '' ? staly : `${staly} — ${wartosc}`;
}

/** Wpisuje wartość w pierwszy niepusty węzeł tekstowy, zostawiając ikonę i przyciski znacznika; fałsz znaczy węzeł bez miejsca na tekst. */
function wpiszTekst(wezel: Element, tekst: string): boolean {
  for (const dziecko of wezel.childNodes) {
    if (dziecko.nodeType !== Node.TEXT_NODE) continue;
    if ((dziecko.nodeValue ?? '').trim() === '') continue;
    dziecko.nodeValue = tekst;
    return true;
  }
  return false;
}

/** Klon węzła wzorcowego, odporny na jego brak w znaczniku. */
function sklonuj<T extends Element>(wezel: T | null): T | null {
  return wezel === null ? null : (wezel.cloneNode(true) as T);
}
