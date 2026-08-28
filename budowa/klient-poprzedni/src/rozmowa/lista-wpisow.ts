import { NAPISY } from './etykiety-rozmowy';
import { utworzWidokWpisu, type WidokWpisu } from './widok-wpisu';
import {
  czyWpisWZapisie,
  WIDOK_ZAPISU_DOMYSLNY,
  type WidokZapisu,
} from './widok-zapisu';
import type { WpisRozmowy } from './wpis-rozmowy';

/** Interfejs listy wpisów historii okna udostępnia element montowany w widoku, zarządzanie wpisami oraz przełączanie trybu widoku transkryptu. */
export interface ListaWpisow {
  /** Element montowany w widoku rozmowy. */
  element: HTMLElement;
  /** Zakłada wpis albo odświeża istniejący — po kluczu wpisu. */
  pokaz(wpis: WpisRozmowy): void;
  /** Liczba wpisów w liście. */
  liczba(): number;
  /** Przestawia całą listę na inny tryb widoku transkryptu. */
  ustawWidokZapisu(widok: WidokZapisu): void;
  /** Tryb widoku transkryptu, w którym lista pracuje w tej chwili. */
  widokZapisu(): WidokZapisu;
  /** Zdejmuje wszystkie wpisy z listy przy czyszczeniu rozmowy ulotnej. */
  wyczysc(): void;
}

/** Funkcja tworzy historię rozmowy jako największy element okna komunikacji, rozpoznając wpisy po kluczu i utrzymując wybrany tryb widoku transkryptu niezależnie od pamięci wpisów. */
export function utworzListeWpisow(): ListaWpisow {
  const element = document.createElement('section');
  // Odstęp, wyściółkę i tło obszaru niesie klasa z biblioteki wspólnej.
  element.className = 'dn-rozmowa-historia dc-historia';
  element.setAttribute('aria-label', NAPISY.historia);
  element.setAttribute('aria-live', 'polite');

  const pusty = utworzPustyStan();
  element.append(pusty.element);

  const widoki = new Map<string, WidokWpisu>();
  const wpisy = new Map<string, WpisRozmowy>();
  let widok: WidokZapisu = WIDOK_ZAPISU_DOMYSLNY;

  /** Czy dolna krawędź listy jest w polu widzenia. */
  function naKoncu(): boolean {
    const zapas = element.scrollHeight - element.scrollTop - element.clientHeight;
    return zapas < 64;
  }

  function pokaz(wpis: WpisRozmowy): void {
    const znany = widoki.get(wpis.id);
    const bylNaKoncu = naKoncu();
    wpisy.set(wpis.id, wpis);

    if (znany !== undefined) {
      znany.aktualizuj(wpis);
    } else {
      const nowy = utworzWidokWpisu(wpis, widok);
      widoki.set(wpis.id, nowy);
      element.append(nowy.element);
    }
    odswiezZapis();

    if (bylNaKoncu) element.scrollTop = element.scrollHeight;
  }

  /** Przestawia widoki na tryb bieżący, ukrywając wpisy odrzucone przez tryb zamiast ich usuwać. */
  function odswiezZapis(): void {
    let widoczne = 0;
    for (const [id, wpisu] of widoki) {
      const dane = wpisy.get(id);
      const wZapisie = dane === undefined || czyWpisWZapisie(widok, dane);
      wpisu.element.hidden = !wZapisie;
      if (wZapisie) widoczne += 1;
    }
    pusty.ustaw(widok, widoki.size);
    pusty.element.hidden = widoczne > 0;
  }

  function ustawWidokZapisu(nowy: WidokZapisu): void {
    widok = nowy;
    for (const wpisu of widoki.values()) wpisu.ustawWidokZapisu(nowy);
    odswiezZapis();
  }

  function wyczysc(): void {
    for (const wpisu of widoki.values()) wpisu.element.remove();
    widoki.clear();
    wpisy.clear();
    odswiezZapis();
  }

  return {
    element,
    pokaz,
    liczba: () => widoki.size,
    ustawWidokZapisu,
    widokZapisu: () => widok,
    wyczysc,
  };
}

/** Interfejs stanu pustej historii opisuje komunikat wyświetlany zamiast pustej powierzchni listy wpisów. */
interface PustyStan {
  /** Element montowany na wierzchu listy. */
  element: HTMLElement;
  /** Dobiera zdania do przyczyny pustki: brak wątku albo tryb bez treści. */
  ustaw(widok: WidokZapisu, wpisow: number): void;
}

/** Funkcja tworzy stan pusty historii z osobnym zdaniem dla wątku, który się nie zaczął, i dla trybu Streszczenie bez zamkniętych tur. */
function utworzPustyStan(): PustyStan {
  const element = document.createElement('div');
  element.className = 'dn-pusty-stan dc-historia__pusto';

  const tytul = document.createElement('p');
  tytul.className = 'dn-pusty-stan-tytul';

  const opis = document.createElement('p');
  opis.className = 'dn-pusty-stan-opis';

  element.append(tytul, opis);

  function ustaw(widok: WidokZapisu, wpisow: number): void {
    const streszczenieBezTur = widok === 'streszczenie' && wpisow > 0;
    tytul.textContent = streszczenieBezTur ? NAPISY.streszczeniePustoTytul : NAPISY.pustoTytul;
    opis.textContent = streszczenieBezTur ? NAPISY.streszczeniePustoOpis : NAPISY.pustoOpis;
  }

  ustaw(WIDOK_ZAPISU_DOMYSLNY, 0);
  return { element, ustaw };
}
