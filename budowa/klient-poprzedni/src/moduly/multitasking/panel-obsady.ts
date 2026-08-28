import { WindowRole, type Window, type WindowCreateRequest } from '../../../../shared/contract';
import { przyciskAkcji as przycisk, pozycjaWykazu, wykaz } from '../../modele/kontrolki-formularza';
import { utworzSekcjePanelu } from '../../powloka/sekcje-panelu';
import type { ZrodloSekcjiPaneli } from '../../powloka/zrodlo-sekcji-paneli';
import { PRZEDROSTEK, wierszOpisu } from './kontrolki';
import { utworzWykazNadanRol } from './wykaz-nadan-rol';
import {
  czyRolaNadana,
  dopnijWiezWykonawcy,
  przypnijPodKoordynatora,
  utworzPasekWcielenia,
  zdanieOZalozeniu,
} from './nadanie-rol';
import { czyWykonawca, LICZBA_WYKONAWCOW, type Obsada } from './obsada-rol';
import { utworzStanTresci } from './stany-okna';
import { utrwalWskazanieAnalityka, wskazanieAnalitykaZSesji } from './wskazanie-analityka';
import type { StanMultitaskingu } from './stan-multitaskingu';
import type { ZrodloOkien } from './zrodlo-okien';

// Panel obsady zakłada i nadaje role czterem oknom sceny na wzór okna już istniejącego w sesji.

/** Interfejs opisuje panel obsady ról: element osadzany w scenie wraz z funkcją odświeżania okien i wskazania analityka. */
export interface PanelObsady {
  element: HTMLElement;
  /** Odczytuje okna sesji i wskazanie analityka. */
  odswiez(): void;
}

export interface OpcjeObsady {
  zrodlo: ZrodloOkien;
  stan: StanMultitaskingu;
  /** Wzorzec powłoki układu sekcji — panel podaje własne sekcje, nie mechanikę. */
  sekcje: ZrodloSekcjiPaneli;
}

/** Stała podaje klucz wskazania okna analityka na poziomie sesji, wspólny dla zapisu i odczytu tego wskazania. */
export const KLUCZ_ANALITYKA = 'multitasking.analityk';

/** Identyfikator panelu w oknie koordynatora — adres dla komend sekcji. Napis stały, bo układ ma przeżyć zamknięcie karty. */
const PANEL_OBSADY = 'obsada-rol';

export function utworzPanelObsady(opcje: OpcjeObsady): PanelObsady {
  const { zrodlo, stan, sekcje } = opcje;
  let wzorzec: Window | null = null;

  const tresci = utworzStanTresci();
  const potwierdz = (zdanie: string, udane: boolean): void => tresci.potwierdzenie(zdanie, udane);

  const { element, opis, obce, zalozKoordynatora, zalozWykonawce, zalozAnalityka } =
    zlozPowierzchnieObsady();

  // Wcielenie roli jest jedynym wołaczem tej komendy; dotyczy okna wskazanego jako koordynator.
  const wcielenie = utworzPasekWcielenia(zrodlo, stan, potwierdz);

  // Wykaz nadań pyta rdzeń o role i je zdejmuje; po zdjęciu scena czyta okna od nowa zwrotnie.
  const nadania = utworzWykazNadanRol({
    zrodlo,
    stan,
    poZdjeciuRoli: () => {
      void odczytaj();
    },
  });

  // Sekcje panelu idą przez wzorzec powłoki: kolejność trzyma rdzeń, moduł podaje własne sekcje.
  const uklad = utworzSekcjePanelu({
    zrodlo: sekcje,
    panelId: PANEL_OBSADY,
    przedrostek: PRZEDROSTEK,
    meldunek: potwierdz,
    sekcje: [
      { id: 'nadania', tytul: 'Nadania ról w rdzeniu', tresc: nadania.element },
      { id: 'wcielenie', tytul: 'Wcielenie koordynatora', tresc: wcielenie.element },
      { id: 'obce', tytul: 'Okna wykonawcze spoza tej obsady', tresc: obce },
    ],
  });
  // Stan treści stoi pod sekcjami, bo melduje o nich wszystkich naraz, także o odmowie zapisu układu.
  element.append(uklad.element, tresci.element);

  zalozKoordynatora.addEventListener('click', () => {
    void zaloz(WindowRole.Coordinator, 'Coordinator Chat');
  });
  zalozWykonawce.addEventListener('click', () => {
    void zaloz(WindowRole.Executor, `Executor ${stan.obsada().wykonawcy.length + 1}`);
  });
  zalozAnalityka.addEventListener('click', () => {
    void zaloz(WindowRole.Standalone, 'Results Analyzer');
  });

  /** Zakłada okno roli na wzór okna istniejącego w sesji. */
  async function zaloz(rola: WindowRole, tytul: string): Promise<void> {
    if (wzorzec === null) {
      potwierdz('Sesja nie ma ani jednego okna — nie ma wzorca kanału modelu.', false);
      return;
    }
    const koordynator = stan.obsada().koordynator;
    if (rola === WindowRole.Executor && koordynator === null) {
      potwierdz('Wykonawca bez koordynatora nie ma dokąd zgłaszać końca tury.', false);
      return;
    }
    tresci.ladowanie(`Zakładanie okna roli ${rola}…`);
    const wynik = await zrodlo.zalozOkno(
      zadanieZalozeniaOkna(wzorzec, stan.sesja(), rola, tytul, koordynator),
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      tresci.blad('Rdzeń odmówił założenia okna roli.', wynik.blad);
      return;
    }
    const zalozone = wynik.wynik;
    const nadana = czyRolaNadana(zalozone, rola, koordynator);
    // Wskazanie analityka utrwala się tylko, gdy rola naprawdę powstała, nie przy innej roli.
    const oUtrwaleniu =
      nadana && rola === WindowRole.Standalone
        ? ` ${await utrwalWskazanieAnalityka(zrodlo, stan, zalozone.id)}`
        : '';
    await odczytaj();
    // Wykonawca wymaga dwóch kroków; przesłanką jest obsada odczytana po założeniu, nie odpowiedź żądania.
    const oWiezi = await dopnijWiezJesliZginela(zalozone, rola, koordynator);
    potwierdz(
      `${zdanieOZalozeniu(zalozone, rola, koordynator)}${oUtrwaleniu}${oWiezi.zdanie}`,
      nadana || oWiezi.naprawione,
    );
  }

  // Obejście dopina więź wykonawcy, gdy po odczycie okien nie ma go w obsadzie; milczy w innych rolach.
  async function dopnijWiezJesliZginela(
    zalozone: Window,
    rola: WindowRole,
    koordynator: Window | null,
  ): Promise<{ naprawione: boolean; zdanie: string }> {
    if (rola !== WindowRole.Executor || koordynator === null) {
      return { naprawione: false, zdanie: '' };
    }
    if (czyWykonawca(stan.obsada(), zalozone.id)) return { naprawione: false, zdanie: '' };
    const dopiecie = await dopnijWiezWykonawcy(zrodlo, stan, zalozone);
    if (dopiecie.udane) await odczytaj();
    return { naprawione: dopiecie.udane, zdanie: ` ${dopiecie.zdanie}` };
  }

  /** Przypina cudze okno wykonawcze pod koordynatora tej obsady. */
  async function przypnij(okno: Window): Promise<void> {
    if (await przypnijPodKoordynatora(zrodlo, stan, okno, potwierdz)) await odczytaj();
  }

  /** Odczyt okien sesji i wskazania analityka; jedno źródło obsady. */
  async function odczytaj(): Promise<void> {
    const wskazanie = await wskazanieAnalitykaZSesji(zrodlo, stan.sesja());
    if (wskazanie !== '') stan.ustawAnalityka(wskazanie);

    const okna = await zrodlo.okna({ sessionId: stan.sesja() });
    if (!okna.udany || okna.wynik === undefined) {
      tresci.blad('Odczyt okien sesji odmówiony.', okna.blad);
      return;
    }
    wzorzec = okna.wynik[0] ?? null;
    stan.ustawOkna(okna.wynik);
    tresci.pusto('');
    // Nadania czyta się z rdzenia osobno od okien: to dwa pytania o dwie różne rzeczy, nie jeden rejestr.
    nadania.odswiez();
  }

  function odswiez(): void {
    const obsada = stan.obsada();
    opis.replaceChildren(...opisObsady(obsada));
    // Blokady chronią przed skutkiem nieodwracalnym: drugie okno tej roli powstałoby w rdzeniu naprawdę.
    wygasZPowodem(
      zalozKoordynatora,
      obsada.koordynator !== null,
      `Obsada ma już koordynatora (${obsada.koordynator?.id ?? ''}). Drugie okno tej roli powstałoby w rdzeniu naprawdę i zostało w sesji bez użytku.`,
    );
    wygasZPowodem(
      zalozWykonawce,
      obsada.wykonawcy.length >= LICZBA_WYKONAWCOW,
      `Obsada ma komplet ${LICZBA_WYKONAWCOW} wykonawców. Trzecie okno wykonawcy powstałoby w rdzeniu naprawdę, a scena pokazuje dwa.`,
    );
    wygasZPowodem(
      zalozAnalityka,
      obsada.analityk !== null,
      `Analityk jest już wskazany (${obsada.analityk?.id ?? ''}). Drugie okno analityka powstałoby w rdzeniu naprawdę i zostało bez wskazania.`,
    );
    wcielenie.odswiez();
    // Układ sekcji jest własnością okna koordynatora; bez niego adres pusty, sekcje stoją miejscowo.
    uklad.ustawOkno(obsada.koordynator?.id ?? '');
    obce.replaceChildren(
      ...pozycjeOkienObcych(obsada.obce, (okno) => {
        void przypnij(okno);
      }),
    );
  }

  stan.naZmiane(odswiez);
  odswiez();

  return {
    element,
    odswiez: () => {
      void odczytaj();
    },
  };
}

/** Interfejs zestawia kontrolki panelu obsady wraz z miejscami na opis obsady i wykaz okien wykonawczych spoza niej. */
interface PowierzchniaObsady {
  element: HTMLElement;
  opis: HTMLElement;
  obce: HTMLElement;
  zalozKoordynatora: HTMLButtonElement;
  zalozWykonawce: HTMLButtonElement;
  zalozAnalityka: HTMLButtonElement;
}

/** Funkcja składa pasek zakładania, opis obsady i wykaz okien obcych jako czystą konstrukcję bez własnych nasłuchów. */
function zlozPowierzchnieObsady(): PowierzchniaObsady {
  const zalozKoordynatora = przycisk('Załóż koordynatora', 'dn-btn dn-btn--sm');
  const zalozWykonawce = przycisk('Załóż wykonawcę', 'dn-btn dn-btn--sm');
  const zalozAnalityka = przycisk('Załóż analityka', 'dn-btn dn-btn--sm');

  const pasek = document.createElement('div');
  pasek.className = 'dm-obsada__pasek';
  pasek.append(zalozKoordynatora, zalozWykonawce, zalozAnalityka);

  const opis = document.createElement('div');
  opis.className = 'dm-obsada__opis';

  const obce = wykaz('Okna wykonawcze spoza tej obsady', 'dm-wykaz');

  const element = document.createElement('div');
  element.className = 'dm-obsada';
  element.setAttribute('aria-label', 'Obsada ról MultitaskingAI');
  element.append(pasek, opis);

  return { element, opis, obce, zalozKoordynatora, zalozWykonawce, zalozAnalityka };
}

/** Funkcja wygasza kontrolkę wraz z powodem blokady, zapisanym trzema drogami: dla wskaźnika, czytnika ekranu i sprawdzianu. */
function wygasZPowodem(kontrolka: HTMLButtonElement, wygaszona: boolean, powod: string): void {
  kontrolka.disabled = wygaszona;
  if (!wygaszona) {
    kontrolka.removeAttribute('title');
    kontrolka.removeAttribute('aria-description');
    delete kontrolka.dataset['powodBlokady'];
    return;
  }
  kontrolka.title = powod;
  kontrolka.setAttribute('aria-description', powod);
  kontrolka.dataset['powodBlokady'] = 'tak';
}

/**
 * Żądanie `window.create` złożone z pól okna wzorcowego.
 *
 * Czysta funkcja danych: wzorzec, sesję, rolę i koordynatora bierze wprost.
 * Przypięcie do koordynatora dokłada się wyłącznie wykonawcy — inna rola nie
 * należy do pętli.
 */
function zadanieZalozeniaOkna(
  wzorzec: Window,
  sesja: string,
  rola: WindowRole,
  tytul: string,
  koordynator: Window | null,
): WindowCreateRequest {
  return {
    sessionId: sesja,
    moduleId: wzorzec.moduleId,
    modelChannelId: wzorzec.modelChannelId,
    workingDirs: wzorzec.workingDirs,
    executionEnv: wzorzec.executionEnv,
    permissionMode: wzorzec.permissionMode,
    windowRole: rola,
    title: tytul,
    ...(rola === WindowRole.Executor && koordynator !== null
      ? { coordinatorWindowId: koordynator.id }
      : {}),
  };
}

/** Funkcja zamienia obsadę na trzy wiersze opisu: koordynatora, komplet wykonawców i wskazanego analityka. */
function opisObsady(obsada: Obsada): HTMLElement[] {
  return [
    wierszOpisu('Koordynator', obsada.koordynator?.id ?? 'nie założony'),
    wierszOpisu('Wykonawcy', `${obsada.wykonawcy.length} z ${LICZBA_WYKONAWCOW}`),
    wierszOpisu('Analityk', obsada.analityk?.id ?? 'nie wskazany'),
  ];
}

/**
 * Wykaz okien wykonawczych spoza obsady.
 *
 * Wyjęty przez wywołanie zwrotne: sam wykaz jest złożeniem pozycji, a sięgnięcie
 * po `window.update` przy kliknięciu zostaje po stronie panelu.
 */
function pozycjeOkienObcych(
  obce: readonly Window[],
  naPrzypiecie: (okno: Window) => void,
): HTMLElement[] {
  return obce.map((okno) => {
    const { element: wiersz, akcje } = pozycjaWykazu(
      okno.id,
      `koordynator ${okno.coordinatorWindowId ?? 'brak'}`,
      PRZEDROSTEK,
    );
    const kontrolka = przycisk('Przypnij tutaj', 'dn-btn dn-btn--sm dn-btn--zarys');
    kontrolka.addEventListener('click', () => {
      naPrzypiecie(okno);
    });
    akcje.append(kontrolka);
    return wiersz;
  });
}
