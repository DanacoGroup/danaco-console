import type { ErrorInfo, PanelSection } from '../../../shared/contract';
import { zlozPasZdjetych, zlozSekcje, type CzynnosciSekcji } from './czesci-sekcji';
import {
  przelaczZdjecie,
  przelaczZwiniecie,
  przestawSekcje,
  scalUklad,
  tenSamUklad,
  ukladDomyslny,
  type OpisSekcji,
} from './uklad-sekcji';
import type { ZrodloSekcjiPaneli } from './zrodlo-sekcji-paneli';

/**
 * Sekcje panelu okna — wspólny mechanizm zwijania, przestawiania i zdejmowania
 * sekcji dla wszystkich okien operacyjnych. Moduł podaje własne sekcje oraz
 * identyfikator panelu; pętlę „odczytaj układ · zmień · zapisz · przerysuj"
 * niesie ten plik.
 *
 * Układ adresuje parę (okno, panel), więc bez wskazanego okna zmiana zostaje
 * miejscowa i mówi to wprost. Widok przerysowuje się układem oddanym przez
 * rdzeń w polu `sections`, nie układem wysłanym — rozjazd obu jest wtedy
 * widoczny od razu. Zapis, który dałby układ tożsamy z bieżącym, nie idzie.
 */
export interface SekcjePanelu {
  /** Element osadzany w ciele panelu — sekcje wraz z pasem sekcji zdjętych. */
  element: HTMLElement;
  /**
   * Wskazanie okna, do którego panel należy. Pusty napis znaczy „okna nie ma
   * jeszcze" — układ zostaje miejscowy i nie jest utrwalany.
   */
  ustawOkno(idOkna: string): void;
  /** Ponowny odczyt układu z rdzenia i przerysowanie. */
  odswiez(): void;
}

export interface OpcjeSekcjiPanelu {
  zrodlo: ZrodloSekcjiPaneli;
  /** Identyfikator panelu w obrębie okna — trafia w `panelId` obu komend. */
  panelId: string;
  /** Sekcje, które panel naprawdę buduje, w kolejności domyślnej. */
  sekcje: readonly OpisSekcji[];
  /** Przedrostek klas rodziny modułu (`dm`, `dt`…); wygląd pochodzi z biblioteki. */
  przedrostek: string;
  /** Meldunek widoczny dla Operatora — powodzenie i odmowa idą tą samą drogą. */
  meldunek(zdanie: string, udane: boolean): void;
}

export function utworzSekcjePanelu(opcje: OpcjeSekcjiPanelu): SekcjePanelu {
  const { zrodlo, panelId, sekcje, przedrostek, meldunek } = opcje;
  let idOkna = '';
  let uklad: PanelSection[] = ukladDomyslny(sekcje);

  const element = document.createElement('div');
  element.className = `${przedrostek}-sekcje`;
  element.dataset['panel'] = panelId;
  element.setAttribute('aria-label', `Sekcje panelu ${panelId}`);

  const czynnosci: CzynnosciSekcji = {
    przyZwinieciu: (id) => void zapisz(przelaczZwiniecie(uklad, id), 'zwinięcie sekcji'),
    przyPrzesunieciu: (id, kierunek) =>
      void zapisz(przestawSekcje(uklad, id, kierunek), 'kolejność sekcji'),
    przyZdjeciu: (id) => void zapisz(przelaczZdjecie(uklad, id), 'zdjęcie sekcji z widoku'),
    przyPrzywroceniu: (id) => void zapisz(przelaczZdjecie(uklad, id), 'przywrócenie sekcji'),
  };

  /** Rysuje sekcje w kolejności układu; zdjęte trafiają wyłącznie do pasa. */
  function narysuj(): void {
    const opisy = new Map(sekcje.map((opis) => [opis.id, opis]));
    const widoczne: HTMLElement[] = [];
    const zdjete: Array<{ id: string; tytul: string }> = [];
    for (const stan of uklad) {
      const opis = opisy.get(stan.id);
      if (opis === undefined) continue;
      if (stan.hidden === true) {
        zdjete.push({ id: opis.id, tytul: opis.tytul });
        continue;
      }
      widoczne.push(zlozSekcje(opis, stan, przedrostek, czynnosci));
    }
    element.replaceChildren(
      ...widoczne,
      zlozPasZdjetych(zdjete, przedrostek, czynnosci.przyPrzywroceniu),
    );
  }

  /**
   * Zapis układu zamierzonego i przerysowanie tym, co rdzeń oddał.
   *
   * Bez wskazanego okna nie ma czego adresować: zmiana zostaje wyłącznie
   * w tym widoku, a meldunek nazywa przyczynę.
   */
  async function zapisz(zamierzony: PanelSection[], czynnosc: string): Promise<void> {
    if (tenSamUklad(zamierzony, uklad)) return;
    if (idOkna === '') {
      uklad = zamierzony;
      narysuj();
      meldunek(
        `Układ sekcji trzyma WYŁĄCZNIE ten widok — panel nie ma jeszcze okna, do którego rdzeń mógłby go przypiąć (${czynnosc}).`,
        false,
      );
      return;
    }
    const wynik = await zrodlo.zapiszUklad({ windowId: idOkna, panelId, sections: zamierzony });
    if (!wynik.udany || wynik.wynik === undefined) {
      meldunek(
        `Rdzeń odmówił zapisu układu sekcji (${czynnosc}). ${powodOdmowy(wynik.blad)}`,
        false,
      );
      return;
    }
    uklad = scalUklad(wynik.wynik, sekcje);
    narysuj();
    // Zdanie rozróżnia układ oddany przez rdzeń: tożsamy z zamierzonym znaczy
    // „zapisane", inny — „przyjęte, ale nie tak".
    meldunek(
      tenSamUklad(uklad, zamierzony)
        ? `Rdzeń zapisał układ sekcji panelu ${panelId} okna ${idOkna} (${czynnosc}).`
        : `Rdzeń przyjął zapis układu sekcji, ale oddał układ inny niż zamierzony (${czynnosc}) — panel pokazuje układ oddany.`,
      tenSamUklad(uklad, zamierzony),
    );
  }

  /** Odczyt układu z rdzenia; odmowa zostawia układ domyślny i nazywa powód. */
  async function odczytaj(): Promise<void> {
    if (idOkna === '') {
      uklad = scalUklad(uklad, sekcje);
      narysuj();
      return;
    }
    const wynik = await zrodlo.uklad({ windowId: idOkna, panelId });
    if (!wynik.udany || wynik.wynik === undefined) {
      uklad = ukladDomyslny(sekcje);
      narysuj();
      meldunek(
        `Rdzeń odmówił odczytu układu sekcji panelu ${panelId}. Panel stoi w kolejności domyślnej. ${powodOdmowy(wynik.blad)}`,
        false,
      );
      return;
    }
    uklad = scalUklad(wynik.wynik, sekcje);
    narysuj();
  }

  narysuj();

  return {
    element,
    ustawOkno(nowe) {
      if (nowe === idOkna) return;
      idOkna = nowe;
      void odczytaj();
    },
    odswiez() {
      void odczytaj();
    },
  };
}

/** Treść odmowy wraz z kodem kontraktu — ta sama postać co w modułach. */
function powodOdmowy(blad?: ErrorInfo): string {
  if (blad === undefined) return 'Rdzeń nie podał przyczyny.';
  return `Powód: ${blad.message} (kod ${blad.code}).`;
}
